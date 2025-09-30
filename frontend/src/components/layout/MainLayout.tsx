import { Outlet } from 'react-router-dom';
import Navigation from '../navigation/Navigation';

const MainLayout: React.FC = () => {
  return (
    <div className="main-layout">
      <Navigation />
      <main className="main-content">
        <Outlet />
      </main>
    </div>
  );
};

export default MainLayout;