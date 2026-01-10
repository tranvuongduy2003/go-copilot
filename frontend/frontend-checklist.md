# Expert-Level Frontend Setup Checklist

## Vite + Bun + React + TypeScript | Production-Ready SPA

---

# Phase 1: Project Foundation

## 1.1 Project Initialization

- [x] Initialize project using Bun with Vite React TypeScript template
- [x] Configure `bunfig.toml` for Bun-specific settings
- [x] Update `package.json` with correct project metadata (name, version, description)
- [x] Configure `tsconfig.json` with strict TypeScript settings
  - [x] Enable strict mode, noImplicitAny, strictNullChecks
  - [x] Configure path aliases (@/ for src/)
  - [x] Set target to ES2022 for modern JavaScript features
  - [x] Enable skipLibCheck for faster builds
- [x] Create `.gitignore` with appropriate patterns (node_modules, dist, .env)
- [x] Create `.editorconfig` for consistent coding style
- [x] Create `README.md` with project overview and setup instructions
- [x] Initialize git repository with initial commit

## 1.2 Vite Configuration (`vite.config.ts`)

- [x] Configure path aliases to match tsconfig
- [x] Configure development server port and host
- [x] Enable CORS for API communication during development
- [x] Configure proxy for API requests to backend
- [x] Configure build output settings
- [x] Enable source maps for development
- [x] Configure chunk splitting strategy for production
- [x] Set up environment variable handling (VITE\_ prefix)
- [x] Configure preview server for production build testing
- [x] Add build analysis plugin (rollup-plugin-visualizer) for bundle optimization

## 1.3 Environment Configuration

- [x] Create `.env.example` with all required variables documented
- [x] Define `VITE_API_BASE_URL` for backend API endpoint
- [x] Define `VITE_APP_NAME` for application branding
- [x] Define `VITE_APP_VERSION` for version tracking
- [x] Define `VITE_ENABLE_MOCK` for mock API toggle (development)
- [x] Create `.env.development` for development defaults
- [x] Create `.env.production` for production defaults
- [x] Create environment type declarations (`env.d.ts`)
- [x] Implement environment validation on app startup

## 1.4 Package Dependencies

### Core Dependencies

- [x] Install React 18+ and React DOM
- [x] Install React Router DOM v6+ for routing
- [x] Install @tanstack/react-query for server state management
- [x] Install zustand for client state management
- [x] Install axios for HTTP requests
- [x] Install react-hook-form for form management
- [x] Install zod for schema validation
- [x] Install @hookform/resolvers for zod integration with react-hook-form
- [x] Install date-fns for date manipulation
- [x] Install clsx and tailwind-merge for className utilities

### UI Dependencies

- [x] Install Tailwind CSS and required PostCSS plugins
- [x] Install shadcn/ui CLI and initialize
- [x] Install lucide-react for icons
- [x] Install class-variance-authority for component variants
- [x] Install @radix-ui primitives as needed by shadcn components
- [x] Install sonner or react-hot-toast for notifications
- [ ] Install framer-motion for animations (optional)

### Development Dependencies

- [x] Install TypeScript and @types/react, @types/react-dom
- [x] Install ESLint with React and TypeScript plugins
- [x] Install Prettier for code formatting
- [x] Install @typescript-eslint/parser and @typescript-eslint/eslint-plugin
- [x] Install eslint-plugin-react-hooks for hooks linting
- [x] Install eslint-plugin-react-refresh for Fast Refresh compatibility
- [ ] Install husky for git hooks
- [ ] Install lint-staged for pre-commit linting
- [x] Install @testing-library/react for component testing
- [x] Install vitest for unit testing
- [x] Install msw for API mocking in tests
- [ ] Install playwright or cypress for E2E testing (optional)

## 1.5 Code Quality Setup

### ESLint Configuration (`.eslintrc.cjs` or `eslint.config.js`)

- [x] Configure parser for TypeScript
- [x] Enable React and React Hooks plugins
- [x] Enable import sorting rules
- [x] Configure no-unused-vars with TypeScript override
- [x] Enable strict type-checking rules
- [x] Configure path alias resolution for import plugin
- [x] Add custom rules for project conventions

### Prettier Configuration (`.prettierrc`)

- [x] Configure semi, singleQuote, tabWidth, trailingComma
- [x] Configure printWidth (80-120 recommended)
- [x] Configure endOfLine for cross-platform compatibility
- [x] Create `.prettierignore` for build artifacts

### Git Hooks

- [x] Initialize husky with `bunx husky init`
- [x] Create pre-commit hook for lint-staged
- [ ] Create commit-msg hook for conventional commits (optional)
- [x] Configure lint-staged in `package.json` or `.lintstagedrc`
- [x] Run ESLint and Prettier on staged files

## 1.6 Tailwind CSS Configuration

- [x] Initialize Tailwind CSS with `bunx tailwindcss init -p`
- [x] Configure content paths in `tailwind.config.js`
- [x] Extend theme with custom colors matching design system
- [x] Configure custom spacing, typography, and breakpoints
- [x] Add CSS variables for theming (light/dark mode support)
- [x] Configure animation keyframes for custom animations
- [x] Set up Tailwind plugins (forms, typography, aspect-ratio)
- [x] Create base styles in `index.css` with Tailwind directives

## 1.7 shadcn/ui Setup

- [x] Run `bunx shadcn-ui@latest init`
- [x] Configure components.json with correct paths
- [x] Select style (default or new-york)
- [x] Configure base color and CSS variables
- [x] Install essential components: button, input, card, form, label
- [x] Install feedback components: alert, toast, dialog, drawer
- [x] Install data display components: table, badge, avatar
- [x] Install navigation components: dropdown-menu, tabs, navigation-menu
- [x] Install form components: select, checkbox, radio-group, switch
- [x] Create component barrel exports for cleaner imports

