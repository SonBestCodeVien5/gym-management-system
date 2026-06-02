import { useState } from 'react'
import { ROLE_OPTIONS, STATUS_OPTIONS, isObjectId } from './employeeFormatters.js'

function EmployeeFilters({ branches = [], onApply }) {
  const [values, setValues] = useState({ role: '', status: '', branch_id: '' })
  const [error, setError] = useState('')

  function handleSubmit(event) {
    event.preventDefault()

    if (values.branch_id.trim() && !isObjectId(values.branch_id)) {
      setError('Branch ID must be a 24 character ObjectID.')
      return
    }

    setError('')
    onApply({
      role: values.role,
      status: values.status,
      branch_id: values.branch_id.trim(),
    })
  }

  return (
    <form className="resource-form" onSubmit={handleSubmit}>
      <div className="resource-form__grid resource-form__grid--compact">
        <div className="field-group">
          <label htmlFor="employee-role-filter">Role</label>
          <select id="employee-role-filter" value={values.role} onChange={(event) => setValues((current) => ({ ...current, role: event.target.value }))}>
            <option value="">Any role</option>
            {ROLE_OPTIONS.map((role) => <option key={role} value={role}>{role}</option>)}
          </select>
        </div>
        <div className="field-group">
          <label htmlFor="employee-status-filter">Status</label>
          <select id="employee-status-filter" value={values.status} onChange={(event) => setValues((current) => ({ ...current, status: event.target.value }))}>
            <option value="">Any status</option>
            {STATUS_OPTIONS.map((status) => <option key={status} value={status}>{status}</option>)}
          </select>
        </div>
        <div className="field-group">
          <label htmlFor="employee-branch-filter">Branch ID</label>
          <input id="employee-branch-filter" list="employee-filter-branches" value={values.branch_id} onChange={(event) => setValues((current) => ({ ...current, branch_id: event.target.value }))} />
          <datalist id="employee-filter-branches">
            {branches.map((branch) => <option key={branch.id} value={branch.id}>{branch.name} · {branch.branch_code}</option>)}
          </datalist>
        </div>
      </div>
      {error ? <div className="form-alert" role="alert">{error}</div> : null}
      <div className="resource-form__actions">
        <button className="btn-primary" type="submit">Apply filters</button>
        <button className="btn-outline" type="button" onClick={() => { setValues({ role: '', status: '', branch_id: '' }); onApply({}) }}>Clear</button>
      </div>
    </form>
  )
}

export default EmployeeFilters
