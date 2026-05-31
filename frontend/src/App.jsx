import { useCallback, useEffect, useMemo, useState } from 'react'
import AppShell from './components/AppShell.jsx'
import DashboardHome from './components/DashboardHome.jsx'
import LoginView from './components/LoginView.jsx'
import ModulePlaceholder from './components/ModulePlaceholder.jsx'
import RouteGuard from './components/RouteGuard.jsx'
import StatusMessage from './components/StatusMessage.jsx'
import StateBlock from './components/StateBlock.jsx'
import { AuthProvider, useAuth } from './context/AuthContext.jsx'
import { canAccessRoute } from './lib/permissions.js'
import { matchRoute, normalizePath } from './routes/matchRoute.js'
import { APP_HOME_PATH, APP_NAV_ITEMS } from './routes/routeConfig.js'

function readRoute() {
  return normalizePath(window.location.pathname)
}

function AppRoutes() {
  const [route, setRoute] = useState(readRoute)
  const { employee, status } = useAuth()
  const matchedRoute = useMemo(() => matchRoute(route), [route])

  const navigate = useCallback((nextRoute, options = {}) => {
    const normalizedRoute = normalizePath(nextRoute)

    if (window.location.pathname !== normalizedRoute) {
      const method = options.replace ? 'replaceState' : 'pushState'
      window.history[method](null, '', normalizedRoute)
    }

    setRoute(normalizedRoute)
  }, [])

  useEffect(() => {
    const handlePopState = () => setRoute(readRoute())
    window.addEventListener('popstate', handlePopState)
    return () => window.removeEventListener('popstate', handlePopState)
  }, [])

  useEffect(() => {
    if (route !== '/' || status === 'checking') {
      return
    }

    navigate(status === 'authenticated' ? APP_HOME_PATH : '/login', { replace: true })
  }, [navigate, route, status])

  useEffect(() => {
    if (matchedRoute.type !== 'redirect') {
      return
    }

    navigate(matchedRoute.redirectTo, { replace: true })
  }, [matchedRoute, navigate])

  function renderAppRoute() {
    if (matchedRoute.type === 'app-not-found') {
      return (
        <StateBlock
          tone="notFound"
          title="Page not found"
          message="This workspace route is not registered yet."
        />
      )
    }

    if (!canAccessRoute(matchedRoute.route, employee?.role || [])) {
      return (
        <StateBlock
          tone="forbidden"
          title="Access denied"
          message="Your current role cannot open this workspace module."
        />
      )
    }

    if (matchedRoute.route.key === 'dashboard') {
      const dashboardNavItems = APP_NAV_ITEMS.map((item) => ({
        ...item,
        available: item.status === 'ready' && canAccessRoute(item, employee?.role || []),
        visible: canAccessRoute(item, employee?.role || []),
      }))

      return (
        <DashboardHome
          employee={employee}
          navItems={dashboardNavItems}
          activeItem={matchedRoute.route.navKey}
        />
      )
    }

    return <ModulePlaceholder route={matchedRoute.route} params={matchedRoute.params} />
  }

  if (matchedRoute.type === 'root' || matchedRoute.type === 'redirect') {
    return (
      <StatusMessage
        fullPage
        title="Dang kiem tra phien lam viec"
        message="He thong dang tai trang thai dang nhap hien tai."
      />
    )
  }

  if (matchedRoute.type === 'app' || matchedRoute.type === 'app-not-found') {
    return (
      <RouteGuard mode="protected" navigate={navigate}>
        <AppShell navigate={navigate} activeRoute={matchedRoute.route}>
          {renderAppRoute()}
        </AppShell>
      </RouteGuard>
    )
  }

  return (
    <RouteGuard mode="public" navigate={navigate}>
      <LoginView navigate={navigate} />
    </RouteGuard>
  )
}

function App() {
  return (
    <AuthProvider>
      <AppRoutes />
    </AuthProvider>
  )
}

export default App
