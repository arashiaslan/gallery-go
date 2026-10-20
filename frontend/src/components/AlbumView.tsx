import React, { useState, useEffect, useCallback, useRef } from 'react';
import {
  AppBar,
  Box,
  Button,
  IconButton,
  Skeleton,
  Toolbar,
  Tooltip,
  Typography,
} from '@mui/material';
import {
  ArrowBack,
  ChevronLeft,
  ChevronRight,
  Delete,
  DarkMode,
  LightMode,
  Pause,
  PhotoLibrary,
  PlayArrow,
} from '@mui/icons-material';
import { toast } from 'sonner';
import { useBookDetail, useBookPhotos, useDeletePhoto } from '../hooks/useBooks';
import { UploadPhotosDialog } from './UploadPhotosDialog';
import { DeleteConfirmDialog } from './DeleteConfirmDialog';
import { booksApi } from '../lib/api';
import { getErrorMessage } from '../lib/client';
import { useTheme } from '../context/ThemeContext';
import type { Photo } from '../types';

const AUTOPLAY_MS = 3000;

type ViewMode = 'grid' | 'slideshow';

interface AlbumViewProps {
  bookId: string;
  /** Guest view: read-only (no upload/delete). */
  guest?: boolean;
  onBack: () => void;
}

function formatTakenRange(photos: Photo[]): string | null {
  const dates = photos
    .map((p) => (p.taken_at ? p.taken_at.slice(0, 10) : null))
    .filter((d): d is string => Boolean(d))
    .sort();
  if (dates.length === 0) return null;

  const first = new Date(dates[0]);
  const last = new Date(dates[dates.length - 1]);
  const monthYear = (d: Date) => d.toLocaleDateString('en-US', { month: 'short', year: 'numeric' });
  const month = (d: Date) => d.toLocaleDateString('en-US', { month: 'short' });

  if (dates[0] === dates[dates.length - 1]) return monthYear(first);
  if (first.getFullYear() === last.getFullYear()) return `${month(first)} – ${monthYear(last)}`;
  return `${monthYear(first)} – ${monthYear(last)}`;
}

