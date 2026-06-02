export const OBJECT_ID_PATTERN = /^[a-f0-9]{24}$/i

export function isObjectId(value) {
  return OBJECT_ID_PATTERN.test((value || '').trim())
}

export function compactId(value) {
  const text = String(value || '')

  if (!text || text.length <= 12) {
    return text || 'Not set'
  }

  return `${text.slice(0, 6)}...${text.slice(-6)}`
}

export function formatText(value, fallback = 'Not set') {
  return value === undefined || value === null || value === '' ? fallback : value
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

export function parseTags(value) {
  return String(value || '')
    .split(/[\n,]/)
    .map((tag) => tag.trim())
    .filter(Boolean)
}

export function formatTags(tags) {
  return Array.isArray(tags) && tags.length ? tags.join(', ') : 'None'
}

export function toDateTimeLocal(value) {
  if (!value) {
    return ''
  }

  const date = new Date(value)

  if (Number.isNaN(date.getTime())) {
    return ''
  }

  const offset = date.getTimezoneOffset()
  const localDate = new Date(date.getTime() - offset * 60 * 1000)
  return localDate.toISOString().slice(0, 16)
}

export function dateTimeLocalToRfc3339(value) {
  if (!value) {
    return ''
  }

  const date = new Date(value)

  if (Number.isNaN(date.getTime())) {
    return ''
  }

  return date.toISOString()
}

export function buildQuery(params = {}) {
  const searchParams = new URLSearchParams()

  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      searchParams.set(key, value)
    }
  })

  const query = searchParams.toString()
  return query ? `?${query}` : ''
}

export function cleanPayload(values) {
  return Object.fromEntries(
    Object.entries(values).map(([key, value]) => [
      key,
      typeof value === 'string' ? value.trim() : value,
    ]),
  )
}

export function apiErrorText(error, fallback = 'Request failed.') {
  if (!error) {
    return fallback
  }

  return error.code ? `${error.message || fallback} (${error.code})` : error.message || fallback
}
