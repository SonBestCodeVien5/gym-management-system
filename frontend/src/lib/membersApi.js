import { apiRequest } from './api.js'

export function createMember(accessToken, payload) {
  return apiRequest('/api/v1/members', {
    method: 'POST',
    accessToken,
    body: payload,
  })
}

export function getMember(accessToken, memberId) {
  return apiRequest(`/api/v1/members/${encodeURIComponent(memberId)}`, {
    accessToken,
  })
}

export function activateMember(accessToken, memberId, subscriptionId) {
  return apiRequest(`/api/v1/members/${encodeURIComponent(memberId)}/activate`, {
    method: 'PATCH',
    accessToken,
    body: { subscription_id: subscriptionId },
  })
}

export function listMemberSubscriptions(accessToken, memberId) {
  return apiRequest(`/api/v1/members/${encodeURIComponent(memberId)}/subscriptions`, {
    accessToken,
  })
}
