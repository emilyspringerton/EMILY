#include "title_screen_proc_bg.h"

#include <string.h>

int title_proc_bg_init(TitleProcBg *bg, int tex_w, int tex_h, float now_sec) {
    if (bg == NULL) {
        return 0;
    }

    memset(bg, 0, sizeof(*bg));
    bg->refresh_interval_sec = 0.18f; /* tweak for animation speed */
    bg->uv_scale = 1.6f;              /* tweak for visual frequency */
    bg->brightness = 0.28f;           /* tweak for title readability */
    bg->seed = 12.0f;

    if (!proc_tex_create(&bg->texture, tex_w, tex_h)) {
        bg->enabled = 0;
        return 0;
    }

    proc_tex_fill_emily_vibe(&bg->texture, bg->seed, now_sec);
    proc_tex_upload(&bg->texture);

    bg->enabled = 1;
    bg->next_refresh_sec = now_sec + bg->refresh_interval_sec;
    return 1;
}

void title_proc_bg_shutdown(TitleProcBg *bg) {
    if (bg == NULL) {
        return;
    }

    proc_tex_destroy(&bg->texture);
    memset(bg, 0, sizeof(*bg));
}

void title_proc_bg_update(TitleProcBg *bg, float now_sec) {
    if (bg == NULL || !bg->enabled) {
        return;
    }

    if (now_sec < bg->next_refresh_sec) {
        return;
    }

    proc_tex_fill_emily_vibe(&bg->texture, bg->seed, now_sec);
    proc_tex_upload(&bg->texture);
    bg->next_refresh_sec = now_sec + bg->refresh_interval_sec;
}

void title_proc_bg_draw(const TitleProcBg *bg, float viewport_w, float viewport_h) {
    float uv = 1.0f;

    if (bg == NULL || !bg->enabled || bg->texture.gl_id == 0) {
        return;
    }

    uv = bg->uv_scale;

    glEnable(GL_TEXTURE_2D);
    glBindTexture(GL_TEXTURE_2D, bg->texture.gl_id);
    glColor3f(bg->brightness, bg->brightness, bg->brightness);

    glBegin(GL_QUADS);
    glTexCoord2f(0.0f, 0.0f);
    glVertex2f(0.0f, 0.0f);

    glTexCoord2f(uv, 0.0f);
    glVertex2f(viewport_w, 0.0f);

    glTexCoord2f(uv, uv);
    glVertex2f(viewport_w, viewport_h);

    glTexCoord2f(0.0f, uv);
    glVertex2f(0.0f, viewport_h);
    glEnd();

    glBindTexture(GL_TEXTURE_2D, 0);
    glDisable(GL_TEXTURE_2D);
    glColor3f(1.0f, 1.0f, 1.0f);
}
