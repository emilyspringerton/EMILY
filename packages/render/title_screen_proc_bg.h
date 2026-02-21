#ifndef EMILY_RENDER_TITLE_SCREEN_PROC_BG_H
#define EMILY_RENDER_TITLE_SCREEN_PROC_BG_H

#include "proc_tex.h"

typedef struct TitleProcBg {
    ProcTexture texture;
    int enabled;
    float refresh_interval_sec;
    float next_refresh_sec;
    float uv_scale;
    float brightness;
    float seed;
} TitleProcBg;

int title_proc_bg_init(TitleProcBg *bg, int tex_w, int tex_h, float now_sec);
void title_proc_bg_shutdown(TitleProcBg *bg);
void title_proc_bg_update(TitleProcBg *bg, float now_sec);
void title_proc_bg_draw(const TitleProcBg *bg, float viewport_w, float viewport_h);

#endif
