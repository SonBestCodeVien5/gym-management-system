import { apiRequest } from './api.js'

export function checkInAttendance(accessToken, payload) {
  return apiRequest('/api/v1/attendance/checkin', {
    method: 'POST',
    accessToken,
    body: payload,
  })
}

export function reportMissedAttendance(accessToken, payload) {
  return apiRequest('/api/v1/attendance/report', {
    method: 'POST',
    accessToken,
    body: payload,
  })
}

export function createMakeupAttendance(accessToken, payload) {
  return apiRequest('/api/v1/attendance/makeup', {
    method: 'POST',
    accessToken,
    body: payload,
  })
}

export function listSubscriptionAttendance(accessToken, subscriptionId) {
  return apiRequest(`/api/v1/subscriptions/${subscriptionId}/attendance`, { accessToken })
}
