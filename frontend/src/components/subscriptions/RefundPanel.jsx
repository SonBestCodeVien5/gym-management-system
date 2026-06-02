import { useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import { refundSubscription } from '../../lib/subscriptionsApi.js'
import { apiErrorText, formatMoney } from '../../lib/featureHelpers.js'

function RefundPanel({ accessToken, subscription, onChanged }) {
  const [reason, setReason] = useState('')
  const [state, setState] = useState({ status: 'idle', error: null, notice: '', refund: null })
  const canRefund = subscription.status === 'active'

  async function handleSubmit(event) {
    event.preventDefault()

    if (!window.confirm('Refund this subscription?')) {
      return
    }

    setState({ status: 'submitting', error: null, notice: '', refund: null })

    try {
      const response = await refundSubscription(accessToken, subscription.id, reason.trim())
      setState({
        status: 'success',
        error: null,
        notice: 'Subscription refunded.',
        refund: response.refund || null,
      })
      await onChanged()
    } catch (error) {
      setState({ status: 'error', error, notice: '', refund: null })
    }
  }

  return (
    <DataPanel title="Refund" description="Refund applies to active subscriptions and uses backend remaining-session calculation.">
      <form className="resource-form" onSubmit={handleSubmit}>
        <div className="field-group">
          <label htmlFor="refund-reason">Reason</label>
          <textarea id="refund-reason" rows="3" value={reason} onChange={(event) => setReason(event.target.value)} />
        </div>
        <div className="resource-form__actions">
          <button className="btn-primary" type="submit" disabled={!canRefund || state.status === 'submitting'}>
            {state.status === 'submitting' ? 'Refunding' : 'Refund subscription'}
          </button>
        </div>
      </form>

      {state.error ? <div className="form-alert" role="alert">{apiErrorText(state.error, 'Refund failed.')}</div> : null}
      {state.notice ? (
        <div className="form-success" role="status">
          {state.notice}
          {state.refund?.refund_amount ? ` Amount: ${formatMoney(state.refund.refund_amount)}.` : ''}
        </div>
      ) : null}
    </DataPanel>
  )
}

export default RefundPanel
