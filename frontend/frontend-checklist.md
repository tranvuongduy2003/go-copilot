# Expert-Level Frontend Setup Checklist

## Vite + Bun + React + TypeScript | Production-Ready SPA

---

# Phase 1: Project Foundation

## 1.1 Project Initialization

- [ ] Initialize project using Bun with Vite React TypeScript template
- [ ] Configure `bunfig.toml` for Bun-specific settings
- [ ] Update `package.json` with correct project metadata (name, version, description)
- [ ] Configure `tsconfig.json` with strict TypeScript settings
  - [ ] Enable strict mode, noImplicitAny, strictNullChecks
  - [ ] Configure path aliases (@/ for src/)
  - [ ] Set target to ES2022 for modern JavaScript features
  - [ ] Enable skipLibCheck for faster builds
- [ ] Create `.gitignore` with appropriate patterns (node_modules, dist, .env)
- [ ] Create `.editorconfig` for consistent coding style
- [ ] Create `README.md` with project overview and setup instructions
- [ ] Initialize git repository with initial commit

## 1.2 Vite Configuration (`vite.config.ts`)

- [ ] Configure path aliases to match tsconfig
- [ ] Configure development server port and host
- [ ] Enable CORS for API communication during development
- [ ] Configure proxy for API requests to backend
- [ ] Configure build output settings
- [ ] Enable source maps for development
- [ ] Configure chunk splitting strategy for production
- [ ] Set up environment variable handling (VITE_ prefix)
- [ ] Configure preview server for production build testing
- [ ] Add build analysis plugin (rollup-plugin-visualizer) for bundle optimization

## 1.3 Environment Configuration

- [ ] Create `.env.example` with all required variables documented
- [ ] Define `VITE_API_BASE_URL` for backend API endpoint
- [ ] Define `VITE_APP_NAME` for application branding
- [ ] Define `VITE_APP_VERSION` for version tracking
- [ ] Define `VITE_ENABLE_MOCK` for mock API toggle (development)
- [ ] Create `.env.development` for development defaults
- [ ] Create `.env.production` for production defaults
- [ ] Create environment type declarations (`env.d.ts`)
- [ ] Implement environment validation on app startup

## 1.4 Package Dependencies

### Core Dependencies
- [ ] Install React 18+ and React DOM
- [ ] Install React Router DOM v6+ for routing
- [ ] Install @tanstack/react-query for server state management
- [ ] Install zustand for client state management
- [ ] Install axios for HTTP requests
- [ ] Install react-hook-form for form management
- [ ] Install zod for schema validation
- [ ] Install @hookform/resolvers for zod integration with react-hook-form
- [ ] Install date-fns for date manipulation
- [ ] Install clsx and tailwind-merge for className utilities

### UI Dependencies
- [ ] Install Tailwind CSS and required PostCSS plugins
- [ ] Install shadcn/ui CLI and initialize
- [ ] Install lucide-react for icons
- [ ] Install class-variance-authority for component variants
- [ ] Install @radix-ui primitives as needed by shadcn components
- [ ] Install sonner or react-hot-toast for notifications
- [ ] Install framer-motion for animations (optional)

### Development Dependencies
- [ ] Install TypeScript and @types/react, @types/react-dom
- [ ] Install ESLint with React and TypeScript plugins
- [ ] Install Prettier for code formatting
- [ ] Install @typescript-eslint/parser and @typescript-eslint/eslint-plugin
- [ ] Install eslint-plugin-react-hooks for hooks linting
- [ ] Install eslint-plugin-react-refresh for Fast Refresh compatibility
- [ ] Install husky for git hooks
- [ ] Install lint-staged for pre-commit linting
- [ ] Install @testing-library/react for component testing
- [ ] Install vitest for unit testing
- [ ] Install msw for API mocking in tests
- [ ] Install playwright or cypress for E2E testing (optional)

## 1.5 Code Quality Setup

### ESLint Configuration (`.eslintrc.cjs` or `eslint.config.js`)
- [ ] Configure parser for TypeScript
- [ ] Enable React and React Hooks plugins
- [ ] Enable import sorting rules
- [ ] Configure no-unused-vars with TypeScript override
- [ ] Enable strict type-checking rules
- [ ] Configure path alias resolution for import plugin
- [ ] Add custom rules for project conventions

### Prettier Configuration (`.prettierrc`)
- [ ] Configure semi, singleQuote, tabWidth, trailingComma
- [ ] Configure printWidth (80-120 recommended)
- [ ] Configure endOfLine for cross-platform compatibility
- [ ] Create `.prettierignore` for build artifacts

