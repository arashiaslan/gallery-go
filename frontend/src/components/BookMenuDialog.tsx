import React, { useState } from 'react';
import {
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Divider,
  IconButton,
  TextField,
  Typography,
} from '@mui/material';
import { Close, Edit } from '@mui/icons-material';
import { useUpdateBook, useBookPhotos } from '../hooks/useBooks';
import { getErrorMessage } from '../lib/client';
import { booksApi } from '../lib/api';
import { toast } from 'sonner';
import { BOOK_COLORS } from '../lib/palette';
import type { Book } from '../types';

interface BookMenuDialogProps {
  book: Book | null;
  isOpen: boolean;
  onClose: () => void;
  onDelete: () => void;
  onUpload: () => void;
}

interface BookSettingsContentProps {
  book: Book;
  onClose: () => void;
  onDelete: () => void;
  onUpload: () => void;
}

const BookSettingsContent: React.FC<BookSettingsContentProps> = ({
  book,
  onClose,
  onDelete,
  onUpload,
}) => {
  const [editTitle, setEditTitle] = useState(book.title);
  const [selectedColor, setSelectedColor] = useState(book.bg_color || BOOK_COLORS[0]);
  const [isEditingTitle, setIsEditingTitle] = useState(false);
  const [isEditingColor, setIsEditingColor] = useState(false);
  const updateBookMutation = useUpdateBook();
  const { data: photosData } = useBookPhotos(book.id);
  const photos = photosData?.photos ?? [];

  const handleSetCover = async (photoId: string) => {
    try {
      await updateBookMutation.mutateAsync({
        bookId: book.id,
        data: { cover_photo_id: photoId },
      });
      toast.success('Cover updated');
    } catch (error) {
      toast.error(getErrorMessage(error, 'Failed to update cover'));
    }
  };

  const handleClearCover = async () => {
    try {
      await updateBookMutation.mutateAsync({
        bookId: book.id,
        data: { cover_photo_id: '' },
      });
      toast.success('Cover removed');
    } catch (error) {
      toast.error(getErrorMessage(error, 'Failed to update cover'));
    }
  };

  const handleSaveTitle = async () => {
    if (!editTitle.trim()) {
      toast.error('Book title cannot be empty');
      return;
    }

    try {
      await updateBookMutation.mutateAsync({
        bookId: book.id,
        data: { title: editTitle.trim() },
      });
      toast.success('Title updated');
      setIsEditingTitle(false);
    } catch (error) {
      toast.error(getErrorMessage(error, 'Failed to update title'));
    }
  };

  const handleSaveColor = async () => {
    try {
      await updateBookMutation.mutateAsync({
        bookId: book.id,
        data: { bg_color: selectedColor },
      });
      toast.success('Cover color updated');
      setIsEditingColor(false);
    } catch (error) {
      toast.error(getErrorMessage(error, 'Failed to update cover color'));
    }
  };

  return (
    <>
      <DialogTitle sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        Book settings
        <IconButton onClick={onClose} size="small" aria-label="Close">
          <Close fontSize="small" />
        </IconButton>
      </DialogTitle>

      <DialogContent sx={{ display: 'flex', flexDirection: 'column', gap: 2.5 }}>
        {/* Title */}
        <Box>
          <Typography variant="caption" color="text.secondary">
            Title
          </Typography>
          {isEditingTitle ? (
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1, mt: 0.5 }}>
              <TextField
                size="small"
                fullWidth
                value={editTitle}
                onChange={(e) => setEditTitle(e.target.value)}
                disabled={updateBookMutation.isPending}
                autoFocus
              />
              <Box sx={{ display: 'flex', gap: 1 }}>
                <Button
                  size="small"
                  onClick={() => {
                    setEditTitle(book.title);
                    setIsEditingTitle(false);
                  }}
                  disabled={updateBookMutation.isPending}
                >
                  Cancel
                </Button>
                <Button
                  size="small"
                  variant="contained"
                  disableElevation
                  onClick={handleSaveTitle}
                  disabled={updateBookMutation.isPending}
                >
                  Save
                </Button>
              </Box>
            </Box>
          ) : (
            <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mt: 0.5 }}>
              <Typography variant="body2" noWrap sx={{ flex: 1, mr: 1 }}>
                {book.title}
              </Typography>
              <Button size="small" startIcon={<Edit fontSize="small" />} onClick={() => setIsEditingTitle(true)}>
                Edit
              </Button>
            </Box>
          )}
        </Box>

        {/* Cover color */}
        <Box>
          <Typography variant="caption" color="text.secondary">
            Cover color
          </Typography>
          {isEditingColor ? (
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1, mt: 0.5 }}>
              <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, alignItems: 'center' }}>
                {BOOK_COLORS.map((color) => (
                  <Box
                    key={color}
                    component="button"
                    type="button"
                    onClick={() => setSelectedColor(color)}
                    disabled={updateBookMutation.isPending}
                    aria-label={`Color ${color}`}
                    sx={{
                      width: 32,
                      height: 32,
                      borderRadius: '8px',
                      border: '2px solid',
                      borderColor: selectedColor === color ? 'primary.main' : 'transparent',
                      bgcolor: color,
                      cursor: 'pointer',
                      p: 0,
                      '&:hover': { opacity: 0.85 },
                    }}
                  />
                ))}
                <Box
                  component="input"
                  type="color"
                  value={selectedColor}
                  onChange={(e) => setSelectedColor((e.target as HTMLInputElement).value)}
                  disabled={updateBookMutation.isPending}
                  aria-label="Custom color"
                  sx={{
                    width: 40,
                    height: 32,
                    border: 'none',
                    bgcolor: 'transparent',
                    cursor: 'pointer',
                    p: 0,
                  }}
                />
              </Box>
              <Box sx={{ display: 'flex', gap: 1 }}>
                <Button
                  size="small"
                  onClick={() => {
                    setSelectedColor(book.bg_color || BOOK_COLORS[0]);
                    setIsEditingColor(false);
                  }}
                  disabled={updateBookMutation.isPending}
                >
                  Cancel
                </Button>
                <Button
                  size="small"
                  variant="contained"
                  disableElevation
                  onClick={handleSaveColor}
                  disabled={updateBookMutation.isPending}
                >
                  Save
                </Button>
              </Box>
            </Box>
          ) : (
            <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mt: 0.5 }}>
              <Box
                sx={{
                  width: 32,
                  height: 32,
                  borderRadius: '8px',
                  bgcolor: book.bg_color || BOOK_COLORS[0],
                }}
              />
              <Button size="small" startIcon={<Edit fontSize="small" />} onClick={() => setIsEditingColor(true)}>
                Edit
              </Button>
            </Box>
          )}
        </Box>

        {/* Cover photo */}
        <Box>
          <Typography variant="caption" color="text.secondary">
            Cover photo
          </Typography>
          {photos.length === 0 ? (
            <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
              No photos yet
            </Typography>
          ) : (
            <>
              <Box sx={{ display: 'flex', gap: 1, overflowX: 'auto', mt: 0.5, pb: 0.5 }}>
                {photos.map((photo) => (
                  <Box
                    key={photo.id}
                    component="button"
                    type="button"
                    onClick={() => handleSetCover(photo.id)}
                    disabled={updateBookMutation.isPending}
                    aria-label="Set as cover photo"
                    sx={{
                      flexShrink: 0,
                      width: 56,
                      height: 56,
                      p: 0,
                      border: '2px solid',
                      borderColor: book.cover_photo_id === photo.id ? 'primary.main' : 'divider',
                      borderRadius: 1,
                      cursor: 'pointer',
                      overflow: 'hidden',
                      bgcolor: 'transparent',
                    }}
                  >
                    <Box
                      component="img"
                      src={booksApi.getPhotoUrl(book.id, photo.id, 'thumb')}
                      alt=""
                      sx={{ width: '100%', height: '100%', objectFit: 'cover', display: 'block' }}
                    />
                  </Box>
                ))}
              </Box>
              {book.cover_photo_id && (
                <Button
                  size="small"
                  sx={{ mt: 0.5 }}
                  onClick={handleClearCover}
                  disabled={updateBookMutation.isPending}
                >
                  Remove cover
                </Button>
              )}
            </>
          )}
        </Box>

        <Divider />

        <Button
          fullWidth
          onClick={() => {
            onClose();
            onUpload();
          }}
        >
          Upload photos
        </Button>

        <Button
          fullWidth
          color="error"
          onClick={() => {
            onClose();
            onDelete();
          }}
        >
          Delete book
        </Button>
      </DialogContent>

      <DialogActions>
        <Button onClick={onClose}>Close</Button>
      </DialogActions>
    </>
  );
};

export const BookMenuDialog: React.FC<BookMenuDialogProps> = ({
  book,
  isOpen,
  onClose,
  onDelete,
  onUpload,
}) => {
  return (
    <Dialog open={isOpen} onClose={onClose} maxWidth="xs" fullWidth>
      {book && (
        <BookSettingsContent
          key={book.id}
          book={book}
          onClose={onClose}
          onDelete={onDelete}
          onUpload={onUpload}
        />
      )}
    </Dialog>
  );
};
