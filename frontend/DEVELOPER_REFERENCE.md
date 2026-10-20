# Frontend Developer Reference

## Quick Navigation

### Authentication
- `src/context/AuthContext.tsx` - Global auth state & `useAuth()` hook
- `src/pages/Login.tsx` - Login form page
- `src/pages/Register.tsx` - Register form page
- `src/lib/api.ts` - Auth API endpoints

### API & Data
- `src/lib/client.ts` - Axios instance setup
- `src/lib/api.ts` - All API endpoint wrappers
- `src/lib/query-client.ts` - TanStack Query config
- `src/hooks/useBooks.ts` - React Query hooks for books

### Components
- `src/components/ProtectedRoute.tsx` - Protect routes
- `src/components/ErrorBoundary.tsx` - Error handling
- `src/pages/Library.tsx` - Main library page (needs data integration)

### Types
- `src/types/auth.ts` - Auth types (User, LoginRequest, etc.)
- `src/types/index.ts` - App types (Book, etc.)

---

## Common Tasks

### Add Authentication to a Component

```typescript
import { useAuth } from '@/context/AuthContext';

export function MyComponent() {
  const { user, isAuthenticated, logout } = useAuth();
  
  if (!isAuthenticated) return <div>Not logged in</div>;
  
  return (
    <div>
      <p>Welcome, {user?.email}</p>
      <button onClick={logout}>Logout</button>
    </div>
  );
}
```

### Fetch Data with React Query

```typescript
import { useBooks } from '@/hooks/useBooks';

export function BooksList() {
  const { data: books, isLoading, error } = useBooks();
  
  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error loading books</div>;
  
  return (
    <ul>
      {books?.map(book => (
        <li key={book.id}>{book.title}</li>
      ))}
    </ul>
  );
}
```

### Create Book

```typescript
import { useCreateBook } from '@/hooks/useBooks';

export function CreateBookButton() {
  const createBook = useCreateBook();
  
  const handleCreate = async () => {
    try {
      await createBook.mutateAsync({ title: 'New Book' });
      // Success - books list auto-refetches
    } catch (error) {
      console.error('Failed to create:', error);
    }
  };
  
  return (
    <button 
      onClick={handleCreate}
      disabled={createBook.isPending}
    >
      {createBook.isPending ? 'Creating...' : 'Create Book'}
    </button>
  );
}
```

### Call API Directly

```typescript
import { authApi, booksApi } from '@/lib/api';

// Auth
const user = await authApi.login(email, password);
const user = await authApi.register(email, password, passwordConfirm, token);
const user = await authApi.getMe();
await authApi.logout();

// Books
const books = await booksApi.getBooks();
const book = await booksApi.createBook({ title: 'My Book' });
const book = await booksApi.updateBook(id, { title: 'Updated' });
await booksApi.deleteBook(id);
```

---

## API Endpoints

### Authentication Endpoints

```
POST /api/login
  Body: { email: string, password: string }
  Response: { user: User }

POST /api/register
  Body: { 
    email: string, 
    password: string, 
    passwordConfirm: string, 
    invitationToken: string 
  }
  Response: { user: User }

GET /api/me
  Response: { user: User }
  (Called on app startup to check session)

POST /api/logout
  Response: 200 OK
```

### Books Endpoints

```
GET /api/books
  Response: { books: Book[] }

POST /api/books
  Body: { title: string, coverImage?: string }
  Response: { book: Book }

PUT /api/books/:id
  Body: { title?: string, backgroundColor?: string, ... }
  Response: { book: Book }

DELETE /api/books/:id
  Response: 204 No Content
```

---

## Type Definitions

### User
```typescript
interface User {
  id: number;
  email: string;
  name?: string;
  role: 'admin' | 'user';
  createdAt: string;
}
```

### Book (from existing types)
```typescript
interface Book {
  id: number;
  title: string;
  pageCount: number;
  coverImage: string;
  backgroundColor: string;
  spine: { color: string };
}
```

### Auth Requests
```typescript
interface LoginRequest {
  email: string;
  password: string;
}

interface RegisterRequest {
  email: string;
  password: string;
  passwordConfirm: string;
  invitationToken: string;
}
```

---

## Hooks Reference

### useAuth()
```typescript
const {
  user,              // User | null
  isLoading,         // boolean
  isAuthenticated,   // boolean
  login,             // (email, password) => Promise<void>
  register,          // (email, password, passwordConfirm, token) => Promise<void>
  logout,            // () => Promise<void>
  fetchMe,           // () => Promise<void>
} = useAuth();
```

### useBooks()
```typescript
const {
  data: books,       // Book[] | undefined
  isLoading,
  error,
} = useBooks();
```

