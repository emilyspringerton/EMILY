# Golden Docs Index
# Machine-readable manifest for GoldenDocCompiler.
# Format: | name | path | tier | budget | description |
# - name: label used in compiled output (no spaces)
# - path: relative to /home/fatbaby/ (the repo root parent)
# - tier: 1 = must always be in context; 2 = subsystem-level; skip others
# - budget: max chars read (0 = unlimited)
# - description: one line
#
# GoldenDocCompiler reads this file at startup. To add a new golden doc:
# append a row here — no Go change needed (S22-11).
# Rows starting with # are comments. Blank lines ignored.

| name | path | tier | budget | description |
|------|------|------|--------|-------------|
| FATBABY | PRRJECT_FATBABY/docs/northstar/northstar.md | 1 | 8000 | 8-K intelligence engine northstar; 4 layers, phase status |
| FATBABY-EXEC | PRRJECT_FATBABY/docs/northstar/executive_summary.md | 1 | 4000 | FatBaby executive summary; system-as-built orientation |
| EMILY-SPEC | EMILY/emily-prime-spec.md | 1 | 4000 | Emily Prime ↔ FatBaby-Emily integration spec |
| EMILY-NORTH | EMILY/docs/NORTHSTAR.md | 1 | 3000 | Emily Prime 3-agent arch + cron cycle northstar |
| EMILY-BACKLOG | EMILY/GOLDEN.md | 1 | 0 | Compressed backlog; open items, priority order |
| EMIREE | EMILY/emiree-emily-fatbaby.md | 1 | 4000 | Canonical 3-agent governance spec |
| RSI-ENGINE | EMILY/docs/emily-rsi-engine-spec.md | 1 | 0 | RSI loop ground truth; supersedes all earlier RSI docs |
| EMIREE-SPEC | EMILY/docs/emiree-over-agent-spec.md | 1 | 0 | Emiree witch engine; gear system implementation |
| IDUNA | IDUNA/golden.md | 1 | 3000 | IDUNA IAM + Apples implementation spec |
| APPLES | APPLES/docs/NORTHSTAR.md | 1 | 0 | Git-authoritative Apple audit trail; sync protocol |
| MJOLNIR | MJOLNIR/docs/NORTHSTAR.md | 1 | 0 | Android intelligence terminal; FCM + Apple feed + WebView |
| SHANKPIT | SHANKPIT/docs2/NORTHSTAR.md | 1 | 4000 | UDP FPS + DragonsNShit persistent world; Steam EA path |
| EMILYOS | EmilyOS/docs/NORTHSTAR.md | 1 | 0 | Bare-metal policy kernel; SOC 2 compliance northstar |
| PITVIPER | PITVIPER/docs/NORTHSTAR.md | 1 | 0 | SDL2 GPU terminal; 5 milestones; Emily Prime hooks |
| EMILY-CLI | emily.cli/docs/NORTHSTAR.md | 1 | 0 | Operator terminal; human half of the operator↔agent loop |
| GTM | PRRJECT_FATBABY/docs/GTM_FUNNEL.md | 1 | 3000 | Ask Emily GTM funnel; free tier → subscription → Merkle |
| TRAINING | EMILY/docs/TRAINING_PIPELINE.md | 1 | 2000 | Self-improving fine-tuning flywheel; RLHF loop |
| EDIS | EDIS/NORTHSTAR.md | 1 | 3000 | WordPress intelligence product; edis-core + signals + ask-emily plugins |
| IDUNA-NORTH | IDUNA/docs/NORTHSTAR.md | 1 | 3000 | IDUNA IAM + Apples + HEIMDAL northstar; port 8080; sprint lifecycle |
| THE-FIELD | EMILY/docs/THE_FIELD.md | 2 | 3000 | Synthetic consciousness architecture; dual-persona harmonic engine; informs Emiree |
