# Diagram Patterns

Use this when: producing Mermaid or D2 technical architecture views after the audience and question have been selected.

Use [visual-communication.md](visual-communication.md) for Level 0 maps, view selection, accessibility, SVG delivery, infographics, and upgrade overlays. This file owns technical C4, sequence, deployment, and diagrams-as-code source patterns.

## Contents

- [C4 context with Mermaid](#c4-context-with-mermaid)
- [C4 container with Mermaid](#c4-container-with-mermaid)
- [Container view with D2](#container-view-with-d2)
- [Critical workflow with Mermaid](#critical-workflow-with-mermaid)
- [Deployment view with Mermaid](#deployment-view-with-mermaid)
- [Diagram checks](#diagram-checks)

## C4 context with Mermaid

Show people, the system of interest as one box, and external software systems. Do not expose internal services here.

```mermaid
flowchart LR
    accTitle: Commerce platform system context
    accDescr: Customers and operators use the commerce platform, which connects to payment and carrier systems.

    customer["Person: Customer"]
    operator["Person: Operations user"]
    platform["Software System: Commerce Platform<br/>Places and fulfils orders"]
    payments["External System: Payment Provider"]
    carrier["External System: Carrier Network"]

    customer -->|"Browses and places orders over HTTPS"| platform
    operator -->|"Manages exceptions over HTTPS"| platform
    platform -->|"Authorizes payments via REST"| payments
    platform -->|"Books shipments via REST"| carrier
```

## C4 container with Mermaid

Show independently runnable applications and data stores within the system boundary. Name technology only when chosen and relevant.

```mermaid
flowchart LR
    accTitle: Commerce platform container view
    accDescr: The web application calls the order API, which owns order data and publishes fulfilment work.

    customer["Person: Customer"]
    payments["External System: Payment Provider"]

    subgraph platform["Software System: Commerce Platform"]
        web["Container: Web Application<br/>Customer experience"]
        api["Container: Order API<br/>Owns checkout workflow"]
        worker["Container: Fulfilment Worker<br/>Processes placed orders"]
        orders[("Container: Order Database<br/>Authoritative order state")]
        broker[["Container: Message Broker<br/>Durable work distribution"]]
    end

    customer -->|"Uses over HTTPS"| web
    web -->|"Submits orders; JSON/HTTPS"| api
    api -->|"Reads and writes; SQL/TLS"| orders
    api -->|"Authorizes; REST/HTTPS"| payments
    api -->|"Publishes OrderPlaced; at least once"| broker
    broker -->|"Delivers OrderPlaced; at least once"| worker
```

## Container view with D2

Use D2 when compact grouping and portable source files are preferable.

```d2
direction: right

customer: Customer { shape: person }
payments: Payment Provider { shape: rectangle }

platform: Commerce Platform {
  web: Web Application
  api: Order API
  worker: Fulfilment Worker
  orders: Order Database { shape: cylinder }
  broker: Message Broker { shape: queue }

  web -> api: Submit order / JSON HTTPS
  api -> orders: Read and write / SQL TLS
  api -> broker: OrderPlaced / at-least-once
  broker -> worker: OrderPlaced / at-least-once
}

customer -> platform.web: Uses / HTTPS
platform.api -> payments: Authorize / REST HTTPS
```

## Critical workflow with Mermaid

Show the primary path plus timeout, duplicate, or asynchronous failure behavior. Match the guarantees described in prose.

```mermaid
sequenceDiagram
    accTitle: Order placement workflow
    accDescr: The API deduplicates a request, authorizes payment, commits the order and outbox event, then starts asynchronous fulfilment.

    autonumber
    actor Customer
    participant API as Order API
    participant DB as Order Database
    participant Pay as Payment Provider
    participant Bus as Message Broker
    participant Worker as Fulfilment Worker

    Customer->>API: PlaceOrder(command, idempotencyKey)
    API->>DB: Find prior result by idempotency key
    alt Request already completed
        DB-->>API: Existing order result
        API-->>Customer: Existing result
    else New request
        API->>Pay: Authorize(amount, idempotencyKey)
        alt Authorization succeeds
            Pay-->>API: Authorization ID
            API->>DB: Commit order + outbox event
            DB-->>API: Order accepted
            API-->>Customer: 202 Accepted
            DB-->>Bus: Relay OrderPlaced
            Bus-->>Worker: Deliver OrderPlaced
            Worker->>Worker: Process idempotently
        else Timeout or rejection
            Pay-->>API: Failure or unknown result
            API-->>Customer: Retry-safe failure
        end
    end
```

## Deployment view with Mermaid

Use a separate view for infrastructure and failure domains. Do not call cloud resources C4 components.

```mermaid
flowchart TB
    accTitle: Commerce platform deployment
    accDescr: Traffic reaches application and worker replicas across two zones backed by multi-zone data and messaging services.

    internet((Internet))
    edge[Global edge and WAF]

    subgraph region[Region]
        subgraph zoneA[Zone A]
            apiA[Order API replicas]
            workerA[Fulfilment workers]
        end
        subgraph zoneB[Zone B]
            apiB[Order API replicas]
            workerB[Fulfilment workers]
        end
        db[(Multi-zone order database)]
        bus[[Durable message broker]]
    end

    internet --> edge
    edge --> apiA
    edge --> apiB
    apiA --> db
    apiB --> db
    apiA --> bus
    apiB --> bus
    bus --> workerA
    bus --> workerB
```

## Diagram checks

- State the diagram type, scope, audience, and question.
- Apply the accessibility, rendering, inspection, and delivery contract in [visual communication](visual-communication.md#deliver-and-validate-visuals).
- Use one meaning for each shape and line style and explain non-obvious semantics.
- Label relationships with action and protocol, event, or delivery guarantee.
- Draw separate relationships when direction or semantics differ.
- Keep logical views vendor-neutral; put infrastructure products in deployment views.
- Put failure behavior in sequence or state views rather than structural views.
- Store editable source near architecture documentation or ADRs; render in CI when the repository supports it.
