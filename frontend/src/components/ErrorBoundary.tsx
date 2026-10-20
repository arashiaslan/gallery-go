import React from 'react';
import type { ReactNode } from 'react';
import { Box, Button, Container, Typography } from '@mui/material';

interface Props {
  children?: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

export class ErrorBoundary extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Error caught by boundary:', error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return (
        <Container maxWidth="sm" sx={{ minHeight: '100dvh', display: 'flex', alignItems: 'center' }}>
          <Box sx={{ width: '100%', textAlign: 'center', py: 6 }}>
            <Typography variant="h5" sx={{ fontWeight: 600, mb: 1 }}>
              Something went wrong
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
              {this.state.error?.message}
            </Typography>
            <Button variant="contained" disableElevation onClick={() => window.location.reload()}>
              Reload page
            </Button>
          </Box>
        </Container>
      );
    }

    return this.props.children;
  }
}
