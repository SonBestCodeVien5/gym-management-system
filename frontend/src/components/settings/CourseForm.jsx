import { cleanCoursePayload } from './settingsFormatters.js'

const INITIAL_VALUES = {
  title: '',
  level: '',
  allowed_tags: '',
  base_price: '',
  session_count: '',
  description: '',
}

export function courseValuesFromCourse(course) {
  return {
    title: course?.title || '',
    level: course?.level || '',
    allowed_tags: Array.isArray(course?.allowed_tags) ? course.allowed_tags.join(', ') : '',
    base_price: course?.base_price ?? '',
    session_count: course?.session_count ?? '',
    description: course?.description || '',
  }
}

function validate(values) {
  const errors = {}
  const price = Number(values.base_price)
  const sessions = Number(values.session_count)

  if (!values.title.trim()) {
    errors.title = 'Title is required.'
  }

  if (!values.level.trim()) {
    errors.level = 'Level is required.'
  }

  if (!Number.isInteger(price) || price <= 0) {
    errors.base_price = 'Base price must be a positive integer.'
  }

  if (!Number.isInteger(sessions) || sessions <= 0) {
    errors.session_count = 'Session count must be a positive integer.'
  }

  return errors
}

function CourseForm({
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

    onSubmit(cleanCoursePayload(values))
  }

  return (
    <form className="resource-form" onSubmit={handleSubmit}>
      <div className="resource-form__grid">
        <div className="field-group">
          <label htmlFor="course-title">Title</label>
          <input
            id="course-title"
            value={values.title}
            onChange={(event) => updateField('title', event.target.value)}
            aria-invalid={errors.title ? 'true' : undefined}
            aria-describedby={errors.title ? 'course-title-error' : undefined}
          />
          {errors.title ? <span id="course-title-error">{errors.title}</span> : null}
        </div>

        <div className="field-group">
          <label htmlFor="course-level">Level</label>
          <input
            id="course-level"
            value={values.level}
            onChange={(event) => updateField('level', event.target.value)}
            aria-invalid={errors.level ? 'true' : undefined}
            aria-describedby={errors.level ? 'course-level-error' : undefined}
          />
          {errors.level ? <span id="course-level-error">{errors.level}</span> : null}
        </div>

        <div className="field-group">
          <label htmlFor="course-price">Base price</label>
          <input
            id="course-price"
            type="number"
            min="1"
            value={values.base_price}
            onChange={(event) => updateField('base_price', event.target.value)}
            aria-invalid={errors.base_price ? 'true' : undefined}
            aria-describedby={errors.base_price ? 'course-price-error' : undefined}
          />
          {errors.base_price ? <span id="course-price-error">{errors.base_price}</span> : null}
        </div>

        <div className="field-group">
          <label htmlFor="course-sessions">Session count</label>
          <input
            id="course-sessions"
            type="number"
            min="1"
            value={values.session_count}
            onChange={(event) => updateField('session_count', event.target.value)}
            aria-invalid={errors.session_count ? 'true' : undefined}
            aria-describedby={errors.session_count ? 'course-sessions-error' : undefined}
          />
          {errors.session_count ? <span id="course-sessions-error">{errors.session_count}</span> : null}
        </div>

        <div className="field-group field-group--wide">
          <label htmlFor="course-tags">Allowed tags</label>
          <textarea
            id="course-tags"
            rows="3"
            value={values.allowed_tags}
            onChange={(event) => updateField('allowed_tags', event.target.value)}
            placeholder="strength, morning, beginner"
          />
        </div>

        <div className="field-group field-group--wide">
          <label htmlFor="course-description">Description</label>
          <textarea
            id="course-description"
            rows="3"
            value={values.description}
            onChange={(event) => updateField('description', event.target.value)}
          />
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

CourseForm.initialValues = INITIAL_VALUES

export default CourseForm
