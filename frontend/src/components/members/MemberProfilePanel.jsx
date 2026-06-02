import DataPanel from '../DataPanel.jsx'
import {
  formatDateTime,
  formatText,
  getRegistrationStatus,
} from './memberFormatters.js'

function DetailItem({ label, value }) {
  return (
    <div>
      <span>{label}</span>
      <strong>{formatText(value)}</strong>
    </div>
  )
}

function MemberProfilePanel({ member }) {
  const registrationStatus = getRegistrationStatus(member)

  return (
    <DataPanel title="Profile" description="Member identity, registration, and training summary.">
      <div className="member-profile-card">
        <div className="member-profile-card__head">
          <div>
            <span className="panel-label">Member</span>
            <h3>{member.full_name}</h3>
            <p>{member.ccid}</p>
          </div>
          <span className={`status-badge status-badge--${registrationStatus.tone}`}>
            {registrationStatus.label}
          </span>
        </div>

        <div className="member-detail-grid">
          <DetailItem label="Email" value={member.email} />
          <DetailItem label="Phone" value={member.phone} />
          <DetailItem label="Gender" value={member.gender} />
          <DetailItem label="Level" value={member.level} />
          <DetailItem label="Sessions attended" value={member.total_sessions_attended ?? 0} />
          <DetailItem label="Created" value={formatDateTime(member.created_at)} />
          <DetailItem label="Updated" value={formatDateTime(member.updated_at)} />
        </div>
      </div>
    </DataPanel>
  )
}

export default MemberProfilePanel
