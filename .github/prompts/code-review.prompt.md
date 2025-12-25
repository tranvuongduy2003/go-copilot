---
description: Perform thorough code review
---

# Code Review

Perform a thorough code review focusing on quality, security, performance, and design system compliance.

## Review Target

**Files/PR to Review**: {{target}}

## Review Checklist

### Code Quality

- [ ] Code is readable and self-documenting
- [ ] Functions are focused and not too long
- [ ] Naming is clear and consistent
- [ ] No code duplication
- [ ] Error handling is comprehensive
- [ ] Edge cases are handled
- [ ] Comments explain "why" not "what"

### Go Specific

- [ ] Error handling follows project patterns
- [ ] Context is passed correctly
- [ ] Resources are properly closed (defer)
- [ ] Concurrency is handled safely
- [ ] Interfaces are defined at point of use
- [ ] Tests use table-driven pattern

### React Specific

- [ ] Components are appropriately sized
- [ ] Props have proper TypeScript types
- [ ] Hooks follow rules of hooks
- [ ] Memoization is used appropriately
- [ ] Error boundaries are in place
- [ ] Loading states are handled

### Design System Compliance

- [ ] Uses design system colors (no arbitrary colors)
- [ ] Uses standard spacing scale (no arbitrary spacing)
- [ ] Uses correct border radius values
- [ ] Typography follows scale
- [ ] Consistent with existing UI patterns

### Security

- [ ] Input is validated
- [ ] SQL uses parameterized queries
- [ ] No sensitive data in logs
- [ ] Authentication/authorization checks
- [ ] No XSS vulnerabilities
- [ ] CSRF protection (for mutations)

### Performance

- [ ] No N+1 queries
- [ ] Proper pagination
- [ ] Appropriate caching
- [ ] No unnecessary re-renders
- [ ] Expensive computations memoized
- [ ] Large lists virtualized

### Testing

- [ ] Unit tests for business logic
- [ ] Handler/component tests
- [ ] Edge cases tested
- [ ] Error scenarios tested
- [ ] Test names are descriptive

## Review Output Format

```markdown
## Code Review Summary

### Overall Assessment
[Brief summary: Approve / Request Changes / Comment]

### Critical Issues (Must Fix)
1. **[Category]** - `file:line`
   - Issue: [Description]
   - Why: [Explanation]
   - Fix: [Suggested solution]

### Major Issues (Should Fix)
1. **[Category]** - `file:line`
   - Issue: [Description]
   - Suggestion: [How to improve]

### Minor Issues (Consider)
- `file:line`: [Brief description]

### Design System Violations
- `file:line`: Using `bg-purple-500` instead of `bg-primary`

### Security Concerns
- `file:line`: [Description of security issue]

### Performance Concerns
- `file:line`: [Description of performance issue]

### Positive Observations
- [What was done well]

### Questions
- [Any clarifications needed]
```

## Common Issues to Look For

### Go
```go
// Issue: Ignoring error
result, _ := doSomething()

// Issue: Not wrapping error
if err != nil {
    return err // Should wrap with context
}

// Issue: Logging sensitive data
log.Printf("User login: %s password: %s", email, password)

// Issue: SQL injection
query := "SELECT * FROM users WHERE email = '" + email + "'"
```

### React
```tsx
// Issue: Arbitrary colors
<div className="bg-purple-500">

// Issue: Missing dependency
useEffect(() => {
    fetchData(id);
}, []); // id should be in deps

// Issue: Missing key in list
{items.map(item => <Item {...item} />)} // needs key

// Issue: Direct state mutation
items.push(newItem);
setItems(items);
```

### Security
```go
// Issue: No rate limiting on auth endpoint
// Issue: JWT secret too short
// Issue: No input validation
```

## Severity Levels

**Critical**: Must fix before merge
- Security vulnerabilities
- Data loss risks
- Breaking changes without migration

**Major**: Should fix before merge
- Significant bugs
- Performance issues
- Missing error handling

**Minor**: Consider fixing
- Style improvements
- Minor optimizations
- Documentation gaps
