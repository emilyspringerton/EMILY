// emily-agent/mailinglist_watch.go — OBSERVE-phase check for new mailing-list
// signups (general/stinkies/freehoodie), same shape as observeVelocityAlerts
// in velocityalert.go: poll a count, diff against the last known value,
// file an Apple + FCM push to MJOLNIR when it moves.
//
// Founder ask, 2026-07-19: "lets get a push notification on mjolnir on new
// email signups." Built once the underlying count endpoint existed
// (IDUNA GET /api/v1/mailing-list/count, added for free-hoodie.html's live
// spots-remaining counter) — this is the second consumer of that endpoint.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"emily-agent/pkg/fcm"
)

var mailingListWatchSources = []string{"general", "stinkies", "freehoodie"}

type mailingListWatchState struct {
	LastCounts map[string]int `json:"last_counts"`
}

func mailingListWatchStatePath(stateDir string) string {
	return filepath.Join(stateDir, "mailinglist-watch-state.json")
}

func loadMailingListWatchState(stateDir string) mailingListWatchState {
	var s mailingListWatchState
	data, err := os.ReadFile(mailingListWatchStatePath(stateDir))
	if err != nil {
		return mailingListWatchState{LastCounts: map[string]int{}}
	}
	if err := json.Unmarshal(data, &s); err != nil || s.LastCounts == nil {
		return mailingListWatchState{LastCounts: map[string]int{}}
	}
	return s
}

func saveMailingListWatchState(stateDir string, s mailingListWatchState) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(mailingListWatchStatePath(stateDir), data, 0o644)
}

// observeMailingListSignups fetches the current subscriber count for each
// tracked list and fires an Apple + FCM push for any increase since the
// last cycle. Called from the OBSERVE phase, same as observeVelocityAlerts.
// All errors are logged, never fatal — a mailing-list hiccup should never
// take down the cron cycle.
func (ac *AutonomousCycle) observeMailingListSignups(ctx context.Context, cycleNum int) {
	if ac.iduna == nil {
		return
	}

	state := loadMailingListWatchState(ac.cfg.StateDir)
	dirty := false

	for _, source := range mailingListWatchSources {
		countCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		count, err := ac.iduna.MailingListCount(countCtx, source)
		cancel()
		if err != nil {
			log.Printf("[cycle %d] mailing-list count %q: %v", cycleNum, source, err)
			continue
		}

		last, seen := state.LastCounts[source]
		if !seen {
			// First cycle observing this source — record the baseline
			// without firing a push for subscribers who signed up before
			// this watcher existed.
			state.LastCounts[source] = count
			dirty = true
			continue
		}

		if count > last {
			delta := count - last
			msg := fmt.Sprintf("%d new signup(s) on mailing list %q (%d total)", delta, source, count)
			log.Printf("[cycle %d] %s", cycleNum, msg)

			appleCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			if _, err := ac.iduna.PostApple(appleCtx, ApplePayload{
				SourceRepo: "IDUNA",
				RunID:      fmt.Sprintf("mailinglist-signup-%s-%d", source, time.Now().Unix()),
				AppleType:  "observation",
				Title:      fmt.Sprintf("mailing-list: +%d on %s", delta, source),
				Body:       msg,
			}); err != nil {
				log.Printf("[cycle %d] mailing-list apple warn: %v", cycleNum, err)
			}
			cancel()

			if ac.fcmSender != nil {
				pushCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
				deviceToken, err := ac.iduna.GetPushToken(pushCtx, "mjolnir-emily")
				cancel()
				if err == nil && deviceToken != "" {
					_ = ac.fcmSender.Send(ctx, deviceToken, fcm.Message{
						Title:    fmt.Sprintf("📬 +%d signup(s) — %s", delta, source),
						Body:     msg,
						Priority: "normal",
						Data:     map[string]string{"source": source, "count": fmt.Sprintf("%d", count)},
					})
				}
			}
		}

		state.LastCounts[source] = count
		dirty = true
	}

	if dirty {
		saveMailingListWatchState(ac.cfg.StateDir, state)
	}
}
