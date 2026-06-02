import { useCallback, useEffect, useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import PageHeader from '../PageHeader.jsx'
import { useAuth } from '../../context/AuthContext.jsx'
import { listBranches } from '../../lib/branchesApi.js'
import { listCourses } from '../../lib/coursesApi.js'
import { createSession } from '../../lib/sessionsApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import { cleanSessionPayload, isObjectId } from './sessionFormatters.js'

const INITIAL_VALUES = {
  branch_id: '',
  trainer_id: '',
  course_level: '',
  scheduled_at: '',
  duration_min: '60',
  capacity: '10',
  tags: '',
}

function validate(values) {
  const errors = {}
  const duration = Number(values.duration_min)
  const capacity = Number(values.capacity)

  if (!isObjectId(values.branch_id)) {
    errors.branch_id = 'Branch ID must be a 24 character ObjectID.'
  }

  if (!isObjectId(values.trainer_id)) {
    errors.trainer_id = 'Trainer ID must be a 24 character ObjectID.'
  }

  if (!values.course_level.trim()) {
    errors.course_level = 'Course level is required.'
  }

  if (!values.scheduled_at || Number.isNaN(new Date(values.scheduled_at).getTime())) {
    errors.scheduled_at = 'Scheduled date is required.'
  }

  if (!Number.isInteger(duration) || duration <= 0) {
    errors.duration_min = 'Duration must be a positive integer.'
  }

  if (!Number.isInteger(capacity) || capacity <= 0) {
    errors.capacity = 'Capacity must be a positive integer.'
  }

  return errors
}

function SessionCreateView({ navigate }) {
  const { accessToken, employee } = useAuth()
  const [values, setValues] = useState(INITIAL_VALUES)
  const [errors, setErrors] = useState({})
  const [refs, setRefs] = useState({ status: 'loading', branches: [], courses: [], error: null })
  const [submitState, setSubmitState] = useState({ status: 'idle', error: null })

  useEffect(() => {
    if (!employee?.role?.includes('trainer') || !employee.id) {
      return
    }

    setValues((current) => (
      current.trainer_id ? current : { ...current, trainer_id: employee.id }
    ))
  }, [employee])

  const loadRefs = useCallback(async () => {
    try {
      const [branchesResponse, coursesResponse] = await Promise.all([
        listBranches(accessToken),
        listCourses(accessToken),
      ])
      setRefs({
        status: 'success',
        branches: branchesResponse.data || [],
        courses: coursesResponse.data || [],
        error: null,
      })
    } catch (error) {
      setRefs({ status: 'error', branches: [], courses: [], error })
    }
  }, [accessToken])

  useEffect(() => {
    loadRefs()
  }, [loadRefs])

  function updateField(name, value) {
    setValues((current) => ({ ...current, [name]: value }))
    setErrors((current) => ({ ...current, [name]: '' }))
    setSubmitState((current) => ({ ...current, error: null }))
  }

  async function handleSubmit(event) {
    event.preventDefault()
    const nextErrors = validate(values)
    setErrors(nextErrors)

    if (Object.keys(nextErrors).length) {
      return
    }

    setSubmitState({ status: 'submitting', error: null })

    try {
      const response = await createSession(accessToken, cleanSessionPayload(values))
      const sessionId = response.data?.id

      if (sessionId) {
        navigate(`/app/sessions/${sessionId}`, { replace: true })
        return
      }

      setSubmitState({ status: 'success', error: null })
    } catch (error) {
      setSubmitState({ status: 'error', error })
    }
  }

  const levels = Array.from(new Set(refs.courses.map((course) => course.level).filter(Boolean)))

  return (
    <div className="module-page resource-workspace sessions-workspace">
      <PageHeader
        eyebrow="Sessions"
        title="New session"
        description="Create a scheduled class using live branch data and manual trainer ObjectID."
        actions={<button className="btn-outline" type="button" onClick={() => navigate('/app/sessions')}>Sessions</button>}
      />

      {refs.status === 'error' ? <div className="form-alert" role="alert">Reference lists failed to load. Manual ObjectID entry is still available. {apiErrorText(refs.error)}</div> : null}

      <DataPanel title="Session details">
        <form className="resource-form" onSubmit={handleSubmit}>
          <div className="resource-form__grid">
            <div className="field-group">
              <label htmlFor="session-branch">Branch ID</label>
              <input id="session-branch" list="session-create-branches" value={values.branch_id} onChange={(event) => updateField('branch_id', event.target.value)} />
              <datalist id="session-create-branches">
                {refs.branches.map((branch) => <option key={branch.id} value={branch.id}>{branch.name} · {branch.branch_code}</option>)}
              </datalist>
              {errors.branch_id ? <span>{errors.branch_id}</span> : null}
            </div>
            <div className="field-group">
              <label htmlFor="session-trainer">Trainer ID</label>
              <input id="session-trainer" value={values.trainer_id} onChange={(event) => updateField('trainer_id', event.target.value)} />
              {errors.trainer_id ? <span>{errors.trainer_id}</span> : null}
            </div>
            <div className="field-group">
              <label htmlFor="session-level">Course level</label>
              <input id="session-level" list="session-course-levels" value={values.course_level} onChange={(event) => updateField('course_level', event.target.value)} />
              <datalist id="session-course-levels">
                {levels.map((level) => <option key={level} value={level} />)}
              </datalist>
              {errors.course_level ? <span>{errors.course_level}</span> : null}
            </div>
            <div className="field-group">
              <label htmlFor="session-scheduled">Scheduled at</label>
              <input id="session-scheduled" type="datetime-local" value={values.scheduled_at} onChange={(event) => updateField('scheduled_at', event.target.value)} />
              {errors.scheduled_at ? <span>{errors.scheduled_at}</span> : null}
            </div>
            <div className="field-group">
              <label htmlFor="session-duration">Duration minutes</label>
              <input id="session-duration" type="number" min="1" value={values.duration_min} onChange={(event) => updateField('duration_min', event.target.value)} />
              {errors.duration_min ? <span>{errors.duration_min}</span> : null}
            </div>
            <div className="field-group">
              <label htmlFor="session-capacity">Capacity</label>
              <input id="session-capacity" type="number" min="1" value={values.capacity} onChange={(event) => updateField('capacity', event.target.value)} />
              {errors.capacity ? <span>{errors.capacity}</span> : null}
            </div>
            <div className="field-group field-group--wide">
              <label htmlFor="session-tags">Tags</label>
              <textarea id="session-tags" rows="3" value={values.tags} onChange={(event) => updateField('tags', event.target.value)} />
            </div>
          </div>
          {submitState.error ? <div className="form-alert" role="alert">{apiErrorText(submitState.error, 'Session could not be created.')}</div> : null}
          {submitState.status === 'success' ? <div className="form-success" role="status">Session created, but response did not include an ID.</div> : null}
          <div className="resource-form__actions">
            <button className="btn-primary" type="submit" disabled={submitState.status === 'submitting'}>{submitState.status === 'submitting' ? 'Creating' : 'Create session'}</button>
            <button className="btn-outline" type="button" onClick={() => navigate('/app/sessions')}>Cancel</button>
          </div>
        </form>
      </DataPanel>
    </div>
  )
}

export default SessionCreateView
