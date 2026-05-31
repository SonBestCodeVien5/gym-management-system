const STATE_TITLES = {
  empty: 'No data yet',
  error: 'Request failed',
  forbidden: 'Access denied',
  loading: 'Loading',
  notFound: 'Page not found',
  planned: 'Planned module',
}

function StateBlock({ tone = 'empty', title, message, details = null }) {
  return (
    <section className={`state-block state-block--${tone}`} aria-live={tone === 'error' ? 'polite' : undefined}>
      <span>{tone}</span>
      <h2>{title || STATE_TITLES[tone] || 'Status'}</h2>
      {message ? <p>{message}</p> : null}
      {details ? <div className="state-block__details">{details}</div> : null}
    </section>
  )
}

export default StateBlock
