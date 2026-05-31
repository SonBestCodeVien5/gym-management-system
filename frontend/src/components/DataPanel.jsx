function DataPanel({ title, description, children, action = null }) {
  return (
    <section className="data-panel" aria-labelledby={title ? `${title.toLowerCase().replace(/\s+/g, '-')}-title` : undefined}>
      {title || description || action ? (
        <div className="data-panel__head">
          <div>
            {title ? <h3 id={`${title.toLowerCase().replace(/\s+/g, '-')}-title`}>{title}</h3> : null}
            {description ? <p>{description}</p> : null}
          </div>
          {action}
        </div>
      ) : null}
      {children}
    </section>
  )
}

export default DataPanel
