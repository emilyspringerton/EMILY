#!/usr/bin/env bash
# attach-game-ip.sh -- give the game node pool its one public address (Vultr Reserved IP), idempotently.
#
# Every game port (REDGARDEN/GFD/SHANKPIT/Einhorn UDP, tcp-edge TCP) is a hostPort on the node labelled
# emily/pool=game (see PARENA tools/pod_render_host.c, PARENA_K8S_PROVIDER=vultr). This script:
#   1. finds -- or creates -- the Reserved IP labelled $GAME_IP_LABEL in the game node's region,
#   2. attaches it to the game node (re-run after the node is replaced: VKE upgrades/recycles nodes),
#   3. applies game-ip.yaml so the address is configured on the node's interface.
# Point game DNS (A records) at the printed address once; it survives node replacement as long as this is re-run.
#
# env: VULTR_API_KEY (required), GAME_IP_LABEL (default emily-game). Needs kubectl on the VKE cluster, curl, jq.
set -euo pipefail
HERE="$(cd "$(dirname "$0")" && pwd)"
: "${VULTR_API_KEY:?set VULTR_API_KEY}"
LABEL="${GAME_IP_LABEL:-emily-game}"
API=https://api.vultr.com/v2
vapi() { curl -fsS -H "Authorization: Bearer $VULTR_API_KEY" -H 'Content-Type: application/json' "$@"; }

nodes=$(kubectl get nodes -l emily/pool=game -o json)
count=$(jq '.items | length' <<<"$nodes")
if [ "$count" -ne 1 ]; then
  echo "attach-game-ip: want exactly 1 node labelled emily/pool=game, found $count" >&2
  echo "  (create the VKE node pool with label emily/pool=game and node_quantity 1; see README.md)" >&2
  exit 1
fi
instance=$(jq -r '.items[0].spec.providerID' <<<"$nodes" | sed 's#^vultr://##')
region=$(jq -r '.items[0].metadata.labels["topology.kubernetes.io/region"] // empty' <<<"$nodes")
[ -n "$instance" ] && [ -n "$region" ] || { echo "attach-game-ip: node has no providerID/region label" >&2; exit 1; }

rip=$(vapi "$API/reserved-ips?per_page=500" | jq -c --arg l "$LABEL" '[.reserved_ips[] | select(.label == $l)] | .[0] // empty')
if [ -z "$rip" ]; then
  echo "creating Reserved IP '$LABEL' in $region"
  rip=$(vapi -X POST "$API/reserved-ips" -d "$(jq -nc --arg r "$region" --arg l "$LABEL" '{region:$r, ip_type:"v4", label:$l}')" | jq -c .reserved_ip)
fi
id=$(jq -r .id <<<"$rip"); ip=$(jq -r .subnet <<<"$rip"); attached=$(jq -r '.instance_id // ""' <<<"$rip")
[ "$(jq -r .region <<<"$rip")" = "$region" ] || { echo "attach-game-ip: Reserved IP $ip is in another region than the game node ($region)" >&2; exit 1; }

if [ "$attached" != "$instance" ]; then
  [ -n "$attached" ] && { echo "detaching $ip from $attached"; vapi -X POST "$API/reserved-ips/$id/detach" >/dev/null; }
  echo "attaching $ip to $instance"
  vapi -X POST "$API/reserved-ips/$id/attach" -d "$(jq -nc --arg i "$instance" '{instance_id:$i}')" >/dev/null
fi

sed "s/\${GAME_IP}/$ip/" "$HERE/game-ip.yaml" | kubectl apply -f -
echo "game IP: $ip  (A records for game hosts point here)"
