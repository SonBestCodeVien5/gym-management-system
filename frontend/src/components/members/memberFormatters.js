export const OBJECT_ID_PATTERN = /^[a-f0-9]{24}$/i

export function isObjectId(value) {
  return OBJECT_ID_PATTERN.test((value || '').trim())
}

export function compactId(value) {
  if (!value || value.length <= 12) {
    return value || 'Not set'
  }

  return `${value.slice(0, 6)}...${value.slice(-6)}`
}

export function formatDateTime(value) {
  if (!value) {
    return 'Not set'
  }

  const date = new Date(value)

  if (Number.isNaN(date.getTime())) {
    return value
  }

  return new Intl.DateTimeFormat('en', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(date)
}

export function formatMoney(value) {
  const amount = Number(value)

  if (!Number.isFinite(amount)) {
    return 'Not set'
  }

  return new Intl.NumberFormat('vi-VN', {
    style: 'currency',
    currency: 'VND',
    maximumFractionDigits: 0,
  }).format(amount)
}

export function formatText(value, fallback = 'Not set') {
  return value === undefined || value === null || value === '' ? fallback : value
}

export function getRegistrationStatus(member) {
  if (member?.is_suspended) {
    return { label: 'Suspended', tone: 'danger' }
  }

  if (member?.is_registered) {
    return { label: 'Registered', tone: 'success' }
  }

  return { label: 'Pending payment', tone: 'warning' }
}

export function getSubscriptionStatus(status) {
  switch (status) {
    case 'active':
      return { label: 'Active', tone: 'success' }
    case 'pending':
      return { label: 'Pending', tone: 'warning' }
    case 'suspended':
      return { label: 'Suspended', tone: 'neutral' }
    case 'expired':
      return { label: 'Expired', tone: 'danger' }
    case 'refunded':
      return { label: 'Refunded', tone: 'neutral' }
    default:
      return { label: formatText(status, 'Unknown'), tone: 'neutral' }
  }
}

export function cleanMemberPayload(values) {
  return Object.fromEntries(
    Object.entries(values).map(([key, value]) => [key, typeof value === 'string' ? value.trim() : value]),
  )
}
