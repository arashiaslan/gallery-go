import React from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Alert,
  Box,
  Button,
  Chip,
  Container,
  Grid,
  IconButton,
  Toolbar,
  Tooltip,
  Typography,
} from '@mui/material';
import { DarkMode, LightMode, PhotoLibrary } from '@mui/icons-material';
import { useBooks } from '../hooks/useBooks';
import { useTheme } from '../context/ThemeContext';
import { AlbumCard, AlbumCardSkeleton } from '../components/AlbumCard';
import { BookCoverImage } from '../components/BookCoverImage';

const COLLAGE_POSITIONS = [
  { left: 0, top: 30, rotate: '-7deg' },
  { left: 95, top: 0, rotate: '4deg' },
  { left: 190, top: 60, rotate: '9deg' },
];

export const GuestLibrary: React.FC = () => {
  const navigate = useNavigate();
  const { theme, toggleTheme } = useTheme();
  const { data: books = [], isLoading, error, refetch } = useBooks(true);

  const collageBooks = books.slice(0, 3);

  return (
    <Box sx={{ minHeight: '100dvh', display: 'flex', flexDirection: 'column' }}>
      <Toolbar sx={{ justifyContent: 'space-between', position: 'relative', zIndex: 2 }}>
        <Typography variant="subtitle1" color="text.secondary" sx={{ fontWeight: 700, letterSpacing: 1 }}>
          GALLERY
        </Typography>
        <Tooltip title={theme === 'dark' ? 'Ganti ke mode terang' : 'Ganti ke mode gelap'}>
          <IconButton size="small" onClick={toggleTheme} aria-label="Ganti tema warna">
            {theme === 'dark' ? <LightMode fontSize="small" /> : <DarkMode fontSize="small" />}
          </IconButton>
        </Tooltip>
      </Toolbar>

      {/* Hero */}
      <Box sx={{ position: 'relative', overflow: 'hidden', pb: { xs: 5, md: 7 } }}>
        <Box aria-hidden sx={{ position: 'absolute', inset: 0, pointerEvents: 'none' }}>
          <Box
            sx={{
              position: 'absolute',
              width: 440,
              height: 440,
              borderRadius: '50%',
              top: -160,
              right: -80,
              bgcolor: 'primary.main',
              opacity: 0.16,
              filter: 'blur(80px)',
            }}
          />
          <Box
            sx={{
              position: 'absolute',
              width: 340,
              height: 340,
              borderRadius: '50%',
              bottom: -140,
              left: -70,
              bgcolor: theme === 'dark' ? '#f0abfc' : '#38bdf8',
              opacity: 0.14,
              filter: 'blur(80px)',
            }}
          />
          <Box
            sx={{
              position: 'absolute',
              width: 260,
              height: 260,
              borderRadius: '50%',
              top: '30%',
              left: '42%',
              bgcolor: theme === 'dark' ? '#93c5fd' : '#7c3aed',
              opacity: 0.1,
              filter: 'blur(70px)',
            }}
          />
        </Box>

        <Container
          maxWidth="lg"
          sx={{
            position: 'relative',
            display: 'flex',
            flexDirection: { xs: 'column', md: 'row' },
            alignItems: { xs: 'flex-start', md: 'center' },
            gap: { xs: 4, md: 6 },
            pt: { xs: 4, md: 8 },
          }}
        >
          <Box sx={{ flex: 1, minWidth: 0, animation: 'fadeUp 0.55s both' }}>
            <Chip
              label="KHUSUS KELUARGA"
              size="small"
              variant="outlined"
              sx={{
                mb: 2,
                fontWeight: 700,
                letterSpacing: 1.2,
                fontSize: 11,
                color: 'primary.main',
                borderColor: 'primary.main',
              }}
            />
            <Typography
              component="h1"
              variant="h2"
              sx={{
                fontWeight: 800,
                lineHeight: 1.1,
                letterSpacing: -1.5,
                color: 'text.primary',
              }}
            >
              Momen yang berharga.
            </Typography>
            <Typography
              variant="h6"
              color="text.secondary"
              sx={{ mt: 1.5, fontWeight: 400, maxWidth: 560 }}
            >
              Jelajahi semua album di bawah tanpa perlu akun. Buka salah satu untuk
              melihat foto layar penuh.
            </Typography>
            {!isLoading && books.length > 0 && (
              <Typography
                variant="body2"
                sx={{
                  mt: 2.5,
                  display: 'inline-flex',
                  alignItems: 'center',
                  gap: 1,
                  px: 1.5,
                  py: 0.75,
                  borderRadius: 2,
                  bgcolor: 'action.hover',
                  fontWeight: 600,
                  color: 'text.secondary',
                }}
              >
                <PhotoLibrary fontSize="small" />
                {books.length} album
              </Typography>
            )}
          </Box>

          {/* Floating cover collage */}
          {collageBooks.length > 0 && (
            <Box
              aria-hidden
              sx={{
                display: { xs: 'none', md: 'block' },
                position: 'relative',
                width: 370,
                height: 250,
                flexShrink: 0,
                animation: 'fadeIn 0.7s 0.25s both',
              }}
            >
              {collageBooks.map((book, i) => (
                <Box
                  key={book.id}
                  sx={{
                    position: 'absolute',
                    left: COLLAGE_POSITIONS[i].left,
                    top: COLLAGE_POSITIONS[i].top,
                    rotate: COLLAGE_POSITIONS[i].rotate,
                    width: 175,
                    p: 1,
                    pb: 0.75,
                    bgcolor: 'background.paper',
                    borderRadius: 2,
                    boxShadow: '0 12px 32px rgba(0,0,0,0.18)',
                    border: '1px solid',
                    borderColor: 'divider',
                    animation: 'float 5.5s ease-in-out infinite',
                    animationDelay: `${i * 0.8}s`,
                  }}
                >
                  <Box sx={{ borderRadius: 1, overflow: 'hidden' }}>
                    <BookCoverImage book={book} />
                  </Box>
                  <Typography
                    variant="caption"
                    noWrap
                    sx={{ display: 'block', px: 0.5, pt: 0.75, fontWeight: 700 }}
                  >
                    {book.title}
                  </Typography>
                </Box>
              ))}
            </Box>
          )}
        </Container>
      </Box>

      {/* Album */}
      <Container maxWidth="lg" sx={{ flex: 1, pb: 8 }}>
        {isLoading ? (
          <>
            <Box sx={{ display: 'flex', alignItems: 'baseline', gap: 2, mb: 3 }}>
              <Typography variant="h5" sx={{ fontWeight: 800 }}>
                Album
              </Typography>
              <Box sx={{ flex: 1, height: '1px', bgcolor: 'divider' }} />
            </Box>
            <Grid container spacing={2.5}>
              {Array.from({ length: 8 }).map((_, i) => (
                <Grid key={i} size={{ xs: 12, sm: 6, md: 4, lg: 3 }}>
                  <AlbumCardSkeleton />
                </Grid>
              ))}
            </Grid>
          </>
        ) : error ? (
          <Alert
            severity="error"
            action={
              <Button color="inherit" size="small" onClick={() => refetch()}>
                Coba lagi
              </Button>
            }
          >
            Gagal memuat album
          </Alert>
        ) : books.length === 0 ? (
          <Box
            sx={{
              textAlign: 'center',
              py: 14,
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              gap: 1.5,
            }}
          >
            <Box
              sx={{
                width: 64,
                height: 64,
                borderRadius: '50%',
                bgcolor: 'action.hover',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            >
              <PhotoLibrary sx={{ fontSize: 30, color: 'text.secondary' }} />
            </Box>
            <Typography variant="body1" color="text.secondary">
              Belum ada album
            </Typography>
          </Box>
        ) : (
          <>
            <Box
              sx={{
                display: 'flex',
                alignItems: 'baseline',
                gap: 2,
                mb: 3,
                animation: 'fadeUp 0.5s 0.1s both',
              }}
            >
              <Typography variant="h5" sx={{ fontWeight: 800 }}>
                Album
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {books.length} album
              </Typography>
              <Box sx={{ flex: 1, height: '1px', bgcolor: 'divider' }} />
            </Box>
            <Grid container spacing={2.5}>
              {books.map((book, i) => (
                <Grid key={book.id} size={{ xs: 12, sm: 6, md: 4, lg: 3 }}>
                  <AlbumCard
                    book={book}
                    index={i}
                    onClick={() => navigate(`/books/${book.id}`)}
                  />
                </Grid>
              ))}
            </Grid>
          </>
        )}
      </Container>
    </Box>
  );
};
