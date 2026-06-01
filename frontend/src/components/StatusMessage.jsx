import loadingIcon from '../assets/brand/loading-spinner-icon.svg'
import BrandMark from './BrandMark.jsx'

function StatusMessage({ title, message, tone = 'neutral', fullPage = false, action }) {
  return (
    <div className={fullPage ? 'status-page grid-bg' : 'status-wrap'} aria-live="polite">
      <section className={`status-message status-message--${tone}`}>
        <div className="status-brand-row">
          <BrandMark variant="compact" label="Iron Forge Staff" meta="Staff" />
          <img className="status-spinner" src={loadingIcon} alt="" aria-hidden="true" />
        </div>
        <h1>{title}</h1>
        {message ? <p>{message}</p> : null}
        {action ? <div className="status-action">{action}</div> : null}
      </section>
    </div>
  )
}

export default StatusMessage
