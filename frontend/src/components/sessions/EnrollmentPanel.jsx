import { useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import { enrollSubscription } from '../../lib/sessionsApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import { isObjectId } from './sessionFormatters.js'

function EnrollmentPanel({ accessToken, sessionId, onChanged }) {
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
      await enrollSubscription(accessToken, sessionId, subscriptionId.trim())
      setState({ status: 'success', error: null, notice: 'Subscription enrolled.' })
      setSubscriptionId('')
      await onChanged()
    } catch (mutationError) {
      setState({ status: 'error', error: mutationError, notice: '' })
    }
  }

  return (
    <DataPanel title="Enroll subscription">
      <form className="resource-form" onSubmit={handleSubmit}>
        <div className="field-group">
          <label htmlFor="enroll-subscription">Subscription ID</label>
          <input id="enroll-subscription" value={subscriptionId} onChange={(event) => { setSubscriptionId(event.target.value); setError('') }} />
          {error ? <span>{error}</span> : null}
        </div>
        <div className="resource-form__actions">
          <button className="btn-primary" type="submit" disabled={state.status === 'submitting'}>{state.status === 'submitting' ? 'Enrolling' : 'Enroll'}</button>
        </div>
      </form>
      {state.error ? <div className="form-alert" role="alert">{apiErrorText(state.error, 'Enrollment failed.')}</div> : null}
      {state.notice ? <div className="form-success" role="status">{state.notice}</div> : null}
    </DataPanel>
  )
}

export default EnrollmentPanel
