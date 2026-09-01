# Indexing & Query Access Patterns

Use this reference to design high-performance indexing strategies aligned with specific application access patterns.

---

## 1. Index Type Selection

| Index Type | Typical Use Case | Example (PostgreSQL) |
| :--- | :--- | :--- |
| **B-Tree (Default)** | Exact match (`=`), ranges (`<`, `>`, `BETWEEN`), sorting (`ORDER BY`). | `CREATE INDEX idx_orders_created ON orders(created_at DESC);` |
| **GIN (Generalized Inverted)** | Full-text search, arrays, JSONB containment (`@>`, `?`, `?|`). | `CREATE INDEX idx_audit_meta ON audit_logs USING gin (payload jsonb_path_ops);` |
| **Partial Index** | Highly skewed queries (e.g. only active records, unhandled jobs). | `CREATE INDEX idx_jobs_pending ON jobs(priority DESC) WHERE status = 'pending';` |
| **Covering Index (`INCLUDE`)** | Index-Only Scans avoiding table heap lookups. | `CREATE INDEX idx_users_lookup ON users(email) INCLUDE (full_name, status);` |

---

## 2. Composite Index Column Ordering Rule

When creating multi-column composite indexes `(col1, col2, col3)`:
1. **Equality columns first** (`WHERE col1 = 'val'`)
2. **Range / Sort columns last** (`WHERE col2 > 100 ORDER BY col3`)

*Example:* For query `SELECT * FROM orders WHERE tenant_id = ? AND status = ? AND created_at > ? ORDER BY created_at DESC`:
```sql
CREATE INDEX idx_orders_tenant_status_date ON orders(tenant_id, status, created_at DESC);
```
