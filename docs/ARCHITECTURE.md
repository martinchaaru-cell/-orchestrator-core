# Architecture

## Components

- **Orchestrator (Go)**: REST API, task queue, agent management
- **Agent (Go)**: Executes tasks locally
- **Tasks**: Bash scripts
- **Dashboard (React)**: Web UI on Vercel
- **Redis**: Task queue storage
