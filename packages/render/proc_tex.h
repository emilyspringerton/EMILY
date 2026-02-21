#ifndef EMILY_RENDER_PROC_TEX_H
#define EMILY_RENDER_PROC_TEX_H

#include <stddef.h>

#include <SDL2/SDL_opengl.h>

typedef struct ProcTexture {
    int width;
    int height;
    GLuint gl_id;
    unsigned char *pixels;
    int valid;
} ProcTexture;

int proc_tex_create(ProcTexture *t, int w, int h);
void proc_tex_destroy(ProcTexture *t);
void proc_tex_upload(ProcTexture *t);
void proc_tex_fill_emily_vibe(ProcTexture *t, float seed, float t_sec);

#endif
