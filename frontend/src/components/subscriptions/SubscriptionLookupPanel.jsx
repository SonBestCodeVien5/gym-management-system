import { useState } from 'react'
import { isObjectId } from './subscriptionFormatters.js'

function SubscriptionLookupPanel({ navigate, initialValue = '', label = 'Subscription ID' }) {
  const [subscriptionId, setSubscriptionId] = useState(initialValue)
  const [error, setError] = useState('')

  function handleSubmit(event) {
    event.preventDefault()
    const nextId = subscriptionId.trim()

    if (!nextId) {
      setError('Enter a subscription ObjectID.')
      return
    }

    if (!isObjectId(nextId)) {
      setError('Subscription ID must be a 24 character ObjectID.')
      return
    }

    navigate(`/app/subscriptions/${nextId}`)
  }

  return (
    <form className="resource-form resource-form--inline" onSubmit={handleSubmit}>
      <div className="field-group">
        <label htmlFor="subscription-lookup-id">{label}</label>
        <input
          id="subscription-lookup-id"
          value={subscriptionId}
          onChange={(event) => {
            setSubscriptionId(event.target.value)
            setError('')
          }}
          placeholder="24 character ObjectID"
          autoComplete="off"
          aria-invalid={error ? 'true' : undefined}
        />
        {error ? <span>{error}</span> : null}
      </div>
      <button className="btn-primary" type="submit">Open subscription</button>
    </form>
  )
}

export default SubscriptionLookupPanel
