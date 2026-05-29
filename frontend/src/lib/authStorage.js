const ACCESS_TOKEN_KEY = 'gym.accessToken'
const REFRESH_TOKEN_KEY = 'gym.refreshToken'

function canUseStorage() {
  return typeof window !== 'undefined' && Boolean(window.localStorage)
}

export function readStoredTokens() {
  if (!canUseStorage()) {
    return { accessToken: null, refreshToken: null }
  }

  try {
    return {
      accessToken: window.localStorage.getItem(ACCESS_TOKEN_KEY),
      refreshToken: window.localStorage.getItem(REFRESH_TOKEN_KEY),
    }
  } catch {
    return { accessToken: null, refreshToken: null }
  }
}

export function saveTokens(accessToken, refreshToken) {
  if (!canUseStorage()) {
    return
  }

  try {
    window.localStorage.setItem(ACCESS_TOKEN_KEY, accessToken)
    window.localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken)
  } catch {
    // Storage failures should not expose token values or block the UI.
  }
}

export function clearStoredTokens() {
  if (!canUseStorage()) {
    return
  }

  try {
    window.localStorage.removeItem(ACCESS_TOKEN_KEY)
    window.localStorage.removeItem(REFRESH_TOKEN_KEY)
  } catch {
    // Ignore storage cleanup failures; auth state is still cleared in memory.
  }
}
