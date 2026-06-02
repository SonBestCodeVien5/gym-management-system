import {
  compactId,
  dateTimeLocalToRfc3339,
  formatDateTime,
  formatMoney,
  formatTags,
  formatText,
  isObjectId,
  toDateTimeLocal,
} from '../../lib/featureHelpers.js'

export {
  compactId,
  dateTimeLocalToRfc3339,
  formatDateTime,
  formatMoney,
  formatTags,
  formatText,
  isObjectId,
  toDateTimeLocal,
}

export function subscriptionStatusMeta(status) {
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

export function cleanSubscriptionPayload(values) {
  const payload = {
    member_id: values.member_id.trim(),
    course_id: values.course_id.trim(),
    home_branch_id: values.home_branch_id.trim(),
    start_date: dateTimeLocalToRfc3339(values.start_date),
    end_date: dateTimeLocalToRfc3339(values.end_date),
    session_per_week: Number(values.session_per_week),
    discount_type: values.discount_type || 'none',
    discount_value: Number(values.discount_value || 0),
    promo_code: values.promo_code.trim(),
  }

  if (payload.discount_type === 'none') {
    payload.discount_value = 0
    payload.promo_code = ''
  }

  return payload
}

export function cleanSuspensionPayload(values) {
  return {
    start_date: dateTimeLocalToRfc3339(values.start_date),
    end_date: dateTimeLocalToRfc3339(values.end_date),
    frozen_session: Number(values.frozen_session || 0),
    reason: values.reason.trim(),
  }
}
