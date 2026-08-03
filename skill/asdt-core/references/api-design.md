# API & System Design — Reference

Criteria for judging an API surface, a scaling story, and the blast radius of a structural change. Optional reference for Architect and Developer. This is judgment, not procedure — apply what the change actually touches.

## API surface

- Plural nouns for collections (`/users`), singular for singletons (`/config`), lowercase kebab-case for multi-word paths (`/audit-logs`).
- No verbs in URLs: `POST /orders/{id}/cancel`, never `/cancelOrder/{id}`. The verb is the HTTP method.
- Nest only when the child is always scoped to the parent, and never deeper than two levels. Past that, flatten and use query params.
- Verb semantics are a contract: GET safe and idempotent, PUT full replace and idempotent, DELETE idempotent, POST and PATCH neither. Use POST for actions that do not map to CRUD.
- Versioning: URL path (`/v1/users`) for externally consumed APIs, header negotiation for internal service-to-service, never a query param. Decide the `v1 → v2` migration story before shipping `v1`.
- One error shape everywhere: machine-readable `code` in `SCREAMING_SNAKE_CASE`, human-safe `message` with no internals, field-level `details`, and an always-present `request_id`.
- Status codes carry meaning: `400` malformed, `401` unauthenticated, `403` unauthorized, `404` absent, `409` conflict, `422` semantically invalid, `429` rate limited, `500` server error with no stack trace.
- Pagination: cursor-based for anything that may exceed a few thousand rows, offset/limit only for small stable sets that need page numbers, keyset for stable-sorted data.
- Contract first: write the spec before the handler, generate types from it, and validate requests at the handler boundary. Never trust unvalidated input.

State each endpoint in `api_surface` with its path, method, request and response shape, and failure modes. Auth requirements ride the operation text (`GET /users/{id} — requires users:read`) with matching `401`/`403` cases.

## Scaling

Look for these before calling a design done:

- **N+1 queries** — any loop that calls the data layer, or an ORM walk over `item.related_data`. Fix with eager loading, batch loading, or a join.
- **Missing indexes** — for each query, check the WHERE / JOIN / ORDER BY columns and whether a composite index covers them most-selective-first. Any table expected past 10k rows needs this analysis before shipping.
- **Synchronous chains** — latency stacks along `A → B → C`. Chains longer than 3 hops with p95 over 500ms are a High risk; name which calls are on the critical path and which can be parallelized or deferred.
- **Unbounded work** — expensive operations on unauthenticated endpoints, uploads without limits, recursion driven by user input.

Scaling strategy is a choice with a cost, so name it and say which modules it affects: vertical (fast, low ceiling, no code change), horizontal (requires stateless design), read replicas (read-heavy, eventual consistency acceptable), sharding (large data, operationally expensive), event-driven (bursty, eventual consistency).

Caching is a consistency decision before it is a performance one. Pick the layer — in-process for rarely-changed lookup data, distributed for state shared across instances, HTTP/ETag for public reads, CDN for static assets — and state the TTL and the invalidation approach (TTL-based, write-through, or cache-aside) in the same breath. Never introduce a cache without saying how it is invalidated.

Going async buys latency and pays in eventual consistency: fire-and-forget with `202 Accepted`, a background job returning a job id, or a domain event with independent consumers. Say which guarantee the feature can live with, and reflect it in the response shape.

Estimate load for any new endpoint or job — expected RPS at p50 and p99, bytes per operation, 30-day and 1-year growth. If the numbers do not exist, record it in `open_items` and design for 10x current load.

## Blast radius

Before proposing a structural change, map what is already coupled to the affected area: direct callers (which modules import it), data-shape consumers (same tables, event schemas, API contracts), configuration dependents (config keys, env vars, feature flags), and test dependents.

| Blast radius | Definition | What it demands |
|---|---|---|
| **Contained** | ≤ 1 package, no external contract changes | Proceed; note it in the design |
| **Moderate** | 2–4 packages or one external contract | Document the migration path |
| **Wide** | 5+ packages or multiple external contracts | High risk; phased migration required |
| **Cross-boundary** | Public API, event schema, or shared DB table | Critical risk; backward-compat analysis required |

For any interface or data contract change: can old and new coexist, and for how long? Is there a version field, and if not, should there be? Are there consumers outside this repository? An unanswered one of these is an `open_items` entry, not a silent assumption.

When a breaking change is unavoidable, the path is expand → migrate → contract: introduce the new interface beside the old, move every caller, then delete the old one. State the phases explicitly in the consequences of the decision, because the contract phase is the one everyone forgets.