### Git Hooks
- [ ] Initialize husky with `bunx husky init`
- [ ] Create pre-commit hook for lint-staged
- [ ] Create commit-msg hook for conventional commits (optional)
- [ ] Configure lint-staged in `package.json` or `.lintstagedrc`
- [ ] Run ESLint and Prettier on staged files

## 1.6 Tailwind CSS Configuration

- [ ] Initialize Tailwind CSS with `bunx tailwindcss init -p`
- [ ] Configure content paths in `tailwind.config.js`
- [ ] Extend theme with custom colors matching design system
- [ ] Configure custom spacing, typography, and breakpoints
- [ ] Add CSS variables for theming (light/dark mode support)
- [ ] Configure animation keyframes for custom animations
- [ ] Set up Tailwind plugins (forms, typography, aspect-ratio)
- [ ] Create base styles in `index.css` with Tailwind directives

## 1.7 shadcn/ui Setup

- [ ] Run `bunx shadcn-ui@latest init`
- [ ] Configure components.json with correct paths
- [ ] Select style (default or new-york)
- [ ] Configure base color and CSS variables
- [ ] Install essential components: button, input, card, form, label
- [ ] Install feedback components: alert, toast, dialog, drawer
- [ ] Install data display components: table, badge, avatar
- [ ] Install navigation components: dropdown-menu, tabs, navigation-menu
- [ ] Install form components: select, checkbox, radio-group, switch
- [ ] Create component barrel exports for cleaner imports

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

- [ ] Create folder structure as defined above
- [ ] Move existing components to appropriate locations
- [ ] Create barrel exports (index.ts) for each major folder
- [ ] Configure path aliases for new folders in tsconfig and vite.config

## 2.2 API Layer Setup (`src/lib/api/`)

### 2.2.1 Axios Client Configuration
- [ ] Create axios instance with base configuration
- [ ] Configure base URL from environment variable
- [ ] Set default headers (Content-Type, Accept)
- [ ] Configure timeout (30 seconds recommended)
- [ ] Create request interceptor for authentication
  - [ ] Automatically attach access token to requests
  - [ ] Skip token for public endpoints
- [ ] Create response interceptor for error handling
  - [ ] Handle 401 errors with token refresh logic
  - [ ] Handle network errors gracefully
  - [ ] Implement request retry for transient failures
- [ ] Create response interceptor for data transformation
  - [ ] Unwrap successful responses for cleaner data access

### 2.2.2 Token Refresh Implementation
- [ ] Implement token refresh mechanism
- [ ] Queue failed requests during refresh
- [ ] Retry queued requests after successful refresh
- [ ] Redirect to login on refresh failure
- [ ] Prevent multiple simultaneous refresh requests
- [ ] Handle refresh token expiration

### 2.2.3 API Error Handling
- [ ] Define ApiError class with typed error structure
  - [ ] Include status code, message, and error details
  - [ ] Include original error for debugging
- [ ] Create error type guards for specific error types
- [ ] Map backend error codes to user-friendly messages
- [ ] Create error boundary for API errors

### 2.2.4 API Endpoints Organization
- [ ] Create endpoints constants file
- [ ] Group endpoints by feature/resource
- [ ] Use string templates for parameterized URLs
- [ ] Document each endpoint with expected request/response

## 2.3 React Query Setup (`src/lib/api/query-client.ts`)

### 2.3.1 Query Client Configuration
- [ ] Create QueryClient instance with default options
- [ ] Configure default staleTime (5 minutes for most data)
- [ ] Configure default gcTime (cacheTime) (10 minutes)
- [ ] Configure retry logic (3 retries with exponential backoff)
- [ ] Configure refetchOnWindowFocus (false for most apps)
- [ ] Configure refetchOnMount behavior
- [ ] Set up query client provider in app providers

### 2.3.2 Query Key Factory
- [ ] Create query key factory for consistent key management
- [ ] Define keys by feature (auth, users, etc.)
- [ ] Include parameters in keys for cache granularity
- [ ] Export typed query keys for type safety
- [ ] Document query key patterns

### 2.3.3 Custom Query Hooks Pattern
- [ ] Define pattern for query hooks (useQuery wrapper)
- [ ] Define pattern for mutation hooks (useMutation wrapper)
- [ ] Include loading, error, and success states
- [ ] Include automatic cache invalidation on mutations
- [ ] Include optimistic updates where appropriate

## 2.4 State Management with Zustand (`src/stores/`)

