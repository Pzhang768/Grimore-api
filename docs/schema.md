# Database Schema

## Entity Relationships

```
User (1) ──────────────── (many) Team
                                   │
                          (many) TeamAgent  ← composite PK: (TeamID, AgentType)
                                   │
Team (1) ──────────────── (many) Run
                                   │
                          (many) RunEvent
                                   │
                          (many) Deliverable
```

A **User** owns one or more **Teams**. Each Team configures an ordered list of **TeamAgents** — the agent pipeline. When a pipeline is triggered, a **Run** is created for that Team. As the Run executes, each agent appends **RunEvents** (log entries) and produces **Deliverables** (structured outputs).

## Tables

### users
| Column | Type | Constraints |
|---|---|---|
| id | uuid | PK |
| email | text | unique, not null |
| created_at | timestamptz | |

### teams
| Column | Type | Constraints |
|---|---|---|
| id | uuid | PK |
| user_id | uuid | FK → users.id RESTRICT, not null, index |
| name | text | not null |
| created_at | timestamptz | |

### team_agents
| Column | Type | Constraints |
|---|---|---|
| team_id | uuid | PK (composite), FK → teams.id RESTRICT |
| agent_type | text | PK (composite), not null |
| context | jsonb | |
| position | int | not null |

Composite PK `(team_id, agent_type)` enforces one agent type per team.

### runs
| Column | Type | Constraints |
|---|---|---|
| id | uuid | PK |
| team_id | uuid | FK → teams.id RESTRICT, not null, index |
| status | text | not null |
| iteration | int | not null, default 0 |
| created_at | timestamptz | |
| updated_at | timestamptz | |

### run_events
| Column | Type | Constraints |
|---|---|---|
| id | uuid | PK |
| run_id | uuid | FK → runs.id RESTRICT, not null, index |
| agent_type | text | not null |
| event_type | text | not null |
| content | text | |
| created_at | timestamptz | |

### deliverables
| Column | Type | Constraints |
|---|---|---|
| id | uuid | PK |
| run_id | uuid | FK → runs.id RESTRICT, not null |
| type | text | not null |
| content | jsonb | not null, default '{}' |
| created_at | timestamptz | |
| updated_at | timestamptz | |

Composite index `idx_deliverables_run_type` on `(run_id, type)`.

## Notes

* All PKs are UUIDs generated in Go via `BeforeCreate` hooks.
* All FK constraints use `ON DELETE RESTRICT` — deleting a parent record is refused if children exist. The application must clean up explicitly.
* `AutoMigrate` only adds columns and indexes; destructive changes (column drops, renames) require manual migration scripts.
