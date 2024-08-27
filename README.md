# Debt Tracker

## System Architecture

```mermaid
graph TD
    subgraph Frontend
        ReactApp[React App]
    end

    subgraph Backend
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

    ReactApp -->|HTTP/gRPC| Expenses
    Expenses -->|User Data| Postgres
    Expenses -->|Rate Limit Check| Throttler
    Throttler -->|RPC| Expenses
    Throttler -->|Rate Limit Data| Redis
    Expenses -->|Mailing Job| AMQP
    Mailer -->|Fetch Job| AMQP
    Mailer -->|Send Email| CloudEmail[Cloud Email Service]
    Expenses -->|Event| Kafka
    Notifier -->|Consume Events| Kafka
    Notifier -->|Store Notifications| MongoDB
    Notifier -->|Send Notifications| WebSocket