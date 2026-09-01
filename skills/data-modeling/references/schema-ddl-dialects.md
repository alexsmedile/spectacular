# Schema DDL Dialects & ORM Mappings

Use this reference when generating executable physical schemas or ORM definitions from logical ER diagrams.

---

## 1. Type Translation Matrix

| Logical Type | PostgreSQL | MySQL 8.0+ | SQLite | Prisma | Drizzle ORM |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **UUID** | `UUID` | `VARCHAR(36)` / `BINARY(16)` | `TEXT` | `String @id @default(uuid())` | `uuid().primaryKey()` |
| **Auto-Int** | `BIGSERIAL` / `INT GENERATED` | `BIGINT AUTO_INCREMENT` | `INTEGER PRIMARY KEY AUTOINCREMENT` | `Int @id @default(autoincrement())` | `serial().primaryKey()` |
| **String** | `VARCHAR(n)` / `TEXT` | `VARCHAR(n)` / `TEXT` | `TEXT` | `String` | `varchar({ length: n })` / `text()` |
| **Decimal/Currency** | `NUMERIC(12,2)` | `DECIMAL(12,2)` | `REAL` | `Decimal @db.Decimal(12,2)` | `decimal({ precision: 12, scale: 2 })` |
| **Timestamp (UTC)** | `TIMESTAMPTZ` | `DATETIME(6)` / `TIMESTAMP` | `TEXT` (ISO8601) | `DateTime @default(now())` | `timestamp({ withTimezone: true })` |
| **JSON/Doc** | `JSONB` | `JSON` | `TEXT` | `Json` | `jsonb()` |
| **Boolean** | `BOOLEAN` | `TINYINT(1)` / `BOOLEAN` | `INTEGER` (0/1) | `Boolean` | `boolean()` |

---

## 2. PostgreSQL DDL Standards

```sql
-- Always use strict schemas, foreign keys, timestamps, and indexes
CREATE TABLE customers (
    customer_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE orders (
    order_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(customer_id) ON DELETE RESTRICT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    total_amount NUMERIC(12, 2) NOT NULL CHECK (total_amount >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_orders_status ON orders(status);
```

---

## 3. Drizzle ORM Standard (TypeScript)

```typescript
import { pgTable, uuid, varchar, numeric, timestamp, jsonb } from 'drizzle-orm/pg-core';
import { relations } from 'drizzle-orm';

export const customers = pgTable('customers', {
  id: uuid('id').primaryKey().defaultRandom(),
  email: varchar('email', { length: 255 }).notNull().unique(),
  fullName: varchar('full_name', { length: 255 }).notNull(),
  metadata: jsonb('metadata').default({}).notNull(),
  createdAt: timestamp('created_at', { withTimezone: true }).defaultNow().notNull(),
});

export const orders = pgTable('orders', {
  id: uuid('id').primaryKey().defaultRandom(),
  customerId: uuid('customer_id').notNull().references(() => customers.id, { onDelete: 'restrict' }),
  status: varchar('status', { length: 50 }).default('pending').notNull(),
  totalAmount: numeric('total_amount', { precision: 12, scale: 2 }).notNull(),
  createdAt: timestamp('created_at', { withTimezone: true }).defaultNow().notNull(),
});

export const customersRelations = relations(customers, ({ many }) => ({
  orders: many(orders),
}));
```

---

## 4. Prisma Schema Standard

```prisma
model Customer {
  id        String   @id @default(uuid())
  email     String   @unique
  fullName  String   @map("full_name")
  metadata  Json     @default("{}")
  createdAt DateTime @default(now()) @map("created_at")
  orders    Order[]

  @@map("customers")
}

model Order {
  id          String   @id @default(uuid())
  customerId  String   @map("customer_id")
  status      String   @default("pending")
  totalAmount Decimal  @map("total_amount") @db.Decimal(12, 2)
  createdAt   DateTime @default(now()) @map("created_at")
  customer    Customer @relation(fields: [customerId], references: [id], onDelete: Restrict)

  @@index([customerId])
  @@map("orders")
}
```
