# API Design Guidelines

## Purpose

RESTful design principles and contract-first approach. Applied during the `system-design` step of the Architect workflow when the change touches an API surface.

## Resource Naming

- Use **plural nouns** for collections: `/users`, `/orders`, `/reports`.
- Use **singular noun** for singletons: `/config`, `/status`.
- Use **lowercase kebab-case** for multi-word resources: `/audit-logs`, `/password-reset`.
- Nest resources only when the child is always scoped to the parent: `/users/{id}/addresses`.
- Limit nesting to **two levels**. Deeper nesting → prefer flat resources with query params.
- Never use verbs in URLs: use `/orders/{id}/cancel` NOT `/cancelOrder/{id}`. The verb is the HTTP method.

## HTTP Verb Semantics

| Verb | Semantics | Body | Idempotent | Safe |
|------|-----------|------|-----------|------|
| GET | Fetch resource(s) | No | Yes | Yes |
| POST | Create resource or trigger action | Yes | No | No |
| PUT | Replace resource entirely | Yes | Yes | No |
| PATCH | Partial update | Yes | No | No |
| DELETE | Remove resource | No | Yes | No |

Use POST for actions that don't map cleanly to CRUD (e.g., `POST /orders/{id}/cancel`).

## Versioning Strategies

| Strategy | URL example | When to use |
|----------|-------------|-------------|
| **URL path** | `/v1/users` | Public APIs; clear deprecation story |
| **Header** | `Accept: application/vnd.api+json;version=2` | Internal APIs; client controls negotiation |
| **Query param** | `/users?version=2` | Avoid — pollutes resource address |

Preference: URL path versioning for externally consumed APIs. Header versioning for internal service-to-service.

Always plan for `v1` → `v2` migration before writing `v1`. Document the versioning strategy chosen in the ADR.

## Error Response Shape

Consistent error shapes across all endpoints:

```json
{
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "User with ID 42 does not exist.",
    "details": [],
    "request_id": "req_abc123"
  }
}
```

Rules:
- `code`: machine-readable constant in `SCREAMING_SNAKE_CASE`.
- `message`: human-readable, safe to display (no internal details).
- `details`: array of field-level errors for validation failures.
- `request_id`: always present; aids debugging.

HTTP status codes:
- `400` — client sent invalid input
- `401` — unauthenticated
- `403` — authenticated but not authorized
- `404` — resource not found
- `409` — conflict (e.g., duplicate creation)
- `422` — validation error (semantically invalid input)
- `429` — rate limited
- `500` — server error (never expose stack traces)

## Pagination Patterns

| Pattern | Use when |
|---------|----------|
| **Cursor-based** | Large datasets, real-time data, infinite scroll |
| **Offset/limit** | Small, stable datasets; page numbers needed in UI |
| **Keyset** | Ordered data with stable sort key |

Prefer cursor-based for anything that may grow beyond a few thousand records.

Response envelope for paginated collections:
```json
{
  "data": [],
  "pagination": {
    "next_cursor": "opaque_string_or_null",
    "has_more": true
  }
}
```

## OpenAPI-First Approach

1. Write the OpenAPI spec (`openapi.yaml`) before writing any handler code.
2. Use code generation to produce request/response types from the spec.
3. Validate requests against the spec at the handler boundary — never trust unvalidated input.
4. Document all endpoints, including error responses, in the spec.

In the `architect/system-design` payload, list every endpoint as an `api_surface[]` entry: the resource path in `operation`, the HTTP verb in `method`, the request shape in `input`, the response shape in `success_response`, and each failure mode in `error_cases`. The schema has no dedicated auth field — state the auth requirement in the `operation` text (e.g. `GET /users/{id} — requires authenticated caller with users:read`) and add the matching `401` / `403` entries to `error_cases`.
