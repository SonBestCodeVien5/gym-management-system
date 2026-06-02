import { capacityLabel, compactId, formatDateTime, formatTags } from './sessionFormatters.js'

function SessionList({ sessions, navigate }) {
  return (
    <div className="resource-list">
      {sessions.map((session) => (
        <article className="resource-row" key={session.id}>
          <div>
            <strong>{session.course_level}</strong>
            <span>{formatDateTime(session.scheduled_at)} · {session.duration_min} min · capacity {capacityLabel(session)}</span>
            <small>{formatTags(session.tags)}</small>
          </div>
          <div>
            <span>Branch {compactId(session.branch_id)}</span>
            <button className="btn-outline" type="button" onClick={() => navigate(`/app/sessions/${session.id}`)}>Open</button>
          </div>
        </article>
      ))}
    </div>
  )
}

export default SessionList
