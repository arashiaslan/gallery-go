import React, { type ReactNode } from 'react';
import { Box, Card, CardActionArea, Skeleton, Typography } from '@mui/material';
import { BookCoverImage } from './BookCoverImage';
import type { Book } from '../types';

interface AlbumCardProps {
  book: Book;
  onClick: () => void;
  /** Extra controls rendered on top of the card (e.g. settings menu). */
  action?: ReactNode;
  /** Position in the grid, used for the staggered entrance. */
  index?: number;
}

export const AlbumCard: React.FC<AlbumCardProps> = ({ book, onClick, action, index = 0 }) => (
  <Card
    sx={{
      position: 'relative',
      overflow: 'hidden',
      animation: 'fadeUp 0.45s both',
      animationDelay: `${Math.min(index, 12) * 45}ms`,
      transition: 'box-shadow 0.25s ease, transform 0.25s ease',
      '&:hover': { boxShadow: 6, transform: 'translateY(-3px)' },
    }}
  >
    <CardActionArea
      onClick={onClick}
      sx={{
        position: 'relative',
        display: 'block',
        aspectRatio: '4 / 3',
        overflow: 'hidden',
        '&:hover .album-cover': { transform: 'scale(1.05)' },
      }}
    >
      <Box
        className="album-cover"
        sx={{ position: 'absolute', inset: 0, transition: 'transform 0.4s ease' }}
      >
        <BookCoverImage book={book} />
      </Box>
      <Typography
        variant="subtitle2"
        sx={{
          position: 'absolute',
          left: 0,
          right: 0,
          bottom: 0,
          px: 1.75,
          py: 1.25,
          color: 'common.white',
          fontWeight: 700,
          whiteSpace: 'nowrap',
          overflow: 'hidden',
          textOverflow: 'ellipsis',
          background: 'linear-gradient(to top, rgba(0,0,0,0.72) 0%, rgba(0,0,0,0.28) 55%, transparent 100%)',
          zIndex: 1,
          pointerEvents: 'none',
        }}
      >
        {book.title}
      </Typography>
    </CardActionArea>
    {action}
  </Card>
);

export const AlbumCardSkeleton: React.FC = () => (
  <Card variant="outlined" sx={{ overflow: 'hidden' }}>
    <Skeleton variant="rectangular" sx={{ aspectRatio: '4 / 3', height: 'auto' }} />
    <Typography variant="subtitle2" sx={{ px: 1.75, py: 1.5 }}>
      <Skeleton width="60%" />
    </Typography>
  </Card>
);
