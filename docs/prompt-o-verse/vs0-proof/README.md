# VS0 proof-of-concept — first real generation

The first real artifact for `PROMPT-O-VERSE-NORTH` (`EMILY/docs/NORTHSTAR_PROMPT_O_VERSE.md`) —
proof that the core mechanic (top-level prompt → generated data, §2 of the northstar) actually
works, not just that a viable backend exists on paper.

**Generated:** 2026-08-17
**Backend:** Vertex AI (`aiplatform.googleapis.com`), model `gemini-2.5-flash-image` ("Nano Banana")
**Auth:** this box's existing `gcloud` ADC (`garybifrost@gmail.com`), region `us-central1` —
no dedicated `GEMINI_API_KEY` required; see the northstar's §5 for the full path (the AI Studio
API-key flow and the Gemini CLI's OAuth-personal tier were both tried first and both confirmed
dead — Vertex AI via the account's own already-provisioned `gcloud` credentials is what worked)
**Prompt (EZ-tier):** "A 1990s glossy baseball rookie card portrait photo, studio lighting, blue
background."
**Output:** `1990s-rookie-card-portrait.png`, 1024×1024, real inline PNG data from the API
response, not post-processed.

Not scoped or claimed here: this is a single proof call, not VS0's real taxonomy-building pass
(§7/§8 of the northstar — pick a backend, then run the generate→observe→categorize→isolate-prompt
loop for real against many prompts). This file exists to prove the backend question is answerable,
not to be the taxonomy itself.