---

# Phase 2: Application Architecture

## 2.1 Folder Structure Reorganization

```
src/
├── app/                    # Application setup
│   ├── providers/          # React context providers
│   ├── router/             # Router configuration
│   └── App.tsx             # Root component
├── components/
│   ├── ui/                 # shadcn/ui components
│   ├── common/             # Shared/reusable components
│   └── layout/             # Layout components
├── features/               # Feature modules
│   ├── auth/
│   ├── users/
│   └── dashboard/
├── hooks/                  # Shared custom hooks
├── lib/                    # Utilities and configurations
│   ├── api/                # API client setup
│   ├── utils/              # Utility functions
│   └── validations/        # Zod schemas
├── stores/                 # Zustand stores
├── types/                  # Shared TypeScript types
├── styles/                 # Global styles
└── constants/              # Application constants
```

- [x] Create folder structure as defined above
- [x] Move existing components to appropriate locations
- [x] Create barrel exports (index.ts) for each major folder
- [x] Configure path aliases for new folders in tsconfig and vite.config

## 2.2 API Layer Setup (`src/lib/api/`)

### 2.2.1 Axios Client Configuration

- [x] Create axios instance with base configuration
- [x] Configure base URL from environment variable
- [x] Set default headers (Content-Type, Accept)
- [x] Configure timeout (30 seconds recommended)
- [x] Create request interceptor for authentication
  - [x] Automatically attach access token to requests
  - [x] Skip token for public endpoints
- [x] Create response interceptor for error handling
  - [x] Handle 401 errors with token refresh logic
  - [x] Handle network errors gracefully
  - [x] Implement request retry for transient failures
- [x] Create response interceptor for data transformation
  - [x] Unwrap successful responses for cleaner data access

### 2.2.2 Token Refresh Implementation

- [x] Implement token refresh mechanism
- [x] Queue failed requests during refresh
- [x] Retry queued requests after successful refresh
- [x] Redirect to login on refresh failure
- [x] Prevent multiple simultaneous refresh requests
- [x] Handle refresh token expiration

### 2.2.3 API Error Handling

- [x] Define ApiError class with typed error structure
  - [x] Include status code, message, and error details
  - [x] Include original error for debugging
- [x] Create error type guards for specific error types
- [x] Map backend error codes to user-friendly messages
- [x] Create error boundary for API errors

### 2.2.4 API Endpoints Organization

- [x] Create endpoints constants file
- [x] Group endpoints by feature/resource
- [x] Use string templates for parameterized URLs
- [x] Document each endpoint with expected request/response

## 2.3 React Query Setup (`src/lib/api/query-client.ts`)

### 2.3.1 Query Client Configuration

- [x] Create QueryClient instance with default options
- [x] Configure default staleTime (5 minutes for most data)
- [x] Configure default gcTime (cacheTime) (10 minutes)
- [x] Configure retry logic (3 retries with exponential backoff)
- [x] Configure refetchOnWindowFocus (false for most apps)
- [x] Configure refetchOnMount behavior
- [x] Set up query client provider in app providers

### 2.3.2 Query Key Factory

- [x] Create query key factory for consistent key management
- [x] Define keys by feature (auth, users, etc.)
- [x] Include parameters in keys for cache granularity
- [x] Export typed query keys for type safety
- [x] Document query key patterns

### 2.3.3 Custom Query Hooks Pattern

- [x] Define pattern for query hooks (useQuery wrapper)
- [x] Define pattern for mutation hooks (useMutation wrapper)
- [x] Include loading, error, and success states
- [x] Include automatic cache invalidation on mutations
- [x] Include optimistic updates where appropriate

## 2.4 State Management with Zustand (`src/stores/`)

### 2.4.1 Store Architecture

- [x] Define store slicing strategy (one store per domain)
- [x] Create auth store for authentication state
- [x] Create UI store for global UI state (sidebar, modals)
- [x] Create user preferences store (theme, language)
- [x] Implement persist middleware for relevant stores
- [x] Configure storage adapter (localStorage)
- [x] Define persisted vs non-persisted state

### 2.4.2 Auth Store Implementation

- [x] Define auth state interface
  - [x] Include user object (nullable)
  - [x] Include authentication status (loading, authenticated, unauthenticated)
  - [x] Include tokens (if storing in memory)
- [x] Define auth actions
  - [x] setUser action
  - [x] clearAuth action
  - [x] updateUser action
- [x] Implement selectors for derived state
  - [x] isAuthenticated selector
  - [x] hasPermission selector (for RBAC)
  - [x] hasRole selector

### 2.4.3 UI Store Implementation

- [x] Define UI state interface
  - [x] Include sidebar collapsed state
  - [x] Include active modal state
  - [x] Include global loading state
  - [x] Include notification queue (if not using toast library)
- [x] Define UI actions for state mutations
- [x] Implement reset action for cleanup

## 2.5 Form Management Setup

### 2.5.1 Zod Schema Definitions (`src/lib/validations/`)

- [x] Create base validation schemas for common types
  - [x] Email schema with proper regex
  - [x] Password schema with strength requirements
  - [x] UUID schema for IDs
  - [x] Pagination schema for list queries
- [x] Create auth-related schemas
  - [x] Login schema (email, password)
  - [x] Register schema (email, password, confirmPassword, fullName)
  - [x] Forgot password schema
  - [x] Reset password schema
  - [x] Change password schema
- [x] Create user-related schemas
  - [x] Create user schema
  - [x] Update user schema (partial)
  - [x] User filter schema for list queries
- [x] Export TypeScript types inferred from schemas

### 2.5.2 Form Components Pattern

- [x] Create form field wrapper component
  - [x] Integrate with react-hook-form Controller
  - [x] Display validation errors from form state
  - [x] Support all input types (text, password, select, etc.)
