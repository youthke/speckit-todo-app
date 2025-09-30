import { useNavigate, NavigateOptions } from 'react-router-dom';
import { ROUTES, RouteKey } from '../routes/routeConfig';

export const useTypedNavigate = () => {
  const navigate = useNavigate();

  return {
    navigateTo: (key: RouteKey, options?: NavigateOptions) => {
      navigate(ROUTES[key], options);
    },
    navigateToPath: (path: string, options?: NavigateOptions) => {
      navigate(path, options);
    },
    goBack: () => navigate(-1),
    goForward: () => navigate(1),
  };
};