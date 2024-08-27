# Debt Tracker

## System Architecture

```mermaid
graph TD
    style ReactApp fill:#1f78b4,stroke:#333,stroke-width:2px,color:#fff
    style ReverseProxy fill:#1f78b4,stroke:#333,stroke-width:2px,color:#fff
    style Expenses fill:#33a02c,stroke:#333,stroke-width:2px,color:#fff
    style Throttler fill:#ff7f00,stroke:#333,stroke-width:2px,color:#fff
    style Mailer fill:#6a3d9a,stroke:#333,stroke-width:2px,color:#fff
    style Notifier fill:#b15928,stroke:#333,stroke-width:2px,color:#fff
    style Postgres fill:#a6cee3,stroke:#333,stroke-width:2px,color:#333
    style Redis fill:#fb9a99,stroke:#333,stroke-width:2px,color:#333
    style MongoDB fill:#fdbf6f,stroke:#333,stroke-width:2px,color:#333
    style Kafka fill:#cab2d6,stroke:#333,stroke-width:2px,color:#333
    style AMQP fill:#ffff99,stroke:#333,stroke-width:2px,color:#333
    style CloudEmail fill:#e31a1c,stroke:#333,stroke-width:2px,color:#fff

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
    linkStyle 0 stroke:#1f78b4,stroke-width:2px
    linkStyle 1 stroke:#33a02c,stroke-width:2px
    linkStyle 2 stroke:#33a02c,stroke-width:2px
    linkStyle 3 stroke:#ff7f00,stroke-width:2px
    linkStyle 4 stroke:#ff7f00,stroke-width:2px
    linkStyle 5 stroke:#6a3d9a,stroke-width:2px
    linkStyle 6 stroke:#6a3d9a,stroke-width:2px
    linkStyle 7 stroke:#6a3d9a,stroke-width:2px
    linkStyle 8 stroke:#cab2d6,stroke-width:2px
    linkStyle 9 stroke:#b15928,stroke-width:2px
    linkStyle 10 stroke:#b15928,stroke-width:2px
    linkStyle 11 stroke:#1f78b4,stroke-width:2px
    linkStyle 12 stroke:#b15928,stroke-width:2px
    linkStyle 13 stroke:#b15928,stroke-width:2px
    linkStyle 14 stroke:#b15928,stroke-width:2px
