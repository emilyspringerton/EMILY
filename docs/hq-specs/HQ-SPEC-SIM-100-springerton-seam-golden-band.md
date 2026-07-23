# HQ-SPEC-SIM-100 — The Springerton Seam & the GOLDEN BAND Animation Layer

**Status:** DRAFT v0 — pending Emily Prime review
**Custody:** Emily OS
**Governs:** The sim↔real boundary and the animation system that feeds both sides
**Prime values, stated first because they govern everything below:** SAFETY and TRANSPARENCY. Autonomy is unlimited inside simulation; nothing crosses into atoms without a human gate; every crossing is provable after the fact.
**Working name:** GOLDEN BAND — the animation layer. (A band is a timeline; the Logistics Architect of the Golden Band finally has her department. Rename if the resonance is wrong — but it isn't.)

---

## 1. Northstar Statement

One canonical motion asset, three consumers:

1. **Game runtime** — SHANKPIT (and successors) play it back as skeletal animation.
2. **Reward compiler** — RL training in the physics backbone treats it as the imitation target: the policy is graded on how faithfully physics can honor the animator's intent.
3. **Hardware deployment** — the trained policy, not the clip, drives physical characters for EINHORN Parks & Cruises attractions.

The animator's clip is the spec; physics is the implementation; the policy is the compiled artifact. This is the Disney BDX insight made into house doctrine: charm is authored by humans, robustness is learned by machines, and the reward function is the contract between them.

## 2. Positioning: Between Unreal and Godot, Bare Metal

**What that means concretely:**

- **Footprint like Godot, semantics like Unreal.** Full skeletal animation stack — blending, state machines, layered/masked playback, retargeting, root motion, IK passes — but no editor monolith, no scene-graph religion, no GC pauses in the hot path.
- **Bare metal aesthetic:** core in C (linking cleanly into SHANKPIT's existing C/OpenGL 2.1 runtime), tooling and pipeline in Go (house language), data-driven everything. The runtime is a library, not a framework: `gb_init()`, `gb_sample(rig, clip, t)`, `gb_blend(...)`, and you own the loop.
- **Determinism is load-bearing, not aspirational.** Animation sampling is locked to SHANKPIT's 64-tick fixed timestep; render-rate smoothness comes from the existing state interpolation layer, never from re-sampling. Same clip + same seed + same tick = bit-identical pose, on every platform. This is what makes a game replay, a training rollout, and an incident reconstruction the same operation.
- **No engine dependency in the asset.** GOLDEN BAND assets know nothing about Unreal, Godot, or SHANKPIT. Adapters exist at the edges (a Godot GDExtension and an Unreal plugin are explicitly on the roadmap for external licensing/hiring reasons), but the canonical format is ours.

## 3. The Canonical Motion Asset (`.gband`)

Content-addressed, versioned, two-part:

- **Manifest (JSON):** skeleton reference (by hash), channel list, duration, tick rate, authorship provenance (who/what authored it — human animator, mocap session, or generative tool, *labeled honestly*), intent tags (`gait`, `gesture`, `idle`, `showpiece`), loop points, and safety annotations (see §7).
- **Channel data (flat binary):** per-joint rotation/translation tracks, uniform sampling at authoring rate. No curves-with-twelve-interpolation-modes cleverness; resample at import. Simple enough to parse in a hundred lines of C — that's the bare metal test.

**Skeleton spec:** a canonical rig per character family — joint hierarchy, limits, and (for hardware-bound characters) actuator metadata: torque limits, velocity limits, gear backlash notes. The skeleton is where game fantasy and hardware truth are forced to negotiate.

**Retargeting maps** are first-class versioned assets: canonical rig → game character mesh, canonical rig → robot morphology. The hardware retarget includes a **feasibility pass** — any frame demanding infeasible joint velocity/torque is flagged at authoring time, not discovered on a test stand. Animators get feasibility feedback inside their tool loop, the way Imagineering's animators iterate against the physical envelope.

Everything content-addressed; promotions of assets to `golden` status are Apples-logged `ApplePublished` events. The motion library is an auditable corpus, not a folder of files.

## 4. The Reward Compiler

The glue itself. Input: a `.gband` clip + a reward profile. Output: a versioned reward function artifact consumed by the training backbone.

**Reward terms (weighted, all standard imitation-learning practice — DeepMimic/AMP lineage):**
- Pose tracking (joint orientation error vs. reference)
- Velocity tracking
- End-effector position tracking (hands/feet/prop points — where charm lives)
- Root motion tracking (trajectory + heading)
- Regularizers: energy/torque penalties, action smoothness, self-collision, joint-limit margins
- Task terms layered on top when the motion must serve gameplay or show logic (reach the mark, face the guest, survive the shove)

**Domain randomization profiles** are versioned config, not code: mass/inertia jitter, friction ranges, motor latency, sensor noise, terrain variation, external perturbation schedules. Each deployment target (SHANKPIT NPC, test-stand biped, park character) has its own DR profile; the profile hash travels with every trained policy. Zero-shot transfer isn't magic — it's DR breadth purchased with sim compute, and the receipt is the profile hash.

**Training backbone:** GPU-parallel physics (Isaac Lab / MJX / Genesis — decided per SIM-100 predecessor discussion, adapter-isolated the same way ELP isolates rails). SHANKPIT is the *observation head and evolution orchestrator*: GameEvolutionEngine proposes reward-weight and DR-profile variations as structured JSON improvement proposals; frozen eval suites grade them; no-regression gate before promotion. Same machine as the EPS grader and the recon matcher. CAST heads spectate training runs live — connection failures invisible to training, per CAST doctrine.

## 5. Dual Deployment Paths

**Path A — Games (ungated autonomy):**
Trained policies become NPC brains: physics-animated characters in SHANKPIT that stumble, recover, and improvise instead of playing canned clips. The educational angle compounds — edu_vm scripts can perturb characters and students watch policies cope. Games are also the *personality proving ground*: a character that can't be charming at 64 ticks/sec in SHANKPIT doesn't earn a body.

**Path B — Attractions (hard-gated):**
The same policy artifact, retargeted to hardware, crosses the Springerton Seam only through the gate protocol (§6). Character continuity is the point: the fiction universe (Tyler, the Syndicate, the Golden Band itself) supplies the cast; the motion library supplies their souls; the parks supply the bodies. Long horizon, honestly held — Parks & Cruises is a northstar the way Column is in FIN-099: it disciplines the architecture now, it does not appear on a sprint board.

## 6. The Springerton Seam — Gate Protocol

The seam is the customs checkpoint between simulation and atoms. Rules:

1. **Deployment bundle = content hash** over {policy weights, skeleton, retarget map, DR profile, reward artifact, sim config, eval results}. The hash is what gets approved. Any component changes → new hash → new approval. (Identical mechanism to KAREN's payment approval binding; `actuation_approved` is `payment_approved` for atoms.)
2. **Human biometric approval** via the companion app for any physical deployment. No standing approvals, no batch waivers — Invariant 2 doctrine, inherited verbatim.
3. **Graduated actuation ladder:** sim → hardware-in-the-loop bench (actuators, no body) → tethered test stand → untethered in controlled space → staff-facing → guest-facing. Each rung is a separate gate with its own eval suite. No rung skipping, even for demos. *Especially* for demos.
4. **Telemetry return flow:** every physical run streams back through the seam as events; sim-vs-real deviation is a first-class metric and a recon-style exception queue. Reality grades the simulator; the simulator's fidelity improvements pass their own no-regression gate.
5. **Löbian clause, physical edition:** no component may certify a deployment bundle that its own outputs are part of. Eval suites are frozen and versioned independently of the training loop that they grade. The system never grades its own homework at the seam.

**Amendment (append-only), 2026-07-23:** the five rules above are this document's own description
of a propose→grade→gate→promote loop, written before the loop itself was extracted into a shared
kernel. `HQ-SPEC-PRIME-101-norn-loop-kernel.md` (`pkg/norn` — Artifact/Proposer/Oracle/Gate/Registry)
is now that kernel, and §6's "Instantiation Table" carries this seam as its own row ("Deployment
bundles (Seam)", tier `biometric`, reality root "Physical telemetry"). This section is not
rewritten — the gate protocol above remains the authoritative description of what the seam
requires — but going forward it should be read as *this domain's instantiation of NORN*, not a
bespoke loop: rule 1 (content-hash bundle identity) is NORN's `Artifact`; rule 4 (telemetry
return flow) is NORN's `Oracle`, reality-rooted per PRIME-101's own terms; rule 5 (the Löbian
clause) is exactly what `CheckLineage`'s content-hash ancestry walk enforces mechanically in
`pkg/norn` today (see NORN's S141-02/03 migrations for two other domains already governed the same
way). When seam tooling is actually built (Build Sequence step 6, hardware-in-the-loop bench), it
should register as a NORN instantiation rather than re-implement gate/promotion logic from
scratch. Per PRIME-101 §8 build-sequence item 5.

## 7. Safety & Transparency Architecture

Guest-facing robots are the highest-stakes thing EINHORN will ever ship. The framing that holds:

- **Standards as the skeleton, not the paperwork:** ISO 10218 / ISO/TS 15066 (collaborative robot force/speed limits) for anything near humans; IEC 61508 functional-safety thinking for the control chain; **ASTM F24 (F2291 et al.)** — the amusement-ride standards body — for attractions specifically. Design to these from day one so certification evidence is a *byproduct of the event log*, not a retrofit project. Transparency is the product: our audit story should be the reason a park operator or classification society trusts us.
- **Hardware-enforced envelopes:** speed/force/torque limits live below the policy in a supervisor layer the learned system cannot override — a deterministic C component, formally reviewable, boring on purpose. The policy proposes motion; the envelope disposes.
- **Characterful failure:** every character has an authored `safe_pose` clip and degradation ladder (full performance → reduced envelope → in-character settle → e-stop). Failure modes are animated content, authored with the same care as showpieces — a character that shuts down gracefully in-fiction is both safer and more magical than one that freezes.
- **Provenance for guests and regulators alike:** every deployed character can answer "what version are you, who approved you, when, against which evals" as a query. Safety annotations in the manifest (max intended speed, guest-proximity class) are machine-checked against the deployment target's class at gate time.
- **Honest labeling of authorship:** generative-assisted motion is tagged as such in provenance. When a guest asks if a character's wave was made by a person, EINHORN knows the answer and isn't embarrassed by it.

## 8. Build Sequence

1. **`.gband` format + C sampler + Go pipeline tools** (import from BVH/glTF, resample, hash, validate). A hundred-line parser is the acceptance test.
2. **SHANKPIT integration:** skeletal playback at fixed timestep with existing interpolation; one character, one idle, one walk. Rotoscope-shader lineage says the render side is already comfortable with stylization.
3. **Reward compiler v0 + training backbone adapter:** one character learns to walk in physics sim tracking the authored walk. Grade against the clip; CAST the rollouts.
4. **GameEvolutionEngine hookup:** reward-weight/DR proposals through the frozen-eval no-regression loop.
5. **Physics-driven NPC in SHANKPIT** — Path A shipped; personality proving ground open.
6. **Hardware-in-the-loop bench** — first rung of the actuation ladder, cheapest possible articulated testbed (a single 2–3 DOF limb is enough to make the seam real).
7. **Retarget feasibility tooling for animators** — close the authoring loop before scaling the cast.

## 9. Open Questions

- **GOLDEN BAND as the name** — proposed above; Springerton's department by right. Confirm or overrule.
- **Authoring tool reality:** animators author in Blender/Maya/mocap and export, or do we eventually owe them a native tool? (v0: import-only. A native tool is a product decision, not a pipeline necessity.)
- **First character:** who earns the first body? The answer is a casting decision from the fiction universe, and it should be made like one.
- **Cruise-specific compliance** (classification societies, maritime electrical standards) — parked until Parks precedes Cruises or vice versa; flag with counsel alongside the FIN-099 licensing questions.
- **Backbone selection** (Isaac Lab vs MJX vs Genesis) — decide by running build-sequence step 3 on two of them; the adapter keeps it reversible.

---

*The pen is still owed. CLEAN BUILDS FIRST.*
