# Documentation Index - Gallery Project

Complete guide to all documentation files. Start here to navigate the project.

---

## 🎯 Start Here

### For Project Overview
1. **[README.md](README.md)** - Quick project overview
2. **[PROJECT_STATUS.md](PROJECT_STATUS.md)** - Current status & progress
3. **[MILESTONES_SUMMARY.md](MILESTONES_SUMMARY.md)** - All 3 milestones overview

### For Setup & Running
1. **[SETUP.md](SETUP.md)** - Development environment setup
2. **[QUICK_START_M3.md](QUICK_START_M3.md)** - Fast frontend setup (M3)

### For Understanding Features
1. **[SPEC.md](SPEC.md)** - Full project specification & design
2. **[API.md](API.md)** - API endpoints reference

---

## 📋 By Milestone

### Milestone 1: Backend Skeleton
**Status**: ✅ COMPLETE

| Document | Purpose |
|----------|---------|
| [M1.md](M1.md) | M1 requirements |
| [M1_COMPLETION.md](M1_COMPLETION.md) | Detailed M1 report |
| [MILESTONE_1_SUMMARY.md](MILESTONE_1_SUMMARY.md) | M1 overview |

**Contains**:
- Backend structure
- Database schema
- Authentication system
- CRUD endpoints
- Tests & verification

**Run**: `cd backend && make dev`

---

### Milestone 2: Photo Upload & Processing
**Status**: ✅ COMPLETE

| Document | Purpose |
|----------|---------|
| [M2.md](M2.md) | M2 requirements |
| [M2_COMPLETION.md](M2_COMPLETION.md) | Detailed M2 report |

**Contains**:
- Photo upload endpoint
- Image processing pipeline
- EXIF extraction
- Media serving with caching
- Worker pool implementation

**Run**: `cd backend && make dev`

---

### Milestone 3: Frontend React
**Status**: ✅ COMPLETE

| Document | Purpose |
|----------|---------|
| [M3.md](M3.md) | M3 requirements |
| [M3_COMPLETION.md](M3_COMPLETION.md) | Detailed M3 report |
| [MILESTONE_3_SUMMARY.md](MILESTONE_3_SUMMARY.md) | M3 overview |
| [QUICK_START_M3.md](QUICK_START_M3.md) | Quick setup guide |
| [VISUAL_GUIDE_M3.md](VISUAL_GUIDE_M3.md) | UI specifications |

**Contains**:
- React + Vite setup
- BookCover component
- Coverflow carousel
- Library page
- Dummy data
- Accessibility features

**Run**: `cd frontend && npm run dev`

---

### Milestone 4: Page Flip (Future)
| Document | Purpose |
|----------|---------|
| [M4.md](M4.md) | M4 requirements |

**Will Include**:
- Page flip animation
- API integration
- Photo grids
- State management

---

### Milestone 5: Lightbox (Future)
| Document | Purpose |
|----------|---------|
| [M5.md](M5.md) | M5 requirements |

---

### Milestone 6: Polish (Future)
| Document | Purpose |
|----------|---------|
| [M6.md](M6.md) | M6 requirements |

---

## 🏗️ Technical Documentation

### Backend
| Document | Focus |
|----------|-------|
| [API.md](API.md) | All HTTP endpoints with examples |
| [SETUP.md](SETUP.md) | Backend development setup |
| `backend/README.md` | Backend-specific guide |
| [M1_COMPLETION.md](M1_COMPLETION.md) | Auth & database details |
| [M2_COMPLETION.md](M2_COMPLETION.md) | Image processing details |

**Key Files**:
- `backend/cmd/server/main.go` - Entry point
- `backend/internal/auth/` - Authentication
- `backend/internal/handlers/` - HTTP endpoints
- `backend/internal/store/` - Database layer
- `backend/internal/media/` - Image processing

---

### Frontend
| Document | Focus |
|----------|-------|
| [QUICK_START_M3.md](QUICK_START_M3.md) | Frontend setup |
| [VISUAL_GUIDE_M3.md](VISUAL_GUIDE_M3.md) | UI design specs |
| `frontend/README.md` | Frontend-specific guide |
| [M3_COMPLETION.md](M3_COMPLETION.md) | Component details |

