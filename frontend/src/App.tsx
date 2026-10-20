import { lazy, Suspense, useMemo, type ReactNode } from 'react';
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
  useLocation,
} from 'react-router-dom';
import {
  ThemeProvider as MuiThemeProvider,
  CssBaseline,
  Box,
  CircularProgress,
  Fade,
} from '@mui/material';
import { AuthProvider } from './context/AuthContext';
import { ThemeProvider, useTheme } from './context/ThemeContext';
import { ProtectedRoute } from './components/ProtectedRoute';
import { ErrorBoundary } from './components/ErrorBoundary';
import { buildTheme } from './theme';

const GuestLibrary = lazy(() => import('./pages/GuestLibrary').then((m) => ({ default: m.GuestLibrary })));
const AlbumRoute = lazy(() => import('./pages/AlbumRoute').then((m) => ({ default: m.AlbumRoute })));
const Login = lazy(() => import('./pages/Login').then((m) => ({ default: m.Login })));
const Register = lazy(() => import('./pages/Register').then((m) => ({ default: m.Register })));
const Library = lazy(() => import('./pages/Library').then((m) => ({ default: m.Library })));

function RouteFallback() {
  return (
    <Box sx={{ minHeight: '100dvh', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
      <CircularProgress />
    </Box>
  );
}

function RouteFade({ children }: { children: ReactNode }) {
  const location = useLocation();
  return (
    <Fade key={location.pathname} in timeout={300}>
      <Box>{children}</Box>
    </Fade>
  );
}

function AppRoutes() {
  const location = useLocation();

  return (
    <RouteFade>
      <Suspense fallback={<RouteFallback />}>
        <Routes location={location}>
          {/* Home is the gallery itself — no sign-in prompt anywhere. */}
          <Route path="/" element={<GuestLibrary />} />
          <Route path="/books/:id" element={<AlbumRoute />} />

          {/* Admin login: only reachable by typing /admin. */}
          <Route path="/admin" element={<Login />} />
          <Route path="/login" element={<Navigate to="/" replace />} />
          <Route path="/register" element={<Register />} />

          {/* Admin dashboard */}
          <Route
            path="/library"
            element={
              <ProtectedRoute>
                <Library />
              </ProtectedRoute>
            }
          />

          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </Suspense>
    </RouteFade>
  );
}

function ThemedApp() {
  const { theme } = useTheme();
  const muiTheme = useMemo(() => buildTheme(theme), [theme]);

  return (
    <MuiThemeProvider theme={muiTheme}>
      <CssBaseline />
      <Router>
        <AuthProvider>
          <AppRoutes />
        </AuthProvider>
      </Router>
    </MuiThemeProvider>
  );
}

function App() {
  return (
    <ErrorBoundary>
      <ThemeProvider>
        <ThemedApp />
      </ThemeProvider>
    </ErrorBoundary>
  );
}

export default App;
