function MemberTable({ members }) {
  if (!members.length) {
    return <p className="dashboard-empty">No recent member registrations.</p>
  }

  return (
    <div className="table-scroll">
      <table className="member-table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Level</th>
            <th>Joined</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {members.map((member) => (
            <tr key={member.id || `${member.name}-${member.phone}`}>
              <td data-label="Name">
                <strong>{member.name}</strong>
                <span>{member.phone}</span>
              </td>
              <td data-label="Level">
                <span className={`status-badge status-badge--${member.levelTone}`}>{member.level}</span>
              </td>
              <td data-label="Joined">{member.joinedAt}</td>
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
