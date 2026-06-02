import { useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import { resetEmployeePassword } from '../../lib/employeesApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'

function PasswordResetPanel({ accessToken, employeeId }) {
  const [values, setValues] = useState({ password: '', confirm: '' })
  const [error, setError] = useState('')
  const [state, setState] = useState({ status: 'idle', error: null, notice: '' })

  async function handleSubmit(event) {
    event.preventDefault()

    if (values.password.length < 8) {
      setError('Password must be at least 8 characters.')
      return
    }

    if (values.password !== values.confirm) {
      setError('Confirmation must match.')
      return
    }

    if (!window.confirm('Reset this employee password and revoke active refresh tokens?')) {
      return
    }

    setState({ status: 'submitting', error: null, notice: '' })

    try {
      await resetEmployeePassword(accessToken, employeeId, values.password)
      setValues({ password: '', confirm: '' })
      setState({ status: 'success', error: null, notice: 'Password reset completed. Active refresh tokens were revoked by the backend.' })
    } catch (mutationError) {
      setState({ status: 'error', error: mutationError, notice: '' })
    }
  }

  return (
    <DataPanel title="Password reset" description="The entered password is only sent to the API and is not stored in frontend state after success.">
      <form className="resource-form" onSubmit={handleSubmit}>
        <div className="resource-form__grid">
          <div className="field-group">
            <label htmlFor="reset-password">New password</label>
            <input id="reset-password" type="password" value={values.password} onChange={(event) => { setValues((current) => ({ ...current, password: event.target.value })); setError('') }} />
          </div>
          <div className="field-group">
            <label htmlFor="reset-confirm">Confirm password</label>
            <input id="reset-confirm" type="password" value={values.confirm} onChange={(event) => { setValues((current) => ({ ...current, confirm: event.target.value })); setError('') }} />
          </div>
        </div>
        {error ? <div className="form-alert" role="alert">{error}</div> : null}
        {state.error ? <div className="form-alert" role="alert">{apiErrorText(state.error, 'Password reset failed.')}</div> : null}
        {state.notice ? <div className="form-success" role="status">{state.notice}</div> : null}
        <div className="resource-form__actions">
          <button className="btn-primary" type="submit" disabled={state.status === 'submitting'}>{state.status === 'submitting' ? 'Resetting' : 'Reset password'}</button>
        </div>
      </form>
    </DataPanel>
  )
}

export default PasswordResetPanel
