# Gallery Frontend - Milestone 3

## Overview

React + Vite + TypeScript frontend for the Family Gallery application. Implements the Library page with Coverflow carousel and 3D book covers using CSS perspective transforms.

## Features

### ✅ Completed
- React 19 + Vite + TypeScript setup with strict mode
- Tailwind CSS 4 + PostCSS for styling
- **BookCover component**: 2:3 ratio, dark spine (12% width), rounded right corners (20px), settings button, layered shadow effects
- **CoverflowCarousel component**: 
  - Embla carousel integration with smooth scrolling
  - CSS 3D perspective transforms (container: 1200px)
  - Active book: scale 1, rotateY 0°
  - Neighboring books: scale ~0.8, rotateY ±30° (mobile) / ±70° (desktop)
  - 60fps transforms via requestAnimationFrame + direct style updates (not React state)
  - Input support: swipe, drag, keyboard (←/→), horizontal wheel
  - Smart snap to nearest book
  - Mobile: neighbors clipped 15-20% at edges
  - Desktop (≥768px): fan-like display with more books visible
- **Library page**:
  - Dynamic gradient background that transitions to active book's color
  - Header: left button (book icon), 3 right buttons (toggle page count, search, menu)
  - Title + optional page count display
  - 4 action buttons below carousel (menu, share, delete, add)
  - 2 glass-morphic buttons at bottom-right (fullscreen, sparkles/search)
- **Dummy data**: 8 books with picsum placeholder images and unique background gradients
- **Accessibility**: Respects prefers-reduced-motion (disables rotations and animations)
- **Responsive design**: Mobile-first, touch-friendly (44px+ tap targets)

## Tech Stack

- **React** 19.3.1 with TypeScript 7
- **Vite** 8.3.4 for fast dev/build
- **Tailwind CSS** 4.3.3
- **Embla Carousel** 8.6.0
- **Framer Motion** 14.4.0 (for future page animations)
- **Lucide React** 1.59.0 (icons)
- **Plus Jakarta Sans** font from Google Fonts

## Installation & Setup

```bash
cd frontend
npm install
npm run dev
```

Open [http://localhost:5173](http://localhost:5173) to view the Library page.

## Project Structure

```
src/
├── components/
│   ├── BookCover.tsx        # Individual book cover with spine & effects
│   └── CoverflowCarousel.tsx # Main carousel with 3D transforms
├── pages/
│   └── Library.tsx          # Main Library page layout
├── types/
│   └── index.ts            # TypeScript interfaces (Book, BookCarouselState)
├── lib/
│   ├── dummy-books.ts      # 8 sample books with images & colors
│   └── utils.ts            # Helper functions (cn)
├── App.tsx                  # Root component
├── index.css               # Global styles + CSS 3D rules
└── main.tsx                # React entry point
```

## Key Components

### BookCover
Renders a single book with:
- Dark spine on left (12% width)
- Rounded right corners (20px border-radius)
- Cover image (2:3 aspect ratio)
- Settings button (sliders icon) in top-right
- Layered shadows for depth (3D thick-page effect)

### CoverflowCarousel
3D carousel powered by embla-carousel:
- **Perspective**: 1200px container perspective
- **Active book**: center position, full size, full brightness
- **Neighboring books**: scaled down, rotated away, dimmed
- **Desktop (≥768px)**: fan-like display showing 5-7 slivers
- **Mobile (<768px)**: neighbors peek 15-20% from edges
- **Smooth 60fps**: updates via requestAnimationFrame + direct DOM style changes
- **Input modes**:
  - Touch swipe/drag
  - Mouse drag
  - Keyboard: ← → arrows
  - Horizontal wheel scroll
  - Snap on release

### Library Page
Main page layout:
- Gradient background (shifts with active book)
- Header with action buttons
- Centered book carousel
- Action buttons row
- Glass-morphic buttons for additional features

## Styling Notes

- **Dark theme**: Periwinkle (#8E8DB8) on mobile, Navy (#3F4A78) reference colors
- **Glass-morphism**: backdrop-blur + semi-transparent white borders
- **Shadows**: Multi-layer shadows for depth (book spine, cover, overall shadow)
- **Font**: Plus Jakarta Sans with weights 400/500/600/700

## Performance Optimizations

1. **60fps transforms**: CSS transforms + requestAnimationFrame bypass React re-renders
2. **Direct style updates**: Embla scroll progress → direct `element.style.transform` (no state)
3. **Will-change**: `will-change: transform` on carousel slides
4. **Lazy loading**: Placeholder images from picsum (can be replaced with optimized images)

## Accessibility

- **Keyboard navigation**: Arrow keys to switch books
- **Reduced motion**: Detects `prefers-reduced-motion` and disables rotations
- **Tap targets**: All buttons ≥44px for mobile
- **ARIA labels**: Meaningful labels for buttons
- **Semantic HTML**: Proper heading hierarchy

## Known Limitations / TODOs

- No backend integration yet (M4 will add page flip animation)
- Placeholder images from external picsum (can be replaced with local assets)
- Glass buttons on mobile stack vertically (desktop: horizontal)
- No page-flip animation (M4 milestone)

## Building for Production

```bash
npm run build      # Creates optimized bundle in dist/
npm run preview    # Preview production build locally
```

## Environment

No environment variables required for M3 (uses dummy data).

## Next Steps (M4+)

- Page flip animation with react-pageflip (StPageFlip)
- Photo grid layouts (1-6 images per page)
- Lightbox for photo viewing
- Backend API integration
- Authentication

---

**Status**: Milestone 3 ✅ Complete - Library page with Coverflow carousel ready for development!
