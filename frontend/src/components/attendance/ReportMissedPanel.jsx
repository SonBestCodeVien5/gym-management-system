import { useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import { reportMissedAttendance } from '../../lib/attendanceApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import { cleanAttendancePayload, isObjectId } from './attendanceFormatters.js'

const INITIAL_VALUES = { subscription_id: '', branch_id: '', date: '' }

function validate(values) {
  const errors = {}

  if (!isObjectId(values.subscription_id)) {
    errors.subscription_id = 'Subscription ID must be a 24 character ObjectID.'
  }

  if (!isObjectId(values.branch_id)) {
    errors.branch_id = 'Branch ID must be a 24 character ObjectID.'
  }

  return errors
}

function ReportMissedPanel({ accessToken, branches = [], subscriptionId = '', onSuccess, branchListId = 'report-branch-options' }) {
  const [values, setValues] = useState({ ...INITIAL_VALUES, subscription_id: subscriptionId })
  const [errors, setErrors] = useState({})
  const [state, setState] = useState({ status: 'idle', error: null, notice: '' })

  function updateField(name, value) {
    setValues((current) => ({ ...current, [name]: value }))
    setErrors((current) => ({ ...current, [name]: '' }))
    setState((current) => ({ ...current, error: null, notice: '' }))
  }

  async function handleSubmit(event) {
    event.preventDefault()
    const nextErrors = validate(values)
    setErrors(nextErrors)

    if (Object.keys(nextErrors).length) {
      return
    }

    setState({ status: 'submitting', error: null, notice: '' })

    try {
      await reportMissedAttendance(accessToken, cleanAttendancePayload(values))
      setState({ status: 'success', error: null, notice: 'Missed attendance reported.' })
      onSuccess?.()
    } catch (error) {
      setState({ status: 'error', error, notice: '' })
    }
  }

  return (
    <DataPanel title="Report missed">
      <form className="resource-form" onSubmit={handleSubmit}>
        <div className="resource-form__grid">
          <div className="field-group">
            <label htmlFor="report-sub">Subscription ID</label>
            <input id="report-sub" value={values.subscription_id} onChange={(event) => updateField('subscription_id', event.target.value)} />
            {errors.subscription_id ? <span>{errors.subscription_id}</span> : null}
          </div>
          <div className="field-group">
            <label htmlFor="report-branch">Branch ID</label>
            <input id="report-branch" list={branchListId} value={values.branch_id} onChange={(event) => updateField('branch_id', event.target.value)} />
            {errors.branch_id ? <span>{errors.branch_id}</span> : null}
          </div>
          <div className="field-group">
            <label htmlFor="report-date">Missed date</label>
            <input id="report-date" type="datetime-local" value={values.date} onChange={(event) => updateField('date', event.target.value)} />
          </div>
        </div>
        <datalist id={branchListId}>
          {branches.map((branch) => <option key={branch.id} value={branch.id}>{branch.name} · {branch.branch_code}</option>)}
        </datalist>
        <div className="resource-form__actions">
          <button className="btn-primary" type="submit" disabled={state.status === 'submitting'}>{state.status === 'submitting' ? 'Reporting' : 'Report missed'}</button>
        </div>
      </form>
      {state.error ? <div className="form-alert" role="alert">{apiErrorText(state.error, 'Report failed.')}</div> : null}
      {state.notice ? <div className="form-success" role="status">{state.notice}</div> : null}
    </DataPanel>
  )
}

export default ReportMissedPanel
