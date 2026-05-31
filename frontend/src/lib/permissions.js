export const ROLE_LABELS = {
  admin: 'Admin',
  manager: 'Manager',
  trainer: 'Trainer',
  receptionist: 'Receptionist',
}

export function hasAnyRole(employeeRoles = [], allowedRoles = []) {
  if (!allowedRoles.length) {
    return true
  }

  return allowedRoles.some((role) => employeeRoles.includes(role))
}

export function canAccessRoute(route, employeeRoles = []) {
  return hasAnyRole(employeeRoles, route?.roles || [])
}

export function formatRoles(roles = []) {
  if (!roles.length) {
    return 'All staff'
  }

  return roles.map((role) => ROLE_LABELS[role] || role).join(', ')
}
