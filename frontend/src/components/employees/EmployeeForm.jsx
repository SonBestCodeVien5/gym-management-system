import RoleSelector from './RoleSelector.jsx'
import {
  LEVEL_OPTIONS,
  STATUS_OPTIONS,
  cleanEmployeePayload,
  isObjectId,
} from './employeeFormatters.js'

export const EMPTY_EMPLOYEE_VALUES = {
  employee_id: '',
  full_name: '',
  email: '',
  password: '',
  role: [],
  level: '',
  phone: '',
  branch_id: '',
  status: 'active',
}

export function validateEmployee(values, { requirePassword = false } = {}) {
  const errors = {}

  if (!values.employee_id.trim()) {
    errors.employee_id = 'Employee ID is required.'
  }

  if (!values.full_name.trim()) {
    errors.full_name = 'Full name is required.'
  }

  if (!values.email.trim()) {
    errors.email = 'Email is required.'
  }

  if (requirePassword && values.password.length < 8) {
    errors.password = 'Password must be at least 8 characters.'
  }

  if (!values.role.length) {
    errors.role = 'Select at least one role.'
  }

  if (values.role.includes('trainer') && !values.level.trim()) {
    errors.level = 'Trainer role requires a level.'
  }

  const branchIds = values.branch_id
    .split(/[\n,]/)
    .map((branchId) => branchId.trim())
    .filter(Boolean)

  if (branchIds.some((branchId) => !isObjectId(branchId))) {
    errors.branch_id = 'Each branch ID must be a 24 character ObjectID.'
  }

  return errors
}

function EmployeeForm({
  values,
  setValues,
  errors,
  setErrors,
  branches = [],
  onSubmit,
  submitLabel,
  submittingLabel,
  status,
  requirePassword = false,
  onCancel,
}) {
  function updateField(name, value) {
    setValues((current) => ({ ...current, [name]: value }))
    setErrors((current) => ({ ...current, [name]: '' }))
  }

  function handleSubmit(event) {
    event.preventDefault()
    const nextErrors = validateEmployee(values, { requirePassword })
    setErrors(nextErrors)

    if (Object.keys(nextErrors).length) {
      return
    }

    onSubmit(cleanEmployeePayload(values, { includePassword: requirePassword }))
  }

  return (
    <form className="resource-form" onSubmit={handleSubmit}>
      <div className="resource-form__grid">
        <div className="field-group">
          <label htmlFor="employee-code">Employee ID</label>
          <input id="employee-code" value={values.employee_id} onChange={(event) => updateField('employee_id', event.target.value)} />
          {errors.employee_id ? <span>{errors.employee_id}</span> : null}
        </div>
        <div className="field-group">
          <label htmlFor="employee-name">Full name</label>
          <input id="employee-name" value={values.full_name} onChange={(event) => updateField('full_name', event.target.value)} />
          {errors.full_name ? <span>{errors.full_name}</span> : null}
        </div>
        <div className="field-group">
          <label htmlFor="employee-email">Email</label>
          <input id="employee-email" type="email" value={values.email} onChange={(event) => updateField('email', event.target.value)} />
          {errors.email ? <span>{errors.email}</span> : null}
        </div>
        {requirePassword ? (
          <div className="field-group">
            <label htmlFor="employee-password">Initial password</label>
            <input id="employee-password" type="password" value={values.password} onChange={(event) => updateField('password', event.target.value)} />
            {errors.password ? <span>{errors.password}</span> : null}
          </div>
        ) : null}
        <div className="field-group">
          <label htmlFor="employee-phone">Phone</label>
          <input id="employee-phone" value={values.phone} onChange={(event) => updateField('phone', event.target.value)} />
        </div>
        <div className="field-group">
          <label htmlFor="employee-status">Status</label>
          <select id="employee-status" value={values.status} onChange={(event) => updateField('status', event.target.value)}>
            {STATUS_OPTIONS.map((statusOption) => <option key={statusOption} value={statusOption}>{statusOption}</option>)}
          </select>
        </div>
        <div className="field-group">
          <label htmlFor="employee-level">Trainer level</label>
          <select id="employee-level" value={values.level} onChange={(event) => updateField('level', event.target.value)}>
            <option value="">No level</option>
            {LEVEL_OPTIONS.map((level) => <option key={level} value={level}>{level}</option>)}
          </select>
          {errors.level ? <span>{errors.level}</span> : null}
        </div>
        <div className="field-group field-group--wide">
          <label htmlFor="employee-branches">Branch IDs</label>
          <input
            id="employee-branches"
            list="employee-branch-options"
            value={values.branch_id}
            onChange={(event) => updateField('branch_id', event.target.value)}
            placeholder="One or more ObjectIDs, separated by commas"
          />
          <datalist id="employee-branch-options">
            {branches.map((branch) => <option key={branch.id} value={branch.id}>{branch.name} · {branch.branch_code}</option>)}
          </datalist>
          {errors.branch_id ? <span>{errors.branch_id}</span> : null}
        </div>
        <div className="field-group field-group--wide">
          <RoleSelector value={values.role} onChange={(nextRoles) => updateField('role', nextRoles)} error={errors.role} />
        </div>
      </div>

      <div className="resource-form__actions">
        <button className="btn-primary" type="submit" disabled={status === 'submitting'}>{status === 'submitting' ? submittingLabel : submitLabel}</button>
        {onCancel ? <button className="btn-outline" type="button" onClick={onCancel}>Cancel</button> : null}
      </div>
    </form>
  )
}

export default EmployeeForm
