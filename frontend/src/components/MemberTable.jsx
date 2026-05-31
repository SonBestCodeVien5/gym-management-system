function MemberTable({ members }) {
  if (!members.length) {
    return <p className="dashboard-empty">No new member registrations in this sample.</p>
  }

  return (
    <div className="table-scroll">
      <table className="member-table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Plan</th>
            <th>Trainer</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {members.map((member) => (
            <tr key={`${member.name}-${member.phone}`}>
              <td data-label="Name">
                <strong>{member.name}</strong>
                <span>{member.phone}</span>
              </td>
              <td data-label="Plan">
                <span className={`status-badge status-badge--${member.planTone}`}>{member.plan}</span>
              </td>
              <td data-label="Trainer">{member.trainer}</td>
              <td data-label="Status">
                <span className={`status-badge status-badge--${member.statusTone}`}>{member.status}</span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default MemberTable
