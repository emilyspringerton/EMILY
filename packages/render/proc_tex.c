#include "proc_tex.h"

#include <math.h>
#include <stdlib.h>
#include <string.h>

static float clamp01(float v) {
    if (v < 0.0f) {
        return 0.0f;
    }
    if (v > 1.0f) {
        return 1.0f;
    }
    return v;
}

static float lerpf(float a, float b, float t) {
    return a + (b - a) * t;
}

static float smoothstep01(float x) {
    x = clamp01(x);
    return x * x * (3.0f - 2.0f * x);
}

static float hash2f(int x, int y, float seed) {
    float p = (float)x * 127.1f + (float)y * 311.7f + seed * 74.7f;
    float s = sinf(p) * 43758.5453f;
    return s - floorf(s);
}

static float value_noise_2d(float x, float y, float seed) {
    int ix = (int)floorf(x);
    int iy = (int)floorf(y);
    float fx = x - (float)ix;
    float fy = y - (float)iy;

    float v00 = hash2f(ix, iy, seed);
    float v10 = hash2f(ix + 1, iy, seed);
    float v01 = hash2f(ix, iy + 1, seed);
    float v11 = hash2f(ix + 1, iy + 1, seed);

    float sx = smoothstep01(fx);
    float sy = smoothstep01(fy);

    float nx0 = lerpf(v00, v10, sx);
    float nx1 = lerpf(v01, v11, sx);
    return lerpf(nx0, nx1, sy);
}

static float fbm_2d(float x, float y, float seed, int octaves) {
    float sum = 0.0f;
    float amp = 0.5f;
    float freq = 1.0f;
    int i;

    for (i = 0; i < octaves; ++i) {
        sum += amp * value_noise_2d(x * freq, y * freq, seed + (float)i * 13.0f);
        freq *= 2.02f;
        amp *= 0.5f;
    }

    return sum;
}

int proc_tex_create(ProcTexture *t, int w, int h) {
    size_t bytes;

    if (t == NULL || w <= 0 || h <= 0) {
        return 0;
    }

    memset(t, 0, sizeof(*t));
    bytes = (size_t)w * (size_t)h * 4u;
    t->pixels = (unsigned char *)malloc(bytes);
    if (t->pixels == NULL) {
        return 0;
    }

    t->width = w;
    t->height = h;
    memset(t->pixels, 0, bytes);

    glGenTextures(1, &t->gl_id);
    if (t->gl_id == 0) {
        free(t->pixels);
        memset(t, 0, sizeof(*t));
        return 0;
    }

    glBindTexture(GL_TEXTURE_2D, t->gl_id);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_REPEAT);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_REPEAT);
    glBindTexture(GL_TEXTURE_2D, 0);

    t->valid = 1;
    return 1;
}

void proc_tex_destroy(ProcTexture *t) {
    if (t == NULL) {
        return;
    }

    if (t->gl_id != 0) {
        glDeleteTextures(1, &t->gl_id);
    }
    if (t->pixels != NULL) {
        free(t->pixels);
    }

    memset(t, 0, sizeof(*t));
}

void proc_tex_upload(ProcTexture *t) {
    if (t == NULL || t->pixels == NULL || t->gl_id == 0) {
        return;
    }

    glBindTexture(GL_TEXTURE_2D, t->gl_id);
    glPixelStorei(GL_UNPACK_ALIGNMENT, 1);
    glTexImage2D(GL_TEXTURE_2D,
                 0,
                 GL_RGBA,
                 t->width,
                 t->height,
                 0,
                 GL_RGBA,
                 GL_UNSIGNED_BYTE,
                 t->pixels);
    glBindTexture(GL_TEXTURE_2D, 0);
}

void proc_tex_fill_emily_vibe(ProcTexture *t, float seed, float t_sec) {
    int x, y;
    float drift = t_sec * 0.06f;

    if (t == NULL || t->pixels == NULL || t->width <= 0 || t->height <= 0) {
        return;
    }

    for (y = 0; y < t->height; ++y) {
        for (x = 0; x < t->width; ++x) {
            float u = (float)x / (float)t->width;
            float v = (float)y / (float)t->height;

            float fog = fbm_2d(u * 2.2f + drift, v * 2.2f - drift * 0.35f, seed + 2.0f, 4);
            float grain = value_noise_2d(u * 38.0f + drift * 3.0f,
                                         v * 38.0f - drift * 2.0f,
                                         seed + 11.0f);

            float diag = u * 5.5f + v * 11.0f + drift * 1.2f;
            float glitch_mask = value_noise_2d(u * 12.0f + drift * 0.8f,
                                               v * 12.0f + 3.1f,
                                               seed + 29.0f);
            float stripe = fabsf(sinf(diag * 3.1415926f));
            float glitch = (stripe > 0.93f && glitch_mask > 0.74f) ? 1.0f : 0.0f;

            float base_r = 0.05f + fog * 0.10f + grain * 0.04f;
            float base_g = 0.30f + fog * 0.30f + grain * 0.05f;
            float base_b = 0.36f + fog * 0.33f + grain * 0.06f;

            float magenta_amt = glitch * 0.22f;
            float r = clamp01(base_r + magenta_amt);
            float g = clamp01(base_g - magenta_amt * 0.35f);
            float b = clamp01(base_b + magenta_amt * 0.48f);

            int idx = (y * t->width + x) * 4;
            t->pixels[idx + 0] = (unsigned char)(r * 255.0f);
            t->pixels[idx + 1] = (unsigned char)(g * 255.0f);
            t->pixels[idx + 2] = (unsigned char)(b * 255.0f);
            t->pixels[idx + 3] = 255;
        }
    }
}
