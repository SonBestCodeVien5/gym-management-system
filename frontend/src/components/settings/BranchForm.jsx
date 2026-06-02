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
  const lngText = values.lng.toString().trim()
  const latText = values.lat.toString().trim()
  const lng = Number(values.lng)
  const lat = Number(values.lat)

  ;['branch_code', 'name', 'address', 'province'].forEach((field) => {
    if (!values[field].trim()) {
      errors[field] = 'This field is required.'
    }
  })

  if (!lngText) {
    errors.lng = 'Longitude is required.'
  } else if (!Number.isFinite(lng) || lng < -180 || lng > 180) {
    errors.lng = 'Longitude must be between -180 and 180.'
  }

  if (!latText) {
    errors.lat = 'Latitude is required.'
  } else if (!Number.isFinite(lat) || lat < -90 || lat > 90) {
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
  function fieldErrorProps(name) {
    return errors[name]
      ? {
        'aria-invalid': 'true',
        'aria-describedby': `branch-${name.replace(/_/g, '-')}-error`,
      }
      : {}
  }

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
          <input id="branch-code" value={values.branch_code} onChange={(event) => updateField('branch_code', event.target.value)} {...fieldErrorProps('branch_code')} />
          {errors.branch_code ? <span id="branch-branch-code-error">{errors.branch_code}</span> : null}
        </div>

        <div className="field-group">
          <label htmlFor="branch-name">Name</label>
          <input id="branch-name" value={values.name} onChange={(event) => updateField('name', event.target.value)} {...fieldErrorProps('name')} />
          {errors.name ? <span id="branch-name-error">{errors.name}</span> : null}
        </div>

        <div className="field-group">
          <label htmlFor="branch-province">Province</label>
          <input id="branch-province" value={values.province} onChange={(event) => updateField('province', event.target.value)} {...fieldErrorProps('province')} />
          {errors.province ? <span id="branch-province-error">{errors.province}</span> : null}
        </div>

        <div className="field-group">
          <label htmlFor="branch-manager">Manager ID</label>
          <input
            id="branch-manager"
            value={values.manager_id}
            onChange={(event) => updateField('manager_id', event.target.value)}
            placeholder="Optional ObjectID"
            {...fieldErrorProps('manager_id')}
          />
          {errors.manager_id ? <span id="branch-manager-id-error">{errors.manager_id}</span> : null}
        </div>

        <div className="field-group">
          <label htmlFor="branch-lng">Longitude</label>
          <input id="branch-lng" type="number" step="any" value={values.lng} onChange={(event) => updateField('lng', event.target.value)} {...fieldErrorProps('lng')} />
          {errors.lng ? <span id="branch-lng-error">{errors.lng}</span> : null}
        </div>

        <div className="field-group">
          <label htmlFor="branch-lat">Latitude</label>
          <input id="branch-lat" type="number" step="any" value={values.lat} onChange={(event) => updateField('lat', event.target.value)} {...fieldErrorProps('lat')} />
          {errors.lat ? <span id="branch-lat-error">{errors.lat}</span> : null}
        </div>

        <div className="field-group field-group--wide">
          <label htmlFor="branch-address">Address</label>
          <textarea id="branch-address" rows="3" value={values.address} onChange={(event) => updateField('address', event.target.value)} {...fieldErrorProps('address')} />
          {errors.address ? <span id="branch-address-error">{errors.address}</span> : null}
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