### 2.4.1 Store Architecture
- [ ] Define store slicing strategy (one store per domain)
- [ ] Create auth store for authentication state
- [ ] Create UI store for global UI state (sidebar, modals)
- [ ] Create user preferences store (theme, language)
- [ ] Implement persist middleware for relevant stores
- [ ] Configure storage adapter (localStorage)
- [ ] Define persisted vs non-persisted state

### 2.4.2 Auth Store Implementation
- [ ] Define auth state interface
  - [ ] Include user object (nullable)
  - [ ] Include authentication status (loading, authenticated, unauthenticated)
  - [ ] Include tokens (if storing in memory)
- [ ] Define auth actions
  - [ ] setUser action
  - [ ] clearAuth action
  - [ ] updateUser action
- [ ] Implement selectors for derived state
  - [ ] isAuthenticated selector
  - [ ] hasPermission selector (for RBAC)
  - [ ] hasRole selector

### 2.4.3 UI Store Implementation
- [ ] Define UI state interface
  - [ ] Include sidebar collapsed state
  - [ ] Include active modal state
  - [ ] Include global loading state
  - [ ] Include notification queue (if not using toast library)
- [ ] Define UI actions for state mutations
- [ ] Implement reset action for cleanup

## 2.5 Form Management Setup

### 2.5.1 Zod Schema Definitions (`src/lib/validations/`)
- [ ] Create base validation schemas for common types
  - [ ] Email schema with proper regex
  - [ ] Password schema with strength requirements
  - [ ] UUID schema for IDs
  - [ ] Pagination schema for list queries
- [ ] Create auth-related schemas
  - [ ] Login schema (email, password)
  - [ ] Register schema (email, password, confirmPassword, fullName)
  - [ ] Forgot password schema
  - [ ] Reset password schema
  - [ ] Change password schema
- [ ] Create user-related schemas
  - [ ] Create user schema
  - [ ] Update user schema (partial)
  - [ ] User filter schema for list queries
- [ ] Export TypeScript types inferred from schemas

### 2.5.2 Form Components Pattern
- [ ] Create form field wrapper component
  - [ ] Integrate with react-hook-form Controller
  - [ ] Display validation errors from form state
  - [ ] Support all input types (text, password, select, etc.)
- [ ] Create reusable form components
  - [ ] FormInput component
  - [ ] FormSelect component
  - [ ] FormCheckbox component
  - [ ] FormTextarea component
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

- [ ] Create folder structure for auth feature
- [ ] Create barrel export (index.ts) for public API

## 3.2 Auth Types (`types/auth.types.ts`)

- [ ] Define User interface matching backend UserDTO
  - [ ] Include id, email, fullName, status, roles, permissions
- [ ] Define LoginRequest interface
- [ ] Define LoginResponse interface with tokens and user
- [ ] Define RegisterRequest interface
- [ ] Define TokenPair interface (accessToken, refreshToken, expiresIn)
- [ ] Define AuthState interface for store
- [ ] Define Permission and Role types for RBAC
- [ ] Define ForgotPasswordRequest interface
- [ ] Define ResetPasswordRequest interface
- [ ] Define ChangePasswordRequest interface

## 3.3 Auth API (`api/auth.api.ts`)

- [ ] Implement login API function
  - [ ] Accept credentials, return token pair and user
- [ ] Implement register API function
- [ ] Implement logout API function
  - [ ] Call backend logout endpoint
- [ ] Implement refresh token API function
- [ ] Implement forgot password API function
- [ ] Implement reset password API function
- [ ] Implement change password API function
- [ ] Implement get current user API function (GET /auth/me)
- [ ] Implement get sessions API function
- [ ] Implement revoke session API function

## 3.4 Auth React Query Hooks (`api/auth.queries.ts`)

### 3.4.1 Auth Mutations
- [ ] Create useLogin mutation hook
  - [ ] On success: store tokens, update auth store, redirect
  - [ ] On error: display error message
- [ ] Create useRegister mutation hook
  - [ ] On success: auto-login or redirect to login
- [ ] Create useLogout mutation hook
  - [ ] On success: clear auth store, clear query cache, redirect
- [ ] Create useForgotPassword mutation hook
  - [ ] On success: display success message
- [ ] Create useResetPassword mutation hook
  - [ ] On success: redirect to login
- [ ] Create useChangePassword mutation hook
  - [ ] On success: display success message, optionally logout

### 3.4.2 Auth Queries
- [ ] Create useCurrentUser query hook
  - [ ] Fetch current user data on app mount
  - [ ] Update auth store with user data
  - [ ] Handle unauthorized (token expired)
