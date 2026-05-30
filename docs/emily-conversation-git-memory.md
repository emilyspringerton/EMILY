# Emily Conversation Git Memory

Emily can persist each chat turn into a self-documenting Git repository. The store keeps a JSONL session log for fast replay and renders the same turns into markdown so Git history becomes the durable audit trail and knowledge base.

## Repository Layout

By default, `CONVERSATION_DIR` points at `./fartco-memory`. Inside that Git repository Emily writes:

```text
fartco-memory/
├── conversations/
│   ├── YYYY/
│   │   └── MM/
│   │       └── YYYY-MM-DD-HH-MM-topic-slug.md
│   ├── topics/
│   │   └── git.md
│   └── index.md
└── sessions/
    └── session-id.jsonl
```

## Runtime Flags

- `CONVERSATION_DIR`: Git repository root for memory storage. Defaults to `./fartco-memory`.
- `GIT_COMMIT`: Set to `false` to write files without committing. Defaults to `true`.
- `GIT_PUSH`: Set to `true` to run `git push` after each successful memory commit. Defaults to `false` so local development is not blocked by remote configuration.

## Markdown Transcript Format

Each conversation markdown file includes:

- date, updated time, session ID, participants, topic, and status metadata;
- alternating `CEO` and `Emily` sections for every turn;
- turn metadata for validation and tool calls when present;
- inferred decisions and action items as checklist sections;
- tags used to populate topic pages.

## Search

Use the HTTP search endpoint to find stored markdown transcripts:

```bash
curl "http://localhost:8080/search?q=kubernetes"
```

The endpoint returns JSON results with the markdown path, title, timestamp, and a short snippet.
