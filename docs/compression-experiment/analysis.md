# Ops Docs Multilingual Compression Experiment
## Step 1: Candidate Analysis + Step 2: Bilingual Test Version
### Status: Steps 1–2 complete. Step 3 (A/B test) blocked on ANTHROPIC_API_KEY.

---

## Step 1: Candidate Identification

### Hot-path docs (loaded into haiku at runtime)

| Doc | Lines | Approx tokens | Usage |
|---|---|---|---|
| `EMILY/GOLDEN.md` | 72 | ~576 (well under 1200 target) | `fable.go`: both `RunFableAdvice` and `RunFableExecute` load this into claude-haiku user prompt |

All other docs in `EMILY/docs/` are **design documents, not runtime context**. They are not loaded by emily-agent at runtime. The large files (emily-cron-autonomous-evolution.md ~1400 lines, emily-training-layer.md ~1200 lines, etc.) are reference material only.

### Non-hot-path candidates for future optimization

If GOLDEN.md grows past budget, or if additional haiku calls are introduced that load docs, the secondary candidates are:

| Doc | Lines | Approx tokens | Currently loaded? |
|---|---|---|---|
| `EMILY/emiree.md` | 521 | ~4000 | No — emiree state is persisted JSON, not the .md file |
| `PRRJECT_FATBABY/docs/northstar/northstar.md` | 670 | ~5400 | No |
| `SHANKPIT/docs2/NORTHSTAR.md` | 196 | ~1600 | No |

**emiree.md** already demonstrates the bilingual encoding pattern (七卷 Chinese poetry + Sanskrit encoding of the witch-math equations). It serves as the template for this experiment.

### Conclusion: Primary target is GOLDEN.md

GOLDEN.md at 72 lines / ~576 tokens is already well under the 1200-token design budget. However:
- The budget is a ceiling, not the actual usage. If backlog sections grow, GOLDEN.md grows.
- A bilingual compressed version could establish a new ceiling or reduce API costs at scale.
- The experiment is worth running to validate the technique before it becomes urgent.

---

## Step 2: Bilingual Test Version

See `GOLDEN_BILINGUAL_TEST.md` in this directory.

Compression strategy:
- Section headers → 中文 header (4–6 characters) + English abbreviation
- Repo list → 中文 abbreviations only (repos are known to haiku from training)
- Status counts → keep as-is (numerals are compact)
- Item descriptions → compressed to core noun phrase in 中文 (2–6 characters) + keep English key terms (API names, function names)
- RSI directive → 中文 summary (saves ~1/3 of the tokens in that section)

Expected reduction: ~30–40% token count reduction. Estimated compressed size: ~400 tokens vs ~576.

---

## Step 3: A/B Test (BLOCKED — needs ANTHROPIC_API_KEY)

Test script ready at: `EMILY/scripts/compression-abtest.sh`

Test protocol:
1. Load GOLDEN.md → send to haiku with 3 comprehension questions
2. Load GOLDEN_BILINGUAL_TEST.md → same 3 comprehension questions
3. Compare response accuracy (manual review or automated scoring)
4. If scores equal: GOLDEN_BILINGUAL_TEST.md → replace GOLDEN.md
5. If scores degrade: keep GOLDEN.md, document failure mode

Comprehension questions:
- Q1: "Which sections have open items? List section numbers and counts."
- Q2: "What is the RSI prime directive pipeline in order?"
- Q3: "Name the repos in this system."

---

## Step 4: Deploy (dependent on Step 3)

If A/B test passes: update `emily backlog compress` command to output bilingual format by default.
