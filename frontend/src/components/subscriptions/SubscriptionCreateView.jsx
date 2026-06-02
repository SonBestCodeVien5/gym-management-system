import { useCallback, useEffect, useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import PageHeader from '../PageHeader.jsx'
import { useAuth } from '../../context/AuthContext.jsx'
import { listBranches } from '../../lib/branchesApi.js'
import { listCourses } from '../../lib/coursesApi.js'
import { createSubscription } from '../../lib/subscriptionsApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import {
  cleanSubscriptionPayload,
  formatMoney,
  isObjectId,
} from './subscriptionFormatters.js'

const INITIAL_VALUES = {
  member_id: '',
  course_id: '',
  home_branch_id: '',
  start_date: '',
  end_date: '',
  session_per_week: '3',
  discount_type: 'none',
  discount_value: '0',
  promo_code: '',
}

function validate(values) {
  const errors = {}
  const start = new Date(values.start_date)
  const end = new Date(values.end_date)
  const weekly = Number(values.session_per_week)
  const discountValue = Number(values.discount_value || 0)

  ;['member_id', 'course_id', 'home_branch_id'].forEach((field) => {
    if (!isObjectId(values[field])) {
      errors[field] = 'Must be a 24 character ObjectID.'
    }
  })

  if (!values.start_date || Number.isNaN(start.getTime())) {
    errors.start_date = 'Start date is required.'
  }

  if (!values.end_date || Number.isNaN(end.getTime())) {
    errors.end_date = 'End date is required.'
  }

  if (values.start_date && values.end_date && start > end) {
    errors.end_date = 'End date must be after start date.'
  }

  if (!Number.isInteger(weekly) || weekly <= 0) {
    errors.session_per_week = 'Weekly sessions must be a positive integer.'
  }

  if (!Number.isFinite(discountValue) || discountValue < 0) {
    errors.discount_value = 'Discount value must be zero or more.'
  }

  return errors
}

function SubscriptionCreateView({ navigate }) {
  const { accessToken } = useAuth()
  const [values, setValues] = useState(INITIAL_VALUES)
  const [errors, setErrors] = useState({})
  const [refs, setRefs] = useState({ status: 'loading', courses: [], branches: [], error: null })
  const [submitState, setSubmitState] = useState({ status: 'idle', error: null })

  const loadRefs = useCallback(async () => {
    setRefs((current) => ({ ...current, status: 'loading', error: null }))

    try {
      const [coursesResponse, branchesResponse] = await Promise.all([
        listCourses(accessToken),
        listBranches(accessToken),
      ])
      setRefs({
        status: 'success',
        courses: coursesResponse.data || [],
        branches: branchesResponse.data || [],
        error: null,
      })
    } catch (error) {
      setRefs({ status: 'error', courses: [], branches: [], error })
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
      const response = await createSubscription(accessToken, cleanSubscriptionPayload(values))
      const subscriptionId = response.data?.id

      if (subscriptionId) {
        navigate(`/app/subscriptions/${subscriptionId}`, { replace: true })
        return
      }

      setSubmitState({ status: 'success', error: null })
    } catch (error) {
      setSubmitState({ status: 'error', error })
    }
  }

  return (
    <div className="module-page resource-workspace subscriptions-workspace">
      <PageHeader
        eyebrow="Subscriptions"
        title="New subscription"
        description="Create a pending subscription. Offline payment activation remains on member detail."
        actions={<button className="btn-outline" type="button" onClick={() => navigate('/app/subscriptions')}>Subscriptions</button>}
      />

      {refs.status === 'error' ? (
        <div className="form-alert" role="alert">
          Reference lists failed to load. Manual ObjectID entry is still available. {apiErrorText(refs.error)}
        </div>
      ) : null}

      <DataPanel title="Subscription details" description="Dates are converted to RFC3339 before sending to the API.">
        <form className="resource-form" onSubmit={handleSubmit}>
          <div className="resource-form__grid">
            <div className="field-group">
              <label htmlFor="sub-member">Member ID</label>
              <input id="sub-member" value={values.member_id} onChange={(event) => updateField('member_id', event.target.value)} />
              {errors.member_id ? <span>{errors.member_id}</span> : null}
            </div>

            <div className="field-group">
              <label htmlFor="sub-course">Course ID</label>
              <input id="sub-course" list="course-options" value={values.course_id} onChange={(event) => updateField('course_id', event.target.value)} />
              <datalist id="course-options">
                {refs.courses.map((course) => <option key={course.id} value={course.id}>{course.title} · {formatMoney(course.base_price)}</option>)}
              </datalist>
              {errors.course_id ? <span>{errors.course_id}</span> : null}
            </div>

            <div className="field-group">
              <label htmlFor="sub-branch">Home branch ID</label>
              <input id="sub-branch" list="branch-options" value={values.home_branch_id} onChange={(event) => updateField('home_branch_id', event.target.value)} />
              <datalist id="branch-options">
                {refs.branches.map((branch) => <option key={branch.id} value={branch.id}>{branch.name} · {branch.branch_code}</option>)}
              </datalist>
              {errors.home_branch_id ? <span>{errors.home_branch_id}</span> : null}
            </div>

            <div className="field-group">
              <label htmlFor="sub-weekly">Sessions per week</label>
              <input id="sub-weekly" type="number" min="1" value={values.session_per_week} onChange={(event) => updateField('session_per_week', event.target.value)} />
              {errors.session_per_week ? <span>{errors.session_per_week}</span> : null}
            </div>

            <div className="field-group">
              <label htmlFor="sub-start">Start date</label>
              <input id="sub-start" type="datetime-local" value={values.start_date} onChange={(event) => updateField('start_date', event.target.value)} />
              {errors.start_date ? <span>{errors.start_date}</span> : null}
            </div>

            <div className="field-group">
              <label htmlFor="sub-end">End date</label>
              <input id="sub-end" type="datetime-local" value={values.end_date} onChange={(event) => updateField('end_date', event.target.value)} />
              {errors.end_date ? <span>{errors.end_date}</span> : null}
            </div>

            <div className="field-group">
              <label htmlFor="discount-type">Discount type</label>
              <select id="discount-type" value={values.discount_type} onChange={(event) => updateField('discount_type', event.target.value)}>
                <option value="none">none</option>
                <option value="percent">percent</option>
                <option value="fixed">fixed</option>
              </select>
            </div>

            <div className="field-group">
              <label htmlFor="discount-value">Discount value</label>
              <input id="discount-value" type="number" min="0" value={values.discount_value} onChange={(event) => updateField('discount_value', event.target.value)} disabled={values.discount_type === 'none'} />
              {errors.discount_value ? <span>{errors.discount_value}</span> : null}
            </div>

            <div className="field-group">
              <label htmlFor="promo-code">Promo code</label>
              <input id="promo-code" value={values.promo_code} onChange={(event) => updateField('promo_code', event.target.value)} disabled={values.discount_type === 'none'} />
            </div>
          </div>

          {submitState.error ? <div className="form-alert" role="alert">{apiErrorText(submitState.error, 'Subscription could not be created.')}</div> : null}
          {submitState.status === 'success' ? <div className="form-success" role="status">Subscription created, but response did not include an ID.</div> : null}

          <div className="resource-form__actions">
            <button className="btn-primary" type="submit" disabled={submitState.status === 'submitting'}>
              {submitState.status === 'submitting' ? 'Creating' : 'Create subscription'}
            </button>
            <button className="btn-outline" type="button" onClick={() => navigate('/app/subscriptions')}>Cancel</button>
          </div>
        </form>
      </DataPanel>
    </div>
  )
}

export default SubscriptionCreateView
