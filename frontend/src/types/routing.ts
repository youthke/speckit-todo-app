/**
 * Routing Type Definitions
 *
 * This file contains all TypeScript type definitions related to routing,
 * navigation, and route configuration in the React Router implementation.
 */

import { NavigateOptions } from 'react-router-dom';
import { ROUTES } from '../routes/routeConfig';

/**
 * Route keys from ROUTES configuration
 */
export type RouteKey = keyof typeof ROUTES;

/**
 * Route path values from ROUTES configuration
 */
export type RoutePath = typeof ROUTES[RouteKey];

/**
 * Navigation metadata for route items
 */
export interface NavItem {
  /** Display label for navigation link */
  label: string;
  /** Route key from ROUTES config */
  to: RouteKey;
  /** Icon name or component (optional) */
  icon?: string;
  /** Whether this nav item requires authentication */
  protected?: boolean;
}

/**
 * Route metadata configuration
 */
export interface RouteMeta {
  /** Human-readable title for the route */
  title: string;
  /** Whether this route requires authentication */
  requiresAuth: boolean;
  /** Parent route key (for breadcrumb trails) */
  parent?: RouteKey;
  /** Description for SEO or documentation */
  description?: string;
}

/**
 * Complete route configuration with metadata
 */
export interface RouteConfig {
  /** Route key from ROUTES */
  key: RouteKey;
  /** Route path pattern */
  path: RoutePath;
  /** Route metadata */
  meta: RouteMeta;
  /** Child routes (for nested routing) */
  children?: RouteConfig[];
}

/**
 * Typed navigation function signature
 */
export interface TypedNavigate {
  /** Navigate to a route by key */
  navigateTo: (key: RouteKey, options?: NavigateOptions) => void;
  /** Navigate to a path by string */
  navigateToPath: (path: string, options?: NavigateOptions) => void;
  /** Go back in history */
  goBack: () => void;
  /** Go forward in history */
  goForward: () => void;
}

/**
 * Protected route props
 */
export interface ProtectedRouteProps {
  /** Redirect path when user is not authenticated */
  redirectTo?: string;
}

/**
 * Navigation component props
 */
export interface NavigationProps {
  /** Optional class name for styling */
  className?: string;
  /** Whether to show navigation in mobile view */
  mobileVisible?: boolean;
}

/**
 * Breadcrumb item for navigation trails
 */
export interface BreadcrumbItem {
  /** Display label */
  label: string;
  /** Route key to navigate to */
  to?: RouteKey;
  /** Whether this is the current/active item */
  active: boolean;
}

/**
 * Route guard result
 */
export interface RouteGuardResult {
  /** Whether navigation is allowed */
  allowed: boolean;
  /** Redirect path if navigation is not allowed */
  redirectTo?: string;
  /** Reason for denial (for logging) */
  reason?: string;
}

/**
 * Navigation state for OAuth redirects
 */
export interface NavigationState {
  /** Original intended destination */
  from?: string;
  /** Timestamp of navigation attempt */
  timestamp?: number;
}