- [ ] Create useSessions query hook for active sessions list
- [ ] Create useRevokeSession mutation hook

## 3.5 Auth Components

### 3.5.1 Login Form Component
- [ ] Integrate react-hook-form with zod resolver
- [ ] Include email input with validation
- [ ] Include password input with show/hide toggle
- [ ] Include "Remember me" checkbox (optional)
- [ ] Include "Forgot password" link
- [ ] Include submit button with loading state
- [ ] Display form-level errors from API
- [ ] Display field-level validation errors
- [ ] Handle form submission with useLogin hook
- [ ] Redirect to dashboard on success

### 3.5.2 Register Form Component
- [ ] Include email input with validation
- [ ] Include full name input with validation
- [ ] Include password input with strength indicator
- [ ] Include confirm password input with match validation
- [ ] Include terms acceptance checkbox
- [ ] Include submit button with loading state
- [ ] Display validation errors
- [ ] Handle form submission with useRegister hook
- [ ] Include link to login page

### 3.5.3 Forgot Password Form Component
- [ ] Include email input with validation
- [ ] Include submit button with loading state
- [ ] Display success message after submission
- [ ] Include link back to login

### 3.5.4 Reset Password Form Component
- [ ] Extract reset token from URL query params
- [ ] Include new password input with strength indicator
- [ ] Include confirm password input
- [ ] Include submit button with loading state
- [ ] Handle invalid/expired token error
- [ ] Redirect to login on success

### 3.5.5 Auth Guard Component
- [ ] Check authentication status
- [ ] Redirect to login if not authenticated
- [ ] Show loading state while checking auth
- [ ] Render children if authenticated
- [ ] Accept required permissions/roles props
- [ ] Check permissions if specified
- [ ] Redirect to unauthorized page if permission denied

## 3.6 Auth Pages

### 3.6.1 Login Page
- [ ] Create page layout with centered form
- [ ] Include app logo/branding
- [ ] Include LoginForm component
- [ ] Include link to register page
- [ ] Include social login options (if applicable)
- [ ] Redirect if already authenticated

### 3.6.2 Register Page
- [ ] Create page layout with centered form
- [ ] Include app logo/branding
- [ ] Include RegisterForm component
- [ ] Include link to login page
- [ ] Redirect if already authenticated

### 3.6.3 Forgot Password Page
- [ ] Create page layout with centered form
- [ ] Include ForgotPasswordForm component
- [ ] Include link back to login

### 3.6.4 Reset Password Page
- [ ] Create page layout with centered form
- [ ] Validate token presence in URL
- [ ] Include ResetPasswordForm component
- [ ] Handle token validation errors

## 3.7 Token Management

- [ ] Implement secure token storage strategy
  - [ ] Option 1: Memory only (most secure, no persistence)
  - [ ] Option 2: localStorage (convenient, less secure)
  - [ ] Option 3: HttpOnly cookies (requires backend support)
- [ ] Implement token persistence (if using localStorage)
  - [ ] Store encrypted or obfuscated tokens
  - [ ] Clear tokens on logout
- [ ] Implement token refresh on app mount
- [ ] Check token expiration before requests
- [ ] Implement automatic token refresh before expiration
  - [ ] Set up refresh interval or intercept 401 responses

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

- [ ] Create folder structure for users feature
- [ ] Create barrel export for public API

## 4.2 User Types (`types/user.types.ts`)

- [ ] Define User interface matching backend DTO
- [ ] Define UserStatus enum (pending, active, inactive, banned)
- [ ] Define CreateUserRequest interface
- [ ] Define UpdateUserRequest interface
- [ ] Define UserFilter interface for list queries
- [ ] Define PaginatedUsers interface for list response
- [ ] Define UserRole interface
- [ ] Define UserPermission interface

## 4.3 User API (`api/users.api.ts`)

- [ ] Implement getUsers API function with pagination and filters
- [ ] Implement getUser API function by ID
- [ ] Implement createUser API function
- [ ] Implement updateUser API function
- [ ] Implement deleteUser API function
- [ ] Implement activateUser API function
- [ ] Implement deactivateUser API function
- [ ] Implement getUserRoles API function
- [ ] Implement assignRole API function
- [ ] Implement revokeRole API function

## 4.4 User React Query Hooks (`api/users.queries.ts`)

### 4.4.1 User Queries
- [ ] Create useUsers query hook with pagination
  - [ ] Accept page, limit, filters as parameters
  - [ ] Return paginated data with metadata
  - [ ] Enable keepPreviousData for smooth pagination
