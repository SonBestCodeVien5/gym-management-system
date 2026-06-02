function RevenueBars({ items }) {
  if (!items.length) {
    return <p className="dashboard-empty">No revenue recorded for this range.</p>
  }

  const maxMagnitude = Math.max(...items.map((item) => Math.abs(Number(item.value || 0))), 0)
  if (maxMagnitude <= 0) {
    return <p className="dashboard-empty">No net revenue recorded for this range.</p>
  }

  const peak = items.reduce(
    (highest, item) => (Math.abs(item.amount) > Math.abs(highest.amount) ? item : highest),
    items[0],
  )
  const average = items.reduce((total, item) => total + item.value, 0) / items.length
  const averageLabel = `${average < 0 ? '-' : ''}${Math.abs(average).toFixed(1)}M VND`

  return (
    <>
      <div className="revenue-bars" aria-label="Live net revenue for the last 7 days">
        {items.map((item) => {
          const signedValue = Number(item.value || 0)
          const height = Math.max(6, Math.round((Math.abs(signedValue) / maxMagnitude) * 46))
          const positionStyle = signedValue < 0 ? { top: '50%' } : { bottom: '50%' }

          return (
            <div className="revenue-bar-wrap" key={item.day}>
              <div className="revenue-bar-track">
                <span className="revenue-bar-baseline" aria-hidden="true" />
                <div
                  className={`revenue-bar revenue-bar--${item.tone}${signedValue < 0 ? ' revenue-bar--negative' : ''}`}
                  style={{ ...positionStyle, height: `${height}%` }}
                  role="img"
                  aria-label={`${item.day}: ${item.displayValue}`}
                  title={`${item.day}: ${item.displayValue}`}
                />
              </div>
              <span>{item.day}</span>
            </div>
          )
        })}
      </div>

      <div className="chart-footnote">
        <span>
          Largest move: {peak.day} ({peak.displayValue})
        </span>
        <span>Avg/day: {averageLabel}</span>
      </div>
    </>
  )
}

export default RevenueBars
