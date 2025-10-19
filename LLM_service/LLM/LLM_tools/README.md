# LLM Tools Overview
Таже диаграмма в mermaid-diagram-2025-09-29-112412.png
```mermaid
flowchart TD
    subgraph Agent Runtime
        A[LLM Agent]
    end

    subgraph Package
        R[ToolRegistry]
        B[BaseTool]
        L[LOGGER]
        U[Environment Helpers]
        C1[CSV Analyzer]
        CP[Postgres Connector]
        CC[ClickHouse Connector]
        GP[Postgres DDL Generator]
        GC[ClickHouse DDL Generator]
        GA[Airflow DAG Generator]
        FM[Markdown Formatter]
    end

    subgraph External Systems
        ENV[(.env / environment)]
        DB1[(PostgreSQL)]
        DB2[(ClickHouse)]
        Files[(CSV Files)]
    end

    A -->|"list/get tools"| R
    R -->|"instance"| B
    B -->|"validate & execute"| C1
    B -->|"validate & execute"| CP
    B -->|"validate & execute"| CC
    B -->|"execute"| GP
    B -->|"execute"| GC
    B -->|"execute"| GA
    B -->|"execute"| FM

    C1 --> Files
    CP -->|"connect"| DB1
    CC -->|"connect"| DB2
    CP --> U
    CC --> U
    U --> ENV

    C1 --> L
    CP --> L
    CC --> L
    GP --> L
    GC --> L
    GA --> L
    FM --> L

    C1 -->|"ToolOutput"| A
    CP -->|"ToolOutput"| A
    CC -->|"ToolOutput"| A
    GP -->|"ToolOutput"| A
    GC -->|"ToolOutput"| A
    GA -->|"ToolOutput"| A
    FM -->|"ToolOutput"| A
```

## Environment Configuration

Connection tools automatically pull missing parameters from environment variables. Populate these keys (e.g. in `.env`) to avoid passing credentials explicitly:

- PostgreSQL:
  - `POSTGRES_HOST`
  - `POSTGRES_PORT`
  - `POSTGRES_DB`
  - `POSTGRES_USER`
  - `POSTGRES_PASSWORD`
  - `POSTGRES_USE_SSL`
  - `POSTGRES_TIMEOUT`

- ClickHouse:
  - `CLICKHOUSE_HOST`
  - `CLICKHOUSE_PORT`
  - `CLICKHOUSE_DB`
  - `CLICKHOUSE_USER`
  - `CLICKHOUSE_PASSWORD`
  - `CLICKHOUSE_USE_SSL`
  - `CLICKHOUSE_TIMEOUT`

Missing values in tool calls are filled from these variables via `LLM_tools.utils.env`.


