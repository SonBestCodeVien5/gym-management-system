import {
  compactId,
  formatDateTime,
  formatText,
  isObjectId,
} from '../../lib/featureHelpers.js'

export { compactId, formatDateTime, formatText, isObjectId }

export const ROLE_OPTIONS = ['admin', 'manager', 'trainer', 'receptionist']
export const STATUS_OPTIONS = ['active', 'inactive']
export const LEVEL_OPTIONS = ['basic', 'advanced', 'professional']

export function roleLabel(roles) {
  return Array.isArray(roles) && roles.length ? roles.join(', ') : 'No role'
}

export function cleanEmployeePayload(values, { includePassword = false } = {}) {
  const payload = {
    employee_id: values.employee_id.trim(),
    full_name: values.full_name.trim(),
    email: values.email.trim(),
    role: values.role,
    level: values.level.trim(),
    phone: values.phone.trim(),
    branch_id: values.branch_id
      .split(/[\n,]/)
      .map((branchId) => branchId.trim())
      .filter(Boolean),
    status: values.status || 'active',
  }

  if (includePassword) {
    payload.password = values.password
  }

  return payload
}

export function employeeValuesFromEmployee(employee) {
  return {
    employee_id: employee?.employee_id || '',
    full_name: employee?.full_name || '',
    email: employee?.email || '',
    password: '',
    role: Array.isArray(employee?.role) ? employee.role : [],
    level: employee?.level || '',
    phone: employee?.phone || '',
    branch_id: Array.isArray(employee?.branch_id) ? employee.branch_id.join(', ') : '',
    status: employee?.status || 'active',
  }
}
