# Security Review — Reference

OWASP Top 10 (2021) detection patterns plus STRIDE framing. Optional reference for Security, and for any role touching auth, data handling, or external integrations. Work the checklists against what the change actually touches.

## A01 — Broken Access Control

- Direct object references with predictable IDs (`/orders/42`) and no ownership check
- Admin-only endpoints without role verification; force-browsing to hidden-but-unprotected pages
- Authorization check missing after the authentication middleware
- CORS allowing untrusted origins
- Go: missing middleware on router groups, `context.Value` abuse bypassing checks. Node: missing `next()`-based checks in the chain. React: client-side route guards with no server-side counterpart — security theater.

Enforce authorization in every server-side handler, centrally, and log denials.

## A02 — Cryptographic Failures

- HTTP anywhere data travels; passwords in plain text or under MD5/SHA-1
- Keys hardcoded in source or committed config; unauthenticated symmetric encryption (AES-CBC with no MAC)
- PII or tokens in URLs, logs, or browser history

Use bcrypt (cost ≥ 12) or argon2, TLS 1.2+, PII encrypted at rest, secrets from env or a manager — never committed.

## A03 — Injection

- String concatenation into SQL; `fmt.Sprintf` into a query
- Dynamic command execution with user-supplied args; `eval()`
- Template injection; LDAP, XPath, NoSQL injection; log injection

Parameterized queries everywhere, ORM query builders over raw interpolation, validation at the input boundary.

## A04 — Insecure Design

- No threat model for the feature
- Business logic flaws (a coupon applicable more than once)
- No rate limiting designed for sensitive operations; no lockout or CAPTCHA on brute-forceable endpoints

## A05 — Security Misconfiguration

- Default credentials; debug endpoints or verbose errors in production
- Missing `Content-Security-Policy`, `X-Frame-Options`, `X-Content-Type-Options`, `Strict-Transport-Security`
- `Access-Control-Allow-Origin: *` on authenticated endpoints; directory listing enabled
- Stack traces returned to clients on 500

## A06 — Vulnerable Components

- Dependencies with known CVEs (`govulncheck`, `npm audit`, `trivy`)
- Unmaintained dependencies (no commits or patches in 2+ years); components running as root

## A07 — Authentication Failures

- Weak passwords permitted; no lockout or rate limit on login
- Predictable session tokens; sessions not invalidated on logout
- Insecure forgot-password flow (predictable or non-expiring tokens); no MFA on sensitive operations
- JWT: verify the `alg` header, reject `none`, validate `exp` and `iss`. Server-side invalidation needs a denylist or short expiry.

## A08 — Integrity Failures

- Deserializing untrusted data; updates fetched over HTTP without integrity checks
- CI/CD modifiable by external pull requests without review; plugins from untrusted sources

## A09 — Logging & Monitoring Failures

- Authentication failures not logged; no audit trail for deletions, permission changes, payments
- Logs modifiable or deletable by application users; no alerting on repeated failures
- PII or secrets in log output

## A10 — SSRF

- Webhook registration accepting arbitrary URLs; URL preview or screenshot generation; import-from-URL
- Internal service URLs reachable from user input; DNS rebinding via user-controlled hostnames

Allowlist schemes and hosts, block RFC 1918 ranges and localhost, prefer a dedicated egress proxy.

## STRIDE framing

The six questions to ask of any trust boundary. Use them when the OWASP categories do not obviously cover the feature's shape.

| Category | The question | Typical tell |
|---|---|---|
| **Spoofing** | Can an attacker claim to be someone they are not? | Forgeable tokens, unverified webhook senders, trusted client-supplied identity |
| **Tampering** | Can data be modified that should not be? | Unencrypted transit, no integrity checks at rest, mass assignment |
| **Repudiation** | Can a user deny an action with no way to prove otherwise? | Missing or subject-deletable audit logs, no immutable state-transition record |
| **Information disclosure** | Can an attacker read what they should not? | Verbose errors, over-returning endpoints, IDOR, shared-layer caching of authenticated responses |
| **Denial of service** | Can resources be exhausted? | Unauthenticated expensive endpoints, no rate limit, ReDoS, no outbound timeouts |
| **Elevation of privilege** | Can someone do more than their role allows? | AuthN without AuthZ, horizontal and vertical escalation, path traversal, insecure deserialization |

A trust boundary is any point where data crosses between trust zones: internet → application, application → database, service → service, browser → server, admin UI → application. For each one, state what crosses it, in which direction, and what validation happens there.

When a threat has several attack paths, name them explicitly rather than in aggregate — "bypass authentication" splits into stealing a session token (XSS, MITM), forging one (weak JWT secret, algorithm confusion), and credential stuffing (no rate limit on login), and those three have different fixes. Carry the highest-risk paths into `risks` as `{risk, mitigation}`.