**Key Files**:
- `frontend/src/App.tsx` - Main app
- `frontend/src/pages/Library.tsx` - Library page
- `frontend/src/components/CoverflowCarousel.tsx` - Carousel
- `frontend/src/components/BookCover.tsx` - Book design
- `frontend/src/lib/dummy-books.ts` - Sample data

---

### Database
| Document | Focus |
|----------|-------|
| [M1_COMPLETION.md](M1_COMPLETION.md) | Schema details (see "Database Schema") |

**Schema**:
- `users` - Authentication
- `books` - Albums
- `photos` - Images
- `sessions` - Login state
- `invites` - Signup links

---

## 📊 Navigation Guide

### "I want to..."

#### ...understand the full project
→ Start with [SPEC.md](SPEC.md), then [PROJECT_STATUS.md](PROJECT_STATUS.md)

#### ...set up development environment
→ Read [SETUP.md](SETUP.md)

#### ...run the backend
→ Go to [backend](backend/) folder, read [M1_COMPLETION.md](M1_COMPLETION.md)

#### ...run the frontend
→ Go to [frontend](frontend/) folder, read [QUICK_START_M3.md](QUICK_START_M3.md)

#### ...understand API endpoints
→ Read [API.md](API.md)

#### ...see code details for M1
→ Read [M1_COMPLETION.md](M1_COMPLETION.md)

#### ...see code details for M2
→ Read [M2_COMPLETION.md](M2_COMPLETION.md)

#### ...see code details for M3
→ Read [M3_COMPLETION.md](M3_COMPLETION.md)

#### ...deploy the project
→ Check [SETUP.md](SETUP.md) Docker section

#### ...understand security
→ Check [M1_COMPLETION.md](M1_COMPLETION.md) "Security Features"

#### ...learn about image processing
→ Read [M2_COMPLETION.md](M2_COMPLETION.md)

#### ...understand the carousel
→ Read [VISUAL_GUIDE_M3.md](VISUAL_GUIDE_M3.md)

---

## 🎯 Quick Reference

### Commands

**Backend**
```bash
cd backend
make dev                # Start server
make test               # Run tests
make build              # Build binary
make clean              # Clean
```

**Frontend**
```bash
cd frontend
npm run dev             # Start dev server
npm run build           # Production build
npm run preview         # Preview build
```

**Full Stack**
```bash
docker-compose up       # Start everything
```

### Key Endpoints

**Auth**
- `POST /api/auth/login` - Login
- `POST /api/auth/register?token=...` - Register
- `POST /api/auth/logout` - Logout

**Books**
- `GET /api/books` - List
- `POST /api/books` - Create
- `GET /api/books/:id` - Get
- `PATCH /api/books/:id` - Update
- `DELETE /api/books/:id` - Delete

**Photos**
- `POST /api/books/:id/photos` - Upload
- `GET /api/books/:id/pages` - Get all
- `PATCH /api/photos/:id` - Update caption
- `DELETE /api/photos/:id` - Delete

**Media**
- `GET /media/:book/:photo/:size` - Download (thumb/medium/original)

---

## 📈 Progress Tracking

| Milestone | Status | Completion | Documentation |
|-----------|--------|------------|----------------|
| M1 | ✅ Done | 100% | [M1_COMPLETION.md](M1_COMPLETION.md) |
| M2 | ✅ Done | 100% | [M2_COMPLETION.md](M2_COMPLETION.md) |
| M3 | ✅ Done | 100% | [M3_COMPLETION.md](M3_COMPLETION.md) |
| M4 | ⏳ Plan | 0% | [M4.md](M4.md) |
| M5 | ⏳ Plan | 0% | [M5.md](M5.md) |
| M6 | ⏳ Plan | 0% | [M6.md](M6.md) |

---

## 🔍 File Organization

### Root Documentation
```
gallery/
├── README.md                       # Overview
├── SPEC.md                         # Full specification
├── API.md                          # Endpoint reference
├── SETUP.md                        # Dev environment
├── PROJECT_STATUS.md               # Current status
├── MILESTONES_SUMMARY.md           # All milestones
├── DOCUMENTATION_INDEX.md          # This file
│
├── M1.md                           # M1 requirements
├── M1_COMPLETION.md                # M1 details
├── MILESTONE_1_SUMMARY.md          # M1 summary
│
├── M2.md                           # M2 requirements
├── M2_COMPLETION.md                # M2 details
│
├── M3.md                           # M3 requirements
├── M3_COMPLETION.md                # M3 details
├── MILESTONE_3_SUMMARY.md          # M3 summary
├── QUICK_START_M3.md               # M3 quick start
├── VISUAL_GUIDE_M3.md              # M3 design
│
├── M4.md                           # M4 requirements
├── M5.md                           # M5 requirements
├── M6.md                           # M6 requirements
│
└── .gitignore
```

