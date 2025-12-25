---
description: Add comprehensive tests for existing code
agent: "Testing Specialist"
---

# Add Tests

Add comprehensive tests for existing code.

## Target

**File/Function to Test**: {{target}}

**Test Type**:
- [ ] Unit tests
- [ ] Integration tests
- [ ] Component tests (React)
- [ ] E2E tests

## Instructions

### For Go Code

1. Analyze the code to understand:
   - Input parameters
   - Return values
   - Side effects
   - Error conditions
   - Edge cases

2. Create test file: `{{filename}}_test.go`

3. Use table-driven tests:
```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr error
    }{
        {
            name:    "success case",
            input:   validInput,
            want:    expectedOutput,
            wantErr: nil,
        },
        {
            name:    "error case",
            input:   invalidInput,
            want:    zero,
            wantErr: ErrExpected,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionName(tt.input)

            if tt.wantErr != nil {
                require.Error(t, err)
                assert.ErrorIs(t, err, tt.wantErr)
                return
            }

            require.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

4. Test coverage goals:
   - Happy path (success cases)
   - Error conditions
   - Edge cases (empty, nil, max values)
   - Boundary conditions

### For React Code

1. Analyze the component:
   - Props and their types
   - User interactions
   - State changes
   - API calls
   - Conditional rendering

2. Create test file: `{{filename}}.test.tsx`

3. Use React Testing Library patterns:
```tsx
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { Component } from './component';

describe('Component', () => {
  it('renders correctly', () => {
    render(<Component />);
    expect(screen.getByRole('button')).toBeInTheDocument();
  });

  it('handles user interaction', async () => {
    const user = userEvent.setup();
    const onClick = vi.fn();

    render(<Component onClick={onClick} />);

    await user.click(screen.getByRole('button'));

    expect(onClick).toHaveBeenCalled();
  });

  it('handles async operations', async () => {
    render(<Component />);

    await waitFor(() => {
      expect(screen.getByText('Loaded')).toBeInTheDocument();
    });
  });
});
```

4. Test coverage:
   - Rendering with different props
   - User interactions
   - Loading states
   - Error states
   - Accessibility

## Test Scenarios to Cover

### Go Services
- [ ] Success with valid input
- [ ] Validation errors
- [ ] Repository errors
- [ ] Not found scenarios
- [ ] Conflict scenarios
- [ ] Authorization checks

### Go Handlers
- [ ] Valid request returns correct status
- [ ] Invalid JSON returns 400
- [ ] Validation errors return details
- [ ] Not found returns 404
- [ ] Authentication required returns 401
- [ ] Authorization denied returns 403

### React Components
- [ ] Renders with required props
- [ ] Renders with optional props
- [ ] Handles click events
- [ ] Handles form submission
- [ ] Shows loading state
- [ ] Shows error state
- [ ] Shows empty state
- [ ] Is keyboard accessible

### React Hooks
- [ ] Returns correct initial state
- [ ] Fetches data on mount
- [ ] Handles loading state
- [ ] Handles error state
- [ ] Invalidates cache on mutation

## Verification

After adding tests:
```bash
# Go
go test ./... -v
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# React
npm test
npm run test:coverage
```

## Output

Provide:
1. Test file(s) created
2. Coverage report summary
3. Any issues found while writing tests
4. Suggestions for improving testability
