const DEFAULT_API_BASE_URL = 'http://localhost:8080'

export const API_BASE_URL = (
  import.meta.env.VITE_API_BASE_URL || DEFAULT_API_BASE_URL
).replace(/\/+$/, '')

export function createClientError(code, message, details = {}, status = 0) {
  return {
    code,
    message,
    details: details && typeof details === 'object' ? details : {},
    status,
  }
}

async function readJson(response) {
  const text = await response.text()

  if (!text) {
    return null
  }

  try {
    return JSON.parse(text)
  } catch {
    return null
  }
}

function normalizeApiError(payload, status) {
  if (payload?.error) {
    return createClientError(
      payload.error.code || 'API_ERROR',
      payload.error.message || 'Request failed',
      payload.error.details || {},
      status,
    )
  }

  return createClientError(
    'API_ERROR',
    payload?.message || `Request failed with status ${status}`,
    {},
    status,
  )
}

export async function apiRequest(path, options = {}) {
  const {
    method = 'GET',
    body,
    accessToken,
    headers = {},
  } = options
  const requestHeaders = new Headers(headers)

  requestHeaders.set('Accept', 'application/json')

  if (body !== undefined) {
    requestHeaders.set('Content-Type', 'application/json')
  }

  if (accessToken) {
    requestHeaders.set('Authorization', `Bearer ${accessToken}`)
  }

  let response

  try {
    response = await fetch(`${API_BASE_URL}${path}`, {
      method,
      headers: requestHeaders,
      body: body === undefined ? undefined : JSON.stringify(body),
    })
  } catch {
    throw createClientError(
      'NETWORK_ERROR',
      'Cannot reach API server. Check that the backend is running.',
    )
  }

  const payload = await readJson(response)

  if (!response.ok) {
    throw normalizeApiError(payload, response.status)
  }

  return payload || { message: '', data: null }
}

export function loginRequest(email, password) {
  return apiRequest('/api/v1/auth/login', {
    method: 'POST',
    body: { email, password },
  })
}

export function currentEmployeeRequest(accessToken) {
  return apiRequest('/api/v1/auth/me', { accessToken })
}

export function refreshRequest(refreshToken) {
  return apiRequest('/api/v1/auth/refresh', {
    method: 'POST',
    body: { refresh_token: refreshToken },
  })
}

export function logoutRequest(refreshToken) {
  return apiRequest('/api/v1/auth/logout', {
    method: 'POST',
    body: { refresh_token: refreshToken },
  })
}
