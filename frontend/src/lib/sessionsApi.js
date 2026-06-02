import { apiRequest } from './api.js'
import { buildQuery } from './featureHelpers.js'

export function listSessions(accessToken, params = {}) {
  return apiRequest(`/api/v1/sessions${buildQuery(params)}`, { accessToken })
}

export function createSession(accessToken, payload) {
  return apiRequest('/api/v1/sessions', {
    method: 'POST',
    accessToken,
    body: payload,
  })
}

export function getSession(accessToken, sessionId) {
  return apiRequest(`/api/v1/sessions/${sessionId}`, { accessToken })
}

export function enrollSubscription(accessToken, sessionId, subscriptionId) {
  return apiRequest(`/api/v1/sessions/${sessionId}/enroll`, {
    method: 'POST',
    accessToken,
    body: { subscription_id: subscriptionId },
  })
}

export function checkInSessionSubscription(accessToken, sessionId, subscriptionId) {
  return apiRequest(`/api/v1/sessions/${sessionId}/checkin`, {
    method: 'POST',
    accessToken,
    body: { subscription_id: subscriptionId },
  })
}
