import React, { useState } from 'react';
import {
  Box,
  Button,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
  Typography,
} from '@mui/material';
import { useCreateBook } from '../hooks/useBooks';
import { getErrorMessage } from '../lib/client';
import { toast } from 'sonner';
import { BOOK_COLORS } from '../lib/palette';

interface CreateBookDialogProps {
  isOpen: boolean;
  onClose: () => void;
}

export const CreateBookDialog: React.FC<CreateBookDialogProps> = ({ isOpen, onClose }) => {
  const [title, setTitle] = useState('');
  const [bgColor, setBgColor] = useState(BOOK_COLORS[0]);
  const createBookMutation = useCreateBook();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) {
      toast.error('Book title cannot be empty');
      return;
    }

    try {
      await createBookMutation.mutateAsync({
        title: title.trim(),
        bg_color: bgColor,
      });
      toast.success('Book created');
      setTitle('');
      setBgColor(BOOK_COLORS[0]);
      onClose();
    } catch (error) {
      toast.error(getErrorMessage(error, 'Failed to create book'));
    }
  };

  return (
    <Dialog open={isOpen} onClose={onClose} maxWidth="xs" fullWidth>
      <DialogTitle>New book</DialogTitle>
      <Box component="form" onSubmit={handleSubmit}>
        <DialogContent sx={{ display: 'flex', flexDirection: 'column', gap: 2.5 }}>
          <TextField
            id="book-title"
            label="Title"
            size="small"
            fullWidth
            required
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Summer Trip 2026"
            disabled={createBookMutation.isPending}
            autoFocus
          />

          <Box>
            <Typography variant="caption" color="text.secondary">
              Cover color
            </Typography>
            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mt: 1 }}>
              {BOOK_COLORS.map((color) => (
                <Box
                  key={color}
                  component="button"
                  type="button"
                  onClick={() => setBgColor(color)}
                  disabled={createBookMutation.isPending}
                  aria-label={`Color ${color}`}
                  sx={{
                    width: 32,
                    height: 32,
                    borderRadius: '8px',
                    border: '2px solid',
                    borderColor: bgColor === color ? 'primary.main' : 'transparent',
                    bgcolor: color,
                    cursor: 'pointer',
                    p: 0,
                    '&:hover': { opacity: 0.85 },
                  }}
                />
              ))}
            </Box>
          </Box>
        </DialogContent>

        <DialogActions>
          <Button onClick={onClose} disabled={createBookMutation.isPending}>
            Cancel
          </Button>
          <Button type="submit" variant="contained" disabled={createBookMutation.isPending} disableElevation>
            {createBookMutation.isPending ? <CircularProgress size={18} color="inherit" /> : 'Create'}
          </Button>
        </DialogActions>
      </Box>
    </Dialog>
  );
};
