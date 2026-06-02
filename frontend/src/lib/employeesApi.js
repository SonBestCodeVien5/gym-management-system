import { apiRequest } from './api.js'
import { buildQuery } from './featureHelpers.js'

export function listEmployees(accessToken, params = {}) {
  return apiRequest(`/api/v1/employees${buildQuery(params)}`, { accessToken })
}

export function getEmployee(accessToken, employeeId) {
  return apiRequest(`/api/v1/employees/${employeeId}`, { accessToken })
}

export function createEmployee(accessToken, payload) {
  return apiRequest('/api/v1/employees', {
    method: 'POST',
    accessToken,
    body: payload,
  })
}

export function updateEmployee(accessToken, employeeId, payload) {
  return apiRequest(`/api/v1/employees/${employeeId}`, {
    method: 'PATCH',
    accessToken,
    body: payload,
  })
}

export function resetEmployeePassword(accessToken, employeeId, password) {
  return apiRequest(`/api/v1/employees/${employeeId}/password`, {
    method: 'PATCH',
    accessToken,
    body: { password },
  })
}