- [ ] Create useUser query hook by ID
  - [ ] Enable caching with appropriate staleTime
- [ ] Create useUserRoles query hook

### 4.4.2 User Mutations
- [ ] Create useCreateUser mutation
  - [ ] Invalidate users list cache on success
- [ ] Create useUpdateUser mutation
  - [ ] Invalidate specific user and list cache
  - [ ] Implement optimistic update (optional)
- [ ] Create useDeleteUser mutation
  - [ ] Invalidate users list cache on success
- [ ] Create useActivateUser mutation
- [ ] Create useDeactivateUser mutation
- [ ] Create useAssignRole mutation
- [ ] Create useRevokeRole mutation

## 4.5 User Components

### 4.5.1 User Table Component
- [ ] Display users in responsive table
- [ ] Include columns: name, email, status, roles, created, actions
- [ ] Implement column sorting
- [ ] Implement row selection (if bulk actions needed)
- [ ] Include action buttons (edit, delete, activate/deactivate)
- [ ] Display loading skeleton during fetch
- [ ] Handle empty state with appropriate message
- [ ] Implement pagination controls

### 4.5.2 User Form Component
- [ ] Support create and edit modes
- [ ] Include all user fields with validation
- [ ] Include role assignment (multi-select)
- [ ] Include status selection (for edit mode)
- [ ] Handle form submission with appropriate mutation
- [ ] Display loading state during submission
- [ ] Display success/error feedback

### 4.5.3 User Filters Component
- [ ] Include search input (debounced)
- [ ] Include status filter dropdown
- [ ] Include role filter dropdown
- [ ] Include date range filter (optional)
- [ ] Include clear filters button
- [ ] Sync filters with URL query params

### 4.5.4 User Card Component
- [ ] Display user info in card format (for grid view)
- [ ] Include avatar, name, email, status badge
- [ ] Include quick action buttons
- [ ] Support click to navigate to detail

### 4.5.5 User Actions Component
- [ ] Dropdown menu for user actions
- [ ] Include edit action
- [ ] Include activate/deactivate action (conditional)
- [ ] Include delete action with confirmation
- [ ] Include view sessions action
- [ ] Check permissions before showing actions

## 4.6 User Pages

### 4.6.1 Users List Page
- [ ] Include page header with title and create button
- [ ] Include UserFilters component
- [ ] Include UserTable component
- [ ] Include pagination component
- [ ] Support view toggle (table/grid) (optional)
- [ ] Protect with required permission

### 4.6.2 User Detail Page
- [ ] Fetch user data by ID from URL params
- [ ] Display user information
- [ ] Display user roles and permissions
- [ ] Display user sessions (optional)
- [ ] Include edit button linking to edit form
- [ ] Include back button to list
- [ ] Handle user not found error

### 4.6.3 User Create/Edit Page
- [ ] Include page header with appropriate title
- [ ] Include UserForm component
- [ ] Handle create vs edit based on route
- [ ] Pre-populate form in edit mode
- [ ] Navigate to list on success
- [ ] Protect with required permission

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

- [ ] Create folder structure for roles feature
- [ ] Create barrel export for public API

## 5.2 Role Types (`types/role.types.ts`)

- [ ] Define Role interface matching backend DTO
- [ ] Define Permission interface matching backend DTO
- [ ] Define CreateRoleRequest interface
- [ ] Define UpdateRoleRequest interface
- [ ] Define RolePermissions interface

## 5.3 Role API Functions

- [ ] Implement getRoles API function
- [ ] Implement getRole API function by ID
- [ ] Implement createRole API function
- [ ] Implement updateRole API function
- [ ] Implement deleteRole API function
- [ ] Implement getRolePermissions API function
- [ ] Implement setRolePermissions API function
- [ ] Implement getPermissions API function (all permissions)

## 5.4 Role React Query Hooks

- [ ] Create useRoles query hook
- [ ] Create useRole query hook by ID
- [ ] Create usePermissions query hook (all available)
- [ ] Create useCreateRole mutation
- [ ] Create useUpdateRole mutation
- [ ] Create useDeleteRole mutation
- [ ] Create useSetRolePermissions mutation

## 5.5 Role Components

### 5.5.1 Role Table Component
- [ ] Display roles in table format
- [ ] Include columns: name, display name, users count, permissions count, actions
- [ ] Indicate system roles (non-deletable)
- [ ] Include action buttons (edit, delete, manage permissions)

