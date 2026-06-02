const RADIUS = 35
const CIRCUMFERENCE = 2 * Math.PI * RADIUS

function PlanDonut({ items }) {
  const total = items.reduce((sum, item) => sum + item.value, 0)

  if (!items.length || total === 0) {
    return <p className="dashboard-empty">No live plan distribution yet.</p>
  }

  let offset = 0

  return (
    <div className="plan-donut">
      <svg
        className="plan-donut__chart"
        width="116"
        height="116"
        viewBox="0 0 100 100"
        role="img"
        aria-label="Live subscription distribution by plan"
      >
        <circle cx="50" cy="50" r={RADIUS} fill="none" stroke="var(--color-border)" strokeWidth="13" />
        {items.map((item) => {
          const length = (item.value / total) * CIRCUMFERENCE
          const dashOffset = -offset
          offset += length

          return (
            <circle
              key={item.label}
              cx="50"
              cy="50"
              r={RADIUS}
              fill="none"
              stroke={item.color}
              strokeWidth="13"
              strokeDasharray={`${length} ${CIRCUMFERENCE - length}`}
              strokeDashoffset={dashOffset}
              transform="rotate(-90 50 50)"
            />
          )
        })}
        <text x="50" y="48" textAnchor="middle" className="plan-donut__total">
          {total.toLocaleString('en-US')}
        </text>
        <text x="50" y="61" textAnchor="middle" className="plan-donut__label">
          subs
        </text>
      </svg>

      <div className="plan-donut__legend">
        {items.map((item) => {
          const percent = Math.round((item.value / total) * 100)

          return (
            <div className="legend-item" key={item.label}>
              <span className="legend-dot" style={{ background: item.color }} />
              <span>
                {item.label} - <strong>{item.value.toLocaleString('en-US')}</strong> ({percent}%)
              </span>
            </div>
          )
        })}
      </div>
    </div>
  )
}

export default PlanDonut