- [x] Create reusable form components
  - [x] FormInput component
  - [x] FormSelect component
  - [x] FormCheckbox component
  - [x] FormTextarea component
  - [ ] FormDatePicker component (if needed)

---

# Phase 3: Authentication Feature (`src/features/auth/`)

## 3.1 Auth Feature Structure

```
features/auth/
├── api/
│   ├── auth.api.ts         # API functions
│   └── auth.queries.ts     # React Query hooks
├── components/
│   ├── login-form.tsx
│   ├── register-form.tsx
│   ├── forgot-password-form.tsx
│   ├── reset-password-form.tsx
│   └── auth-guard.tsx
├── hooks/
│   └── use-auth.ts
├── pages/
│   ├── login.page.tsx
│   ├── register.page.tsx
│   ├── forgot-password.page.tsx
│   └── reset-password.page.tsx
├── types/
│   └── auth.types.ts
└── index.ts
```

- [x] Create folder structure for auth feature
- [x] Create barrel export (index.ts) for public API

## 3.2 Auth Types (`types/auth.types.ts`)

- [x] Define User interface matching backend UserDTO
  - [x] Include id, email, fullName, status, roles, permissions
- [x] Define LoginRequest interface
- [x] Define LoginResponse interface with tokens and user
- [x] Define RegisterRequest interface
- [x] Define TokenPair interface (accessToken, refreshToken, expiresIn)
- [x] Define AuthState interface for store
- [x] Define Permission and Role types for RBAC
- [x] Define ForgotPasswordRequest interface
- [x] Define ResetPasswordRequest interface
- [x] Define ChangePasswordRequest interface

## 3.3 Auth API (`api/auth.api.ts`)

- [x] Implement login API function
  - [x] Accept credentials, return token pair and user
- [x] Implement register API function
- [x] Implement logout API function
  - [x] Call backend logout endpoint
- [x] Implement refresh token API function
- [x] Implement forgot password API function
- [x] Implement reset password API function
- [x] Implement change password API function
- [x] Implement get current user API function (GET /auth/me)
- [x] Implement get sessions API function
- [x] Implement revoke session API function

## 3.4 Auth React Query Hooks (`api/auth.queries.ts`)

### 3.4.1 Auth Mutations

- [x] Create useLogin mutation hook
  - [x] On success: store tokens, update auth store, redirect
  - [x] On error: display error message
- [x] Create useRegister mutation hook
  - [x] On success: auto-login or redirect to login
- [x] Create useLogout mutation hook
  - [x] On success: clear auth store, clear query cache, redirect
- [x] Create useForgotPassword mutation hook
  - [x] On success: display success message
- [x] Create useResetPassword mutation hook
  - [x] On success: redirect to login
- [x] Create useChangePassword mutation hook
  - [x] On success: display success message, optionally logout

### 3.4.2 Auth Queries

- [x] Create useCurrentUser query hook
  - [x] Fetch current user data on app mount
  - [x] Update auth store with user data
  - [x] Handle unauthorized (token expired)
- [x] Create useSessions query hook for active sessions list
- [x] Create useRevokeSession mutation hook

## 3.5 Auth Components

### 3.5.1 Login Form Component

- [x] Integrate react-hook-form with zod resolver
- [x] Include email input with validation
- [x] Include password input with show/hide toggle
- [x] Include "Remember me" checkbox (optional)
- [x] Include "Forgot password" link
- [x] Include submit button with loading state
- [x] Display form-level errors from API
- [x] Display field-level validation errors
- [x] Handle form submission with useLogin hook
- [x] Redirect to dashboard on success

### 3.5.2 Register Form Component

- [x] Include email input with validation
- [x] Include full name input with validation
- [x] Include password input with strength indicator
- [x] Include confirm password input with match validation
- [x] Include terms acceptance checkbox
- [x] Include submit button with loading state
- [x] Display validation errors
- [x] Handle form submission with useRegister hook
- [x] Include link to login page

### 3.5.3 Forgot Password Form Component

- [x] Include email input with validation
- [x] Include submit button with loading state
- [x] Display success message after submission
- [x] Include link back to login

### 3.5.4 Reset Password Form Component

- [x] Extract reset token from URL query params
- [x] Include new password input with strength indicator
- [x] Include confirm password input
- [x] Include submit button with loading state
- [x] Handle invalid/expired token error
- [x] Redirect to login on success

### 3.5.5 Auth Guard Component

- [x] Check authentication status
- [x] Redirect to login if not authenticated
- [x] Show loading state while checking auth
- [x] Render children if authenticated
- [x] Accept required permissions/roles props
- [x] Check permissions if specified
- [x] Redirect to unauthorized page if permission denied

## 3.6 Auth Pages

### 3.6.1 Login Page

- [x] Create page layout with centered form
- [x] Include app logo/branding
- [x] Include LoginForm component
- [x] Include link to register page
- [ ] Include social login options (if applicable)
- [x] Redirect if already authenticated

### 3.6.2 Register Page

- [x] Create page layout with centered form
- [x] Include app logo/branding
- [x] Include RegisterForm component
- [x] Include link to login page
- [x] Redirect if already authenticated

### 3.6.3 Forgot Password Page

- [x] Create page layout with centered form
- [x] Include ForgotPasswordForm component
- [x] Include link back to login

### 3.6.4 Reset Password Page

- [x] Create page layout with centered form
- [x] Validate token presence in URL
- [x] Include ResetPasswordForm component
- [x] Handle token validation errors

## 3.7 Token Management

- [x] Implement secure token storage strategy
  - [ ] Option 1: Memory only (most secure, no persistence)
  - [x] Option 2: localStorage (convenient, less secure)
  - [ ] Option 3: HttpOnly cookies (requires backend support)
