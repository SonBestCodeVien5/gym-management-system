import { useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import PageHeader from '../PageHeader.jsx'
import { createMember } from '../../lib/membersApi.js'
import { useAuth } from '../../context/AuthContext.jsx'
import { cleanMemberPayload } from './memberFormatters.js'

const INITIAL_VALUES = {
  ccid: '',
  full_name: '',
  email: '',
  phone: '',
  gender: '',
  level: '',
}

function validate(values) {
  const errors = {}

  if (!values.ccid.trim()) {
    errors.ccid = 'CCID is required.'
  }

  if (!values.full_name.trim()) {
    errors.full_name = 'Full name is required.'
  }

  return errors
}

function MemberCreateView({ navigate }) {
  const { accessToken } = useAuth()
  const [values, setValues] = useState(INITIAL_VALUES)
  const [errors, setErrors] = useState({})
  const [submitError, setSubmitError] = useState(null)
  const [status, setStatus] = useState('idle')

  function updateField(name, value) {
    setValues((current) => ({ ...current, [name]: value }))
    setErrors((current) => ({ ...current, [name]: '' }))
    setSubmitError(null)
  }

  async function handleSubmit(event) {
    event.preventDefault()
    const nextErrors = validate(values)
    setErrors(nextErrors)

    if (Object.keys(nextErrors).length) {
      return
    }

    setStatus('submitting')
    setSubmitError(null)

    try {
      const response = await createMember(accessToken, cleanMemberPayload(values))
      const memberId = response.data?.id

      if (memberId) {
        navigate(`/app/members/${memberId}`, { replace: true })
        return
      }

      setStatus('success')
    } catch (error) {
      setStatus('error')
      setSubmitError(error)
    }
  }

  return (
    <div className="module-page members-workspace">
      <PageHeader
        eyebrow="Members"
        title="New member"
        description="Create a member profile. Subscription creation remains in the later subscription workspace."
        actions={(
          <button className="btn-outline" type="button" onClick={() => navigate('/app/members')}>
            Back to members
          </button>
        )}
      />

      <DataPanel title="Member details" description="Only CCID and full name are required by the current API.">
        <form className="member-form" onSubmit={handleSubmit}>
          <div className="member-form__grid">
            <div className="field-group">
              <label htmlFor="member-ccid">CCID</label>
              <input
                id="member-ccid"
                type="text"
                value={values.ccid}
                onChange={(event) => updateField('ccid', event.target.value)}
                aria-invalid={errors.ccid ? 'true' : undefined}
                aria-describedby={errors.ccid ? 'member-ccid-error' : undefined}
              />
              {errors.ccid ? <span id="member-ccid-error">{errors.ccid}</span> : null}
            </div>

            <div className="field-group">
              <label htmlFor="member-full-name">Full name</label>
              <input
                id="member-full-name"
                type="text"
                value={values.full_name}
                onChange={(event) => updateField('full_name', event.target.value)}
                aria-invalid={errors.full_name ? 'true' : undefined}
                aria-describedby={errors.full_name ? 'member-full-name-error' : undefined}
              />
              {errors.full_name ? <span id="member-full-name-error">{errors.full_name}</span> : null}
            </div>

            <div className="field-group">
              <label htmlFor="member-email">Email</label>
              <input
                id="member-email"
                type="email"
                value={values.email}
                onChange={(event) => updateField('email', event.target.value)}
              />
            </div>

            <div className="field-group">
              <label htmlFor="member-phone">Phone</label>
              <input
                id="member-phone"
                type="tel"
                value={values.phone}
                onChange={(event) => updateField('phone', event.target.value)}
              />
            </div>

            <div className="field-group">
              <label htmlFor="member-gender">Gender</label>
              <input
                id="member-gender"
                type="text"
                value={values.gender}
                onChange={(event) => updateField('gender', event.target.value)}
                placeholder="male, female, or local convention"
              />
            </div>

            <div className="field-group">
              <label htmlFor="member-level">Level</label>
              <input
                id="member-level"
                type="text"
                value={values.level}
                onChange={(event) => updateField('level', event.target.value)}
                placeholder="basic, advanced, professional"
              />
            </div>
          </div>

          {submitError ? (
            <div className="form-alert" role="alert">
              {submitError.message || 'Member could not be created.'}
              {submitError.code ? <span>{submitError.code}</span> : null}
            </div>
          ) : null}

          {status === 'success' ? (
            <div className="form-success" role="status" aria-live="polite">
              Member created, but the API response did not include an ID.
            </div>
          ) : null}

          <div className="member-form__actions">
            <button className="btn-primary" type="submit" disabled={status === 'submitting'}>
              {status === 'submitting' ? 'Creating member' : 'Create member'}
            </button>
            <button className="btn-outline" type="button" onClick={() => navigate('/app/members')}>
              Cancel
            </button>
          </div>
        </form>
      </DataPanel>
    </div>
  )
}

export default MemberCreateView
