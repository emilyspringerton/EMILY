# VS1 Toothbrush Brief — EINHORN_INDUSTRIAL
## Two Brands: BASIC + ULTRA

**Status:** DESIGN LOCKED — vendor selection gated (S139-02)  
**Owner:** Emily Springerton  
**Updated:** 2026-06-27

---

## Brand Architecture

EINHORN_INDUSTRIAL operates two consumer brands in oral care. They do not reference each other.

| Brand | Tier | Price | Positioning |
|-------|------|-------|-------------|
| **STINKIES COMMISSAIRE** | Entry | $9 | The best toothbrush that does nothing extra |
| **ULTRA** | Ultra-premium | $68 | The last toothbrush you buy |

---

## STINKIES COMMISSAIRE

### What It Is

STINKIES COMMISSAIRE is a toothbrush brand with no pretension. It does the job. No logo tour, no wellness copy, no ingredient list. Black and white. Recyclable cardboard. Ships fast.

It exists because the $3 drugstore toothbrush is bad and the $40 toothbrush is theater. STINKIES COMMISSAIRE is the gap. The name says everything.

### Handle

- **Material:** Matte black polypropylene, recycled-content where available
- **Weight:** 18–22g
- **Length:** 185mm
- **Grip:** Flat spine, slight taper at neck — no rubberized grip, no thumb indent
- **Print:** "STINKIES COMMISSAIRE" in white, Helvetica Neue, 6pt, on spine — nothing else

### Head

- **Bristles:** Nylon, 0.20mm diameter, soft
- **Color:** White
- **Replace interval:** 90 days

### Packaging

- One-color (black on kraft cardboard). "BASIC TOOTHBRUSH." No adjectives.
- 100% recyclable

### Price + Margin

| | |
|---|---|
| Retail | $9 |
| Unit cost target | ≤ $1.80 at 500-unit MOQ |
| Gross margin target | ≥ 80% |
| Subscription | 4-pack replacement heads / 3 months — $12 |

---

## ULTRA

### What It Is

ULTRA is the final answer to the toothbrush question. Solid brass handle, cold-forged, PVD matte black. DuPont Tynex bristles. Replaceable head. You buy it once. You replace the head.

No features. No sonic vibration. No app. A weighted brass instrument that you use to clean your teeth.

### Handle

- **Material:** Solid brass, cold-forged
- **Coating:** PVD matte black — no chrome, no gloss, no color
- **Weight:** 40g ± 1g
- **Length:** 192mm
- **Grip:** Diamond knurl, 0.5mm pitch, mid-section only
- **Engraving:** "ULTRA" on spine, 0.3mm recessed, no fill. Nothing else.
- **Head mount:** Replaceable press-fit neck, 15mm OD

### Head

- **Frame:** Matte black nylon
- **Bristles:** DuPont Tynex®, 0.15mm, medium
- **Bristle color:** Black
- **Replace interval:** 90 days

### Packaging

- Matte black rigid box. "ULTRA" debossed on lid. No other text on exterior.
- Interior: precision-cut black foam. One card inside: "replace the head every 90 days."

### Price + Margin

| | |
|---|---|
| Retail | $68 |
| Unit cost target | ≤ $18 at 250-unit MOQ |
| Gross margin target | ≥ 73% |
| Replacement head (2-pack) | $22 / 90 days → $88/yr recurring |

---

## QC Criteria

### STINKIES COMMISSAIRE
1. Bristle pull test: 50N, no loss
2. Handle flex test: no crack at 15° deflection
3. Print legibility: 100% of units pass 1m readability ("STINKIES COMMISSAIRE" must be readable)

### ULTRA
1. Handle weight: 40g ± 1g (weigh all units in batch)
2. PVD adhesion: no delamination after 100 acetone wipes
3. Knurl depth: 0.45–0.55mm (digital caliper)
4. Bristle pull test: 50N, no loss
5. Engraving depth: 0.25–0.35mm
6. Photograph 5 units under D50, commit to `EMILY/docs/merch/vs1_qc/ultra/`

---

## Subscription Model

Both brands run on 90-day replacement head cycles. Emily Prime monitors inventory and reorder threshold via Supply Chain AGI (S136). Below threshold → auto-draft PO → HEIMDAL approval sprint.

| | STINKIES COMMISSAIRE | ULTRA |
|---|---|---|
| Replacement unit | 4-pack heads | 2-pack heads |
| Subscription price | $12 / 90 days | $22 / 90 days |
| Annual recurring | $48/yr | $88/yr |
| Reorder trigger | ≤ 80 units | ≤ 40 units |

---

## Brand Statement

**STINKIES COMMISSAIRE:** It's a toothbrush.  
**ULTRA:** It's the last one.
