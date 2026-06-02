import { useState } from 'react'
import { nearbyBranches } from '../../lib/branchesApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import { compactId, formatCoordinates } from './settingsFormatters.js'

const INITIAL_VALUES = {
  lng: '',
  lat: '',
  max_distance: '5000',
  limit: '10',
}

function validate(values) {
  const errors = {}
  const lngText = values.lng.toString().trim()
  const latText = values.lat.toString().trim()
  const lng = Number(values.lng)
  const lat = Number(values.lat)
  const maxDistance = Number(values.max_distance)
  const limit = Number(values.limit)

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

  if (values.max_distance && (!Number.isInteger(maxDistance) || maxDistance <= 0)) {
    errors.max_distance = 'Distance must be a positive integer.'
  }

  if (values.limit && (!Number.isInteger(limit) || limit < 1 || limit > 100)) {
    errors.limit = 'Limit must be from 1 to 100.'
  }

  return errors
}

function NearbyBranchesPanel({ accessToken, navigate }) {
  const [values, setValues] = useState(INITIAL_VALUES)
  const [errors, setErrors] = useState({})
  const [state, setState] = useState({ status: 'idle', data: [], error: null })

  function updateField(name, value) {
    setValues((current) => ({ ...current, [name]: value }))
    setErrors((current) => ({ ...current, [name]: '' }))
    setState((current) => ({ ...current, error: null }))
  }

  function fieldErrorProps(name) {
    return errors[name]
      ? {
        'aria-invalid': 'true',
        'aria-describedby': `nearby-${name.replace(/_/g, '-')}-error`,
      }
      : {}
  }

  async function handleSubmit(event) {
    event.preventDefault()
    const nextErrors = validate(values)
    setErrors(nextErrors)

    if (Object.keys(nextErrors).length) {
      return
    }

    setState({ status: 'loading', data: [], error: null })

    try {
      const response = await nearbyBranches(accessToken, {
        lng: values.lng,
        lat: values.lat,
        max_distance: values.max_distance,
        limit: values.limit,
      })
      setState({ status: 'success', data: response.data || [], error: null })
    } catch (error) {
      setState({ status: 'error', data: [], error })
    }
  }

  return (
    <div className="resource-stack">
      <form className="resource-form" onSubmit={handleSubmit}>
        <div className="resource-form__grid resource-form__grid--compact">
          <div className="field-group">
            <label htmlFor="nearby-lng">Longitude</label>
            <input id="nearby-lng" type="number" step="any" value={values.lng} onChange={(event) => updateField('lng', event.target.value)} {...fieldErrorProps('lng')} />
            {errors.lng ? <span id="nearby-lng-error">{errors.lng}</span> : null}
          </div>
          <div className="field-group">
            <label htmlFor="nearby-lat">Latitude</label>
            <input id="nearby-lat" type="number" step="any" value={values.lat} onChange={(event) => updateField('lat', event.target.value)} {...fieldErrorProps('lat')} />
            {errors.lat ? <span id="nearby-lat-error">{errors.lat}</span> : null}
          </div>
          <div className="field-group">
            <label htmlFor="nearby-distance">Max distance</label>
            <input id="nearby-distance" type="number" min="1" value={values.max_distance} onChange={(event) => updateField('max_distance', event.target.value)} {...fieldErrorProps('max_distance')} />
            {errors.max_distance ? <span id="nearby-max-distance-error">{errors.max_distance}</span> : null}
          </div>
          <div className="field-group">
            <label htmlFor="nearby-limit">Limit</label>
            <input id="nearby-limit" type="number" min="1" max="100" value={values.limit} onChange={(event) => updateField('limit', event.target.value)} {...fieldErrorProps('limit')} />
            {errors.limit ? <span id="nearby-limit-error">{errors.limit}</span> : null}
          </div>
        </div>
        <div className="resource-form__actions">
          <button className="btn-primary" type="submit" disabled={state.status === 'loading'}>
            {state.status === 'loading' ? 'Searching' : 'Search nearby'}
          </button>
        </div>
      </form>

      {state.status === 'error' ? (
        <div className="form-alert" role="alert">{apiErrorText(state.error, 'Nearby search failed.')}</div>
      ) : null}

      {state.status === 'success' ? (
        state.data.length ? (
          <div className="resource-list">
            {state.data.map((branch) => (
              <article className="resource-row" key={branch.id}>
                <div>
                  <strong>{branch.name}</strong>
                  <span>{branch.branch_code} · {formatCoordinates(branch.location)}</span>
                </div>
                <div>
                  <span>{branch.distance_meters ? `${Math.round(branch.distance_meters)} m` : compactId(branch.id)}</span>
                  <button className="btn-outline" type="button" onClick={() => navigate(`/app/settings/branches/${branch.id}`)}>
                    Open
                  </button>
                </div>
              </article>
            ))}
          </div>
        ) : (
          <div className="form-success" role="status">No nearby branches matched this query.</div>
        )
      ) : null}
    </div>
  )
}

export default NearbyBranchesPanel
