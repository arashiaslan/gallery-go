import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Alert,
  Box,
  Button,
  Container,
  Fab,
  Grid,
  IconButton,
  Toolbar,
  Tooltip,
  Typography,
} from '@mui/material';
import { Add, MoreVert } from '@mui/icons-material';
import { UploadPhotosDialog } from '../components/UploadPhotosDialog';
import { CreateBookDialog } from '../components/CreateBookDialog';
import { BookMenuDialog } from '../components/BookMenuDialog';
import { DeleteConfirmDialog } from '../components/DeleteConfirmDialog';
import { ProfileMenu } from '../components/ProfileMenu';
import { AlbumCard, AlbumCardSkeleton } from '../components/AlbumCard';
import { useBooks, useDeleteBook } from '../hooks/useBooks';
import { useAuth } from '../context/AuthContext';
import { getErrorMessage } from '../lib/client';
import { toast } from 'sonner';
import type { Book } from '../types';

export const Library: React.FC = () => {
  const navigate = useNavigate();
  const { logout } = useAuth();
  const { data: books = [], isLoading, error, refetch } = useBooks();
  const deleteBookMutation = useDeleteBook();

  const [selectedBookId, setSelectedBookId] = useState<string | null>(null);
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);
  const [isUploadOpen, setIsUploadOpen] = useState(false);
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState<Book | null>(null);

  // Keep the settings dialog in sync with the query cache (title, cover, color).
  const selectedBook = books.find((b) => b.id === selectedBookId) ?? null;

  const handleLogout = async () => {
    try {
      await logout();
      toast.success('Signed out');
      navigate('/');
    } catch {
      toast.error('Sign out failed');
    }
  };

  const handleDeleteBook = async () => {
    if (!deleteTarget) return;

    try {
      await deleteBookMutation.mutateAsync(deleteTarget.id);
      toast.success('Book deleted');
      setDeleteTarget(null);
    } catch (error) {
      toast.error(getErrorMessage(error, 'Failed to delete book'));
    }
  };

  return (
    <Box sx={{ minHeight: '100dvh', display: 'flex', flexDirection: 'column' }}>
      <Toolbar sx={{ justifyContent: 'space-between' }}>
        <Typography variant="subtitle1" color="text.secondary" sx={{ fontWeight: 700, letterSpacing: 1 }}>
          GALLERY
        </Typography>
        <ProfileMenu onLogout={handleLogout} />
      </Toolbar>

      <Container maxWidth="lg" sx={{ flex: 1, pb: 10 }}>
        {isLoading ? (
          <Grid container spacing={2.5}>
            {Array.from({ length: 8 }).map((_, i) => (
              <Grid key={i} size={{ xs: 12, sm: 6, md: 4, lg: 3 }}>
                <AlbumCardSkeleton />
              </Grid>
            ))}
          </Grid>
        ) : error ? (
          <Alert
            severity="error"
            action={
              <Button color="inherit" size="small" onClick={() => refetch()}>
                Try again
              </Button>
            }
          >
            Failed to load books
          </Alert>
        ) : books.length === 0 ? (
          <Box sx={{ textAlign: 'center', py: 10 }}>
            <Typography variant="h6" sx={{ fontWeight: 600 }}>
              No books yet
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5, mb: 3 }}>
              Create your first album to start adding photos
            </Typography>
            <Button variant="contained" disableElevation startIcon={<Add />} onClick={() => setIsCreateOpen(true)}>
              New book
            </Button>
          </Box>
        ) : (
          <Grid container spacing={2.5}>
            {books.map((book, i) => (
              <Grid key={book.id} size={{ xs: 12, sm: 6, md: 4, lg: 3 }}>
                <AlbumCard
                  book={book}
                  index={i}
                  onClick={() => navigate(`/books/${book.id}`)}
                  action={
                    <Tooltip title="Book settings">
                      <IconButton
                        size="small"
                        aria-label={`Settings for ${book.title}`}
                        onClick={(e) => {
                          e.stopPropagation();
                          setSelectedBookId(book.id);
                          setIsSettingsOpen(true);
                        }}
                        sx={{
                          position: 'absolute',
                          top: 8,
                          right: 8,
                          bgcolor: 'rgba(0,0,0,0.45)',
                          color: 'common.white',
                          backdropFilter: 'blur(6px)',
                          border: '1px solid rgba(255,255,255,0.25)',
                          '&:hover': { bgcolor: 'rgba(0,0,0,0.65)' },
                          zIndex: 2,
                        }}
                      >
                        <MoreVert fontSize="small" />
                      </IconButton>
                    </Tooltip>
                  }
                />
              </Grid>
            ))}
          </Grid>
        )}
      </Container>

      <Fab
        color="primary"
        aria-label="Create new book"
        onClick={() => setIsCreateOpen(true)}
        sx={{ position: 'fixed', bottom: 24, right: 24 }}
      >
        <Add />
      </Fab>

      {selectedBook && (
        <>
          <BookMenuDialog
            book={selectedBook}
            isOpen={isSettingsOpen}
            onClose={() => setIsSettingsOpen(false)}
            onDelete={() => {
              setIsSettingsOpen(false);
              setDeleteTarget(selectedBook);
            }}
            onUpload={() => setIsUploadOpen(true)}
          />
          <UploadPhotosDialog
            bookId={selectedBook.id}
            isOpen={isUploadOpen}
            onClose={() => setIsUploadOpen(false)}
          />
        </>
      )}

      <DeleteConfirmDialog
        isOpen={deleteTarget !== null}
        title="Delete book?"
        message={`"${deleteTarget?.title}" and all photos inside it will be permanently deleted.`}
        isLoading={deleteBookMutation.isPending}
        onConfirm={handleDeleteBook}
        onCancel={() => setDeleteTarget(null)}
      />

      <CreateBookDialog isOpen={isCreateOpen} onClose={() => setIsCreateOpen(false)} />
    </Box>
  );
};
