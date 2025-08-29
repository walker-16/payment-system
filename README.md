# payment-system

A resilient, event-driven payment processing system built with a microservices architecture.  
It handles payment orchestration, wallet management, and external gateway communication with a focus on fault tolerance and scalability.

## Services

The project is structured as a monorepo containing the following services:

- **Payment**: Handles payment orchestration, persistence, and status tracking.  
  _(currently under construction)_

- **Wallet**: Responsible for managing user balances and wallet operations.  
  _(pending development)_

- **Processor**: Consumes events and interacts with external gateways or banks to execute payments.  
  _(pending development)_

---

## Project Structure

## Getting Started (Payment Service)

### Prerequisites

- [Go 1.24.4+](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/) running locally
- [Docker](https://docs.docker.com/get-docker/) (recommended for running Postgres)

### Environment Variables

The **Payment** service requires the following environment variables:

| Variable               | Description                                  | Default / Example                                                    |
|-------------------------|----------------------------------------------|----------------------------------------------------------------------|
| `LOG_LEVEL`             | Logging level                                | `INFO`                                                               |
| `PORT`                  | Port where the service will run              | `8200`                                                               |
| `DB_DNS`                | PostgreSQL connection string                 | `postgres://postgres:postgres@localhost:15432/local`                 |
| `DB_MAX_CONNS`          | Maximum number of DB connections             | `10`                                                                 |
| `DB_MIN_CONNS`          | Minimum number of DB connections             | `0`                                                                  |
| `DB_MAX_CONN_IDLE_TIME` | Maximum idle time for DB connections         | `30m`                                                                |
| `DB_MAX_CONN_LIFETIME`  | Maximum lifetime for DB connections          | `1h`                                                                 |

---

Or using VS Code with the provided launch.json configuration:

```json
{
  "name": "Launch Package",
  "type": "go",
  "request": "launch",
  "mode": "auto",
  "program": "${workspaceFolder}/services/payment/cmd",
  "env": {
    "LOG_LEVEL": "INFO",
    "PORT": "8200",
    "DB_DNS": "postgres://postgres:postgres@localhost:15432/local",
    "DB_MAX_CONNS": "10",
    "DB_MIN_CONNS": "0",
    "DB_MAX_CONN_IDLE_TIME": "30m",
    "DB_MAX_CONN_LIFETIME": "1h"
  },
  "args": []
}