### Backend Documentation
```
backend/
├── README.md                       # Backend overview
├── Makefile
├── go.mod / go.sum
├── .env.example
│
└── internal/README files (if needed)
```

### Frontend Documentation
```
frontend/
├── README.md                       # Frontend overview
├── package.json
├── vite.config.ts
└── tsconfig.json
```

---

## 🎓 Learning Path

### New to the Project?
1. Read [README.md](README.md) (5 min)
2. Read [SPEC.md](SPEC.md) (20 min)
3. Read [PROJECT_STATUS.md](PROJECT_STATUS.md) (10 min)
4. Pick a component to explore

### Want to Understand Backend?
1. Read [SETUP.md](SETUP.md) - Development setup
2. Read [M1_COMPLETION.md](M1_COMPLETION.md) - Backend details
3. Run `cd backend && make dev`
4. Test endpoints with curl (see [API.md](API.md))

### Want to Understand Frontend?
1. Read [QUICK_START_M3.md](QUICK_START_M3.md) - Quick start
2. Read [M3_COMPLETION.md](M3_COMPLETION.md) - Component details
3. Read [VISUAL_GUIDE_M3.md](VISUAL_GUIDE_M3.md) - Design specs
4. Run `cd frontend && npm run dev`

### Want to Deploy?
1. Read [SETUP.md](SETUP.md) - Docker section
2. Review [PROJECT_STATUS.md](PROJECT_STATUS.md) - Deployment section
3. Follow Dockerfile structure

---

## 📚 External Resources

### Backend (Go)
- [Go Official Docs](https://golang.org/doc/)
- [chi Router Docs](https://github.com/go-chi/chi)
- [SQLite](https://www.sqlite.org/)
- [Goose Migrations](https://github.com/pressly/goose)

### Frontend (React)
- [React Documentation](https://react.dev)
- [Vite Guide](https://vitejs.dev)
- [TypeScript Docs](https://www.typescriptlang.org/)
- [Tailwind CSS](https://tailwindcss.com)

### Deployment
- [Docker Docs](https://docs.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

---

## ✅ Verification Checklist

- [ ] Read README.md
- [ ] Understand SPEC.md
- [ ] Run backend: `cd backend && make dev`
- [ ] Run tests: `cd backend && make test`
- [ ] Run frontend: `cd frontend && npm run dev`
- [ ] Test API endpoints (see [API.md](API.md))
- [ ] Explore codebase
- [ ] Review security (see [M1_COMPLETION.md](M1_COMPLETION.md))

---

## 🚀 Quick Links

| Need | Link |
|------|------|
| Quick Overview | [README.md](README.md) |
| Full Details | [SPEC.md](SPEC.md) |
| Current Status | [PROJECT_STATUS.md](PROJECT_STATUS.md) |
| Backend Setup | [SETUP.md](SETUP.md) |
| Frontend Setup | [QUICK_START_M3.md](QUICK_START_M3.md) |
| API Reference | [API.md](API.md) |
| M1 Details | [M1_COMPLETION.md](M1_COMPLETION.md) |
| M2 Details | [M2_COMPLETION.md](M2_COMPLETION.md) |
| M3 Details | [M3_COMPLETION.md](M3_COMPLETION.md) |
| All Milestones | [MILESTONES_SUMMARY.md](MILESTONES_SUMMARY.md) |

---

## 📞 Need Help?

1. **Setup issues?** → [SETUP.md](SETUP.md)
2. **API questions?** → [API.md](API.md)
3. **Understand code?** → [M*_COMPLETION.md](M1_COMPLETION.md)
4. **Design questions?** → [VISUAL_GUIDE_M3.md](VISUAL_GUIDE_M3.md)
5. **Full spec?** → [SPEC.md](SPEC.md)

---

**Last Updated**: 2024-10-20  
**Status**: All documentation complete  
**Version**: M1-M3 stable
