import {
  compactId,
  dateTimeLocalToRfc3339,
  formatDateTime,
  formatText,
  isObjectId,
  toDateTimeLocal,
} from '../../lib/featureHelpers.js'

export {
  compactId,
  dateTimeLocalToRfc3339,
  formatDateTime,
  formatText,
  isObjectId,
  toDateTimeLocal,
}

export function attendanceStatusLabel(status) {
  switch (status) {
    case 'attended':
      return 'Attended'
    case 'reported_missed':
      return 'Reported missed'
    case 'makeup':
      return 'Makeup'
    default:
      return formatText(status, 'Unknown')
  }
}

export function cleanAttendancePayload(values, status) {
  const payload = {
    subscription_id: values.subscription_id.trim(),
    branch_id: values.branch_id.trim(),
  }

  if (status) {
    payload.status = status
  }

  const date = dateTimeLocalToRfc3339(values.date)
  if (date) {
    payload.date = date
  }

  const makeupFor = dateTimeLocalToRfc3339(values.is_makeup_for)
  if (makeupFor) {
    payload.is_makeup_for = makeupFor
  }

  return payload
}
