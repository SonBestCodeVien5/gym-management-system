import { useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import { checkInSessionSubscription } from '../../lib/sessionsApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import { isObjectId } from './sessionFormatters.js'

function SessionCheckInPanel({ accessToken, sessionId, onChanged }) {
  const [subscriptionId, setSubscriptionId] = useState('')
  const [error, setError] = useState('')
  const [state, setState] = useState({ status: 'idle', error: null, notice: '' })

  async function handleSubmit(event) {
    event.preventDefault()

    if (!isObjectId(subscriptionId)) {
      setError('Subscription ID must be a 24 character ObjectID.')
      return
    }

    setState({ status: 'submitting', error: null, notice: '' })

    try {
      await checkInSessionSubscription(accessToken, sessionId, subscriptionId.trim())
      setState({ status: 'success', error: null, notice: 'Session check-in recorded.' })
      setSubscriptionId('')
      await onChanged()
    } catch (mutationError) {
      setState({ status: 'error', error: mutationError, notice: '' })
    }
  }

  return (
    <DataPanel title="Session check-in">
      <form className="resource-form" onSubmit={handleSubmit}>
        <div className="field-group">
          <label htmlFor="session-checkin-subscription">Subscription ID</label>
          <input id="session-checkin-subscription" value={subscriptionId} onChange={(event) => { setSubscriptionId(event.target.value); setError('') }} />
          {error ? <span>{error}</span> : null}
        </div>
        <div className="resource-form__actions">
          <button className="btn-primary" type="submit" disabled={state.status === 'submitting'}>{state.status === 'submitting' ? 'Checking in' : 'Check in'}</button>
        </div>
      </form>
      {state.error ? <div className="form-alert" role="alert">{apiErrorText(state.error, 'Session check-in failed.')}</div> : null}
      {state.notice ? <div className="form-success" role="status">{state.notice}</div> : null}
    </DataPanel>
  )
}

export default SessionCheckInPanel
