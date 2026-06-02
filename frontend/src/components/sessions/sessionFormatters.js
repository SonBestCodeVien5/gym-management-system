import {
  compactId,
  dateTimeLocalToRfc3339,
  formatDateTime,
  formatTags,
  formatText,
  isObjectId,
  parseTags,
} from '../../lib/featureHelpers.js'

export {
  compactId,
  dateTimeLocalToRfc3339,
  formatDateTime,
  formatTags,
  formatText,
  isObjectId,
  parseTags,
}

export function capacityLabel(session) {
  return `${session.enrolled_count || 0}/${session.capacity || 0}`
}

export function cleanSessionPayload(values) {
  return {
    branch_id: values.branch_id.trim(),
    trainer_id: values.trainer_id.trim(),
    course_level: values.course_level.trim(),
    scheduled_at: dateTimeLocalToRfc3339(values.scheduled_at),
    duration_min: Number(values.duration_min),
    capacity: Number(values.capacity),
    tags: parseTags(values.tags),
  }
}