export const AlbumView: React.FC<AlbumViewProps> = ({ bookId, guest = false, onBack }) => {
  const [view, setView] = useState<ViewMode>('grid');
  const [index, setIndex] = useState(0);
  const [autoplay, setAutoplay] = useState(false);
  const [isUploadOpen, setIsUploadOpen] = useState(false);
  const [isConfirmDeleteOpen, setIsConfirmDeleteOpen] = useState(false);
  const [failedCoverId, setFailedCoverId] = useState<string | null>(null);
  const activeThumbRef = useRef<HTMLButtonElement | null>(null);
  const { theme, toggleTheme } = useTheme();

  const { data: book, isLoading: bookLoading } = useBookDetail(bookId, guest);
  const { data: photosData, isLoading: photosLoading } = useBookPhotos(bookId, guest);
  const deletePhotoMutation = useDeletePhoto();

  const photos = photosData?.photos ?? [];
  const total = photos.length;
  const safeIndex = total > 0 ? Math.min(index, total - 1) : 0;
  const current = total > 0 ? photos[safeIndex] : null;

  const currentCoverId = book?.cover_photo_id ?? null;
  const heroFailed = currentCoverId !== null && failedCoverId === currentCoverId;

  const handleBackOneLevel = useCallback(() => {
    if (view === 'slideshow') {
      setView('grid');
      setAutoplay(false);
    } else {
      onBack();
    }
  }, [view, onBack]);

  const openSlideshow = (photoIndex: number) => {
    setIndex(photoIndex);
    setView('slideshow');
  };

  const goNext = useCallback(() => {
    setIndex((i) => (total > 0 ? (i + 1) % total : 0));
  }, [total]);

  const goPrev = useCallback(() => {
    setIndex((i) => (total > 0 ? (i - 1 + total) % total : 0));
  }, [total]);

  // Autoplay slideshow
  useEffect(() => {
    if (!autoplay || total < 2) return;
    const timer = window.setInterval(goNext, AUTOPLAY_MS);
    return () => window.clearInterval(timer);
  }, [autoplay, total, goNext]);

  // Keep the active filmstrip thumb in view
  useEffect(() => {
    if (view === 'slideshow') {
      activeThumbRef.current?.scrollIntoView({ behavior: 'smooth', block: 'nearest', inline: 'center' });
    }
  }, [safeIndex, view]);

  // Keyboard: arrows navigate slideshow, Escape goes one level back
  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      if (view === 'slideshow') {
        if (e.key === 'ArrowRight') {
          e.preventDefault();
          goNext();
        } else if (e.key === 'ArrowLeft') {
          e.preventDefault();
          goPrev();
        } else if (e.key === 'Escape') {
          handleBackOneLevel();
        }
      } else if (e.key === 'Escape') {
        onBack();
      }
    };
    window.addEventListener('keydown', onKeyDown);
    return () => window.removeEventListener('keydown', onKeyDown);
  }, [view, goNext, goPrev, handleBackOneLevel, onBack]);

  const handleConfirmDelete = async () => {
    if (!current || guest) return;

    const wasLast = total === 1 || safeIndex === total - 1;
    try {
      await deletePhotoMutation.mutateAsync({ photoId: current.id, bookId });
      toast.success('Photo deleted');
      if (wasLast) {
        if (total === 1) {
          setView('grid');
        } else {
          setIndex(Math.max(0, safeIndex - 1));
        }
      }
      setIsConfirmDeleteOpen(false);
    } catch (error) {
      toast.error(getErrorMessage(error, 'Failed to delete photo'));
    }
  };

  const title = book?.title || 'Album';
  const isLoading = bookLoading || photosLoading;
  const coverUrl =
    book?.cover_photo_id ? booksApi.getPhotoUrl(bookId, book.cover_photo_id, 'medium') : null;
  const takenRange = formatTakenRange(photos);
  const metaLine = [total > 0 ? `${total} photo${total === 1 ? '' : 's'}` : null, takenRange]
    .filter(Boolean)
    .join(' · ');

  return (
    <Box
      sx={{
        minHeight: '100dvh',
        display: 'flex',
        flexDirection: 'column',
        height: view === 'slideshow' && !isLoading ? '100dvh' : 'auto',
      }}
    >
      <AppBar
        position="sticky"
        color="transparent"
        elevation={0}
        enableColorOnDark
        sx={{
          bgcolor: (t) =>
            t.palette.mode === 'dark' ? 'rgba(17,24,39,0.78)' : 'rgba(249,250,251,0.78)',
          backdropFilter: 'blur(12px)',
          borderBottom: 1,
          borderColor: 'divider',
          zIndex: (t) => t.zIndex.appBar,
        }}
      >
        <Toolbar variant="dense" sx={{ gap: 0.5 }}>
          <Tooltip title={view === 'slideshow' ? 'Back to grid' : 'Back'}>
            <IconButton edge="start" size="small" onClick={handleBackOneLevel} aria-label="Back">
              <ArrowBack fontSize="small" />
            </IconButton>
          </Tooltip>

          <Typography variant="subtitle2" noWrap sx={{ flex: 1, mx: 1, fontWeight: 600 }}>
            {title}
          </Typography>

          {view === 'slideshow' && (
            <Typography
              variant="body2"
              color="text.secondary"
              sx={{ mr: 0.5, minWidth: 48, textAlign: 'right' }}
            >
              {total > 0 ? `${safeIndex + 1} / ${total}` : '0'}
            </Typography>
          )}

          {view === 'grid' && (
            <Tooltip title={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}>
              <IconButton size="small" onClick={toggleTheme} aria-label="Toggle color theme">
                {theme === 'dark' ? <LightMode fontSize="small" /> : <DarkMode fontSize="small" />}
              </IconButton>
            </Tooltip>
          )}

          {!guest && (
            <Tooltip title="Upload photos">
              <IconButton size="small" onClick={() => setIsUploadOpen(true)} aria-label="Upload photos">
                <PhotoLibrary fontSize="small" />
              </IconButton>
            </Tooltip>
          )}

          {view === 'slideshow' && (
            <>
              <Tooltip title={autoplay ? 'Pause slideshow' : 'Play slideshow'}>
                <span>
                  <IconButton
                    size="small"
                    onClick={() => setAutoplay((a) => !a)}
                    disabled={total < 2}
                    aria-label={autoplay ? 'Pause slideshow' : 'Play slideshow'}
                  >
                    {autoplay ? <Pause fontSize="small" /> : <PlayArrow fontSize="small" />}
                  </IconButton>
                </span>
              </Tooltip>

              {!guest && (
                <Tooltip title="Delete photo">
                  <span>
                    <IconButton
                      size="small"
                      color="error"
                      disabled={!current || deletePhotoMutation.isPending}
                      onClick={() => setIsConfirmDeleteOpen(true)}
                      aria-label="Delete photo"
                    >
                      <Delete fontSize="small" />
                    </IconButton>
                  </span>
                </Tooltip>
              )}
            </>
          )}
        </Toolbar>
      </AppBar>

      {isLoading ? (
        <Box sx={{ p: { xs: 2, sm: 2.5, md: 3 }, maxWidth: 1440, mx: 'auto', width: '100%' }}>
          <Skeleton variant="rounded" sx={{ height: { xs: 220, sm: 280, md: 340 }, mb: 3 }} />
          <Box sx={{ columnCount: { xs: 2, sm: 3, md: 4, lg: 5 }, columnGap: 1.5 }}>
            {[160, 220, 140, 260, 180, 200, 150, 240, 170, 210].map((h, i) => (
              <Skeleton
                key={i}
                variant="rounded"
                sx={{ height: h, mb: 1.5, animationDelay: `${i * 80}ms` }}
              />
            ))}
          </Box>
        </Box>
      ) : view === 'grid' ? (
        <>
          {/* Hero: full-bleed cover with title + meta */}
          <Box
            sx={{
              position: 'relative',
              height: { xs: 240, sm: 300, md: 360 },
              bgcolor: book?.bg_color || 'grey.900',
              overflow: 'hidden',
              flexShrink: 0,
            }}
          >
            {coverUrl && !heroFailed && (
              <Box
                component="img"
                src={coverUrl}
                alt=""
                onError={() => setFailedCoverId(book?.cover_photo_id ?? null)}
                sx={{
                  position: 'absolute',
                  inset: 0,
                  width: '100%',
                  height: '100%',
                  objectFit: 'cover',
                  animation: 'fadeIn 0.6s ease both',
                }}
              />
            )}
            <Box
              sx={{
                position: 'absolute',
                inset: 0,
                background:
                  'linear-gradient(to top, rgba(0,0,0,0.78) 0%, rgba(0,0,0,0.38) 45%, rgba(0,0,0,0.16) 100%)',
              }}
            />
            <Box
              sx={{
                position: 'absolute',
                left: 0,
                right: 0,
                bottom: 0,
                px: { xs: 3, md: 4 },
                pb: { xs: 3, md: 4 },
                color: 'common.white',
                animation: 'fadeUp 0.55s 0.1s both',
              }}
            >
              <Typography
                variant="h4"
                sx={{
                  fontWeight: 800,
                  textShadow: '0 1px 14px rgba(0,0,0,0.4)',
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  whiteSpace: 'nowrap',
                }}
              >
                {title}
              </Typography>
              {metaLine && (
                <Typography
                  variant="body2"
                  sx={{ mt: 0.5, opacity: 0.9, fontWeight: 500, letterSpacing: 0.3 }}
                >
                  {metaLine}
                </Typography>
              )}
            </Box>
          </Box>

          {/* Masonry photo grid — native aspect ratios */}
          <Box sx={{ p: { xs: 2, sm: 2.5, md: 3 }, maxWidth: 1440, mx: 'auto', width: '100%' }}>
            {total === 0 ? (
              <Box
                sx={{
                  textAlign: 'center',
                  py: 12,
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
                  No photos in this album
                </Typography>
                {!guest && (
                  <Button
                    variant="contained"
                    disableElevation
                    sx={{ mt: 0.5 }}
                    onClick={() => setIsUploadOpen(true)}
                  >
                    Upload photos
                  </Button>
                )}
              </Box>
            ) : (
              <Box sx={{ columnCount: { xs: 2, sm: 3, md: 4, lg: 5 }, columnGap: 1.5 }}>
                {photos.map((photo, photoIndex) => (
                  <Box
                    key={photo.id}
                    role="button"
                    tabIndex={0}
                    aria-label={`Open photo ${photoIndex + 1}`}
                    onClick={() => openSlideshow(photoIndex)}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') openSlideshow(photoIndex);
                    }}
                    sx={{
                      mb: 1.5,
                      breakInside: 'avoid',
                      position: 'relative',
                      borderRadius: 1.5,
                      overflow: 'hidden',
                      cursor: 'pointer',
                      display: 'block',
                      animation: 'fadeUp 0.45s both',
                      animationDelay: `${Math.min(photoIndex, 15) * 40}ms`,
                      '&:hover img': { transform: 'scale(1.06)' },
                      '&:focus-visible': { outline: 2, outlineColor: 'primary.main' },
                    }}
                  >
                    <Box
                      component="img"
                      src={booksApi.getPhotoUrl(bookId, photo.id, 'thumb')}
                      alt={photo.caption || photo.filename || `Photo ${photoIndex + 1}`}
                      loading="lazy"
                      sx={{
                        width: '100%',
                        height: 'auto',
                        display: 'block',
                        aspectRatio:
                          photo.width > 0 && photo.height > 0
                            ? `${photo.width} / ${photo.height}`
                            : '1',
                        objectFit: 'cover',
                        transition: 'transform 0.4s ease',
                      }}
                    />
                  </Box>
                ))}
              </Box>
            )}
          </Box>
        </>
      ) : (
        <Box
          sx={{
            flex: 1,
            position: 'relative',
            display: 'flex',
            flexDirection: 'column',
            bgcolor: '#0f172a',
            minHeight: 0,
          }}
        >
          {/* Autoplay progress bar */}
          {autoplay && total > 1 && (
            <Box
              sx={{
                position: 'absolute',
                top: 0,
                left: 0,
                right: 0,
                height: 3,
                zIndex: 2,
                bgcolor: 'rgba(255,255,255,0.15)',
              }}
            >
              <Box
                key={safeIndex}
                sx={{
                  height: '100%',
                  bgcolor: 'primary.main',
                  animation: `progress ${AUTOPLAY_MS}ms linear forwards`,
                }}
              />
            </Box>
          )}

          {/* Current photo */}
          <Box
            sx={{
              flex: 1,
              position: 'relative',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              minHeight: 0,
              overflow: 'hidden',
              p: { xs: 1, sm: 2 },
            }}
          >
            {current && (
              <Box
                component="img"
                key={current.id}
                src={booksApi.getPhotoUrl(bookId, current.id, 'medium')}
                alt={current.caption || current.filename || 'Photo'}
                sx={{
                  maxWidth: '100%',
                  maxHeight: '100%',
                  objectFit: 'contain',
                  display: 'block',
                  userSelect: 'none',
                  animation: 'fadeIn 0.5s ease both, kenburns 9s ease-out both',
                }}
              />
            )}

            {total > 1 && (
              <>
                <IconButton
                  onClick={goPrev}
                  aria-label="Previous photo"
                  sx={{
                    position: 'absolute',
                    left: 12,
                    bgcolor: 'rgba(255,255,255,0.12)',
                    color: 'common.white',
                    backdropFilter: 'blur(8px)',
                    border: '1px solid rgba(255,255,255,0.2)',
                    '&:hover': { bgcolor: 'rgba(255,255,255,0.24)' },
                  }}
                >
                  <ChevronLeft />
                </IconButton>
                <IconButton
                  onClick={goNext}
                  aria-label="Next photo"
                  sx={{
                    position: 'absolute',
                    right: 12,
                    bgcolor: 'rgba(255,255,255,0.12)',
                    color: 'common.white',
                    backdropFilter: 'blur(8px)',
                    border: '1px solid rgba(255,255,255,0.2)',
                    '&:hover': { bgcolor: 'rgba(255,255,255,0.24)' },
                  }}
                >
                  <ChevronRight />
                </IconButton>
              </>
            )}
          </Box>

          {/* Filmstrip */}
          {total > 1 && (
            <Box
              sx={{
                display: 'flex',
                gap: 1,
                px: 1.5,
                py: 1.25,
                overflowX: 'auto',
                flexShrink: 0,
                borderTop: '1px solid rgba(255,255,255,0.1)',
                bgcolor: 'rgba(0,0,0,0.35)',
                backdropFilter: 'blur(8px)',
                scrollBehavior: 'smooth',
                '&::-webkit-scrollbar': { height: 4 },
                '&::-webkit-scrollbar-thumb': { bgcolor: 'rgba(255,255,255,0.25)', borderRadius: 2 },
              }}
            >
              {photos.map((photo, photoIndex) => (
                <Box
                  key={photo.id}
                  ref={photoIndex === safeIndex ? activeThumbRef : undefined}
                  component="button"
                  onClick={() => setIndex(photoIndex)}
                  aria-label={`Go to photo ${photoIndex + 1}`}
                  sx={{
                    flex: '0 0 auto',
                    width: 52,
                    height: 52,
                    p: 0,
                    border: 'none',
                    borderRadius: 1,
                    overflow: 'hidden',
                    cursor: 'pointer',
                    opacity: photoIndex === safeIndex ? 1 : 0.45,
                    outlineWidth: 2,
                    outlineStyle: 'solid',
                    outlineColor: photoIndex === safeIndex ? 'primary.main' : 'transparent',
                    outlineOffset: 1,
                    transition: 'opacity 0.2s ease',
                    '&:hover': { opacity: 1 },
                  }}
                >
                  <Box
                    component="img"
                    src={booksApi.getPhotoUrl(bookId, photo.id, 'thumb')}
                    alt=""
                    sx={{
                      width: '100%',
                      height: '100%',
                      objectFit: 'cover',
                      display: 'block',
                    }}
                  />
                </Box>
              ))}
            </Box>
          )}
        </Box>
      )}

      {!guest && (
        <>
          <UploadPhotosDialog
            bookId={bookId}
            isOpen={isUploadOpen}
            onClose={() => setIsUploadOpen(false)}
          />
          <DeleteConfirmDialog
            isOpen={isConfirmDeleteOpen}
            title="Delete photo"
            message="This photo will be permanently removed from the album."
            isLoading={deletePhotoMutation.isPending}
            onConfirm={handleConfirmDelete}
            onCancel={() => setIsConfirmDeleteOpen(false)}
          />
        </>
      )}
    </Box>
  );
};