### 5.5.2 Role Form Component
- [ ] Support create and edit modes
- [ ] Include name input (disabled for system roles)
- [ ] Include display name input
- [ ] Include description textarea
- [ ] Include permission selector component
- [ ] Handle form submission

### 5.5.3 Permission Selector Component
- [ ] Display permissions grouped by resource
- [ ] Support multi-select with checkboxes
- [ ] Include select all / deselect all per group
- [ ] Show current selection count
- [ ] Support search/filter permissions

### 5.5.4 Permission List Component
- [ ] Display permissions in readable format
- [ ] Group by resource
- [ ] Show resource:action format with description

## 5.6 Role Pages

### 5.6.1 Roles List Page
- [ ] Include page header with create button
- [ ] Include RoleTable component
- [ ] Protect with roles:list permission

### 5.6.2 Role Detail/Edit Page
- [ ] Fetch role data by ID
- [ ] Include RoleForm component
- [ ] Handle system role restrictions
- [ ] Protect with roles:read permission

---

# Phase 6: Layout & Navigation (`src/components/layout/`)

## 6.1 Layout Components

### 6.1.1 Root Layout Component
- [ ] Define overall page structure
- [ ] Include header/navbar
- [ ] Include sidebar (for authenticated routes)
- [ ] Include main content area
- [ ] Include footer (optional)
- [ ] Handle responsive layout (mobile sidebar drawer)

### 6.1.2 Header Component
- [ ] Include app logo/branding
- [ ] Include navigation links (if not using sidebar)
- [ ] Include user menu (avatar, dropdown)
- [ ] Include notifications indicator (optional)
- [ ] Include theme toggle (dark/light mode)
- [ ] Handle responsive behavior (hamburger menu)

### 6.1.3 Sidebar Component
- [ ] Display navigation menu items
- [ ] Support nested/grouped menu items
- [ ] Highlight active route
- [ ] Support collapsed state
- [ ] Filter menu items based on permissions
- [ ] Include user info section (optional)
- [ ] Include logout button

### 6.1.4 User Menu Component
- [ ] Display user avatar and name
- [ ] Include dropdown with menu items
  - [ ] Profile link
  - [ ] Settings link
  - [ ] Change password link
  - [ ] Logout button
- [ ] Handle loading state for user data

### 6.1.5 Footer Component (optional)
- [ ] Include copyright information
- [ ] Include version number
- [ ] Include useful links

## 6.2 Navigation Configuration

- [ ] Create navigation config object
- [ ] Define menu items with labels, icons, paths
- [ ] Define required permissions per menu item
- [ ] Support nested menu groups
- [ ] Create hook to filter navigation by permissions

## 6.3 Page Layout Components

### 6.3.1 Page Header Component
- [ ] Include page title
- [ ] Include breadcrumbs
- [ ] Include action buttons slot
- [ ] Support back button

### 6.3.2 Page Container Component
- [ ] Apply consistent padding and max-width
- [ ] Support full-width variant
- [ ] Include loading state support

### 6.3.3 Card Layout Component
- [ ] Wrapper for page sections
- [ ] Include optional title and description
- [ ] Include action buttons slot

---

# Phase 7: Routing (`src/app/router/`)

## 7.1 Router Configuration

- [ ] Create router using createBrowserRouter
- [ ] Define route structure with nested routes
- [ ] Configure error boundaries per route level
- [ ] Configure loading states for lazy routes

## 7.2 Route Definitions

### 7.2.1 Public Routes
- [ ] Define login route (/login)
- [ ] Define register route (/register)
- [ ] Define forgot password route (/forgot-password)
- [ ] Define reset password route (/reset-password)
- [ ] Redirect authenticated users away from public routes

### 7.2.2 Protected Routes
- [ ] Define dashboard route (/)
- [ ] Define users routes (/users, /users/:id, /users/create)
- [ ] Define roles routes (/roles, /roles/:id, /roles/create)
- [ ] Define profile route (/profile)
- [ ] Define settings route (/settings)
- [ ] Wrap protected routes with AuthGuard

### 7.2.3 Error Routes
- [ ] Define 404 not found route
- [ ] Define 403 unauthorized route
- [ ] Define 500 error route
- [ ] Configure catch-all route for 404

## 7.3 Route Guards

- [ ] Implement AuthGuard component for protected routes
- [ ] Implement PermissionGuard component for permission-based access
- [ ] Implement RoleGuard component for role-based access
- [ ] Handle loading state during auth check
- [ ] Preserve intended URL for post-login redirect

## 7.4 Lazy Loading

