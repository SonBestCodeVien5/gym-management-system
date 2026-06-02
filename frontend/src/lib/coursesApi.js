import { apiRequest } from './api.js'

export function listCourses(accessToken) {
  return apiRequest('/api/v1/courses', { accessToken })
}

export function getCourse(accessToken, courseId) {
  return apiRequest(`/api/v1/courses/${courseId}`, { accessToken })
}

export function createCourse(accessToken, payload) {
  return apiRequest('/api/v1/courses', {
    method: 'POST',
    accessToken,
    body: payload,
  })
}

export function updateCourse(accessToken, courseId, payload) {
  return apiRequest(`/api/v1/courses/${courseId}`, {
    method: 'PATCH',
    accessToken,
    body: payload,
  })
}

export function deleteCourse(accessToken, courseId) {
  return apiRequest(`/api/v1/courses/${courseId}`, {
    method: 'DELETE',
    accessToken,
  })
}
