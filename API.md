# Gallery API Reference (M1)

Base URL: `http://localhost:8080/api`

All endpoints return JSON. Authentication uses HttpOnly session cookies.

## Guest (public, read-only — no session required)

- `GET /api/guest/books` — all albums (same shape as `GET /api/books`)
- `GET /api/guest/books/{id}` — album detail
- `GET /api/guest/books/{id}/pages` — photo list
- `GET /media/{book_id}/{photo_id}/{size}` — image files (`thumb` | `medium` | `original`), now public

All mutations (`POST`/`PATCH`/`DELETE`) still require a session (401 otherwise).

## Authentication

### POST /auth/login
Login dengan email & password.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "data": {
    "id": "abc123",
    "email": "user@example.com",
    "role": "member"
  }
}
```

**Sets cookie:** `session_token=...` (HttpOnly, expires 30 days)

**Errors:**
- 400: Missing email/password
- 401: Invalid credentials
- 429: Too many attempts (max 10 per IP per 15 minutes)

---

### POST /auth/register?token=INVITE_TOKEN
Register dengan invite token.

**Request:**
```json
{
  "email": "newuser@example.com",
  "password": "password123"
}
```

**Response (201):**
```json
{
  "data": {
    "id": "xyz789",
    "email": "newuser@example.com",
    "role": "member"
  }
}
```

**Sets cookie:** `session_token=...` (same as login)

**Errors:**
- 400: Missing token/email/password, invite already used, invite expired
- 404: Invalid token
- 409: Email already registered (would be 400 in current implementation)
- 429: Too many attempts (max 10 per IP per 15 minutes)

---

### POST /auth/logout
Logout dan hapus session.

**Request:** (empty body)

**Response (200):**
```json
{
  "data": {
    "message": "Logged out"
  }
}
```

**Clears cookie:** `session_token`

---

### GET /me
Get current user info. Requires authentication.

**Request:** (empty body)

**Response (200):**
```json
{
  "data": {
    "id": "abc123",
    "email": "user@example.com",
    "role": "member"
  }
}
```

**Errors:**
- 401: Not authenticated

---

## Invites (Admin Only)

### POST /invites
Create signup invite link. **Admin role required.**

**Request:**
```json
{
  "email": "invitee@example.com"
}
```

**Response (201):**
```json
{
  "data": {
    "id": "inv123",
    "email": "invitee@example.com",
    "token": "abc123def456...",
    "expires_at": "2024-10-27T12:00:00Z"
  }
}
```

Share the token with the invitee: `POST /auth/register?token={token}`

**Errors:**
- 400: Missing email
- 403: Not admin

---

## Books

All endpoints require authentication.

### GET /books
List user's books.

**Response (200):**
```json
{
  "data": [
    {
      "id": "book1",
      "title": "Summer 2024",
      "owner_id": "abc123",
      "cover_photo_id": null,
      "bg_color": "#8E8DB8",
      "created_at": "2024-10-20T10:00:00",
      "updated_at": "2024-10-20T10:00:00"
    }
  ]
}
```

Returns empty array if no books.

---

### POST /books
Create new book.

**Request:**
```json
{
  "title": "My Album",
  "bg_color": "#8E8DB8"
}
```

`bg_color` is optional (defaults to `#8E8DB8`).

**Response (201):**
```json
{
  "data": {
    "id": "book_abc123",
    "title": "My Album",
    "owner_id": "user_xyz",
    "cover_photo_id": null,
    "bg_color": "#8E8DB8",
    "created_at": "2024-10-20T10:00:00",
    "updated_at": "2024-10-20T10:00:00"
  }
}
```

**Errors:**
- 400: Missing title
- 401: Not authenticated

---

### GET /books/:id
Get book details.

**Response (200):**
```json
{
  "data": {
    "id": "book1",
    "title": "Summer 2024",
    "owner_id": "abc123",
    "cover_photo_id": null,
    "bg_color": "#8E8DB8",
    "created_at": "2024-10-20T10:00:00",
    "updated_at": "2024-10-20T10:00:00"
  }
}
```

**Errors:**
- 401: Not authenticated
- 403: Not owner (unless admin)
- 404: Book not found

---

### PATCH /books/:id
Update book metadata.

**Request:** (any combination of these fields)
```json
{
  "title": "Updated Title",
  "bg_color": "#FF5733",
  "cover_photo_id": "photo123"
}
```

**Response (200):** Updated book object

**Errors:**
- 401: Not authenticated
- 403: Not owner (unless admin)
- 404: Book not found

---

### DELETE /books/:id
Delete book and all its photos.

**Response (200):**
```json
{
  "data": {
    "message": "Book deleted"
  }
}
```

**Errors:**
- 401: Not authenticated
- 403: Not owner (unless admin)
- 404: Book not found

---

## Health Check

### GET /health
Simple health check.

**Response (200):**
```json
{
  "status": "ok"
}
```

---

## Error Responses

All errors follow this format:

```json
{
  "error": "Error description"
}
```

Common HTTP status codes:
- **200**: Success
- **201**: Created
- **400**: Bad request (validation error)
- **401**: Unauthorized (not authenticated)
- **403**: Forbidden (authenticated but not allowed)
- **404**: Not found
- **500**: Internal server error

---

## cURL Examples

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{"email":"admin@example.com","password":"admin123"}'
```

### Create Book
```bash
curl -X POST http://localhost:8080/api/books \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"title":"My Album","bg_color":"#8E8DB8"}'
```

### List Books
```bash
curl http://localhost:8080/api/books -b cookies.txt
```

### Get Book
```bash
curl http://localhost:8080/api/books/book1 -b cookies.txt
```

### Update Book
```bash
curl -X PATCH http://localhost:8080/api/books/book1 \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"title":"Updated Title"}'
```

### Delete Book
```bash
curl -X DELETE http://localhost:8080/api/books/book1 -b cookies.txt
```

### Get Current User
```bash
curl http://localhost:8080/api/me -b cookies.txt
```

### Logout
```bash
curl -X POST http://localhost:8080/api/auth/logout -b cookies.txt
```

### Create Invite (Admin)
```bash
curl -X POST http://localhost:8080/api/invites \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"email":"newuser@example.com"}'
```

### Register with Invite
```bash
curl -X POST "http://localhost:8080/api/auth/register?token=INVITE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"newuser@example.com","password":"password123"}'
```

---

## Notes

- Sessions expire after a fixed 30 days (sliding renewal is not implemented)
- Login and register are rate limited to 10 attempts per IP per 15 minutes (HTTP 429)
- Photo upload (not documented here yet): max 20 files per request, each up to `MAX_UPLOAD_MB` (default 25MB), request body capped at `(MAX_UPLOAD_MB + 1MB) × 20`
- Invite links expire after 7 days
- One invite can only be used once
- Admin can view/edit any book; members can only manage their own
- All timestamps are in ISO 8601 format