- [x] Implement token persistence (if using localStorage)
  - [x] Store encrypted or obfuscated tokens
  - [x] Clear tokens on logout
- [x] Implement token refresh on app mount
- [x] Check token expiration before requests
- [x] Implement automatic token refresh before expiration
  - [x] Set up refresh interval or intercept 401 responses

---

# Phase 4: User Management Feature (`src/features/users/`)

## 4.1 User Feature Structure

```
features/users/
├── api/
│   ├── users.api.ts
│   └── users.queries.ts
├── components/
│   ├── user-table.tsx
│   ├── user-form.tsx
│   ├── user-card.tsx
│   ├── user-filters.tsx
│   └── user-actions.tsx
├── hooks/
│   └── use-user-permissions.ts
├── pages/
│   ├── users-list.page.tsx
│   ├── user-detail.page.tsx
│   └── user-create.page.tsx
├── types/
│   └── user.types.ts
└── index.ts
```

- [x] Create folder structure for users feature
- [x] Create barrel export for public API

## 4.2 User Types (`types/user.types.ts`)

- [x] Define User interface matching backend DTO
- [x] Define UserStatus enum (pending, active, inactive, banned)
- [x] Define CreateUserRequest interface
- [x] Define UpdateUserRequest interface
- [x] Define UserFilter interface for list queries
- [x] Define PaginatedUsers interface for list response
- [x] Define UserRole interface
- [x] Define UserPermission interface

## 4.3 User API (`api/users.api.ts`)

- [x] Implement getUsers API function with pagination and filters
- [x] Implement getUser API function by ID
- [x] Implement createUser API function
- [x] Implement updateUser API function
- [x] Implement deleteUser API function
- [x] Implement activateUser API function
- [x] Implement deactivateUser API function
- [x] Implement getUserRoles API function
- [x] Implement assignRole API function
- [x] Implement revokeRole API function

## 4.4 User React Query Hooks (`api/users.queries.ts`)

### 4.4.1 User Queries

- [x] Create useUsers query hook with pagination
  - [x] Accept page, limit, filters as parameters
  - [x] Return paginated data with metadata
  - [x] Enable keepPreviousData for smooth pagination
- [x] Create useUser query hook by ID
  - [x] Enable caching with appropriate staleTime
- [x] Create useUserRoles query hook

### 4.4.2 User Mutations

- [x] Create useCreateUser mutation
  - [x] Invalidate users list cache on success
- [x] Create useUpdateUser mutation
  - [x] Invalidate specific user and list cache
  - [x] Implement optimistic update (optional)
- [x] Create useDeleteUser mutation
  - [x] Invalidate users list cache on success
- [x] Create useActivateUser mutation
- [x] Create useDeactivateUser mutation
- [x] Create useAssignRole mutation
- [x] Create useRevokeRole mutation

## 4.5 User Components

### 4.5.1 User Table Component

- [x] Display users in responsive table
- [x] Include columns: name, email, status, roles, created, actions
- [x] Implement column sorting
- [x] Implement row selection (if bulk actions needed)
- [x] Include action buttons (edit, delete, activate/deactivate)
- [x] Display loading skeleton during fetch
- [x] Handle empty state with appropriate message
- [x] Implement pagination controls

### 4.5.2 User Form Component

- [x] Support create and edit modes
- [x] Include all user fields with validation
- [x] Include role assignment (multi-select)
- [x] Include status selection (for edit mode)
- [x] Handle form submission with appropriate mutation
- [x] Display loading state during submission
- [x] Display success/error feedback

### 4.5.3 User Filters Component

- [x] Include search input (debounced)
- [x] Include status filter dropdown
- [x] Include role filter dropdown
- [ ] Include date range filter (optional)
- [x] Include clear filters button
- [x] Sync filters with URL query params

### 4.5.4 User Card Component

- [x] Display user info in card format (for grid view)
- [x] Include avatar, name, email, status badge
- [x] Include quick action buttons
- [x] Support click to navigate to detail

### 4.5.5 User Actions Component

- [x] Dropdown menu for user actions
- [x] Include edit action
- [x] Include activate/deactivate action (conditional)
- [x] Include delete action with confirmation
- [ ] Include view sessions action
- [x] Check permissions before showing actions

## 4.6 User Pages

### 4.6.1 Users List Page

- [x] Include page header with title and create button
- [x] Include UserFilters component
- [x] Include UserTable component
- [x] Include pagination component
- [ ] Support view toggle (table/grid) (optional)
- [x] Protect with required permission

### 4.6.2 User Detail Page

- [x] Fetch user data by ID from URL params
- [x] Display user information
- [x] Display user roles and permissions
- [ ] Display user sessions (optional)
- [x] Include edit button linking to edit form
- [x] Include back button to list
- [x] Handle user not found error

### 4.6.3 User Create/Edit Page

- [x] Include page header with appropriate title
- [x] Include UserForm component
- [x] Handle create vs edit based on route
- [x] Pre-populate form in edit mode
- [x] Navigate to list on success
- [x] Protect with required permission

---

# Phase 5: Role & Permission Management (`src/features/roles/`)

## 5.1 Role Feature Structure

```
features/roles/
├── api/
│   ├── roles.api.ts
│   ├── permissions.api.ts
│   └── roles.queries.ts
├── components/
│   ├── role-table.tsx
│   ├── role-form.tsx
│   ├── permission-list.tsx
│   └── permission-selector.tsx
├── pages/
│   ├── roles-list.page.tsx
│   └── role-detail.page.tsx
├── types/
│   └── role.types.ts
└── index.ts
```

- [x] Create folder structure for roles feature
- [x] Create barrel export for public API

## 5.2 Role Types (`types/role.types.ts`)

- [x] Define Role interface matching backend DTO
- [x] Define Permission interface matching backend DTO
- [x] Define CreateRoleRequest interface
- [x] Define UpdateRoleRequest interface
- [x] Define RolePermissions interface

