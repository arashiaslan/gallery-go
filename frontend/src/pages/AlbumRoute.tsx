import React from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Box, CircularProgress } from '@mui/material';
import { AlbumView } from '../components/AlbumView';
import { useAuth } from '../context/AuthContext';

// Single album route: owner (logged in) gets upload/delete, guests get a
// read-only view. No login prompt is ever shown here.
export const AlbumRoute: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return (
      <Box sx={{ height: '100dvh', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <CircularProgress size={28} />
      </Box>
    );
  }

  if (!id) return null;

  return (
    <AlbumView
      bookId={id}
      guest={!isAuthenticated}
      onBack={() => navigate(isAuthenticated ? '/library' : '/')}
    />
  );
};
