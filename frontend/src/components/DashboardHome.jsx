import { API_BASE_URL } from '../lib/api.js'

const ROLE_LABELS = {
  admin: 'Admin',
  manager: 'Manager',
  trainer: 'Trainer',
  receptionist: 'Receptionist',
}

function DashboardHome({ employee, navItems, activeItem }) {
  const roles = employee.role || []
  const branchCount = employee.branch_id?.length || 0
  const availableModules = navItems.filter((item) => item.available)
  const upcomingModules = navItems.filter((item) => item.visible && !item.ready)
  const activeModule = navItems.find((item) => item.key === activeItem)

  return (
    <div className="dashboard-grid">
      <section className="dashboard-hero">
        <p className="section-eyebrow">Phien lam viec</p>
        <h2>{activeModule?.label || 'Dashboard'}</h2>
        <p>
          Dang nhap thanh cong voi vai tro staff. Cac module nghiep vu se duoc mo theo tung feat tiep
          theo.
        </p>
      </section>

      <section className="info-panel">
        <span className="panel-label">Nhan vien</span>
        <h3>{employee.full_name}</h3>
        <dl className="detail-list">
          <div>
            <dt>Ma nhan vien</dt>
            <dd>{employee.employee_id || 'Chua gan'}</dd>
          </div>
          <div>
            <dt>Email</dt>
            <dd>{employee.email}</dd>
          </div>
          <div>
            <dt>Chi nhanh</dt>
            <dd>{branchCount}</dd>
          </div>
        </dl>
      </section>

      <section className="info-panel">
        <span className="panel-label">Quyen truy cap</span>
        <div className="chip-row">
          {roles.map((role) => (
            <span className="role-chip" key={role}>
              {ROLE_LABELS[role] || role}
            </span>
          ))}
        </div>
        <p className="panel-copy">
          Module hien tai duoc tinh theo role guard backend. Branch-scope authorization chua nam trong
          FE cycle nay.
        </p>
      </section>

      <section className="metric-panel">
        <span>Module san sang</span>
        <strong>{availableModules.length}</strong>
        <p>{availableModules.map((item) => item.label).join(', ') || 'Chua co'}</p>
      </section>

      <section className="metric-panel">
        <span>Dang cho feat sau</span>
        <strong>{upcomingModules.length}</strong>
        <p>{upcomingModules.map((item) => item.label).join(', ') || 'Khong co'}</p>
      </section>

      <section className="info-panel info-panel--wide">
        <span className="panel-label">API session</span>
        <h3>Bearer token active</h3>
        <p className="panel-copy">
          Frontend dang dung access token trong header Authorization va khoi phuc user bang
          `/api/v1/auth/me`.
        </p>
        <div className="api-line">
          <span>Base URL</span>
          <code>{API_BASE_URL}</code>
        </div>
      </section>
    </div>
  )
}

export default DashboardHome
