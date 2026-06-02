import { useState } from 'react'
import { isObjectId } from './attendanceFormatters.js'

function SubscriptionHistoryLookup({ navigate, initialValue = '' }) {
  const [subscriptionId, setSubscriptionId] = useState(initialValue)
  const [error, setError] = useState('')

  function handleSubmit(event) {
    event.preventDefault()
    const nextId = subscriptionId.trim()

    if (!isObjectId(nextId)) {
      setError('Enter a valid subscription ObjectID.')
      return
    }

    navigate(`/app/subscriptions/${nextId}/attendance`)
  }

  return (
    <form className="resource-form resource-form--inline" onSubmit={handleSubmit}>
      <div className="field-group">
        <label htmlFor="attendance-history-id">Subscription ID</label>
        <input id="attendance-history-id" value={subscriptionId} onChange={(event) => { setSubscriptionId(event.target.value); setError('') }} />
        {error ? <span>{error}</span> : null}
      </div>
      <button className="btn-primary" type="submit">Open history</button>
    </form>
  )
}

export default SubscriptionHistoryLookup
