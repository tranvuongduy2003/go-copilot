---
description: Plan and implement a new feature end-to-end following DDD + CQRS patterns
---

# New Feature Implementation

I need to implement a new feature. Help me plan and build it end-to-end.

## Feature Details

**Feature Name**: {{featureName}}

**Description**: {{description}}

**Requirements**:
{{requirements}}

## Your Task

### Phase 1: Analysis
1. Search the codebase to understand existing patterns
2. Identify which files need to be created or modified
3. Check for similar implementations to follow

### Phase 2: Planning
Create a detailed implementation plan including:
- Database changes needed (migrations)
- Backend components (domain aggregate, application command/query handlers, infrastructure repository, interface handler)
- Frontend components (types, API client, hooks, components, pages)
- Tests to write

### Phase 3: Implementation

#### Backend (if applicable)
1. Create database migration using Goose in `backend/migrations/sql/`
2. Define domain aggregate in `backend/internal/domain/{aggregate}/`
   - Entity with private fields + getters
   - Repository interface (port)
   - Domain errors
3. Implement CQRS handlers in `backend/internal/application/`
   - Command handlers in `command/`
   - Query handlers in `query/`
   - DTOs in `dto/`
4. Implement repository adapter in `backend/internal/infrastructure/persistence/postgres/`
5. Create HTTP handler in `backend/internal/interfaces/http/handler/`
6. Register routes

#### Frontend (if applicable)
1. Define TypeScript types in `frontend/src/types/`
2. Create API client functions in `frontend/src/lib/api/`
3. Create React Query hooks in `frontend/src/hooks/`
4. Build components in `frontend/src/components/features/`
5. Create page component in `frontend/src/pages/`
6. Add routing

### Phase 4: Testing
1. Write unit tests for service layer
2. Write handler tests
3. Write component tests
4. Write integration tests if needed

### Phase 5: Verification
1. Run all tests
2. Check for linting errors
3. Verify the feature works end-to-end

## Guidelines

- Follow the project's coding standards
- Use the design system colors and spacing (no arbitrary values)
- Write comprehensive error handling
- Add appropriate logging
- Consider security implications
- Write tests alongside implementation

## Output

Provide:
1. A summary of changes made
2. List of files created/modified
3. Any manual steps needed
4. Testing instructions
