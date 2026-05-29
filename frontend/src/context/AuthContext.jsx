import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
} from 'react'
import {
  createClientError,
  currentEmployeeRequest,
  loginRequest,
  logoutRequest,
  refreshRequest,
} from '../lib/api.js'
import {
  clearStoredTokens,
  readStoredTokens,
  saveTokens,
} from '../lib/authStorage.js'

const AuthContext = createContext(null)

const INITIAL_STATE = {
  status: 'checking',
  employee: null,
  accessToken: null,
  refreshToken: null,
  error: null,
  notice: null,
}

const SESSION_EXPIRED_NOTICE = 'Phien dang nhap da het han. Vui long dang nhap lai.'

function requireTokenPair(data) {
  const accessToken = data?.access_token
  const refreshToken = data?.refresh_token

  if (!accessToken || !refreshToken) {
    throw createClientError(
      'INVALID_RESPONSE',
      'API response did not include a complete token pair.',
    )
  }

  return { accessToken, refreshToken }
}

function requireEmployee(data) {
  if (!data?.employee && !data?.id) {
    throw createClientError(
      'INVALID_RESPONSE',
      'API response did not include employee data.',
    )
  }

  return data.employee || data
}

export function AuthProvider({ children }) {
  const [authState, setAuthState] = useState(INITIAL_STATE)
  const restoreStarted = useRef(false)

  const setAnonymous = useCallback((error = null, notice = null) => {
    setAuthState({
      status: 'anonymous',
      employee: null,
      accessToken: null,
      refreshToken: null,
      error,
      notice,
    })
  }, [])

  useEffect(() => {
    if (restoreStarted.current) {
      return
    }

    restoreStarted.current = true

    async function restoreSession() {
      const storedTokens = readStoredTokens()

      if (!storedTokens.accessToken) {
        setAnonymous()
        return
      }

      try {
        const currentEmployeeResponse = await currentEmployeeRequest(storedTokens.accessToken)
        setAuthState({
          status: 'authenticated',
          employee: requireEmployee(currentEmployeeResponse.data),
          accessToken: storedTokens.accessToken,
          refreshToken: storedTokens.refreshToken,
          error: null,
          notice: null,
        })
        return
      } catch (currentEmployeeError) {
        if (currentEmployeeError.code !== 'UNAUTHORIZED' || !storedTokens.refreshToken) {
          if (currentEmployeeError.code === 'UNAUTHORIZED') {
            clearStoredTokens()
            setAnonymous(currentEmployeeError, SESSION_EXPIRED_NOTICE)
            return
          }

          setAnonymous(currentEmployeeError)
          return
        }

        try {
          const refreshResponse = await refreshRequest(storedTokens.refreshToken)
          const nextTokens = requireTokenPair(refreshResponse.data)
          saveTokens(nextTokens.accessToken, nextTokens.refreshToken)

          const retryResponse = await currentEmployeeRequest(nextTokens.accessToken)
          setAuthState({
            status: 'authenticated',
            employee: requireEmployee(retryResponse.data),
            accessToken: nextTokens.accessToken,
            refreshToken: nextTokens.refreshToken,
            error: null,
            notice: null,
          })
        } catch (refreshError) {
          clearStoredTokens()
          setAnonymous(refreshError, SESSION_EXPIRED_NOTICE)
        }
      }
    }

    restoreSession()
  }, [setAnonymous])

  const login = useCallback(async ({ email, password }) => {
    try {
      const response = await loginRequest(email.trim(), password)
      const tokens = requireTokenPair(response.data)
      const employee = requireEmployee(response.data)

      saveTokens(tokens.accessToken, tokens.refreshToken)
      setAuthState({
        status: 'authenticated',
        employee,
        accessToken: tokens.accessToken,
        refreshToken: tokens.refreshToken,
        error: null,
        notice: null,
      })

      return employee
    } catch (error) {
      setAnonymous(error)
      throw error
    }
  }, [setAnonymous])

  const logout = useCallback(async () => {
    const refreshToken = authState.refreshToken || readStoredTokens().refreshToken

    try {
      if (refreshToken) {
        await logoutRequest(refreshToken)
      }
    } catch {
      // Local logout should remain reliable even if the API is offline.
    } finally {
      clearStoredTokens()
      setAnonymous()
    }
  }, [authState.refreshToken, setAnonymous])

  const clearAuthMessage = useCallback(() => {
    setAuthState((current) => ({
      ...current,
      error: null,
      notice: null,
    }))
  }, [])

  const value = useMemo(
    () => ({
      ...authState,
      login,
      logout,
      clearAuthMessage,
    }),
    [authState, clearAuthMessage, login, logout],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)

  if (!context) {
    throw new Error('useAuth must be used within AuthProvider')
  }

  return context
}
