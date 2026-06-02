import { apiRequest } from './api.js'
import { buildQuery } from './featureHelpers.js'

export function listBranches(accessToken) {
  return apiRequest('/api/v1/branches', { accessToken })
}

export function getBranch(accessToken, branchId) {
  return apiRequest(`/api/v1/branches/${branchId}`, { accessToken })
}

export function createBranch(accessToken, payload) {
  return apiRequest('/api/v1/branches', {
    method: 'POST',
    accessToken,
    body: payload,
  })
}

export function updateBranch(accessToken, branchId, payload) {
  return apiRequest(`/api/v1/branches/${branchId}`, {
    method: 'PATCH',
    accessToken,
    body: payload,
  })
}

export function deleteBranch(accessToken, branchId) {
  return apiRequest(`/api/v1/branches/${branchId}`, {
    method: 'DELETE',
    accessToken,
  })
}

export function nearbyBranches(accessToken, params) {
  return apiRequest(`/api/v1/branches/nearby${buildQuery(params)}`, { accessToken })
}
