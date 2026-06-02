import { useState } from 'react'
import { isObjectId } from './memberFormatters.js'

function MemberLookupPanel({ navigate, initialValue = '' }) {
  const [memberId, setMemberId] = useState(initialValue)
  const [error, setError] = useState('')

  function handleSubmit(event) {
    event.preventDefault()
    const nextMemberId = memberId.trim()

    if (!nextMemberId) {
      setError('Enter a member ObjectID.')
      return
    }

    if (!isObjectId(nextMemberId)) {
      setError('Member ID must be a 24 character hex ObjectID.')
      return
    }

    setError('')
    navigate(`/app/members/${nextMemberId}`)
  }

  return (
    <form className="member-form member-lookup" onSubmit={handleSubmit}>
      <div className="field-group">
        <label htmlFor="member-lookup-id">Member ID</label>
        <input
          id="member-lookup-id"
          type="text"
          value={memberId}
          onChange={(event) => {
            setMemberId(event.target.value)
            setError('')
          }}
          placeholder="24 character ObjectID"
          autoComplete="off"
          aria-invalid={error ? 'true' : undefined}
          aria-describedby={error ? 'member-lookup-error' : undefined}
        />
        {error ? <span id="member-lookup-error">{error}</span> : null}
      </div>

      <button className="btn-primary" type="submit">
        Open profile
      </button>
    </form>
  )
}

export default MemberLookupPanel
