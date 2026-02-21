# Title screen procedural texture integration snippets

Use these snippets in your existing title screen paths.

## Globals/static state

```c
#include "packages/render/title_screen_proc_bg.h"

static TitleProcBg g_title_bg;
```

## Init (after OpenGL context exists)

```c
/* now_sec can come from SDL_GetTicks() * 0.001f */
(void)title_proc_bg_init(&g_title_bg, 256, 256, now_sec);
/* fail-open: if init fails, title screen still renders normally */
```

## Update (title screen tick)

```c
title_proc_bg_update(&g_title_bg, now_sec);
```

## Render (draw before title/logo/menu text)

```c
title_proc_bg_draw(&g_title_bg, (float)screen_w, (float)screen_h);
/* existing title/logo/menu rendering continues unchanged */
```

## Shutdown

```c
title_proc_bg_shutdown(&g_title_bg);
```

## Tuning

- Animation refresh interval: `TitleProcBg::refresh_interval_sec` (default `0.18f`)
- Texture brightness/intensity: `TitleProcBg::brightness` (default `0.28f`)
- UV scale: `TitleProcBg::uv_scale` (default `1.6f`)
- Magenta glitch intensity: `magenta_amt` in `proc_tex_fill_emily_vibe()`
