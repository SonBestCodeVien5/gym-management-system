import { useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import { suspendSubscription, unsuspendSubscription, expireSubscription } from '../../lib/subscriptionsApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import { cleanSuspensionPayload } from './subscriptionFormatters.js'

const INITIAL_VALUES = {
  start_date: '',
  end_date: '',
  frozen_session: '0',
  reason: '',
}

function validate(values) {
  const errors = {}
  const start = new Date(values.start_date)
  const end = new Date(values.end_date)
  const frozenSession = Number(values.frozen_session)

  if (!values.start_date || Number.isNaN(start.getTime())) {
    errors.start_date = 'Start date is required.'
  }

  if (!values.end_date || Number.isNaN(end.getTime())) {
    errors.end_date = 'End date is required.'
  }

  if (values.start_date && values.end_date && start > end) {
    errors.end_date = 'End date must be after start date.'
  }

  if (!Number.isInteger(frozenSession) || frozenSession < 0) {
    errors.frozen_session = 'Frozen sessions must be zero or more.'
  }

  return errors
}

function SubscriptionLifecyclePanel({ accessToken, subscription, onChanged }) {
  const [values, setValues] = useState(INITIAL_VALUES)
  const [errors, setErrors] = useState({})
  const [state, setState] = useState({ status: 'idle', error: null, notice: '' })

  function updateField(name, value) {
    setValues((current) => ({ ...current, [name]: value }))
    setErrors((current) => ({ ...current, [name]: '' }))
    setState((current) => ({ ...current, error: null, notice: '' }))
  }

  async function runMutation(action, successMessage) {
    setState({ status: 'submitting', error: null, notice: '' })

    try {
      await action()
      setState({ status: 'success', error: null, notice: successMessage })
      await onChanged()
    } catch (error) {
      setState({ status: 'error', error, notice: '' })
    }
  }

  function handleSuspend(event) {
    event.preventDefault()
    const nextErrors = validate(values)
    setErrors(nextErrors)

    if (Object.keys(nextErrors).length) {
      return
    }

    runMutation(
      () => suspendSubscription(accessToken, subscription.id, cleanSuspensionPayload(values)),
      'Subscription suspended.',
    )
  }

  const canSuspend = subscription.status === 'active'
  const canUnsuspend = subscription.status === 'suspended'
  const canExpire = ['pending', 'active', 'suspended'].includes(subscription.status)

  return (
    <DataPanel title="Lifecycle actions" description="Actions are enabled by current subscription status.">
      <div className="resource-stack">
        <form className="resource-form" onSubmit={handleSuspend}>
          <div className="resource-form__grid">
            <div className="field-group">
              <label htmlFor="suspend-start">Suspend start</label>
              <input id="suspend-start" type="datetime-local" value={values.start_date} onChange={(event) => updateField('start_date', event.target.value)} />
              {errors.start_date ? <span>{errors.start_date}</span> : null}
            </div>
            <div className="field-group">
              <label htmlFor="suspend-end">Suspend end</label>
              <input id="suspend-end" type="datetime-local" value={values.end_date} onChange={(event) => updateField('end_date', event.target.value)} />
              {errors.end_date ? <span>{errors.end_date}</span> : null}
            </div>
            <div className="field-group">
              <label htmlFor="frozen-session">Frozen sessions</label>
              <input id="frozen-session" type="number" min="0" value={values.frozen_session} onChange={(event) => updateField('frozen_session', event.target.value)} />
              {errors.frozen_session ? <span>{errors.frozen_session}</span> : null}
            </div>
            <div className="field-group">
              <label htmlFor="suspend-reason">Reason</label>
              <input id="suspend-reason" value={values.reason} onChange={(event) => updateField('reason', event.target.value)} />
            </div>
          </div>
          <div className="resource-form__actions">
            <button className="btn-primary" type="submit" disabled={!canSuspend || state.status === 'submitting'}>
              Suspend
            </button>
            <button className="btn-outline" type="button" disabled={!canUnsuspend || state.status === 'submitting'} onClick={() => runMutation(() => unsuspendSubscription(accessToken, subscription.id), 'Subscription unsuspended.')}>
              Unsuspend
            </button>
            <button className="btn-outline" type="button" disabled={!canExpire || state.status === 'submitting'} onClick={() => window.confirm('Expire this subscription?') && runMutation(() => expireSubscription(accessToken, subscription.id), 'Subscription expired.')}>
              Expire
            </button>
          </div>
        </form>

        {state.error ? <div className="form-alert" role="alert">{apiErrorText(state.error, 'Lifecycle action failed.')}</div> : null}
        {state.notice ? <div className="form-success" role="status">{state.notice}</div> : null}
      </div>
    </DataPanel>
  )
}

export default SubscriptionLifecyclePanel
