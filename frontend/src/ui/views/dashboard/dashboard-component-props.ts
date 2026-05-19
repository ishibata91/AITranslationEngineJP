import type { ShellRouteContract, ShellRouteId } from "@ui/stores/shell-state"

export interface AppHeaderProps {
  currentRoute: ShellRouteContract
  isMobileNavOpen: boolean
  routes: ShellRouteContract[]
  selectRoute: (routeId: ShellRouteId) => void
  toggleMobileNav: () => void
}

export interface GlobalNavigationProps {
  currentRoute: ShellRouteContract
  routes: ShellRouteContract[]
  selectRoute: (routeId: ShellRouteId) => void
}

export interface CurrentPageHeroProps {
  currentRoute: ShellRouteContract
  dataTestId?: string
}

export interface DashboardEntryGridProps {
  routes: ShellRouteContract[]
  selectRoute: (routeId: ShellRouteId) => void
}

export interface DashboardEntryCardProps {
  route: ShellRouteContract
  selectRoute: (routeId: ShellRouteId) => void
}
