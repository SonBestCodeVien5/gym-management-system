import {
  compactId,
  formatMoney,
  formatTags,
  formatText,
  isObjectId,
  parseTags,
} from '../../lib/featureHelpers.js'

export { compactId, formatMoney, formatTags, formatText, isObjectId, parseTags }

export function formatCoordinates(location) {
  const coordinates = location?.coordinates

  if (!Array.isArray(coordinates) || coordinates.length < 2) {
    return 'Not set'
  }

  return `${coordinates[0]}, ${coordinates[1]}`
}

export function cleanCoursePayload(values) {
  return {
    title: values.title.trim(),
    level: values.level.trim(),
    allowed_tags: parseTags(values.allowed_tags),
    base_price: Number(values.base_price),
    session_count: Number(values.session_count),
    description: values.description.trim(),
  }
}

export function cleanBranchPayload(values) {
  const payload = {
    branch_code: values.branch_code.trim(),
    name: values.name.trim(),
    address: values.address.trim(),
    province: values.province.trim(),
    location: {
      type: 'Point',
      coordinates: [Number(values.lng), Number(values.lat)],
    },
  }

  if (values.manager_id.trim()) {
    payload.manager_id = values.manager_id.trim()
  }

  return payload
}
