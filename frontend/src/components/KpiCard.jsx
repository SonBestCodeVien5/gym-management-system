function KpiCard({ item }) {
  return (
    <article className={item.accent ? 'kpi-card kpi-card--accent' : 'kpi-card'}>
      <span className="kpi-card__label">{item.label}</span>
      <strong>{item.value}</strong>
      <p className={`kpi-card__delta kpi-card__delta--${item.tone}`}>{item.delta}</p>
    </article>
  )
}

export default KpiCard
