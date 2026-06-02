function ScheduleList({ sessions }) {
  if (!sessions.length) {
    return <p className="dashboard-empty">No classes scheduled for today.</p>
  }

  return (
    <div className="schedule-list">
      {sessions.map((session) => {
        const capacityLabel = `${session.capacity}/${session.maxCapacity}`

        return (
          <article className="schedule-item" key={session.id || `${session.time}-${session.name}`}>
            <time>{session.time}</time>
            <span className={`schedule-dot schedule-dot--${session.tone}`} aria-hidden="true" />
            <div>
              <h4>{session.name}</h4>
              <p>{session.trainer} - {session.room}</p>
            </div>
            <strong className={`schedule-capacity schedule-capacity--${session.tone}`}>
              {capacityLabel}
            </strong>
          </article>
        )
      })}
    </div>
  )
}

export default ScheduleList
