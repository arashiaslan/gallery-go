import { createTheme } from '@mui/material/styles';
import type { Theme } from '@mui/material/styles';

export type ColorMode = 'light' | 'dark';

export const buildTheme = (mode: ColorMode): Theme =>
  createTheme({
    palette: {
      mode,
      primary: {
        main: mode === 'dark' ? '#a5a3d4' : '#5c5a96',
      },
      background: {
        default: mode === 'dark' ? '#111827' : '#f9fafb',
        paper: mode === 'dark' ? '#1f2937' : '#ffffff',
      },
      divider: mode === 'dark' ? 'rgba(255,255,255,0.12)' : 'rgba(15,23,42,0.10)',
    },
    shape: {
      borderRadius: 10,
    },
    typography: {
      fontFamily: '"Plus Jakarta Sans", Roboto, Helvetica, Arial, sans-serif',
      button: {
        textTransform: 'none',
        fontWeight: 600,
      },
    },
    components: {
      MuiButton: {
        defaultProps: {
          disableElevation: true,
        },
      },
      MuiCard: {
        defaultProps: {
          variant: 'outlined',
        },
      },
      MuiCssBaseline: {
        styleOverrides: {
          '@keyframes fadeUp': {
            from: { opacity: 0, transform: 'translateY(14px)' },
            to: { opacity: 1, transform: 'none' },
          },
          '@keyframes fadeIn': {
            from: { opacity: 0 },
            to: { opacity: 1 },
          },
          '@keyframes kenburns': {
            from: { transform: 'scale(1)' },
            to: { transform: 'scale(1.045)' },
          },
          '@keyframes progress': {
            from: { width: '0%' },
            to: { width: '100%' },
          },
          '@keyframes float': {
            '0%, 100%': { transform: 'translateY(0)' },
            '50%': { transform: 'translateY(-10px)' },
          },
          '@media (prefers-reduced-motion: reduce)': {
            '*, *::before, *::after': {
              animationDuration: '0.01ms !important',
              animationIterationCount: '1 !important',
              transitionDuration: '0.01ms !important',
            },
          },
        },
      },
    },
  });