- [ ] Implement React.lazy for feature pages
- [ ] Create Suspense wrapper with loading fallback
- [ ] Group related routes for code splitting
- [ ] Preload critical routes on hover (optional)

---

# Phase 8: Common Components (`src/components/common/`)

## 8.1 Data Display Components

### 8.1.1 Data Table Component
- [ ] Accept columns configuration
- [ ] Accept data array
- [ ] Support sorting by columns
- [ ] Support row selection
- [ ] Support pagination
- [ ] Display loading skeleton
- [ ] Display empty state
- [ ] Support responsive behavior (horizontal scroll or card view)

### 8.1.2 Pagination Component
- [ ] Display current page and total pages
- [ ] Include previous/next buttons
- [ ] Include page number buttons with truncation
- [ ] Include items per page selector
- [ ] Support controlled pagination

### 8.1.3 Status Badge Component
- [ ] Accept status value
- [ ] Map status to color/variant
- [ ] Support custom status mappings

### 8.1.4 Empty State Component
- [ ] Display icon and message
- [ ] Support action button
- [ ] Support different variants (no data, no results, error)

### 8.1.5 Loading Components
- [ ] Spinner component
- [ ] Skeleton component for content placeholders
- [ ] Full page loading overlay
- [ ] Button loading state

## 8.2 Feedback Components

### 8.2.1 Toast/Notification System
- [ ] Configure toast provider (sonner or similar)
- [ ] Create toast utility functions
  - [ ] success(message)
  - [ ] error(message)
  - [ ] warning(message)
  - [ ] info(message)
- [ ] Support action buttons in toasts
- [ ] Configure positioning and duration

### 8.2.2 Confirmation Dialog Component
- [ ] Accept title, message, confirm/cancel labels
- [ ] Support destructive variant (red confirm button)
- [ ] Accept onConfirm and onCancel callbacks
- [ ] Create useConfirm hook for imperative usage

### 8.2.3 Alert Component
- [ ] Support variants (info, success, warning, error)
- [ ] Include icon based on variant
- [ ] Support dismissible option

## 8.3 Form Components

### 8.3.1 Form Field Wrapper
- [ ] Accept label, description, error message
- [ ] Display required indicator
- [ ] Apply consistent spacing
- [ ] Support horizontal/vertical layout

### 8.3.2 Search Input Component
- [ ] Include search icon
- [ ] Support clear button
- [ ] Support debounced onChange
- [ ] Support loading state

### 8.3.3 Password Input Component
- [ ] Include show/hide toggle button
- [ ] Support password strength indicator
- [ ] Inherit standard input props

### 8.3.4 Multi-Select Component
- [ ] Display selected items as tags
- [ ] Support search/filter options
- [ ] Support select all option
- [ ] Handle large lists with virtualization (if needed)

## 8.4 Navigation Components

### 8.4.1 Breadcrumbs Component
- [ ] Accept items array with label and path
- [ ] Display separator between items
- [ ] Make last item non-clickable (current page)
- [ ] Support auto-generation from route (optional)

### 8.4.2 Tabs Component
- [ ] Support controlled and uncontrolled modes
- [ ] Support URL-synced tabs (query param)
- [ ] Include tab panels with lazy loading option

---

# Phase 9: Providers & App Setup (`src/app/`)

## 9.1 Provider Setup (`src/app/providers/`)

### 9.1.1 Query Provider
- [ ] Create QueryClientProvider wrapper
- [ ] Include ReactQueryDevtools (dev only)
- [ ] Configure devtools position and default open state

### 9.1.2 Theme Provider
- [ ] Implement dark/light mode toggle
- [ ] Persist theme preference
- [ ] Sync with system preference option
- [ ] Provide theme context for components

### 9.1.3 Toast Provider
- [ ] Configure toast provider from chosen library
- [ ] Set global toast configuration

### 9.1.4 Auth Provider (optional)
- [ ] Provide auth context if not using Zustand exclusively
- [ ] Initialize auth state on app mount
- [ ] Handle token refresh on mount

### 9.1.5 Combined Providers Component
- [ ] Create AppProviders component combining all providers
- [ ] Order providers correctly (query outside, theme, toast, etc.)
- [ ] Accept children prop

## 9.2 App Initialization

### 9.2.1 App Entry Point (main.tsx)
- [ ] Import global styles
- [ ] Render App component inside React.StrictMode
- [ ] Mount to root DOM element

### 9.2.2 App Component (App.tsx)
- [ ] Wrap with AppProviders
- [ ] Include RouterProvider
- [ ] Include global error boundary
- [ ] Include initial data loading (auth check)

