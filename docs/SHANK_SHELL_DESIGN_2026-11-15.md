# SHANK SHELL — a desktop app shell without Electron (design, 2026-11-15)

Founder real-time, 2026-11-15: *"lets ensure our shank client works from SPGF for a full stripped down real browser parena
first mods first pragmatic chromium is the glue think about an electron killer"*. This is the thinking, before code. "SPGF"
is not a name anywhere in the repositories; the plan below does not depend on it and asks about it at the end.

## The problem with Electron

An Electron app ships its own copy of Chromium **and** Node.js: 100–200 MB per app, its own update cadence (security fixes
arrive only when the app re-ships), and a Node runtime the app mostly does not need. Ten Electron apps are ten browsers.

## The idea

**Use the Chromium that is already on the machine, and nothing else.** Every Windows 10/11 PC has Microsoft Edge (Chromium);
most Linux and macOS machines have Chrome or Chromium. Started with `--app=<url>` and its own profile folder, a Chromium
browser opens a window with no tabs, no address bar and its own taskbar entry. That is an app window.

Everything else is ours, and PARENA:

| Layer | Electron | SHANK SHELL |
|---|---|---|
| Rendering, input, GPU, audio, WebAssembly | bundled Chromium | the installed Chromium (Edge / Chrome / Chromium), `--app` mode |
| App logic | JavaScript + Node | PARENA compiled to WebAssembly (as the WOTAN games already are: no imports, no allocator) |
| Page glue | hand-written JS | PARENA compiled to TypeScript (`bp_host.prn`, `sp_host.prn` ... already exist) |
| Serving the app | `file://` or a custom protocol | a ~1 MB launcher serving the app's folder on 127.0.0.1 (PARENA `net/http_server.prn`) |
| Native access (files, saves) | Node APIs, `ipcMain` | loopback endpoints on the same server, off unless the app's mod grants them, each call carrying a per-launch secret |
| What an app is | `package.json` + `main.js` | a PARENA mod: name, entry page, window size, capabilities |
| Size per app | 100–200 MB | the app's files + one small launcher shared by every app |

"Chromium is the glue": it does the parts not worth rebuilding (a GPU-accelerated canvas, input, audio, a wasm JIT) and
nothing else. Updates to it come from the OS or the browser's own updater, not from us.

## How it runs

1. `shank-shell <app>` reads the app's entry from the mod (`shell_apps.prn`): folder, entry page, window size,
   capabilities.
2. It starts a loopback HTTP server on a random port, serving only that folder (correct MIME types, `application/wasm`
   included; no directory listing; no path outside the folder).
3. It finds a browser: `SHANK_BROWSER`, else Edge then Chrome then Chromium in their standard places on Windows, macOS and
   Linux.
4. It launches the browser with `--app=http://127.0.0.1:<port>/<entry>`, `--user-data-dir=<app's own profile>`,
   `--window-size=<w>,<h>`, `--no-first-run`, `--no-default-browser-check`. The app's storage is its own profile.
5. When the browser window closes, the launcher stops the server and exits.
6. No Chromium-family browser: the launcher opens Pellegrino (PARENA's own browser) on a page that says what to install.

**The long game (the real Electron killer):** Chromium becomes optional once the PARENA host can run the app's WebAssembly
itself, the way the SHANKPHONE already runs BRAWLPIT, shankpit-460 and Bedrock Racers on SHANKPIT's own wasm3 level host
(EMILY/BACKLOG.md S635-4). Then a SHANK SHELL app is the wasm module plus a native PARENA host: no browser at all.

## What already exists to build on

- The WOTAN games as self-contained static bundles (the WOTAN Arcade artifact proves they run from plain static files).
- PARENA `net/http_server.prn`, `net/tcp.prn` (listen/accept), and generated TS hosts for each game.
- SHANKPIT `packages/chromium_bridge/` (starting and driving a Chromium; the phone's web-game stopgap) and
  `packages/simulation/phone_arcade.h` (the wasm3 level host path).

## Risks, named

- **No Chromium on a Linux box:** fall back to Pellegrino's message; package managers install one.
- **Enterprise policy can disable `--app` or force-update Edge:** the app still opens as a normal browser window.
- **An unsigned launcher .exe can trip antivirus:** sign the release build (the monorepo already GPG-signs app releases).
- **Loopback is reachable by other local programs:** a per-launch secret on every native call; native calls off by default.
- **Browser differences:** pin behaviour to Chromium-family only (Edge/Chrome/Chromium), unlike Tauri's per-OS webviews.

## Build order (proposed)

1. `stdlib/shell/shell_apps.prn` (the mod: apps) + `shell_rules.prn` (browser search order, flags, MIME types).
2. The launcher (PARENA-first: the server is `net/http_server.prn`; C only for process spawn and browser lookup), Linux
   and Windows, tested under Xvfb with this box's Chromium.
3. First app: the WOTAN Arcade (already static). Second: the SHANKPIT web client (needs the open question below).
4. Native bridge with per-app capabilities (saves first).
5. Pellegrino fallback page; then the no-browser path through the PARENA wasm host.

## Open question for the founder

**What is "SPGF"?** It is not a repository, file or name in the code. If it is where the shank client should launch from
(a folder, a server, a page), the launcher's first app entry points there.
