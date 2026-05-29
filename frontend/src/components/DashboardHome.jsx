import KpiCard from './KpiCard.jsx'
import MemberTable from './MemberTable.jsx'
import PlanDonut from './PlanDonut.jsx'
import RevenueBars from './RevenueBars.jsx'
import ScheduleList from './ScheduleList.jsx'
import { dashboardData } from './dashboardData.js'

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
  const activeModule = navItems.find((item) => item.key === activeItem)
  const todayLabel = new Intl.DateTimeFormat('en-US', {
    weekday: 'long',
    month: 'short',
    day: 'numeric',
  }).format(new Date())

  return (
    <div className="ops-dashboard" aria-labelledby="dashboard-title">
      <section className="dashboard-hero ops-panel">
        <div>
          <p className="section-eyebrow">Operations overview</p>
          <h2 id="dashboard-title">{activeModule?.label || 'Dashboard'}</h2>
        </div>
        <p>
          Reference-inspired dashboard using sample operational data. Staff identity and access context
          still come from the authenticated API session.
        </p>
        <div className="dashboard-hero__meta">
          <span>{todayLabel}</span>
          <span>{dashboardData.updatedAtLabel}</span>
          <span>{branchCount} assigned branches</span>
        </div>
      </section>

      <div className="kpi-grid" aria-label="Sample dashboard KPIs">
        {dashboardData.kpis.map((item) => (
          <KpiCard item={item} key={item.key} />
        ))}
      </div>

      <div className="dashboard-chart-row">
        <section className="ops-panel chart-panel" aria-labelledby="revenue-title">
          <div className="panel-head">
            <div>
              <span className="panel-label">Revenue</span>
              <h3 id="revenue-title">Last 7 days</h3>
              <p>Sample revenue, in million VND.</p>
            </div>
            <span className="panel-pill">This week</span>
          </div>
          <RevenueBars items={dashboardData.revenue} />
        </section>

        <section className="ops-panel chart-panel" aria-labelledby="plans-title">
          <div className="panel-head">
            <div>
              <span className="panel-label">Membership mix</span>
              <h3 id="plans-title">Plan distribution</h3>
              <p>Sample member count by plan.</p>
            </div>
          </div>
          <PlanDonut items={dashboardData.planDistribution} />
        </section>
      </div>

      <div className="dashboard-bottom-row">
        <section className="ops-panel members-panel" aria-labelledby="latest-members-title">
          <div className="panel-head">
            <div>
              <span className="panel-label">Members</span>
              <h3 id="latest-members-title">Latest registrations</h3>
            </div>
            <span className="panel-pill">Today - 12</span>
          </div>
          <MemberTable members={dashboardData.latestMembers} />
        </section>

        <section className="ops-panel schedule-panel" aria-labelledby="today-sessions-title">
          <div className="panel-head">
            <div>
              <span className="panel-label">Classes</span>
              <h3 id="today-sessions-title">Today schedule</h3>
            </div>
            <span className="panel-pill">{dashboardData.todaySessions.length} classes</span>
          </div>
          <ScheduleList sessions={dashboardData.todaySessions} />
        </section>
      </div>

      <section className="ops-panel staff-context" aria-labelledby="staff-context-title">
        <div>
          <span className="panel-label">Staff context</span>
          <h3 id="staff-context-title">{employee.full_name}</h3>
          <p>
            Live identity from `GET /api/v1/auth/me`. Dashboard metrics above are frontend sample data
            until backend report APIs exist.
          </p>
        </div>

        <dl className="staff-context__details">
          <div>
            <dt>Employee ID</dt>
            <dd>{employee.employee_id || 'Unassigned'}</dd>
          </div>
          <div>
            <dt>Email</dt>
            <dd>{employee.email}</dd>
          </div>
          <div>
            <dt>Ready modules</dt>
            <dd>{availableModules.map((item) => item.label).join(', ') || 'None'}</dd>
          </div>
          <div>
            <dt>Roles</dt>
            <dd>
              <span className="chip-row">
                {roles.map((role) => (
                  <span className="role-chip" key={role}>
                    {ROLE_LABELS[role] || role}
                  </span>
                ))}
              </span>
            </dd>
          </div>
        </dl>
      </section>
    </div>
  )
}

export default DashboardHome
