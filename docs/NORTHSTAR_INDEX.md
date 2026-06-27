# NORTHSTAR — EINHORN INDEX

**Classification:** Tier 1 — strategic platform  
**Owner:** Emily Springerton / EINHORN_INDUSTRIAL  
**Version:** 0.1 | Updated: 2026-06-27  
**Gates:** S138-03 through S138-06

---

## Mission

The EINHORN INDEX is a continuously updated structured knowledge graph of financial, supply chain, and AI ops intelligence. Emily Prime is the primary consumer; the INDEX is eventually exposed as a data product ($49/mo seed tier). The INDEX is EINHORN_INDUSTRIAL's long-term data moat.

---

## Crawl Domain Strategy

| Domain | Sources | Update cadence |
|--------|---------|---------------|
| Supply chain | Supplier directories, shipping data, commodity feeds, industry publications | Hourly (alert events), Daily (background) |
| Financial signals | SEC EDGAR, company press releases, earnings transcripts, investor decks | Daily |
| AI ops | Arxiv (cs.AI, cs.LG), HN threads, GitHub release notes, model API changelogs | Daily |
| General web | Open web via Research engine (S137-01) | On demand |

The crawler is the bootstrap spider from `PRRJECT_FATBABY/internal/spider`. The INDEX pipeline runs as a cron job separate from the Emily Prime RSI loop.

---

## Entity Type Taxonomy

### Node types

| Type | Examples | Canonical ID strategy |
|------|---------|----------------------|
| `company` | Apple Inc, StickersForAll LLC | CIK (SEC) or domain |
| `person` | Tim Cook | Full name + org affiliation |
| `product` | iPhone 17, die-cut vinyl sticker | SKU or product page URL |
| `location` | Shenzhen Manufacturing Zone | ISO 3166 + subdivision |
| `price_point` | $1.60/unit at MOQ 250 | value + currency + unit + date |
| `event` | Q2 2026 earnings miss | event_type + entity_id + date |
| `technology` | GPT-2, RLHF, IDUNA JWT | canonical name + version |

### Edge (relationship) types

| Predicate | Subject → Object | Notes |
|-----------|-----------------|-------|
| `acquired` | company → company | with date |
| `sells` | company → product | with price_point |
| `employs` | company → person | with role |
| `located_in` | company/product → location | |
| `competes_with` | company → company | bidirectional |
| `supplies` | company → company | supply chain edge |
| `priced_at` | product → price_point | versioned by date |
| `cited_in` | entity → research_log Apple | audit provenance |

---

## Graph Schema

### MongoDB: collection `kg_nodes`

```json
{
  "_id": "<sha256 of entity_type:canonical_name>",
  "entity_type": "company",
  "canonical_name": "StickerMule",
  "aliases": ["Sticker Mule", "stickermule.com"],
  "properties": {
    "domain": "stickermule.com",
    "moq": 25,
    "unit_cost_cents_250": 110
  },
  "first_seen": "2026-06-27T00:00:00Z",
  "last_updated": "2026-06-27T12:00:00Z",
  "source_urls": ["https://www.stickermule.com/products/die-cut-stickers"]
}
```

### MongoDB: collection `kg_edges`

```json
{
  "_id": "<sha256 of subject_id:predicate:object_id>",
  "subject_id": "<company-stickermule-id>",
  "predicate": "sells",
  "object_id": "<product-die-cut-vinyl-sticker-id>",
  "confidence": 0.95,
  "source_url": "https://www.stickermule.com/products/die-cut-stickers",
  "extracted_at": "2026-06-27T12:00:00Z"
}
```

---

## Indexing Pipeline

```
Crawl (spider)
    ↓
Document store (raw HTML + text, PRRJECT_FATBABY event store)
    ↓
Entity extraction (claude-haiku, structured JSON output)
    ↓
Dedup + canonicalization (alias resolution, entity merge)
    ↓
MongoDB upsert (kg_nodes + kg_edges)
    ↓
research_log Apple (source provenance audit)
    ↓
research_cache in IDUNA (query-level cache invalidation)
```

Pipeline runs as a goroutine in PRRJECT_FATBABY obs-watcher or as a standalone `einhorn-index` binary. Emily Prime is the query consumer, not the pipeline operator.

---

## Query Interface

### emily-agent KnowledgeQuery tool (S138-05)

```
KnowledgeQuery(entity: "StickerMule", predicate: "sells") → GraphResult[]
KnowledgeQuery(entity: "vinyl sticker", predicate: "") → all edges where subject or object matches
```

Returns: `[]GraphResult{SubjectName, Predicate, ObjectName, Confidence, SourceURL, ExtractedAt}`

Emily uses KnowledgeQuery before Research() for any factual lookup:
1. `KnowledgeQuery` → hit: return, file metadata Apple
2. Miss or low confidence → `Research()` → entity extraction → graph upsert → return

### External API (S138-06): `GET /api/v1/index/query`

```
Authorization: Bearer <IDUNA JWT>
?entity=StickerMule&predicate=sells&limit=10
→ {"results": [...], "query_id": "<uuid>", "billed_at": "2026-06-27T..."}
```

Rate limit: 1000 req/mo at $49/mo seed tier. Each billed query files a billing Apple.

---

## Emily Integration

Emily Prime decision tree for any factual question:

```
Is it in emily-memory? → use memory
  ↓ no
KnowledgeQuery → hit (confidence ≥ 0.8) → use result
  ↓ miss or low confidence
Research() → hit → extract entities → upsert graph → use result
  ↓ miss
Escalate to Emily (human)
```

---

## Long-Term Product Path

| Phase | What | Revenue |
|-------|------|---------|
| INDEX-0 | Emily Prime internal only | None (cost center) |
| INDEX-1 | Seed tier API ($49/mo, 1000 req) | Platform revenue track |
| INDEX-2 | Pro tier + domain specialization | TBD |
| INDEX-3 | Data licensing to third parties | TBD |

INDEX-3 is the strategy document S26 data licensing track.

---

## Success Metrics (INDEX-0)

- ≥ 100 kg_nodes populated by end of S138 sprint
- ≥ 50 kg_edges with confidence ≥ 0.8
- KnowledgeQuery tool used by Emily Prime ≥ 3× per day in RSI loop
- Zero sourceless claims — every KG fact has a research_log Apple
