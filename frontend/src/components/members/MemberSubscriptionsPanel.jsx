import DataPanel from '../DataPanel.jsx'
import StateBlock from '../StateBlock.jsx'
import {
  compactId,
  formatDateTime,
  formatMoney,
  getSubscriptionStatus,
} from './memberFormatters.js'

function MemberSubscriptionsPanel({ subscriptions, status, error, selectedSubscriptionId, onSelectSubscription }) {
  if (status === 'loading') {
    return (
      <DataPanel title="Subscriptions">
        <StateBlock tone="loading" title="Loading subscriptions" message="Fetching member subscription records." />
      </DataPanel>
    )
  }

  if (status === 'error') {
    return (
      <DataPanel title="Subscriptions">
        <StateBlock
          tone="error"
          title="Could not load subscriptions"
          message={error?.message || 'Subscription history could not be loaded.'}
        />
      </DataPanel>
    )
  }

  if (!subscriptions.length) {
    return (
      <DataPanel title="Subscriptions">
        <StateBlock
          tone="empty"
          title="No subscriptions"
          message="This member does not have subscription records yet."
        />
      </DataPanel>
    )
  }

  return (
    <DataPanel
      title="Subscriptions"
      description="Member-scoped subscription records returned by the backend."
    >
      <div className="table-scroll">
        <table className="member-table member-subscriptions-table">
          <thead>
            <tr>
              <th>Subscription</th>
              <th>Status</th>
              <th>Sessions</th>
              <th>Total paid</th>
              <th>Dates</th>
              <th>Payment</th>
            </tr>
          </thead>
          <tbody>
            {subscriptions.map((subscription) => {
              const subscriptionStatus = getSubscriptionStatus(subscription.status)
              const isPending = subscription.status === 'pending'
              const isSelected = selectedSubscriptionId === subscription.id

              return (
                <tr key={subscription.id}>
                  <td data-label="Subscription">
                    <strong title={subscription.id}>{compactId(subscription.id)}</strong>
                    <span>Course {compactId(subscription.course_id)}</span>
                  </td>
                  <td data-label="Status">
                    <span className={`status-badge status-badge--${subscriptionStatus.tone}`}>
                      {subscriptionStatus.label}
                    </span>
                  </td>
                  <td data-label="Sessions">
                    {subscription.remaining_sessions}/{subscription.total_sessions}
                  </td>
                  <td data-label="Total paid">{formatMoney(subscription.total_amount_paid)}</td>
                  <td data-label="Dates">
                    <span>{formatDateTime(subscription.start_date)}</span>
                    <span>{formatDateTime(subscription.end_date)}</span>
                  </td>
                  <td data-label="Payment">
                    {isPending ? (
                      <button
                        className="btn-outline btn-inline"
                        type="button"
                        onClick={() => onSelectSubscription(subscription.id)}
                        aria-pressed={isSelected}
                      >
                        {isSelected ? 'Selected' : 'Select'}
                      </button>
                    ) : (
                      <span className="panel-copy">No action</span>
                    )}
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>
    </DataPanel>
  )
}

export default MemberSubscriptionsPanel