## 5.3 Role API Functions

- [x] Implement getRoles API function
- [x] Implement getRole API function by ID
- [x] Implement createRole API function
- [x] Implement updateRole API function
- [x] Implement deleteRole API function
- [x] Implement getRolePermissions API function
- [x] Implement setRolePermissions API function
- [x] Implement getPermissions API function (all permissions)

## 5.4 Role React Query Hooks

- [x] Create useRoles query hook
- [x] Create useRole query hook by ID
- [x] Create usePermissions query hook (all available)
- [x] Create useCreateRole mutation
- [x] Create useUpdateRole mutation
- [x] Create useDeleteRole mutation
- [x] Create useSetRolePermissions mutation

## 5.5 Role Components

### 5.5.1 Role Table Component

- [x] Display roles in table format
- [x] Include columns: name, display name, users count, permissions count, actions
- [x] Indicate system roles (non-deletable)
- [x] Include action buttons (edit, delete, manage permissions)

### 5.5.2 Role Form Component

- [x] Support create and edit modes
- [x] Include name input (disabled for system roles)
- [x] Include display name input
- [x] Include description textarea
- [x] Include permission selector component
- [x] Handle form submission

### 5.5.3 Permission Selector Component

- [x] Display permissions grouped by resource
- [x] Support multi-select with checkboxes
- [x] Include select all / deselect all per group
- [x] Show current selection count
- [x] Support search/filter permissions

### 5.5.4 Permission List Component

- [x] Display permissions in readable format
- [x] Group by resource
- [x] Show resource:action format with description

## 5.6 Role Pages

### 5.6.1 Roles List Page

- [x] Include page header with create button
- [x] Include RoleTable component
- [x] Protect with roles:list permission

### 5.6.2 Role Detail/Edit Page

- [x] Fetch role data by ID
- [x] Include RoleForm component
- [x] Handle system role restrictions
- [x] Protect with roles:read permission

---

# Phase 6: Layout & Navigation (`src/components/layout/`)

## 6.1 Layout Components

### 6.1.1 Root Layout Component

- [x] Define overall page structure
- [x] Include header/navbar
- [x] Include sidebar (for authenticated routes)
- [x] Include main content area
- [x] Include footer (optional)
- [x] Handle responsive layout (mobile sidebar drawer)

### 6.1.2 Header Component

- [x] Include app logo/branding
- [x] Include navigation links (if not using sidebar)
- [x] Include user menu (avatar, dropdown)
- [ ] Include notifications indicator (optional)
- [x] Include theme toggle (dark/light mode)
- [x] Handle responsive behavior (hamburger menu)

### 6.1.3 Sidebar Component

- [x] Display navigation menu items
- [x] Support nested/grouped menu items
- [x] Highlight active route
- [x] Support collapsed state
- [x] Filter menu items based on permissions
- [x] Include user info section (optional)
- [x] Include logout button

### 6.1.4 User Menu Component

- [x] Display user avatar and name
- [x] Include dropdown with menu items
  - [x] Profile link
  - [x] Settings link
  - [x] Change password link
  - [x] Logout button
- [x] Handle loading state for user data

### 6.1.5 Footer Component (optional)

- [x] Include copyright information
- [x] Include version number
- [x] Include useful links

## 6.2 Navigation Configuration

- [x] Create navigation config object
- [x] Define menu items with labels, icons, paths
- [x] Define required permissions per menu item
- [x] Support nested menu groups
- [x] Create hook to filter navigation by permissions

## 6.3 Page Layout Components

### 6.3.1 Page Header Component

- [x] Include page title
- [x] Include breadcrumbs
- [x] Include action buttons slot
- [x] Support back button

### 6.3.2 Page Container Component

- [x] Apply consistent padding and max-width
- [x] Support full-width variant
- [x] Include loading state support

### 6.3.3 Card Layout Component

- [x] Wrapper for page sections
- [x] Include optional title and description
- [x] Include action buttons slot

---

# Phase 7: Routing (`src/app/router/`)

## 7.1 Router Configuration

- [x] Create router using createBrowserRouter
- [x] Define route structure with nested routes
- [x] Configure error boundaries per route level
- [x] Configure loading states for lazy routes

## 7.2 Route Definitions

### 7.2.1 Public Routes

- [x] Define login route (/login)
- [x] Define register route (/register)
- [x] Define forgot password route (/forgot-password)
- [x] Define reset password route (/reset-password)
- [x] Redirect authenticated users away from public routes

### 7.2.2 Protected Routes

- [x] Define dashboard route (/)
- [x] Define users routes (/users, /users/:id, /users/create)
- [x] Define roles routes (/roles, /roles/:id, /roles/create)
- [x] Define profile route (/profile)
- [x] Define settings route (/settings)
- [x] Wrap protected routes with AuthGuard

### 7.2.3 Error Routes

- [x] Define 404 not found route
- [x] Define 403 unauthorized route
- [x] Define 500 error route
- [x] Configure catch-all route for 404

## 7.3 Route Guards

- [x] Implement AuthGuard component for protected routes
- [x] Implement PermissionGuard component for permission-based access
- [x] Implement RoleGuard component for role-based access
- [x] Handle loading state during auth check
- [x] Preserve intended URL for post-login redirect

## 7.4 Lazy Loading

- [x] Implement React.lazy for feature pages
- [x] Create Suspense wrapper with loading fallback
- [x] Group related routes for code splitting
- [ ] Preload critical routes on hover (optional)

---

# Phase 8: Common Components (`src/components/common/`)

## 8.1 Data Display Components

### 8.1.1 Data Table Component

- [x] Accept columns configuration
- [x] Accept data array
- [x] Support sorting by columns
- [x] Support row selection
- [x] Support pagination
- [x] Display loading skeleton
- [x] Display empty state
- [x] Support responsive behavior (horizontal scroll or card view)

