import React from 'react';
import { NavLink } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import { ROUTES } from '../../routes/routeConfig';
import './Navigation.css';

const Navigation: React.FC = () => {
  const { user, logout, isAuthenticated } = useAuth();

  if (!isAuthenticated) {
    return null;
  }

  const handleLogout = async () => {
    try {
      await logout();
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  return (
    <nav className="navigation">
      <div className="nav-brand">
        <h2>Todo App</h2>
      </div>

      <div className="nav-links">
        <NavLink
          to={ROUTES.DASHBOARD}
          className={({ isActive }) =>
            isActive ? 'nav-link active' : 'nav-link'
          }
        >
          Dashboard
        </NavLink>

        <NavLink
          to={ROUTES.TASKS}
          className={({ isActive }) =>
            isActive ? 'nav-link active' : 'nav-link'
          }
        >
          Tasks
        </NavLink>

        <NavLink
          to={ROUTES.PROFILE}
          className={({ isActive }) =>
            isActive ? 'nav-link active' : 'nav-link'
          }
        >
          Profile
        </NavLink>
      </div>

      <div className="nav-user">
        <span className="user-name">{user?.name}</span>
        <button onClick={handleLogout} className="logout-button">
          Logout
        </button>
      </div>
    </nav>
  );
};

export default Navigation;