### useCreateBook()
```typescript
const {
  mutateAsync,       // (data) => Promise<Book>
  isPending,
  error,
} = useCreateBook();

// Usage:
await mutateAsync({ title: 'New Book' });
```

### useUpdateBook()
```typescript
const {
  mutateAsync,       // ({ bookId, data }) => Promise<Book>
  isPending,
  error,
} = useUpdateBook();

// Usage:
await mutateAsync({ bookId: 1, data: { title: 'Updated' } });
```

### useDeleteBook()
```typescript
const {
  mutateAsync,       // (bookId) => Promise<void>
  isPending,
  error,
} = useDeleteBook();

// Usage:
await mutateAsync(1);
```

---

## Routing

### Routes Available

```
/login              - Login page (public)
/register           - Register page (requires ?token=...)
/library            - Library page (protected)
/                   - Redirects to /library
```

### Protected Route

```typescript
<Route
  path="/library"
  element={
    <ProtectedRoute>
      <Library />
    </ProtectedRoute>
  }
/>
```

When not authenticated, redirects to `/login`.

---

## Error Handling

### Component Level
```typescript
try {
  await someAsyncAction();
} catch (error: any) {
  const message = error.response?.data?.message || 'Something went wrong';
  setError(message);
}
```

### React Query
```typescript
const { data, error } = useBooks();

if (error) {
  return <div>Error: {error.message}</div>;
}
```

### 401 Errors
Automatically handled:
- Axios interceptor detects 401
- AuthContext clears user
- App redirects to /login

---

## Development Tips

### 1. Check Authentication State
```typescript
// In browser DevTools React tab
// Find AuthProvider → see user/isAuthenticated state
```

### 2. Debug API Calls
```typescript
// Browser Network tab shows:
// POST /api/login
// GET /api/me
// etc.

// Browser Console shows:
// Axios interceptor logs (if added)
// Component errors
```

### 3. Vite Proxy
Requests to `/api/*` are proxied to `http://localhost:8080/api/*`
- No CORS issues
- Cookies sent automatically
- Same-origin in dev

### 4. Test Login Without Backend
Will fail with ECONNREFUSED
- Backend must be running on :8080
- Check if `npm run dev` proxy is working

### 5. Clear Session
```typescript
// Browser DevTools → Application → Cookies
// Delete cookies set by backend
// Then refresh page → should redirect to /login
```

---

## Common Errors

### "Cannot POST /api/login"
- ✅ Backend running on http://localhost:8080?
- ✅ Vite proxy configured in vite.config.ts?

### "useAuth must be used within AuthProvider"
- ✅ Make sure AuthProvider wraps your component (in App.tsx)

### "Cannot read property 'email' of null"
- ✅ Check if user is authenticated before accessing user.email
- ✅ Use isAuthenticated or user?.email optional chaining

### Cookies not being sent
- ✅ Backend must set cookies with proper flags
- ✅ Frontend has withCredentials: true in axios
- ✅ Dev: Proxy should work. Prod: Might need HTTPS + proper CORS

---

## Styling

### Tailwind CSS
- Already configured in tailwind.config.js
- Use Tailwind classes in components
- Examples:
  - `bg-gradient-to-br` - Gradient background
  - `text-white` - White text
  - `hover:bg-slate-100` - Hover effects
  - `disabled:opacity-50` - Disabled state

### Color Scheme
- Slate colors for backgrounds/muted tones
- White for primary text
- Red for errors
- Use `opacity-*` for transparency effects

---

## Performance

### TanStack Query Features
- ✅ Automatic caching (5 min stale time)
- ✅ Deduplication (same request twice = one API call)
- ✅ Background refetch on window focus
- ✅ Automatic retry on failure

### Code Splitting
- Consider using `React.lazy()` for large pages
- Routes can be lazy loaded for faster initial load

### Bundle Size
- Current: ~114 KB gzipped
- Monitor with: `npm run build`

---

## Next Tasks for M4 Part 2

1. [ ] Connect Library page to `useBooks()`
2. [ ] Handle empty state (no books)
3. [ ] Add skeleton loading for books carousel
4. [ ] Implement "Create Book" dialog
5. [ ] Implement book image upload
6. [ ] Add delete/edit book actions
7. [ ] Add user menu (profile, logout, manage members)

---

## Useful Resources

- [React Router v6 Docs](https://reactrouter.com/)
- [TanStack Query Docs](https://tanstack.com/query/latest)
- [Axios Docs](https://axios-http.com/)
- [Tailwind CSS Docs](https://tailwindcss.com/)
- [TypeScript React](https://react-typescript-cheatsheet.netlify.app/)
