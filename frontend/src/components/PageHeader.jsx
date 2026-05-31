function PageHeader({ eyebrow, title, description, actions = null }) {
  return (
    <header className="page-header">
      <div>
        {eyebrow ? <span className="panel-label">{eyebrow}</span> : null}
        <h2>{title}</h2>
        {description ? <p>{description}</p> : null}
      </div>
      {actions ? <div className="page-header__actions">{actions}</div> : null}
    </header>
  )
}

export default PageHeader