### 9.2.3 Auth Initialization
- [ ] Check for stored tokens on app mount
- [ ] Validate tokens and fetch current user
- [ ] Handle expired/invalid tokens
- [ ] Set loading state during initialization
- [ ] Redirect based on auth state

---

# Phase 10: Error Handling & Boundaries

## 10.1 Error Boundary Components

### 10.1.1 Global Error Boundary
- [ ] Catch uncaught errors in React tree
- [ ] Display friendly error message
- [ ] Include "Reload" button
- [ ] Log errors to monitoring service (production)
- [ ] Show stack trace in development

### 10.1.2 Route Error Boundary
- [ ] Handle errors within specific routes
- [ ] Display route-specific error UI
- [ ] Include navigation back to safe route
- [ ] Handle 404 and other HTTP errors from loaders

### 10.1.3 Feature Error Boundary
- [ ] Wrap feature components
- [ ] Show feature-specific error UI
- [ ] Allow retry without full page reload

## 10.2 Error Handling Utilities

### 10.2.1 Error Logging Service
- [ ] Create error logging utility
- [ ] Integrate with monitoring service (Sentry, LogRocket)
- [ ] Include user context in error reports
- [ ] Include breadcrumbs/actions before error
- [ ] Filter sensitive information

### 10.2.2 API Error Handler
- [ ] Create central API error handling utility
- [ ] Map error codes to user messages
- [ ] Handle network errors
- [ ] Handle timeout errors
- [ ] Display appropriate toast/notification

## 10.3 Not Found & Unauthorized Pages

### 10.3.1 404 Not Found Page
- [ ] Display friendly message
- [ ] Include illustration/icon
- [ ] Include link to home/dashboard
- [ ] Include search option (optional)

### 10.3.2 403 Unauthorized Page
- [ ] Display permission denied message
- [ ] Explain what permission is needed
- [ ] Include link to request access (optional)
- [ ] Include link to home/dashboard

### 10.3.3 500 Error Page
- [ ] Display generic error message
- [ ] Include error reference ID (for support)
- [ ] Include retry button
- [ ] Include contact support link

---

# Phase 11: Testing

## 11.1 Testing Setup

### 11.1.1 Vitest Configuration
- [ ] Configure vitest.config.ts
- [ ] Set up test environment (jsdom)
- [ ] Configure coverage thresholds
- [ ] Set up path aliases for tests
- [ ] Configure test globals

### 11.1.2 Testing Library Setup
- [ ] Configure @testing-library/react
- [ ] Create custom render function with providers
- [ ] Set up user-event for interactions
- [ ] Configure screen queries

### 11.1.3 MSW Setup for API Mocking
- [ ] Install and configure msw
- [ ] Create mock handlers for auth endpoints
- [ ] Create mock handlers for user endpoints
- [ ] Create mock handlers for role endpoints
- [ ] Set up mock server for tests
- [ ] Configure request handlers reset between tests

## 11.2 Unit Tests

### 11.2.1 Utility Function Tests
- [ ] Test API client utilities
- [ ] Test form validation schemas
- [ ] Test helper functions
- [ ] Test formatters (date, currency, etc.)

### 11.2.2 Hook Tests
- [ ] Test custom hooks with renderHook
- [ ] Test auth hooks (useAuth, useLogin)
- [ ] Test form hooks
- [ ] Test query hooks with mock data

### 11.2.3 Store Tests
- [ ] Test Zustand store actions
- [ ] Test store selectors
- [ ] Test store persistence

## 11.3 Component Tests

### 11.3.1 UI Component Tests
- [ ] Test Button component variants
- [ ] Test Input component with validation
- [ ] Test Form components
- [ ] Test Dialog/Modal components
- [ ] Test Table component

### 11.3.2 Feature Component Tests
- [ ] Test LoginForm submission
- [ ] Test RegisterForm validation
- [ ] Test UserTable with mock data
- [ ] Test UserForm create and edit modes

### 11.3.3 Page Tests
- [ ] Test page rendering with mock data
- [ ] Test page navigation
- [ ] Test protected page access

## 11.4 Integration Tests

### 11.4.1 Auth Flow Tests
- [ ] Test complete login flow
- [ ] Test complete registration flow
- [ ] Test logout flow
- [ ] Test token refresh flow
- [ ] Test protected route access

### 11.4.2 CRUD Flow Tests
- [ ] Test user list with filters
- [ ] Test user creation
- [ ] Test user update
- [ ] Test user deletion
- [ ] Test role management

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
