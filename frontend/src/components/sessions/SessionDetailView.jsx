import { useCallback, useEffect, useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import PageHeader from '../PageHeader.jsx'
import StateBlock from '../StateBlock.jsx'
import { useAuth } from '../../context/AuthContext.jsx'
import { getSession } from '../../lib/sessionsApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import EnrollmentPanel from './EnrollmentPanel.jsx'
import SessionCheckInPanel from './SessionCheckInPanel.jsx'
import {
  capacityLabel,
  compactId,
  formatDateTime,
  formatTags,
  isObjectId,
} from './sessionFormatters.js'

function SessionDetailView({ sessionId, navigate }) {
  const { accessToken } = useAuth()
  const [sessionState, setSessionState] = useState({ status: 'loading', data: null, error: null })

  const loadSession = useCallback(async ({ background = false } = {}) => {
    if (!isObjectId(sessionId)) {
      setSessionState({
        status: 'error',
        data: null,
        error: { code: 'INVALID_ID', message: 'Session ID must be a 24 character ObjectID.' },
      })
      return
    }

    if (!background) {
      setSessionState((current) => ({ ...current, status: 'loading', error: null }))
    }

    try {
      const response = await getSession(accessToken, sessionId)
      setSessionState({ status: 'success', data: response.data, error: null })
    } catch (error) {
      if (background) {
        setSessionState((current) => ({ ...current, error }))
        return
      }

      setSessionState({ status: 'error', data: null, error })
    }
  }, [accessToken, sessionId])

  useEffect(() => {
    loadSession()
  }, [loadSession])

  if (sessionState.status === 'loading') {
    return (
      <div className="module-page resource-workspace">
        <PageHeader eyebrow="Sessions" title="Session detail" description={`Loading ${compactId(sessionId)}.`} />
        <StateBlock tone="loading" title="Loading session" message="Fetching session detail from the API." />
      </div>
    )
  }

  if (sessionState.status === 'error') {
    return (
      <div className="module-page resource-workspace">
        <PageHeader
          eyebrow="Sessions"
          title={sessionState.error?.code === 'NOT_FOUND' ? 'Session not found' : 'Session lookup failed'}
          description="Return to sessions or use a valid ObjectID."
          actions={<button className="btn-outline" type="button" onClick={() => navigate('/app/sessions')}>Sessions</button>}
        />
        <StateBlock tone={sessionState.error?.code === 'NOT_FOUND' ? 'notFound' : 'error'} title="Could not load session" message={apiErrorText(sessionState.error)} />
      </div>
    )
  }

  const session = sessionState.data

  return (
    <div className="module-page resource-workspace sessions-workspace">
      <PageHeader
        eyebrow="Sessions"
        title={`${session.course_level} session`}
        description={`Session ID ${session.id}`}
        actions={<button className="btn-outline" type="button" onClick={() => navigate('/app/sessions')}>Sessions</button>}
      />

      <DataPanel title="Session summary">
        <dl className="detail-grid">
          <div><dt>Scheduled</dt><dd>{formatDateTime(session.scheduled_at)}</dd></div>
          <div><dt>Duration</dt><dd>{session.duration_min} min</dd></div>
          <div><dt>Capacity</dt><dd>{capacityLabel(session)}</dd></div>
          <div><dt>Branch</dt><dd>{compactId(session.branch_id)}</dd></div>
          <div><dt>Trainer</dt><dd>{compactId(session.trainer_id)}</dd></div>
          <div><dt>Tags</dt><dd>{formatTags(session.tags)}</dd></div>
        </dl>
      </DataPanel>

      <DataPanel title="Enrolled subscriptions">
        {session.enrolled_subscription_ids?.length ? (
          <div className="resource-list">
            {session.enrolled_subscription_ids.map((subscriptionId) => (
              <article className="resource-row" key={subscriptionId}>
                <div><strong>{compactId(subscriptionId)}</strong><span>Subscription</span></div>
                <div><button className="btn-outline" type="button" onClick={() => navigate(`/app/subscriptions/${subscriptionId}`)}>Open</button></div>
              </article>
            ))}
          </div>
        ) : (
          <StateBlock tone="empty" title="No enrollments" message="Enroll a subscription below." />
        )}
      </DataPanel>

      <div className="module-page__grid">
        <EnrollmentPanel accessToken={accessToken} sessionId={session.id} onChanged={() => loadSession({ background: true })} />
        <SessionCheckInPanel accessToken={accessToken} sessionId={session.id} onChanged={() => loadSession({ background: true })} />
      </div>
    </div>
  )
}

export default SessionDetailView
