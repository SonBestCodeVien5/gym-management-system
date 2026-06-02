import DataPanel from '../DataPanel.jsx'
import {
  compactId,
  formatDateTime,
  formatMoney,
  formatTags,
  subscriptionStatusMeta,
} from './subscriptionFormatters.js'

function SubscriptionSummaryPanel({ subscription, navigate }) {
  const statusMeta = subscriptionStatusMeta(subscription.status)

  return (
    <DataPanel title="Subscription summary">
      <div className="status-strip">
        <span className={`status-pill status-pill--${statusMeta.tone}`}>{statusMeta.label}</span>
        <span>{subscription.remaining_sessions} of {subscription.total_sessions} sessions left</span>
      </div>

      <dl className="detail-grid">
        <div><dt>Member</dt><dd>{compactId(subscription.member_id)}</dd></div>
        <div><dt>Course</dt><dd>{compactId(subscription.course_id)}</dd></div>
        <div><dt>Home branch</dt><dd>{compactId(subscription.home_branch_id)}</dd></div>
        <div><dt>Weekly limit</dt><dd>{subscription.session_per_week}</dd></div>
        <div><dt>Start</dt><dd>{formatDateTime(subscription.start_date)}</dd></div>
        <div><dt>End</dt><dd>{formatDateTime(subscription.end_date)}</dd></div>
        <div><dt>Subtotal</dt><dd>{formatMoney(subscription.subtotal_amount)}</dd></div>
        <div><dt>Discount</dt><dd>{subscription.discount_type} · {formatMoney(subscription.discount_amount)}</dd></div>
        <div><dt>Paid</dt><dd>{formatMoney(subscription.total_amount_paid)}</dd></div>
        <div><dt>Tags</dt><dd>{formatTags(subscription.allowed_tags)}</dd></div>
      </dl>

      <div className="resource-form__actions">
        <button className="btn-outline" type="button" onClick={() => navigate(`/app/members/${subscription.member_id}`)}>
          Open member
        </button>
        <button className="btn-outline" type="button" onClick={() => navigate(`/app/subscriptions/${subscription.id}/attendance`)}>
          Attendance
        </button>
      </div>
    </DataPanel>
  )
}

export default SubscriptionSummaryPanel
