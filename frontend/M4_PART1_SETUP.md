# M4 Part 1: Setup TanStack Query, API Client, dan Authentication

## ✅ Selesai

### 1. **TanStack Query Setup**
- ✅ `src/lib/query-client.ts` - QueryClient dengan default options
- ✅ `src/lib/client.ts` - Axios instance dengan `withCredentials: true`
- ✅ `src/lib/api.ts` - API endpoints wrapper (`authApi`, `booksApi`)
- ✅ Error handling: Axios interceptors untuk 401 responses
- ✅ Integration di `src/main.tsx` dengan `QueryClientProvider`

### 2. **Vite Proxy (Development)**
- ✅ `vite.config.ts` - Proxy `/api` dan `/media` ke `http://localhost:8080`

### 3. **Halaman Login** (`src/pages/Login.tsx`)
- ✅ UI minimalis dengan gradient background (slate colors)
- ✅ Form: email + password
- ✅ Button: "Masuk", link "Daftar via undangan"
- ✅ Loading state dan error messages
- ✅ Bahasa Indonesia
- ✅ Auto redirect ke /library kalau sudah login
- ✅ Auto redirect ke /register?token=... kalau ada query param

### 4. **Halaman Register** (`src/pages/Register.tsx`)
- ✅ Ambil token dari URL query params (`?token=...`)
- ✅ Form: email + password + confirm password
- ✅ Validation: password minimal 6 char, password harus match
- ✅ Error handling
- ✅ Success redirect ke /library
- ✅ Link kembali ke login

### 5. **Authentication Context** (`src/context/AuthContext.tsx`)
- ✅ Global user state
- ✅ `useAuth()` hook
- ✅ Methods: `login()`, `register()`, `logout()`
- ✅ Auto fetch current user on mount (`/api/me`)
- ✅ Handle 401 responses

### 6. **Protected Route** (`src/components/ProtectedRoute.tsx`)
- ✅ Redirect ke /login kalau belum authenticated
- ✅ Loading state

### 7. **Routing** (`src/App.tsx`)
- ✅ React Router v6 setup
- ✅ Routes: `/login`, `/register`, `/library`
- ✅ Root `/` redirect ke `/library`

### 8. **Dependencies** 
- ✅ `@tanstack/react-query` - ^5.28.0
- ✅ `axios` - ^1.6.2
- ✅ `react-router-dom` - ^6.21.0

### 9. **Hooks** (`src/hooks/useBooks.ts`)
- ✅ `useBooks()` - GET /api/books
- ✅ `useCreateBook()` - POST /api/books
- ✅ `useUpdateBook()` - PUT /api/books/:id
- ✅ `useDeleteBook()` - DELETE /api/books/:id

## API Endpoints Expected

Backend harus menyediakan:

```
POST /api/login
  Request: { email: string, password: string }
  Response: { user: User }

POST /api/register
  Request: { email: string, password: string, passwordConfirm: string, invitationToken: string }
  Response: { user: User }

GET /api/me
  Response: { user: User }

POST /api/logout
  Response: 200 OK

GET /api/books
  Response: { books: Book[] }

POST /api/books
  Request: { title: string, coverImage?: string }
  Response: { book: Book }

PUT /api/books/:id
  Response: { book: Book }

DELETE /api/books/:id
  Response: 204 No Content
```

## User Type

```typescript
interface User {
  id: number;
  email: string;
  name?: string;
  role: 'admin' | 'user';
  createdAt: string;
}
```

## Development

1. **Start backend** di port 8080
2. **Start frontend** (dev server akan proxy `/api` dan `/media`):
   ```bash
   npm run dev
   ```
3. **Build**:
   ```bash
   npm run build
   ```

## Testing Flow

1. ✅ Open http://localhost:5173
2. ✅ Should redirect ke /login (not authenticated)
3. ✅ Login form works (test dengan email/password dari backend)
4. ✅ After login → redirect ke /library
5. ✅ /library protected (unauthorized → redirect ke /login)
6. ✅ Register page: /register?token=xxx (dari invite link)
7. ✅ Logout → redirect ke /login

## Notes

- `withCredentials: true` → cookies auto included in API requests
- 401 response auto handled by axios interceptor
- QueryClient auto invalidate queries after mutations
- AuthProvider perlu di-wrap di root App untuk semua routes bisa akses `useAuth()`
