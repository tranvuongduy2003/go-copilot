# Technical Planner Command

Create technical designs and implementation plans for features following DDD + CQRS patterns.

## Task: $ARGUMENTS

## Planning Workflow

### 1. Requirements Analysis

Before designing, understand:
- What problem are we solving?
- Who are the users?
- What are the constraints?
- What are the non-functional requirements?

### 2. Technical Design Document Template

```markdown
# Technical Design: [Feature Name]

## Overview

### Problem Statement
[What problem does this feature solve?]

### Goals
- Goal 1
- Goal 2

### Non-Goals
- What this feature will NOT do

## Background

### Current State
[How does the system work today?]

### Proposed Solution
[High-level description of the solution]

## Detailed Design

### Architecture

```
[ASCII diagram or description of component interactions]
```

### API Design

#### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/resource | List resources |
| POST | /api/v1/resource | Create resource |

#### Request/Response Examples

```json
// POST /api/v1/resource
// Request
{
  "name": "example"
}

// Response (201 Created)
{
  "data": {
    "id": "uuid",
    "name": "example",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### Database Schema

```sql
CREATE TABLE resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### Data Flow

1. User action triggers API call
2. Handler validates input
3. Command/Query handler processes request
4. Repository persists/retrieves data
5. Response returned to user

## Implementation Plan

### Phase 1: Database and Domain
- [ ] Create migration
- [ ] Define domain entity with private fields + getters
- [ ] Define repository interface

### Phase 2: Backend API
- [ ] Implement command handlers
- [ ] Implement query handlers
- [ ] Create DTOs
- [ ] Implement repository
- [ ] Create HTTP handlers

### Phase 3: Frontend
- [ ] Create TypeScript types
- [ ] Implement API hooks
- [ ] Build React components
- [ ] Add to routing

### Phase 4: Testing
- [ ] Unit tests for handlers
- [ ] HTTP handler tests
- [ ] Frontend component tests

## Risks and Mitigations

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Risk 1 | High | Medium | Mitigation strategy |

## Alternatives Considered

### Alternative 1: [Name]
- **Pros**: ...
- **Cons**: ...
- **Why not chosen**: ...

## Open Questions

- [ ] Question 1?
- [ ] Question 2?
```

### 3. Task Breakdown Template

```markdown
## Epic: [Feature Name]

### Story 1: [User Story]
**Description**: As a [user], I can [action] so that [benefit].

#### Tasks:

1. **Backend: Database**
   - Create migration files
   - Files: `migrations/000001_create_<table>.up.sql`, `migrations/000001_create_<table>.down.sql`

2. **Backend: Domain Layer**
   - Define entity with private fields + getters
   - Define repository interface (port)
   - Define domain errors
   - Files: `internal/domain/<entity>/<entity>.go`, `repository.go`, `errors.go`

3. **Backend: Application Layer (CQRS)**
   - Implement command handlers (Create, Update, Delete)
   - Implement query handlers (Get, List)
   - Define DTOs
   - Files: `internal/application/command/`, `query/`, `dto/`

4. **Backend: Infrastructure Layer**
   - Implement PostgreSQL repository
   - Files: `internal/infrastructure/persistence/repository/<entity>_repository.go`

5. **Backend: Interface Layer**
   - Implement HTTP handlers
   - Files: `internal/interfaces/http/handler/<entity>_handler.go`

6. **Backend: Tests**
   - Unit tests for handlers
   - HTTP handler tests
   - Files: `*_test.go`

7. **Frontend: Types**
   - Define TypeScript interfaces
   - Files: `src/types/<entity>.ts`

8. **Frontend: API**
   - Implement API hooks
   - Files: `src/hooks/use-<entity>.ts`

9. **Frontend: Components**
   - Create UI components
   - Files: `src/components/features/<entity>/`

10. **Frontend: Tests**
    - Component tests
    - Hook tests
```

## Planning Principles

1. **Start with the end in mind** - Define success criteria first
2. **Design for change** - Systems evolve, plan for flexibility
3. **Prefer simple over complex** - Complexity should be justified
4. **Consider operational concerns** - Monitoring, debugging, deployment
5. **Document decisions** - Future developers need context
6. **Identify risks early** - Cheaper to fix in planning than implementation
7. **Break down ruthlessly** - Small tasks are easier to deliver
8. **Define contracts first** - API contracts before implementation
9. **Consider backward compatibility** - Especially for APIs
10. **Plan for testing** - Testing is part of the design

## Boundaries

### Always Do

- Follow DDD + CQRS patterns when designing backend features
- Consider database migrations and schema changes
- Plan for both backend API and frontend integration
- Identify risks and dependencies upfront
- Break tasks into small, deliverable units

### Ask First

- Before proposing architectural changes
- Before suggesting new technologies or dependencies
- When multiple valid approaches exist
- Before planning features that affect authentication/authorization

### Never Do

- Never skip security considerations in designs
- Never propose designs that violate Clean Architecture layers
- Never ignore existing patterns in the codebase
