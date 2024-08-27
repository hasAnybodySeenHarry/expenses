# Debt Tracker

## System Architecture

```mermaid
graph TD
    style ReactApp fill:#4fc3f7,stroke:#0288d1,stroke-width:2px,color:#fff
    style ReverseProxy fill:#4fc3f7,stroke:#0288d1,stroke-width:2px,color:#fff
    style Expenses fill:#81c784,stroke:#388e3c,stroke-width:2px,color:#fff
    style Throttler fill:#ffb74d,stroke:#f57c00,stroke-width:2px,color:#fff
    style Mailer fill:#ba68c8,stroke:#8e24aa,stroke-width:2px,color:#fff
    style Notifier fill:#8d6e63,stroke:#5d4037,stroke-width:2px,color:#fff
    style Postgres fill:#a6cee3,stroke:#1e88e5,stroke-width:2px,color:#333
    style Redis fill:#ef9a9a,stroke:#e53935,stroke-width:2px,color:#333
    style MongoDB fill:#ffcc80,stroke:#fb8c00,stroke-width:2px,color:#333
    style Kafka fill:#ce93d8,stroke:#8e24aa,stroke-width:2px,color:#333
    style AMQP fill:#fff176,stroke:#fbc02d,stroke-width:2px,color:#333
    style CloudEmail fill:#e57373,stroke:#d32f2f,stroke-width:2px,color:#fff

    subgraph Frontend
        ReactApp[React App]
    end

    subgraph Backend
        ReverseProxy[Reverse Proxy]
        Expenses[Expenses Service]
        Throttler[Throttler Service]
        Mailer[Mailer Service]
        Notifier[Notifier Service]
    end

    subgraph Databases
        Postgres[Postgres Database]
        Redis[Redis]
        MongoDB[MongoDB]
    end

    subgraph Messaging
        Kafka[Kafka]
        AMQP[AMQP Proxy]
    end

    subgraph ExternalServices
        CloudEmail[Cloud Email Service]
    end

    %% Define connections with labels and directions
    ReactApp -- HTTP --> ReverseProxy
    ReverseProxy -- HTTP --> Expenses
    Expenses -- User Data --> Postgres
    Throttler -- gRPC --> Expenses
    Throttler -- Rate Limit Buckets --> Redis
    Expenses -- AMQP --> AMQP
    Mailer -- Fetch Job --> AMQP
    Mailer -- Send Email --> CloudEmail
    Expenses -- Send Event --> Kafka
    Notifier -- Consume Events --> Kafka
    Notifier -- Store Notifications --> MongoDB
    ReactApp -- WebSocket --> ReverseProxy
    ReverseProxy -- HTTP --> Notifier
    Notifier -- gRPC --> Expenses
    Notifier -- WebSocket --> ReactApp

    %% Style links
    linkStyle 0 stroke:#0288d1,stroke-width:2px
    linkStyle 1 stroke:#388e3c,stroke-width:2px
    linkStyle 2 stroke:#388e3c,stroke-width:2px
    linkStyle 3 stroke:#f57c00,stroke-width:2px
    linkStyle 4 stroke:#f57c00,stroke-width:2px
    linkStyle 5 stroke:#8e24aa,stroke-width:2px
    linkStyle 6 stroke:#8e24aa,stroke-width:2px
    linkStyle 7 stroke:#8e24aa,stroke-width:2px
    linkStyle 8 stroke:#8e24aa,stroke-width:2px
    linkStyle 9 stroke:#5d4037,stroke-width:2px
    linkStyle 10 stroke:#5d4037,stroke-width:2px
    linkStyle 11 stroke:#0288d1,stroke-width:2px
    linkStyle 12 stroke:#5d4037,stroke-width:2px
    linkStyle 13 stroke:#5d4037,stroke-width:2px
    linkStyle 14 stroke:#5d4037,stroke-width:2px
