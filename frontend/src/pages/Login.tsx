import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Button,
  CircularProgress,
  Container,
  Link,
  TextField,
  Typography,
  Alert,
} from '@mui/material';
import { useAuth } from '../context/AuthContext';
import { getErrorMessage } from '../lib/client';

export const Login: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();
  const { isAuthenticated, login } = useAuth();

  useEffect(() => {
    if (isAuthenticated) {
      navigate('/library');
    }
  }, [isAuthenticated, navigate]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      await login(email, password);
      navigate('/library');
    } catch (err) {
      setError(getErrorMessage(err, 'Invalid email or password'));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Container maxWidth="xs" sx={{ minHeight: '100dvh', display: 'flex', alignItems: 'center' }}>
      <Box sx={{ width: '100%', py: 6 }}>
        <Typography variant="overline" color="text.secondary">
          Gallery
        </Typography>
        <Typography variant="h5" sx={{ fontWeight: 600, mt: 0.5 }}>
          Sign in
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
          Open your photo albums
        </Typography>

        <Box component="form" onSubmit={handleSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          {error && <Alert severity="error">{error}</Alert>}

          <TextField
            id="email"
            label="Email"
            type="email"
            required
            fullWidth
            size="small"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            disabled={isLoading}
            autoFocus
          />
          <TextField
            id="password"
            label="Password"
            type="password"
            required
            fullWidth
            size="small"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            disabled={isLoading}
          />

          <Button type="submit" variant="contained" fullWidth disabled={isLoading} sx={{ mt: 1 }}>
            {isLoading ? <CircularProgress size={20} color="inherit" /> : 'Sign in'}
          </Button>
        </Box>

        <Typography variant="body2" color="text.secondary" sx={{ mt: 3, textAlign: 'center' }}>
          <Link component="button" variant="body2" onClick={() => navigate('/')} underline="hover">
            Back to gallery
          </Link>
        </Typography>
      </Box>
    </Container>
  );
};