### 8.1.2 Pagination Component

- [x] Display current page and total pages
- [x] Include previous/next buttons
- [x] Include page number buttons with truncation
- [x] Include items per page selector
- [x] Support controlled pagination

### 8.1.3 Status Badge Component

- [x] Accept status value
- [x] Map status to color/variant
- [x] Support custom status mappings

### 8.1.4 Empty State Component

- [x] Display icon and message
- [x] Support action button
- [x] Support different variants (no data, no results, error)

### 8.1.5 Loading Components

- [x] Spinner component
- [x] Skeleton component for content placeholders
- [x] Full page loading overlay
- [x] Button loading state

## 8.2 Feedback Components

### 8.2.1 Toast/Notification System

- [x] Configure toast provider (sonner or similar)
- [x] Create toast utility functions
  - [x] success(message)
  - [x] error(message)
  - [x] warning(message)
  - [x] info(message)
- [x] Support action buttons in toasts
- [x] Configure positioning and duration

### 8.2.2 Confirmation Dialog Component

- [x] Accept title, message, confirm/cancel labels
- [x] Support destructive variant (red confirm button)
- [x] Accept onConfirm and onCancel callbacks
- [x] Create useConfirm hook for imperative usage

### 8.2.3 Alert Component

- [x] Support variants (info, success, warning, error)
- [x] Include icon based on variant
- [x] Support dismissible option

## 8.3 Form Components

### 8.3.1 Form Field Wrapper

- [x] Accept label, description, error message
- [x] Display required indicator
- [x] Apply consistent spacing
- [x] Support horizontal/vertical layout

### 8.3.2 Search Input Component

- [x] Include search icon
- [x] Support clear button
- [x] Support debounced onChange
- [x] Support loading state

### 8.3.3 Password Input Component

- [x] Include show/hide toggle button
- [x] Support password strength indicator
- [x] Inherit standard input props

### 8.3.4 Multi-Select Component

- [x] Display selected items as tags
- [x] Support search/filter options
- [x] Support select all option
- [ ] Handle large lists with virtualization (if needed)

## 8.4 Navigation Components

### 8.4.1 Breadcrumbs Component

- [x] Accept items array with label and path
- [x] Display separator between items
- [x] Make last item non-clickable (current page)
- [ ] Support auto-generation from route (optional)

### 8.4.2 Tabs Component

- [x] Support controlled and uncontrolled modes
- [x] Support URL-synced tabs (query param)
- [x] Include tab panels with lazy loading option

---

# Phase 9: Providers & App Setup (`src/app/`)

## 9.1 Provider Setup (`src/app/providers/`)

### 9.1.1 Query Provider

- [x] Create QueryClientProvider wrapper
- [x] Include ReactQueryDevtools (dev only)
- [x] Configure devtools position and default open state

### 9.1.2 Theme Provider

- [x] Implement dark/light mode toggle
- [x] Persist theme preference
- [x] Sync with system preference option
- [x] Provide theme context for components

### 9.1.3 Toast Provider

- [x] Configure toast provider from chosen library
- [x] Set global toast configuration

### 9.1.4 Auth Provider (optional)

- [x] Provide auth context if not using Zustand exclusively
- [x] Initialize auth state on app mount
- [x] Handle token refresh on mount

### 9.1.5 Combined Providers Component

- [x] Create AppProviders component combining all providers
- [x] Order providers correctly (query outside, theme, toast, etc.)
- [x] Accept children prop

## 9.2 App Initialization

### 9.2.1 App Entry Point (main.tsx)

- [x] Import global styles
- [x] Render App component inside React.StrictMode
- [x] Mount to root DOM element

### 9.2.2 App Component (App.tsx)

- [x] Wrap with AppProviders
- [x] Include RouterProvider
- [x] Include global error boundary
- [x] Include initial data loading (auth check)

### 9.2.3 Auth Initialization

- [x] Check for stored tokens on app mount
- [x] Validate tokens and fetch current user
- [x] Handle expired/invalid tokens
- [x] Set loading state during initialization
- [x] Redirect based on auth state

---

# Phase 10: Error Handling & Boundaries

## 10.1 Error Boundary Components

### 10.1.1 Global Error Boundary

- [x] Catch uncaught errors in React tree
- [x] Display friendly error message
- [x] Include "Reload" button
- [ ] Log errors to monitoring service (production)
- [x] Show stack trace in development

### 10.1.2 Route Error Boundary

- [x] Handle errors within specific routes
- [x] Display route-specific error UI
- [x] Include navigation back to safe route
- [x] Handle 404 and other HTTP errors from loaders

### 10.1.3 Feature Error Boundary

- [x] Wrap feature components
- [x] Show feature-specific error UI
- [x] Allow retry without full page reload

## 10.2 Error Handling Utilities

### 10.2.1 Error Logging Service

- [ ] Create error logging utility
- [ ] Integrate with monitoring service (Sentry, LogRocket)
- [ ] Include user context in error reports
- [ ] Include breadcrumbs/actions before error
- [ ] Filter sensitive information

### 10.2.2 API Error Handler

- [x] Create central API error handling utility
- [x] Map error codes to user messages
- [x] Handle network errors
- [x] Handle timeout errors
- [x] Display appropriate toast/notification

## 10.3 Not Found & Unauthorized Pages

### 10.3.1 404 Not Found Page

- [x] Display friendly message
- [x] Include illustration/icon
- [x] Include link to home/dashboard
- [ ] Include search option (optional)

### 10.3.2 403 Unauthorized Page

- [x] Display permission denied message
- [x] Explain what permission is needed
- [ ] Include link to request access (optional)
- [x] Include link to home/dashboard

### 10.3.3 500 Error Page

