import { apiRequest } from './api.js'
import { buildQuery } from './featureHelpers.js'

export function getDashboardSummary(accessToken, params = {}) {
  return apiRequest(`/api/v1/dashboard/summary${buildQuery(params)}`, { accessToken })
}

export function getDashboardRevenue(accessToken, params = {}) {
  return apiRequest(`/api/v1/dashboard/revenue${buildQuery(params)}`, { accessToken })
}

export function getDashboardPlanDistribution(accessToken, params = {}) {
  return apiRequest(`/api/v1/dashboard/plans${buildQuery(params)}`, { accessToken })
}

export function getDashboardRecentMembers(accessToken, params = {}) {
  return apiRequest(`/api/v1/dashboard/members/recent${buildQuery(params)}`, { accessToken })
}

export function getDashboardTodaySessions(accessToken, params = {}) {
  return apiRequest(`/api/v1/dashboard/sessions/today${buildQuery(params)}`, { accessToken })
}
