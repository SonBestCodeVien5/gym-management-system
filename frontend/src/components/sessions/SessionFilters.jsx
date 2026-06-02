import { useState } from 'react'
import { dateTimeLocalToRfc3339, isObjectId } from './sessionFormatters.js'

function SessionFilters({ branches = [], onApply }) {
  const [values, setValues] = useState({ branchId: '', level: '', date: '' })
  const [error, setError] = useState('')

  function handleSubmit(event) {
    event.preventDefault()

    if (values.branchId.trim() && !isObjectId(values.branchId)) {
      setError('Branch ID must be a 24 character ObjectID.')
      return
    }

    setError('')
    onApply({
      branchId: values.branchId.trim(),
      level: values.level.trim(),
      date: dateTimeLocalToRfc3339(values.date),
    })
  }

  return (
    <form className="resource-form" onSubmit={handleSubmit}>
      <div className="resource-form__grid resource-form__grid--compact">
        <div className="field-group">
          <label htmlFor="session-filter-branch">Branch ID</label>
          <input id="session-filter-branch" list="session-branch-options" value={values.branchId} onChange={(event) => setValues((current) => ({ ...current, branchId: event.target.value }))} />
        </div>
        <div className="field-group">
          <label htmlFor="session-filter-level">Level</label>
          <input id="session-filter-level" value={values.level} onChange={(event) => setValues((current) => ({ ...current, level: event.target.value }))} />
        </div>
        <div className="field-group">
          <label htmlFor="session-filter-date">Date</label>
          <input id="session-filter-date" type="datetime-local" value={values.date} onChange={(event) => setValues((current) => ({ ...current, date: event.target.value }))} />
        </div>
      </div>
      <datalist id="session-branch-options">
        {branches.map((branch) => <option key={branch.id} value={branch.id}>{branch.name} · {branch.branch_code}</option>)}
      </datalist>
      {error ? <div className="form-alert" role="alert">{error}</div> : null}
      <div className="resource-form__actions">
        <button className="btn-primary" type="submit">Apply filters</button>
        <button className="btn-outline" type="button" onClick={() => { setValues({ branchId: '', level: '', date: '' }); onApply({}) }}>Clear</button>
      </div>
    </form>
  )
}

export default SessionFilters
