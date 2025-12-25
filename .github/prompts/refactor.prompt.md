---
description: Refactor code for better quality and maintainability
agent: "Fullstack Engineer"
---

# Code Refactoring

Refactor code to improve quality, maintainability, and performance.

## Target

**Code to Refactor**: {{target}}

**Refactoring Goals**:
- [ ] Improve readability
- [ ] Reduce complexity
- [ ] Extract reusable code
- [ ] Improve performance
- [ ] Fix code smells
- [ ] Apply design patterns
- [ ] Improve type safety
- [ ] Enhance error handling

## Analysis Phase

Before refactoring, analyze:

1. **Current state**
   - What does this code do?
   - What are its dependencies?
   - Who calls this code?

2. **Problems identified**
   - Code smells (duplication, long methods, etc.)
   - Performance issues
   - Maintainability concerns
   - Missing error handling
   - Type safety issues

3. **Test coverage**
   - Are there existing tests?
   - Will they catch regressions?

## Refactoring Techniques

### Go Refactoring

**Extract Function**
```go
// Before
func Process(items []Item) {
    for _, item := range items {
        // 20 lines of validation
        // 20 lines of transformation
        // 20 lines of persistence
    }
}

// After
func Process(items []Item) {
    for _, item := range items {
        if err := validate(item); err != nil {
            continue
        }
        transformed := transform(item)
        persist(transformed)
    }
}
```

**Extract Interface**
```go
// Before: Direct dependency
type Service struct {
    db *sql.DB
}

// After: Interface dependency
type Repository interface {
    Find(id string) (*Entity, error)
}

type Service struct {
    repo Repository
}
```

**Simplify Conditionals**
```go
// Before
if err != nil {
    return nil, err
} else {
    return result, nil
}

// After
if err != nil {
    return nil, err
}
return result, nil
```

### React Refactoring

**Extract Component**
```tsx
// Before: Large component
function Dashboard() {
  return (
    <div>
      {/* 50 lines of stats */}
      {/* 50 lines of chart */}
      {/* 50 lines of table */}
    </div>
  );
}

// After: Composed components
function Dashboard() {
  return (
    <div>
      <StatsSection />
      <ChartSection />
      <DataTable />
    </div>
  );
}
```

**Extract Custom Hook**
```tsx
// Before: Logic in component
function UserList() {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchUsers().then(setUsers).finally(() => setLoading(false));
  }, []);

  // Component logic...
}

// After: Custom hook
function useUsers() {
  return useQuery({
    queryKey: ['users'],
    queryFn: fetchUsers,
  });
}

function UserList() {
  const { data: users, isLoading } = useUsers();
  // Component logic...
}
```

**Simplify Props**
```tsx
// Before: Too many props
<Button
  text="Submit"
  onClick={handleClick}
  disabled={isLoading}
  loading={isLoading}
  variant="primary"
  size="medium"
  icon={<Check />}
  iconPosition="left"
/>

// After: Composed/simplified
<Button onClick={handleClick} isLoading={isLoading}>
  <Check className="mr-2" />
  Submit
</Button>
```

## Refactoring Checklist

Before starting:
- [ ] Understand the current behavior
- [ ] Ensure tests exist (add if missing)
- [ ] Tests pass before refactoring

During refactoring:
- [ ] Make small, incremental changes
- [ ] Run tests after each change
- [ ] Commit working states

After refactoring:
- [ ] All tests still pass
- [ ] No new linting errors
- [ ] Code is more readable
- [ ] No behavior changes (unless intended)

## Common Refactoring Targets

### Go
- Long functions (>50 lines)
- Deep nesting (>3 levels)
- Duplicated code
- Missing error handling
- Unexported types that should be exported
- Comments explaining complex code (simplify instead)

### React
- Large components (>200 lines)
- Repeated JSX patterns
- Complex useEffect dependencies
- Prop drilling (use context)
- Inline styles (use Tailwind classes)
- Arbitrary values (use design system)

## Output

Provide:
1. Summary of changes
2. Before/after comparison for key changes
3. Explanation of why each change improves the code
4. Any breaking changes or migration steps
5. Updated tests if needed
