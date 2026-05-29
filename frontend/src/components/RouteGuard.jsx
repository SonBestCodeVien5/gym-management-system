import { useEffect } from 'react'
import { useAuth } from '../context/AuthContext.jsx'
import StatusMessage from './StatusMessage.jsx'

function RouteGuard({ mode, navigate, children }) {
  const { status } = useAuth()

  useEffect(() => {
    if (status === 'checking') {
      return
    }

    if (mode === 'protected' && status === 'anonymous') {
      navigate('/login', { replace: true })
    }

    if (mode === 'public' && status === 'authenticated') {
      navigate('/app', { replace: true })
    }
  }, [mode, navigate, status])

  if (status === 'checking') {
    return (
      <StatusMessage
        fullPage
        title="Dang kiem tra phien lam viec"
        message="He thong dang ket noi API de khoi phuc nhan vien hien tai."
      />
    )
  }

  if (mode === 'protected' && status !== 'authenticated') {
    return null
  }

  if (mode === 'public' && status === 'authenticated') {
    return null
  }

  return children
}

export default RouteGuard
