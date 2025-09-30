import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './providers/AuthProvider';
import ProtectedRoute from './components/routes/ProtectedRoute';
import MainLayout from './components/layout/MainLayout';
import Login from './pages/Login';
import AuthCallback from './pages/AuthCallback';
import Dashboard from './pages/Dashboard';
import TodoListPage from './pages/TodoList';
import Profile from './pages/Profile';
import NotFound from './pages/NotFound';
import { ROUTES } from './routes/routeConfig';
import './App.css';

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          {/* Root redirect */}
          <Route path={ROUTES.HOME} element={<Navigate to={ROUTES.LOGIN} replace />} />

          {/* Public routes */}
          <Route path={ROUTES.LOGIN} element={<Login />} />
          <Route path={ROUTES.AUTH_CALLBACK} element={<AuthCallback />} />

          {/* Protected routes with layout */}
          <Route element={<ProtectedRoute />}>
            <Route element={<MainLayout />}>
              <Route path={ROUTES.DASHBOARD} element={<Dashboard />} />
              <Route path={ROUTES.TASKS} element={<TodoListPage />} />
              <Route path={ROUTES.PROFILE} element={<Profile />} />
            </Route>
          </Route>

          {/* 404 catch-all */}
          <Route path={ROUTES.NOT_FOUND} element={<NotFound />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;