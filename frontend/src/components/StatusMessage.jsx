function StatusMessage({ title, message, tone = 'neutral', fullPage = false, action }) {
  return (
    <div className={fullPage ? 'status-page grid-bg' : 'status-wrap'} aria-live="polite">
      <section className={`status-message status-message--${tone}`}>
        <p className="status-kicker">Iron Forge Staff</p>
        <h1>{title}</h1>
        {message ? <p>{message}</p> : null}
        {action ? <div className="status-action">{action}</div> : null}
      </section>
    </div>
  )
}

export default StatusMessage