- [x] Display generic error message
- [x] Include error reference ID (for support)
- [x] Include retry button
- [x] Include contact support link

---

# Phase 11: Testing

## 11.1 Testing Setup

### 11.1.1 Vitest Configuration

- [x] Configure vitest.config.ts
- [x] Set up test environment (jsdom)
- [x] Configure coverage thresholds
- [x] Set up path aliases for tests
- [x] Configure test globals

### 11.1.2 Testing Library Setup

- [x] Configure @testing-library/react
- [x] Create custom render function with providers
- [x] Set up user-event for interactions
- [x] Configure screen queries

### 11.1.3 MSW Setup for API Mocking

- [x] Install and configure msw
- [x] Create mock handlers for auth endpoints
- [x] Create mock handlers for user endpoints
- [x] Create mock handlers for role endpoints
- [x] Set up mock server for tests
- [x] Configure request handlers reset between tests

## 11.2 Unit Tests

### 11.2.1 Utility Function Tests

- [x] Test API client utilities
- [x] Test form validation schemas
- [x] Test helper functions
- [x] Test formatters (date, currency, etc.)

### 11.2.2 Hook Tests

- [x] Test custom hooks with renderHook
- [x] Test auth hooks (useAuth, useLogin)
- [x] Test form hooks
- [x] Test query hooks with mock data

### 11.2.3 Store Tests

- [x] Test Zustand store actions
- [x] Test store selectors
- [x] Test store persistence

## 11.3 Component Tests

### 11.3.1 UI Component Tests

- [x] Test Button component variants
- [x] Test Input component with validation
- [x] Test Form components
- [x] Test Dialog/Modal components
- [x] Test Table component

### 11.3.2 Feature Component Tests

- [x] Test LoginForm submission
- [x] Test RegisterForm validation
- [x] Test UserTable with mock data
- [x] Test UserForm create and edit modes

### 11.3.3 Page Tests

- [x] Test page rendering with mock data
- [x] Test page navigation
- [x] Test protected page access

## 11.4 Integration Tests

### 11.4.1 Auth Flow Tests

- [x] Test complete login flow
- [x] Test complete registration flow
- [x] Test logout flow
- [x] Test token refresh flow
- [x] Test protected route access

### 11.4.2 CRUD Flow Tests

- [x] Test user list with filters
- [x] Test user creation
- [x] Test user update
- [x] Test user deletion
- [x] Test role management

## 11.5 E2E Tests (Optional)

### 11.5.1 Playwright/Cypress Setup

- [ ] Configure E2E testing framework
- [ ] Set up test database/environment
- [ ] Create test fixtures
- [ ] Configure CI pipeline for E2E

### 11.5.2 Critical Path Tests

- [ ] Test login to dashboard flow
- [ ] Test user management flow
- [ ] Test role assignment flow
- [ ] Test session management flow

---

# Phase 12: Performance Optimization

## 12.1 Bundle Optimization

- [ ] Analyze bundle size with rollup-plugin-visualizer
- [ ] Configure code splitting by route
- [ ] Configure vendor chunk for dependencies
- [ ] Remove unused dependencies
- [ ] Tree-shake unused exports
- [ ] Compress assets with gzip/brotli
- [ ] Configure proper caching headers for assets

## 12.2 React Performance

### 12.2.1 Component Optimization

- [ ] Implement React.memo for expensive components
- [ ] Use useMemo for expensive computations
- [ ] Use useCallback for callback props
- [ ] Avoid inline object/function props
- [ ] Profile with React DevTools

### 12.2.2 List Rendering

- [ ] Use proper keys for list items
- [ ] Implement virtualization for long lists (react-window)
- [ ] Paginate large data sets
- [ ] Implement infinite scroll where appropriate

### 12.2.3 State Management

- [ ] Keep state as local as possible
- [ ] Use selectors to prevent unnecessary rerenders
- [ ] Split stores to minimize subscription scope
- [ ] Use React Query for server state

## 12.3 Network Optimization

### 12.3.1 API Optimization

- [ ] Configure appropriate staleTime for queries
- [ ] Implement request deduplication
- [ ] Use prefetching for predictable navigations
- [ ] Implement optimistic updates for mutations
- [ ] Configure retry logic with backoff

### 12.3.2 Asset Optimization

- [ ] Optimize images (WebP format, proper sizing)
- [ ] Lazy load images below fold
- [ ] Preload critical assets
- [ ] Use CDN for static assets

## 12.4 Loading Performance

### 12.4.1 Initial Load

- [ ] Minimize critical rendering path
- [ ] Inline critical CSS
- [ ] Defer non-critical JavaScript
- [ ] Implement route-based code splitting
- [ ] Add loading indicators for chunks

### 12.4.2 Perceived Performance

- [ ] Implement skeleton screens
- [ ] Use optimistic UI updates
- [ ] Add transition animations
- [ ] Prefetch on hover/focus

---

# Phase 13: Accessibility (a11y)

## 13.1 Semantic HTML

- [ ] Use proper heading hierarchy (h1-h6)
- [ ] Use semantic elements (nav, main, article, section)
- [ ] Use proper list elements (ul, ol, li)
- [ ] Use proper form elements with labels
- [ ] Use button for actions, a for navigation

## 13.2 Keyboard Navigation

- [ ] Ensure all interactive elements are focusable
- [ ] Implement visible focus indicators
- [ ] Support tab navigation order
- [ ] Implement keyboard shortcuts for common actions
- [ ] Handle focus trap in modals/dialogs
- [ ] Restore focus after modal close

## 13.3 Screen Reader Support

- [ ] Add alt text to all images
- [ ] Add aria-labels to icon-only buttons
- [ ] Implement aria-live regions for dynamic content
- [ ] Add proper ARIA roles where needed
- [ ] Test with screen readers (NVDA, VoiceOver)

