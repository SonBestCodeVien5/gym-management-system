import { useCallback, useEffect, useState } from 'react'
import AppShell from './components/AppShell.jsx'
import LoginView from './components/LoginView.jsx'
import RouteGuard from './components/RouteGuard.jsx'
import StatusMessage from './components/StatusMessage.jsx'
import { AuthProvider, useAuth } from './context/AuthContext.jsx'

const KNOWN_ROUTES = new Set(['/login', '/app'])

function readRoute() {
  return KNOWN_ROUTES.has(window.location.pathname) ? window.location.pathname : '/'
}

function AppRoutes() {
  const [route, setRoute] = useState(readRoute)
  const { status } = useAuth()

  const navigate = useCallback((nextRoute, options = {}) => {
    const normalizedRoute = KNOWN_ROUTES.has(nextRoute) ? nextRoute : '/'

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

    navigate(status === 'authenticated' ? '/app' : '/login', { replace: true })
  }, [navigate, route, status])

  if (route === '/') {
    return (
      <StatusMessage
        fullPage
        title="Dang kiem tra phien lam viec"
        message="He thong dang tai trang thai dang nhap hien tai."
      />
    )
  }

  if (route === '/app') {
    return (
      <RouteGuard mode="protected" navigate={navigate}>
        <AppShell navigate={navigate} />
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
