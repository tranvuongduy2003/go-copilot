# Code Review Command

Perform comprehensive code review for quality, security, performance, and design system compliance.

## Task: $ARGUMENTS

## Review Checklist

### 1. Architecture & Design

- [ ] Follows Clean Architecture layers (domain <- application <- infrastructure <- interfaces)
- [ ] Uses CQRS pattern correctly (commands for writes, queries for reads)
- [ ] Domain entities use private fields with getters
- [ ] Repository interfaces defined in domain layer
- [ ] No business logic in HTTP handlers
- [ ] DTOs used for API responses (not domain entities)

### 2. Code Quality

- [ ] Functions are small and focused (single responsibility)
- [ ] Variable and function names are descriptive
- [ ] No magic numbers or strings (use constants)
- [ ] Error messages are clear and actionable
- [ ] Comments explain "why", not "what"
- [ ] No dead code or unused imports

### 3. Error Handling

**Backend (Go)**
- [ ] All errors are handled explicitly
- [ ] Errors are wrapped with context: `fmt.Errorf("context: %w", err)`
- [ ] Custom errors are defined for domain-specific cases
- [ ] No panic for recoverable errors
- [ ] Errors don't leak implementation details to API responses

**Frontend (React)**
- [ ] All async operations have error handling
- [ ] Loading and error states are displayed
- [ ] User-friendly error messages
- [ ] Error boundaries for component errors

### 4. Security

- [ ] No hardcoded secrets or credentials
- [ ] Parameterized queries only (no SQL string concatenation)
- [ ] Input validation at handler level
- [ ] No sensitive data logged (passwords, tokens, PII)
- [ ] Proper authentication/authorization checks
- [ ] No XSS vulnerabilities (no dangerouslySetInnerHTML without sanitization)
- [ ] No CSRF vulnerabilities
- [ ] Secure cookie settings

### 5. Performance

**Backend**
- [ ] Database queries are optimized (indexes, pagination)
- [ ] N+1 queries avoided
- [ ] Connection pooling configured
- [ ] Appropriate caching where needed
- [ ] Context timeouts for external calls

**Frontend**
- [ ] No unnecessary re-renders
- [ ] Large lists use virtualization
- [ ] Images are optimized
- [ ] Code splitting for large components
- [ ] Memoization where appropriate

### 6. Design System Compliance (Frontend)

- [ ] Uses design system colors only (no arbitrary colors)
- [ ] Uses spacing scale (no arbitrary spacing like p-[13px])
- [ ] Uses shadcn/ui components
- [ ] Consistent typography
- [ ] Accessible (proper labels, ARIA attributes, keyboard navigation)

### 7. Testing

- [ ] Unit tests for domain logic
- [ ] Tests for command/query handlers
- [ ] HTTP handler tests
- [ ] Edge cases covered
- [ ] Error paths tested
- [ ] Mocks used appropriately

### 8. TypeScript (Frontend)

- [ ] No `any` types
- [ ] Proper interface definitions
- [ ] Strict mode compliance
- [ ] Types match backend API contracts

## Common Issues to Flag

### Backend

```go
// BAD: Business logic in handler
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    if user.Age < 18 {  // Business logic should be in domain
        // ...
    }
}

// BAD: String concatenation in SQL
query := "SELECT * FROM users WHERE email = '" + email + "'"

// BAD: Ignored error
result, _ := service.DoSomething()

// BAD: Domain entity exposed directly
response.JSON(w, http.StatusOK, user)  // Should use DTO

// BAD: Infrastructure in domain
import "github.com/jackc/pgx/v5"  // In domain layer
```

### Frontend

```tsx
// BAD: Arbitrary color
<button className="bg-purple-500">  // Should use bg-primary

// BAD: Arbitrary spacing
<div className="p-[13px]">  // Should use p-3 or p-4

// BAD: No loading/error states
const { data } = useQuery(...);
return <div>{data.name}</div>;  // Crashes if data is undefined

// BAD: Using any
const handleClick = (e: any) => { ... }

// BAD: Sensitive data in localStorage
localStorage.setItem('authToken', token);
```

## Review Response Format

```markdown
## Summary
[Brief overview of the changes and overall assessment]

## Architecture
[Comments on architectural decisions]

## Security
[Any security concerns]

## Performance
[Performance considerations]

## Code Quality
[General code quality feedback]

## Suggestions
- [ ] Suggestion 1
- [ ] Suggestion 2

## Questions
- Question 1?
- Question 2?
```

## Severity Levels

| Level | Description | Action Required |
|-------|-------------|-----------------|
| **Critical** | Security vulnerability, data loss risk | Must fix before merge |
| **Major** | Bug, architectural violation | Should fix before merge |
| **Minor** | Style, minor improvement | Consider fixing |
| **Nitpick** | Personal preference | Optional |
