function RevenueBars({ items }) {
  if (!items.length) {
    return <p className="dashboard-empty">No revenue sample data.</p>
  }

  const maxValue = Math.max(...items.map((item) => item.value))
  const peak = items.reduce((highest, item) => (item.value > highest.value ? item : highest), items[0])
  const average = items.reduce((total, item) => total + item.value, 0) / items.length

  return (
    <>
      <div className="revenue-bars" aria-label="Revenue for the last 7 days">
        {items.map((item) => {
          const height = Math.max(18, Math.round((item.value / maxValue) * 100))

          return (
            <div className="revenue-bar-wrap" key={item.day}>
              <div
                className={`revenue-bar revenue-bar--${item.tone}`}
                style={{ height: `${height}%` }}
                role="img"
                aria-label={`${item.day}: ${item.value}M VND`}
                title={`${item.day}: ${item.value}M VND`}
              />
              <span>{item.day}</span>
            </div>
          )
        })}
      </div>

      <div className="chart-footnote">
        <span>Peak: {peak.day} - {peak.value.toFixed(1)}M VND</span>
        <span>Avg/day: {average.toFixed(1)}M VND</span>
      </div>
    </>
  )
}

export default RevenueBars
