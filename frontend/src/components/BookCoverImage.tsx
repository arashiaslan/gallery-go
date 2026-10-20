import React, { useState } from 'react';
import { Box } from '@mui/material';
import { Image as ImageIcon } from '@mui/icons-material';
import { booksApi } from '../lib/api';
import type { Book } from '../types';

// Cover image with fallback to the album's color when the thumbnail
// is missing or fails to load (e.g. the cover photo was deleted).
export const BookCoverImage: React.FC<{ book: Book }> = ({ book }) => {
  const [failed, setFailed] = useState(false);

  if (book.cover_photo_id && !failed) {
    return (
      <Box
        component="img"
        src={booksApi.getPhotoUrl(book.id, book.cover_photo_id, 'thumb')}
        alt={book.title}
        onError={() => setFailed(true)}
        sx={{ width: '100%', aspectRatio: '4 / 3', objectFit: 'cover', display: 'block' }}
      />
    );
  }

  return (
    <Box
      sx={{
        width: '100%',
        aspectRatio: '4 / 3',
        bgcolor: book.bg_color,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
      }}
    >
      <ImageIcon sx={{ fontSize: 36, color: 'rgba(255,255,255,0.7)' }} />
    </Box>
  );
};