## 13.4 Color & Contrast

- [ ] Ensure sufficient color contrast (WCAG AA)
- [ ] Don't rely on color alone for information
- [ ] Support high contrast mode
- [ ] Test with color blindness simulators

## 13.5 Testing Accessibility

- [ ] Install axe-core for automated testing
- [ ] Add accessibility tests to component tests
- [ ] Run automated accessibility audits in CI
- [ ] Perform manual accessibility testing

---

# Phase 14: Internationalization (i18n) - Optional

## 14.1 i18n Setup

- [ ] Choose i18n library (react-i18next recommended)
- [ ] Configure i18n with default language
- [ ] Set up language detection (browser, stored preference)
- [ ] Configure fallback language
- [ ] Set up namespace separation

## 14.2 Translation Setup

- [ ] Create translation files structure
- [ ] Define translation keys conventions
- [ ] Set up translation extraction workflow
- [ ] Configure pluralization rules
- [ ] Configure interpolation

## 14.3 Component Integration

- [ ] Replace hardcoded strings with translation keys
- [ ] Handle date/number formatting per locale
- [ ] Handle RTL languages (if needed)
- [ ] Create language switcher component

---

# Phase 15: Production Deployment

## 15.1 Build Configuration

### 15.1.1 Production Build

- [ ] Configure environment variables for production
- [ ] Enable minification and tree shaking
- [ ] Generate source maps for error tracking
- [ ] Configure asset hashing for cache busting
- [ ] Optimize chunk splitting
- [ ] Remove development-only code

### 15.1.2 Build Verification

- [ ] Run production build locally
- [ ] Test production build with preview server
- [ ] Verify environment variables are applied
- [ ] Check bundle sizes
- [ ] Verify no development warnings

## 15.2 Docker Setup

### 15.2.1 Dockerfile

- [ ] Create multi-stage Dockerfile
- [ ] Use Bun image for build stage
- [ ] Use nginx or static file server for runtime
- [ ] Configure nginx for SPA routing
- [ ] Optimize image size
- [ ] Set proper security headers in nginx

### 15.2.2 Docker Compose

- [ ] Create docker-compose for local testing
- [ ] Configure environment variables
- [ ] Set up health checks

## 15.3 CI/CD Pipeline

### 15.3.1 Continuous Integration

- [ ] Configure CI workflow (GitHub Actions, GitLab CI)
- [ ] Run linting on PRs
- [ ] Run unit tests with coverage
- [ ] Run type checking
- [ ] Run build verification
- [ ] Run accessibility audit
- [ ] Fail on coverage below threshold

### 15.3.2 Continuous Deployment

- [ ] Configure deployment workflow
- [ ] Build Docker image
- [ ] Push to container registry
- [ ] Deploy to staging environment
- [ ] Run smoke tests on staging
- [ ] Deploy to production (manual approval)
- [ ] Configure rollback procedure

## 15.4 Environment Configuration

### 15.4.1 Staging Environment

- [ ] Configure staging API URL
- [ ] Enable verbose logging
- [ ] Enable debug tools
- [ ] Configure feature flags for testing

### 15.4.2 Production Environment

- [ ] Configure production API URL
- [ ] Disable debug tools
- [ ] Configure error tracking
- [ ] Configure analytics (if applicable)
- [ ] Configure CSP headers

## 15.5 Monitoring & Analytics

### 15.5.1 Error Tracking

- [ ] Integrate Sentry or similar
- [ ] Configure error grouping
- [ ] Set up alerts for error spikes
- [ ] Include user context in reports
- [ ] Configure source maps upload

### 15.5.2 Performance Monitoring

- [ ] Configure web vitals tracking
- [ ] Set up performance budgets
- [ ] Monitor bundle size changes
- [ ] Track API response times

### 15.5.3 Analytics (Optional)

- [ ] Integrate analytics service
- [ ] Track page views
- [ ] Track key user actions
- [ ] Ensure GDPR compliance

---

# Phase 16: Documentation

## 16.1 Code Documentation

- [ ] Add JSDoc comments to exported functions
- [ ] Document component props with TypeScript
- [ ] Add README files to feature folders
- [ ] Document custom hooks usage
- [ ] Document store structure and actions

## 16.2 Component Documentation

### 16.2.1 Storybook Setup (Optional)

- [ ] Install and configure Storybook
- [ ] Create stories for UI components
- [ ] Document component variants
- [ ] Add interaction tests
- [ ] Deploy Storybook for team reference

## 16.3 Project Documentation

### 16.3.1 README

- [ ] Project overview and purpose
- [ ] Prerequisites and requirements
- [ ] Installation instructions
- [ ] Development workflow
- [ ] Available scripts
- [ ] Environment configuration
- [ ] Deployment instructions
- [ ] Contributing guidelines

### 16.3.2 Architecture Documentation

- [ ] Document folder structure
- [ ] Document state management approach
- [ ] Document API integration patterns
- [ ] Document authentication flow
- [ ] Document RBAC implementation
- [ ] Create architecture diagrams

### 16.3.3 API Documentation

- [ ] Document API client usage
- [ ] Document query hooks
- [ ] Document mutation hooks
- [ ] Provide usage examples

---

# Final Checklist Before Production

- [ ] All features implemented and tested
- [ ] All tests passing with adequate coverage (>80%)
- [ ] No TypeScript errors
- [ ] No ESLint errors or warnings
- [ ] Accessibility audit passing
- [ ] Performance metrics within budget
- [ ] Security headers configured
- [ ] Error tracking configured
- [ ] Environment variables documented
- [ ] Docker build tested
- [ ] CI/CD pipeline working
- [ ] Documentation complete
- [ ] Code review completed
- [ ] Manual QA completed
- [ ] Staging deployment verified
- [ ] Rollback procedure tested
