import brandIcon from '../assets/brand/iron-forge-favicon-if.svg'
import brandLogo from '../assets/brand/iron-forge-logo-horizontal-dark.svg'

function BrandMark({
  variant = 'full',
  label = 'Iron Forge Gym',
  meta,
  className = '',
}) {
  const isCompact = variant === 'compact'
  const classes = ['brand-mark', `brand-mark--${variant}`, className]
    .filter(Boolean)
    .join(' ')

  return (
    <div className={classes}>
      <img
        className="brand-mark__image"
        src={isCompact ? brandIcon : brandLogo}
        alt={label}
      />
      {meta ? <span className="brand-mark__meta">{meta}</span> : null}
    </div>
  )
}

export default BrandMark
