import { cleanBranchPayload, isObjectId } from './settingsFormatters.js'

const INITIAL_VALUES = {
  branch_code: '',
  name: '',
  address: '',
  province: '',
  lng: '',
  lat: '',
  manager_id: '',
}

export function branchValuesFromBranch(branch) {
  const coordinates = branch?.location?.coordinates || []

  return {
    branch_code: branch?.branch_code || '',
    name: branch?.name || '',
    address: branch?.address || '',
    province: branch?.province || '',
    lng: coordinates[0] ?? '',
    lat: coordinates[1] ?? '',
    manager_id: branch?.manager_id && branch.manager_id !== '000000000000000000000000' ? branch.manager_id : '',
  }
}

function validate(values) {
  const errors = {}
  const lng = Number(values.lng)
  const lat = Number(values.lat)

  ;['branch_code', 'name', 'address', 'province'].forEach((field) => {
    if (!values[field].trim()) {
      errors[field] = 'This field is required.'
    }
  })

  if (!Number.isFinite(lng) || lng < -180 || lng > 180) {
    errors.lng = 'Longitude must be between -180 and 180.'
  }

  if (!Number.isFinite(lat) || lat < -90 || lat > 90) {
    errors.lat = 'Latitude must be between -90 and 90.'
  }

  if (values.manager_id.trim() && !isObjectId(values.manager_id)) {
    errors.manager_id = 'Manager ID must be a 24 character ObjectID.'
  }

  return errors
}

function BranchForm({
  values,
  setValues,
  errors,
  setErrors,
  onSubmit,
  submitLabel,
  submittingLabel,
  status,
  onCancel,
}) {
  function updateField(name, value) {
    setValues((current) => ({ ...current, [name]: value }))
    setErrors((current) => ({ ...current, [name]: '' }))
  }

  function handleSubmit(event) {
    event.preventDefault()
    const nextErrors = validate(values)
    setErrors(nextErrors)

    if (Object.keys(nextErrors).length) {
      return
    }

    onSubmit(cleanBranchPayload(values))
  }

  return (
    <form className="resource-form" onSubmit={handleSubmit}>
      <div className="resource-form__grid">
        <div className="field-group">
          <label htmlFor="branch-code">Branch code</label>
          <input id="branch-code" value={values.branch_code} onChange={(event) => updateField('branch_code', event.target.value)} />
          {errors.branch_code ? <span>{errors.branch_code}</span> : null}
        </div>

        <div className="field-group">
          <label htmlFor="branch-name">Name</label>
          <input id="branch-name" value={values.name} onChange={(event) => updateField('name', event.target.value)} />
          {errors.name ? <span>{errors.name}</span> : null}
        </div>

        <div className="field-group">
          <label htmlFor="branch-province">Province</label>
          <input id="branch-province" value={values.province} onChange={(event) => updateField('province', event.target.value)} />
          {errors.province ? <span>{errors.province}</span> : null}
        </div>

        <div className="field-group">
          <label htmlFor="branch-manager">Manager ID</label>
          <input
            id="branch-manager"
            value={values.manager_id}
            onChange={(event) => updateField('manager_id', event.target.value)}
            placeholder="Optional ObjectID"
            aria-invalid={errors.manager_id ? 'true' : undefined}
          />
          {errors.manager_id ? <span>{errors.manager_id}</span> : null}
        </div>

        <div className="field-group">
          <label htmlFor="branch-lng">Longitude</label>
          <input id="branch-lng" type="number" step="any" value={values.lng} onChange={(event) => updateField('lng', event.target.value)} />
          {errors.lng ? <span>{errors.lng}</span> : null}
        </div>

        <div className="field-group">
          <label htmlFor="branch-lat">Latitude</label>
          <input id="branch-lat" type="number" step="any" value={values.lat} onChange={(event) => updateField('lat', event.target.value)} />
          {errors.lat ? <span>{errors.lat}</span> : null}
        </div>

        <div className="field-group field-group--wide">
          <label htmlFor="branch-address">Address</label>
          <textarea id="branch-address" rows="3" value={values.address} onChange={(event) => updateField('address', event.target.value)} />
          {errors.address ? <span>{errors.address}</span> : null}
        </div>
      </div>

      <div className="resource-form__actions">
        <button className="btn-primary" type="submit" disabled={status === 'submitting'}>
          {status === 'submitting' ? submittingLabel : submitLabel}
        </button>
        {onCancel ? (
          <button className="btn-outline" type="button" onClick={onCancel}>
            Cancel
          </button>
        ) : null}
      </div>
    </form>
  )
}

BranchForm.initialValues = INITIAL_VALUES

export default BranchForm
