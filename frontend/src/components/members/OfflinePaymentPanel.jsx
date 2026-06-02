import { useEffect, useMemo, useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import { activateMember } from '../../lib/membersApi.js'
import { isObjectId } from './memberFormatters.js'

function OfflinePaymentPanel({
  accessToken,
  memberId,
  subscriptions,
  selectedSubscriptionId,
  onSelectSubscription,
  onActivated,
}) {
  const pendingSubscriptions = useMemo(
    () => subscriptions.filter((subscription) => subscription.status === 'pending'),
    [subscriptions],
  )
  const [subscriptionId, setSubscriptionId] = useState(selectedSubscriptionId || '')
  const [status, setStatus] = useState('idle')
  const [error, setError] = useState(null)
  const [message, setMessage] = useState('')
  const trimmedSubscriptionId = subscriptionId.trim()
  const hasInvalidSubscriptionId = Boolean(trimmedSubscriptionId) && !isObjectId(trimmedSubscriptionId)
  const subscriptionHelpId = [
    hasInvalidSubscriptionId ? 'activation-subscription-id-error' : '',
    error ? 'activation-error' : '',
  ].filter(Boolean).join(' ') || undefined

  useEffect(() => {
    if (selectedSubscriptionId) {
      setSubscriptionId(selectedSubscriptionId)
      setError(null)
      setMessage('')
    }
  }, [selectedSubscriptionId])

  async function handleSubmit(event) {
    event.preventDefault()
    const nextSubscriptionId = subscriptionId.trim()

    if (!isObjectId(nextSubscriptionId)) {
      setError({ message: 'Subscription ID must be a 24 character hex ObjectID.' })
      return
    }

    setStatus('submitting')
    setError(null)
    setMessage('')

    try {
      await activateMember(accessToken, memberId, nextSubscriptionId)
      setStatus('success')
      setMessage('Offline payment confirmed. Member and subscriptions were refreshed.')
      await onActivated()
    } catch (activationError) {
      setStatus('error')
      setError(activationError)
    }
  }

  function handleSelectChange(event) {
    const nextSubscriptionId = event.target.value
    setSubscriptionId(nextSubscriptionId)
    onSelectSubscription(nextSubscriptionId)
    setError(null)
    setMessage('')
  }

  return (
    <DataPanel
      title="Offline payment"
      description="Confirm payment for a pending subscription that belongs to this member."
    >
      {pendingSubscriptions.length ? (
        <div className="field-group">
          <label htmlFor="pending-subscription">Pending subscription</label>
          <select id="pending-subscription" value={subscriptionId} onChange={handleSelectChange}>
            <option value="">Select a pending subscription</option>
            {pendingSubscriptions.map((subscription) => (
              <option key={subscription.id} value={subscription.id}>
                {subscription.id}
              </option>
            ))}
          </select>
        </div>
      ) : (
        <p className="panel-copy">
          No pending subscription is loaded for this member. You can still enter a known subscription
          ID manually.
        </p>
      )}

      <form className="member-form" onSubmit={handleSubmit}>
        <div className="field-group">
          <label htmlFor="activation-subscription-id">Subscription ID</label>
          <input
            id="activation-subscription-id"
            type="text"
            value={subscriptionId}
            onChange={(event) => {
              setSubscriptionId(event.target.value)
              setError(null)
              setMessage('')
            }}
            placeholder="24 character ObjectID"
            autoComplete="off"
            aria-invalid={error || hasInvalidSubscriptionId ? 'true' : undefined}
            aria-describedby={subscriptionHelpId}
          />
          {hasInvalidSubscriptionId ? (
            <span id="activation-subscription-id-error">
              Subscription ID must be a 24 character hex ObjectID.
            </span>
          ) : null}
        </div>

        {error ? (
          <div className="form-alert" id="activation-error" role="alert">
            {error.message || 'Activation failed.'}
            {error.code ? <span>{error.code}</span> : null}
          </div>
        ) : null}

        {message ? (
          <div className="form-success" role="status" aria-live="polite">
            {message}
          </div>
        ) : null}

        <button
          className="btn-primary"
          type="submit"
          disabled={status === 'submitting' || !isObjectId(subscriptionId)}
        >
          {status === 'submitting' ? 'Confirming payment' : 'Confirm payment'}
        </button>
      </form>
    </DataPanel>
  )
}

export default OfflinePaymentPanel
