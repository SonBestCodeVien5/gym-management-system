import { useCallback, useEffect, useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import PageHeader from '../PageHeader.jsx'
import StateBlock from '../StateBlock.jsx'
import { useAuth } from '../../context/AuthContext.jsx'
import { deleteCourse, getCourse, updateCourse } from '../../lib/coursesApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import CourseForm, { courseValuesFromCourse } from './CourseForm.jsx'
import { compactId, formatMoney, formatTags, isObjectId } from './settingsFormatters.js'

function CourseDetailView({ courseId, navigate }) {
  const { accessToken } = useAuth()
  const [courseState, setCourseState] = useState({ status: 'loading', data: null, error: null })
  const [values, setValues] = useState(CourseForm.initialValues)
  const [errors, setErrors] = useState({})
  const [mutation, setMutation] = useState({ status: 'idle', error: null, notice: '' })

  const loadCourse = useCallback(async () => {
    if (!isObjectId(courseId)) {
      setCourseState({
        status: 'error',
        data: null,
        error: { code: 'INVALID_ID', message: 'Course ID must be a 24 character ObjectID.' },
      })
      return
    }

    setCourseState((current) => ({ ...current, status: 'loading', error: null }))

    try {
      const response = await getCourse(accessToken, courseId)
      setCourseState({ status: 'success', data: response.data, error: null })
      setValues(courseValuesFromCourse(response.data))
    } catch (error) {
      setCourseState({ status: 'error', data: null, error })
    }
  }, [accessToken, courseId])

  useEffect(() => {
    loadCourse()
  }, [loadCourse])

  async function handleUpdate(payload) {
    setMutation({ status: 'submitting', error: null, notice: '' })

    try {
      await updateCourse(accessToken, courseId, payload)
      setMutation({ status: 'success', error: null, notice: 'Course updated.' })
      await loadCourse()
    } catch (error) {
      setMutation({ status: 'error', error, notice: '' })
    }
  }

  async function handleDelete() {
    if (!window.confirm('Delete this course?')) {
      return
    }

    setMutation({ status: 'submitting', error: null, notice: '' })

    try {
      await deleteCourse(accessToken, courseId)
      navigate('/app/settings/courses', { replace: true })
    } catch (error) {
      setMutation({ status: 'error', error, notice: '' })
    }
  }

  if (courseState.status === 'loading') {
    return (
      <div className="module-page resource-workspace">
        <PageHeader eyebrow="Settings" title="Course detail" description={`Loading ${compactId(courseId)}.`} />
        <StateBlock tone="loading" title="Loading course" message="Fetching course detail from the API." />
      </div>
    )
  }

  if (courseState.status === 'error') {
    return (
      <div className="module-page resource-workspace">
        <PageHeader
          eyebrow="Settings"
          title={courseState.error?.code === 'NOT_FOUND' ? 'Course not found' : 'Course lookup failed'}
          description="Return to courses or use a valid ObjectID."
          actions={<button className="btn-outline" type="button" onClick={() => navigate('/app/settings/courses')}>Courses</button>}
        />
        <StateBlock tone={courseState.error?.code === 'NOT_FOUND' ? 'notFound' : 'error'} title="Could not load course" message={apiErrorText(courseState.error)} />
      </div>
    )
  }

  const course = courseState.data

  return (
    <div className="module-page resource-workspace">
      <PageHeader
        eyebrow="Settings"
        title={course.title}
        description={`Course ID ${course.id}`}
        actions={<button className="btn-outline" type="button" onClick={() => navigate('/app/settings/courses')}>Courses</button>}
      />

      <div className="module-page__grid">
        <DataPanel title="Summary">
          <dl className="detail-grid">
            <div><dt>Level</dt><dd>{course.level}</dd></div>
            <div><dt>Base price</dt><dd>{formatMoney(course.base_price)}</dd></div>
            <div><dt>Sessions</dt><dd>{course.session_count}</dd></div>
            <div><dt>Tags</dt><dd>{formatTags(course.allowed_tags)}</dd></div>
          </dl>
        </DataPanel>

        <DataPanel title="Danger zone">
          <p className="panel-copy">Delete is permanent. Backend rejects invalid references when needed.</p>
          <button className="btn-outline" type="button" onClick={handleDelete} disabled={mutation.status === 'submitting'}>
            Delete course
          </button>
        </DataPanel>
      </div>

      <DataPanel title="Edit course">
        <CourseForm
          values={values}
          setValues={setValues}
          errors={errors}
          setErrors={setErrors}
          onSubmit={handleUpdate}
          submitLabel="Save course"
          submittingLabel="Saving"
          status={mutation.status}
        />
        {mutation.error ? <div className="form-alert" role="alert">{apiErrorText(mutation.error, 'Course could not be updated.')}</div> : null}
        {mutation.notice ? <div className="form-success" role="status">{mutation.notice}</div> : null}
      </DataPanel>
    </div>
  )
}

export default CourseDetailView
