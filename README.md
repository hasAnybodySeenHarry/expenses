# Debt Tracker

## System Architecture

```mermaid
graph TD
    style ReactApp fill:#2196F3,stroke:#1976D2,stroke-width:2px,color:#fff,font-size:14px,rx:10,ry:10
    style ReverseProxy fill:#2196F3,stroke:#1976D2,stroke-width:2px,color:#fff,font-size:14px,rx:10,ry:10
    style Expenses fill:#66BB6A,stroke:#388E3C,stroke-width:2px,color:#fff,font-size:14px,rx:10,ry:10
    style Throttler fill:#FFB74D,stroke:#F57C00,stroke-width:2px,color:#fff,font-size:14px,rx:10,ry:10
    style Mailer fill:#AB47BC,stroke:#8E24AA,stroke-width:2px,color:#fff,font-size:14px,rx:10,ry:10
    style Notifier fill:#A1887F,stroke:#6D4C41,stroke-width:2px,color:#fff,font-size:14px,rx:10,ry:10
    style Postgres fill:#B3E5FC,stroke:#039BE5,stroke-width:2px,color:#333,font-size:14px,rx:10,ry:10
    style Redis fill:#FFCDD2,stroke:#E57373,stroke-width:2px,color:#333,font-size:14px,rx:10,ry:10
    style MongoDB fill:#FFE082,stroke:#FFC107,stroke-width:2px,color:#333,font-size:14px,rx:10,ry:10
    style Kafka fill:#CE93D8,stroke:#AB47BC,stroke-width:2px,color:#333,font-size:14px,rx:10,ry:10
    style AMQP fill:#FFF59D,stroke:#FBC02D,stroke-width:2px,color:#333,font-size:14px,rx:10,ry:10
    style CloudEmail fill:#EF9A9A,stroke:#E57373,stroke-width:2px,color:#fff,font-size:14px,rx:10,ry:10

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
    linkStyle 0 stroke:#1976D2,stroke-width:2px
    linkStyle 1 stroke:#1976D2,stroke-width:2px
    linkStyle 2 stroke:#388E3C,stroke-width:2px
    linkStyle 3 stroke:#F57C00,stroke-width:2px
    linkStyle 4 stroke:#F57C00,stroke-width:2px
    linkStyle 5 stroke:#8E24AA,stroke-width:2px
    linkStyle 6 stroke:#8E24AA,stroke-width:2px
    linkStyle 7 stroke:#8E24AA,stroke-width:2px
    linkStyle 8 stroke:#8E24AA,stroke-width:2px
    linkStyle 9 stroke:#6D4C41,stroke-width:2px
    linkStyle 10 stroke:#6D4C41,stroke-width:2px
    linkStyle 11 stroke:#1976D2,stroke-width:2px
    linkStyle 12 stroke:#6D4C41,stroke-width:2px
    linkStyle 13 stroke:#6D4C41,stroke-width:2px
    linkStyle 14 stroke:#6D4C41,stroke-width:2px
