import { apiRequest } from './api.js'

export function createSubscription(accessToken, payload) {
  return apiRequest('/api/v1/subscriptions', {
    method: 'POST',
    accessToken,
    body: payload,
  })
}

export function getSubscription(accessToken, subscriptionId) {
  return apiRequest(`/api/v1/subscriptions/${subscriptionId}`, { accessToken })
}

export function suspendSubscription(accessToken, subscriptionId, payload) {
  return apiRequest(`/api/v1/subscriptions/${subscriptionId}/suspend`, {
    method: 'PATCH',
    accessToken,
    body: payload,
  })
}

export function unsuspendSubscription(accessToken, subscriptionId) {
  return apiRequest(`/api/v1/subscriptions/${subscriptionId}/unsuspend`, {
    method: 'PATCH',
    accessToken,
  })
}

export function expireSubscription(accessToken, subscriptionId) {
  return apiRequest(`/api/v1/subscriptions/${subscriptionId}/expire`, {
    method: 'PATCH',
    accessToken,
  })
}

export function refundSubscription(accessToken, subscriptionId, reason) {
  return apiRequest(`/api/v1/subscriptions/${subscriptionId}/refund`, {
    method: 'POST',
    accessToken,
    body: { reason },
  })
}
