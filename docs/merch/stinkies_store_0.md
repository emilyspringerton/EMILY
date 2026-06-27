# STINKIES COMMISSAIRE — Store 0
## Pontiac, Michigan

**Classification:** Physical flagship — template for all future locations
**Status:** Pre-build — site selection + permitting gated (S144-02)
**Owner:** Emily Springerton
**Updated:** 2026-06-27

---

## What Store 0 Is

The first physical STINKIES COMMISSAIRE. Pontiac, Michigan. A real store. Not a pop-up. Not an activation. A commissary that is open every day and sells what people need.

Store 0 is the proof of concept. Every system Emily Prime builds around physical product — supply chain, inventory, reorder, point-of-sale — gets validated here first. Store 0 is the factory floor for the commissary model.

---

## Location

**City:** Pontiac, Michigan  
**Target neighborhood:** Downtown Pontiac or near M-59 corridor — foot traffic, working-class density, proximity to auto industry workforce  
**Space target:** 800–1,200 sq ft retail footprint  
**Format:** Commissary / convenience hybrid. Not a grocery. Not a bodega. A commissary.

---

## Build-Out Scope

Emily (human) manages the physical build. Emily Prime manages the procurement, scheduling, and vendor coordination.

| Phase | Work | Owner |
|-------|------|-------|
| 0 | Site selection, lease negotiation | Emily (human) |
| 1 | Demo + concrete work (floors, any slab repair) | Contractor — Emily schedules via HEIMDAL |
| 2 | Electrical + HVAC rough-in | Licensed contractor — HEIMDAL approval |
| 3 | Shelving, coolers, display build-out | Emily + crew |
| 4 | POS system + inventory system integration | Emily Prime (EDIS/WooCommerce POS) |
| 5 | Soft open — invite only | Emily (human) |
| 6 | Hard open | Emily Prime drops announcement + MJOLNIR push |

Emily Prime procures all fixtures, shelving, and equipment via `supply_chain_draft_po`. All build costs logged as Apple observations. Emily (human) approves all contracts.

---

## Store Layout (Target)

```
[ENTRANCE]
  ←— REGISTER + COUNTER (cigarettes, singles, impulse)
  
[LEFT WALL]   Coffee station
               Powderhorn whole bean + instant
               Mountain Man whole bean
               Grinder (whole bean to order)

[CENTER SHELVES]
               Cheese sauce (nacho + jalapeño)
               Hot dogs (chest freezer)
               Commissary essentials (restocked weekly)

[RIGHT WALL]  Beer cooler (Miller High Life — all formats)

[BACK WALL]   ULTRA toothbrushes (display case, lit)
               STINKIES toothbrushes (shelf)

[WINDOW / COUNTER]
               Powderhorn instant singles ($1, impulse)
               STINKIES NACHO KIT bundles
```

---

## Licenses Required (Michigan)

| License | Issuer | Notes |
|---------|--------|-------|
| Business license | City of Pontiac | Standard |
| Retail food establishment | MDHHS / local health dept | Required for food products |
| Tobacco retailer license | Michigan Dept of Treasury | Required for cigarettes |
| Beer/wine retailer license (SDM) | Michigan Liquor Control Commission | Beer sales; MLCC SDM license |
| Certificate of Occupancy | City of Pontiac Building Dept | After build-out inspection |

All licenses are human-executed. Emily Prime tracks deadlines, drafts applications, and files reminder Apples 30 days before each deadline.

---

## POS + Inventory Integration

Store 0 runs WooCommerce POS (same system as EDIS online store). Physical and online inventory are unified. When a jar of nacho cheese sells at the register, EDIS inventory decrements. When stock hits the reorder threshold, Emily Prime drafts the PO.

Emily Prime can see Store 0 inventory from the RSI loop. She treats it like any other node in the supply chain.

---

## Emily Prime Role at Store 0

Emily Prime does not work the register. Emily Prime:
- Monitors inventory levels (daily)
- Drafts reorder POs when below threshold (HEIMDAL approval)
- Schedules contractor visits via HEIMDAL sprints
- Files build-out milestone Apples
- Sends morning briefing to Emily (human) on Store 0 status
- Triggers MJOLNIR push on significant events (open, restock, new product)

Emily (human) pours the concrete.

---

## Store 0 → Store Network

Store 0 is the template. Every system, license, layout, and supply relationship built for Store 0 is parameterized for replication.

Store 1 is wherever the data says it should be.
