import { useCallback, useEffect, useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import PageHeader from '../PageHeader.jsx'
import StateBlock from '../StateBlock.jsx'
import { useAuth } from '../../context/AuthContext.jsx'
import { createCourse, deleteCourse, listCourses } from '../../lib/coursesApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import CourseForm from './CourseForm.jsx'
import { compactId, formatMoney, formatTags } from './settingsFormatters.js'

function CoursesPage({ navigate }) {
  const { accessToken } = useAuth()
  const [coursesState, setCoursesState] = useState({ status: 'loading', data: [], error: null })
  const [values, setValues] = useState(CourseForm.initialValues)
  const [errors, setErrors] = useState({})
  const [mutation, setMutation] = useState({ status: 'idle', error: null, notice: '' })

  const loadCourses = useCallback(async () => {
    setCoursesState((current) => ({ ...current, status: 'loading', error: null }))

    try {
      const response = await listCourses(accessToken)
      setCoursesState({ status: 'success', data: response.data || [], error: null })
    } catch (error) {
      setCoursesState({ status: 'error', data: [], error })
    }
  }, [accessToken])

  useEffect(() => {
    loadCourses()
  }, [loadCourses])

  async function handleCreate(payload) {
    setMutation({ status: 'submitting', error: null, notice: '' })

    try {
      const response = await createCourse(accessToken, payload)
      const courseId = response.data?.id
      setValues(CourseForm.initialValues)
      setMutation({ status: 'success', error: null, notice: 'Course created.' })

      if (courseId) {
        navigate(`/app/settings/courses/${courseId}`, { replace: true })
        return
      }

      await loadCourses()
    } catch (error) {
      setMutation({ status: 'error', error, notice: '' })
    }
  }

  async function handleDelete(course) {
    if (!window.confirm(`Delete course ${course.title}?`)) {
      return
    }

    setMutation({ status: 'submitting', error: null, notice: '' })

    try {
      await deleteCourse(accessToken, course.id)
      setMutation({ status: 'success', error: null, notice: 'Course deleted.' })
      await loadCourses()
    } catch (error) {
      setMutation({ status: 'error', error, notice: '' })
    }
  }

  return (
    <div className="module-page resource-workspace settings-workspace">
      <PageHeader
        eyebrow="Settings"
        title="Courses"
        description="Manage package templates used by subscriptions and session tags."
        actions={<button className="btn-outline" type="button" onClick={loadCourses}>Refresh</button>}
      />

      <div className="module-page__grid">
        <DataPanel title="Create course" description="Course updates use the full backend course shape.">
          <CourseForm
            values={values}
            setValues={setValues}
            errors={errors}
            setErrors={setErrors}
            onSubmit={handleCreate}
            submitLabel="Create course"
            submittingLabel="Creating"
            status={mutation.status}
          />
          {mutation.error ? <div className="form-alert" role="alert">{apiErrorText(mutation.error, 'Course could not be saved.')}</div> : null}
          {mutation.notice ? <div className="form-success" role="status">{mutation.notice}</div> : null}
        </DataPanel>

        <DataPanel title="Course data">
          <ul className="feature-list">
            <li>Title, level, base price, and session count are required.</li>
            <li>Allowed tags feed session validation later.</li>
            <li>Duplicate or referenced course rules remain enforced by the backend.</li>
          </ul>
        </DataPanel>
      </div>

      <DataPanel title="Course list">
        {coursesState.status === 'loading' ? (
          <StateBlock tone="loading" title="Loading courses" message="Fetching course packages from the API." />
        ) : null}

        {coursesState.status === 'error' ? (
          <StateBlock tone="error" title="Could not load courses" message={apiErrorText(coursesState.error)} />
        ) : null}

        {coursesState.status === 'success' && !coursesState.data.length ? (
          <StateBlock tone="empty" title="No courses yet" message="Create the first package template above." />
        ) : null}

        {coursesState.status === 'success' && coursesState.data.length ? (
          <div className="resource-list">
            {coursesState.data.map((course) => (
              <article className="resource-row" key={course.id}>
                <div>
                  <strong>{course.title}</strong>
                  <span>{course.level} · {formatMoney(course.base_price)} · {course.session_count} sessions</span>
                  <small>{formatTags(course.allowed_tags)}</small>
                </div>
                <div>
                  <span>{compactId(course.id)}</span>
                  <button className="btn-outline" type="button" onClick={() => navigate(`/app/settings/courses/${course.id}`)}>Open</button>
                  <button className="btn-outline" type="button" onClick={() => handleDelete(course)}>Delete</button>
                </div>
              </article>
            ))}
          </div>
        ) : null}
      </DataPanel>
    </div>
  )
}

export default CoursesPage
