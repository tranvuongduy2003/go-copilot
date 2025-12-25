---
description: Debug and fix a bug with root cause analysis
agent: "Fullstack Engineer"
---

# Bug Fix

Debug and fix a bug with proper root cause analysis.

## Bug Details

**Description**: {{description}}

**Expected Behavior**: {{expected}}

**Actual Behavior**: {{actual}}

**Steps to Reproduce**:
{{steps}}

## Debugging Process

### Phase 1: Understand the Bug

1. **Reproduce the bug**
   - Follow the steps to reproduce
   - Verify the actual behavior matches the report
   - Note any error messages or logs

2. **Identify the scope**
   - Which component/function is affected?
   - When did this start happening?
   - What recent changes might be related?

### Phase 2: Investigate

1. **Trace the code path**
   - Find the entry point
   - Follow the execution flow
   - Identify where behavior diverges from expected

2. **Search for related code**
   - Look for similar patterns
   - Check if the bug exists elsewhere
   - Review recent changes to affected files

3. **Check logs and errors**
   - Backend logs
   - Browser console
   - Network requests

### Phase 3: Root Cause Analysis

Ask these questions:
- WHY did this bug occur?
- WHAT allowed it to happen?
- WHERE are similar bugs likely to exist?
- HOW can we prevent this in the future?

Common root causes:
- Off-by-one errors
- Null/undefined handling
- Race conditions
- Missing validation
- Incorrect assumptions
- State management issues
- API contract mismatches

### Phase 4: Fix the Bug

1. **Write a failing test first**
   - Create a test that reproduces the bug
   - Verify it fails

2. **Implement the fix**
   - Make the minimal change needed
   - Don't add unrelated changes
   - Follow existing patterns

3. **Verify the fix**
   - Test passes now
   - No regression in other tests
   - Manual verification works

### Phase 5: Prevent Recurrence

Consider:
- Adding validation
- Improving error messages
- Adding logging
- Updating documentation
- Adding type safety
- Creating regression tests

## Bug Fix Templates

### Go: Null Pointer Fix
```go
// Before (bug)
func (s *Service) GetUser(id string) (*User, error) {
    user := s.cache.Get(id) // Could be nil
    return user, nil        // Returns nil without error
}

// After (fixed)
func (s *Service) GetUser(id string) (*User, error) {
    user := s.cache.Get(id)
    if user == nil {
        return nil, ErrUserNotFound
    }
    return user, nil
}
```

### Go: Error Handling Fix
```go
// Before (bug - error ignored)
result, _ := json.Marshal(data)

// After (fixed)
result, err := json.Marshal(data)
if err != nil {
    return fmt.Errorf("failed to marshal data: %w", err)
}
```

### React: State Update Fix
```tsx
// Before (bug - stale state)
const handleClick = () => {
    setCount(count + 1);
    setCount(count + 1); // Still uses old count
};

// After (fixed)
const handleClick = () => {
    setCount(prev => prev + 1);
    setCount(prev => prev + 1);
};
```

### React: Dependency Fix
```tsx
// Before (bug - missing dependency)
useEffect(() => {
    fetchData(userId);
}, []); // userId not in deps

// After (fixed)
useEffect(() => {
    fetchData(userId);
}, [userId]);
```

## Output

Provide:

1. **Root Cause**
   - What caused the bug
   - Why it wasn't caught earlier

2. **Fix Description**
   - What was changed
   - Why this fix is correct

3. **Files Modified**
   - List of changed files
   - Brief description of each change

4. **Tests Added**
   - Regression test for this bug
   - Any related tests added

5. **Prevention**
   - How to prevent similar bugs
   - Any follow-up tasks

6. **Verification**
   - How to verify the fix works
   - Test commands to run
