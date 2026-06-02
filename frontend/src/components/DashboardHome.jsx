import { useEffect, useMemo, useState } from 'react'
import KpiCard from './KpiCard.jsx'
import MemberTable from './MemberTable.jsx'
import PlanDonut from './PlanDonut.jsx'
import RevenueBars from './RevenueBars.jsx'
import ScheduleList from './ScheduleList.jsx'
import StateBlock from './StateBlock.jsx'
import { useAuth } from '../context/AuthContext.jsx'
import {
  getDashboardPlanDistribution,
  getDashboardRecentMembers,
  getDashboardRevenue,
  getDashboardSummary,
  getDashboardTodaySessions,
} from '../lib/dashboardApi.js'
import { apiErrorText, compactId, formatDateTime, formatMoney, formatText } from '../lib/featureHelpers.js'

const ROLE_LABELS = {
  admin: 'Admin',
  manager: 'Manager',
  trainer: 'Trainer',
  receptionist: 'Receptionist',
}

const DASHBOARD_ROLES = ['admin', 'manager']
const PLAN_COLORS = ['var(--color-accent)', 'rgba(255, 70, 20, 0.52)', 'var(--color-surface-strong)', '#f0ece4']

const EMPTY_DASHBOARD = {
  summary: null,
  revenue: [],
  plans: [],
  members: [],
  sessions: [],
}

function hasDashboardAccess(roles) {
  return roles.some((role) => DASHBOARD_ROLES.includes(role))
}

function signedNumber(value, unit = '') {
  const number = Number(value || 0)
  const prefix = number > 0 ? '+' : ''
  return `${prefix}${number.toLocaleString('en-US')}${unit}`
}

function deltaTone(value) {
  const number = Number(value || 0)
  if (number > 0) {
    return 'up'
  }
  if (number < 0) {
    return 'down'
  }
  return 'neutral'
}

function compactMoneyDelta(value) {
  const number = Number(value || 0)
  const abs = Math.abs(number)
  const prefix = number > 0 ? '+' : number < 0 ? '-' : ''

  if (abs >= 1000000) {
    return `${prefix}${(abs / 1000000).toFixed(1)}M VND`
  }

  return `${prefix}${formatMoney(abs)}`
}

function formatDateLabel(value) {
  if (!value) {
    return 'Live data'
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return 'Live data'
  }

  return new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}

function buildKpis(summary) {
  if (!summary) {
    return [
      { key: 'members', label: 'Active members', value: '...', delta: 'Loading', tone: 'neutral', accent: true },
      { key: 'revenue', label: 'Net revenue', value: '...', delta: 'Loading', tone: 'neutral' },
      { key: 'checkins', label: 'Today check-ins', value: '...', delta: 'Loading', tone: 'neutral' },
      { key: 'classes', label: 'Classes this week', value: '...', delta: 'Loading', tone: 'neutral' },
    ]
  }

  return [
    {
      key: 'members',
      label: 'Active members',
      value: Number(summary.active_members || 0).toLocaleString('en-US'),
      delta: `${signedNumber(summary.active_members_delta)} vs previous range`,
      tone: deltaTone(summary.active_members_delta),
      accent: true,
    },
    {
      key: 'revenue',
      label: 'Net revenue',
      value: formatMoney(summary.monthly_revenue),
      delta: `${compactMoneyDelta(summary.monthly_revenue_delta)} vs previous range`,
      tone: deltaTone(summary.monthly_revenue_delta),
    },
    {
      key: 'checkins',
      label: 'Today check-ins',
      value: Number(summary.today_checkins || 0).toLocaleString('en-US'),
      delta: `${signedNumber(summary.today_checkins_delta)} vs yesterday`,
      tone: deltaTone(summary.today_checkins_delta),
    },
    {
      key: 'classes',
      label: 'Classes this week',
      value: Number(summary.classes_this_week || 0).toLocaleString('en-US'),
      delta: `${signedNumber(summary.classes_this_week_delta)} vs last week`,
      tone: deltaTone(summary.classes_this_week_delta),
    },
  ]
}

function mapRevenueItems(items) {
  return items.map((item, index) => {
    const amount = Number(item.net_amount || 0)
    return {
      day: item.label || `Day ${index + 1}`,
      value: amount / 1000000,
      amount,
      displayValue: formatMoney(amount),
      tone: amount < 0 ? 'negative' : amount > 0 ? (index === items.length - 1 ? 'accent' : 'mid') : 'base',
    }
  })
}

function mapPlanItems(items) {
  return items.map((item, index) => ({
    label: item.label || compactId(item.course_id),
    value: Number(item.count || 0),
    color: PLAN_COLORS[index % PLAN_COLORS.length],
  }))
}

function mapMembers(items) {
  return items.map((member) => ({
    id: member.id,
    name: formatText(member.full_name, 'Unnamed member'),
    phone: formatText(member.phone, 'No phone'),
    level: formatText(member.level, 'Unassigned'),
    joinedAt: formatDateTime(member.created_at),
    status: member.is_registered ? 'Registered' : 'Pending',
    statusTone: member.is_registered ? 'success' : 'warning',
    levelTone: member.level ? 'warning' : 'neutral',
  }))
}

function mapSessions(items) {
  return items.map((session) => {
    const scheduledAt = new Date(session.scheduled_at)
    const time = Number.isNaN(scheduledAt.getTime())
      ? 'TBD'
      : new Intl.DateTimeFormat('en-US', { hour: '2-digit', minute: '2-digit' }).format(scheduledAt)
    const enrolledCount = Number(session.enrolled_count || 0)
    const capacity = Number(session.capacity || 0)
    const tone = capacity > 0 && enrolledCount >= capacity ? 'full' : enrolledCount >= capacity * 0.8 ? 'warn' : 'ok'

    return {
      id: session.id,
      time,
      name: `${formatText(session.course_level, 'Training')} session`,
      trainer: compactId(session.trainer_id),
      room: compactId(session.branch_id),
      capacity: enrolledCount,
      maxCapacity: capacity,
      tone,
    }
  })
}

function DashboardHome({ employee, navItems, activeItem }) {
  const { accessToken } = useAuth()
  const [expandedMobilePanel, setExpandedMobilePanel] = useState(null)
  const [dashboardState, setDashboardState] = useState({
    status: 'idle',
    data: EMPTY_DASHBOARD,
    errors: {},
    error: null,
  })
  const roles = employee.role || []
  const branchCount = employee.branch_id?.length || 0
  const availableModules = navItems.filter((item) => item.available)
  const activeModule = navItems.find((item) => item.key === activeItem)
  const canUseDashboard = hasDashboardAccess(roles)
  const data = dashboardState.data
  const kpis = useMemo(() => buildKpis(data.summary), [data.summary])
  const revenueItems = useMemo(() => mapRevenueItems(data.revenue), [data.revenue])
  const planItems = useMemo(() => mapPlanItems(data.plans), [data.plans])
  const latestMembers = useMemo(() => mapMembers(data.members), [data.members])
  const todaySessions = useMemo(() => mapSessions(data.sessions), [data.sessions])
  const revenueTotal = revenueItems.reduce((total, item) => total + item.amount, 0)
  const revenuePeak = revenueItems.reduce(
    (highest, item) => (!highest || Math.abs(item.amount) > Math.abs(highest.amount) ? item : highest),
    null,
  )
  const planTotal = planItems.reduce((total, item) => total + item.value, 0)
  const topPlan = planItems.reduce(
    (highest, item) => (!highest || item.value > highest.value ? item : highest),
    null,
  )
  const todayLabel = new Intl.DateTimeFormat('en-US', {
    weekday: 'long',
    month: 'short',
    day: 'numeric',
  }).format(new Date())
  const updatedAtLabel = data.summary?.range?.to ? `Updated ${formatDateLabel(data.summary.range.to)}` : 'Live dashboard'
  const sectionErrorMessages = Object.values(dashboardState.errors).filter(Boolean)

  useEffect(() => {
    if (!accessToken || !canUseDashboard) {
      return
    }

    let active = true

    async function loadDashboard() {
      setDashboardState((current) => ({
        ...current,
        status: current.data.summary ? 'refreshing' : 'loading',
        error: null,
        errors: {},
      }))

      const [summaryResult, revenueResult, plansResult, membersResult, sessionsResult] = await Promise.allSettled([
        getDashboardSummary(accessToken),
        getDashboardRevenue(accessToken),
        getDashboardPlanDistribution(accessToken),
        getDashboardRecentMembers(accessToken, { limit: 5 }),
        getDashboardTodaySessions(accessToken),
      ])

      if (!active) {
        return
      }

      if (summaryResult.status === 'rejected') {
        setDashboardState((current) => ({
          status: current.data.summary ? 'stale' : 'error',
          data: current.data,
          errors: {},
          error: summaryResult.reason,
        }))
        return
      }

      setDashboardState((current) => ({
        status: 'success',
        data: {
          summary: summaryResult.value.data,
          revenue: revenueResult.status === 'fulfilled' ? revenueResult.value.data?.items || [] : current.data.revenue,
          plans: plansResult.status === 'fulfilled' ? plansResult.value.data?.items || [] : current.data.plans,
          members: membersResult.status === 'fulfilled' ? membersResult.value.data?.items || [] : current.data.members,
          sessions: sessionsResult.status === 'fulfilled' ? sessionsResult.value.data?.items || [] : current.data.sessions,
        },
        errors: {
          revenue: revenueResult.status === 'rejected' ? apiErrorText(revenueResult.reason) : null,
          plans: plansResult.status === 'rejected' ? apiErrorText(plansResult.reason) : null,
          members: membersResult.status === 'rejected' ? apiErrorText(membersResult.reason) : null,
          sessions: sessionsResult.status === 'rejected' ? apiErrorText(sessionsResult.reason) : null,
        },
        error: null,
      }))
    }

    loadDashboard()

    return () => {
      active = false
    }
  }, [accessToken, canUseDashboard])

  if (!canUseDashboard) {
    return (
      <StateBlock
        tone="forbidden"
        title="Dashboard access denied"
        message="Live reporting metrics are available to admin and manager roles."
      />
    )
  }

  if (dashboardState.status === 'error') {
    return (
      <StateBlock
        tone="error"
        title="Dashboard unavailable"
        message={apiErrorText(dashboardState.error, 'Cannot load live dashboard metrics.')}
      />
    )
  }

  return (
    <div className="ops-dashboard" aria-labelledby="dashboard-title">
      <section className="ops-panel staff-context" aria-labelledby="staff-context-title">
        <div>
          <span className="panel-label">Staff context</span>
          <h2 id="dashboard-title">{activeModule?.label || 'Dashboard'}</h2>
          <h3 id="staff-context-title">{employee.full_name}</h3>
          <p>
            Live identity from `GET /api/v1/auth/me`. Dashboard metrics below are loaded from the
            backend dashboard API.
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

      <section className="dashboard-hero ops-panel">
        <div>
          <p className="section-eyebrow">Operations overview</p>
          <h2>{activeModule?.label || 'Dashboard'}</h2>
        </div>
        <p>
          Live operational dashboard using backend summary, revenue, plan mix, member, and session
          aggregates.
        </p>
        <div className="dashboard-hero__meta">
          <span>{todayLabel}</span>
          <span>
            {dashboardState.status === 'loading'
              ? 'Loading metrics'
              : dashboardState.status === 'stale'
                ? 'Showing last successful snapshot'
                : updatedAtLabel}
          </span>
          <span>{branchCount} assigned branches</span>
        </div>
      </section>

      {dashboardState.status === 'stale' ? (
        <section className="form-alert" aria-live="polite">
          <strong>Live refresh failed, so the last successful snapshot is still on screen.</strong>
          <span>{apiErrorText(dashboardState.error, 'The previous data remains visible.')}</span>
        </section>
      ) : null}

      {sectionErrorMessages.length ? (
        <section className="form-alert" aria-live="polite">
          <strong>Some dashboard sections could not refresh.</strong>
          <span>{sectionErrorMessages.join(' ')}</span>
        </section>
      ) : null}

      <div className="kpi-grid" aria-label="Live dashboard KPIs" aria-busy={dashboardState.status === 'loading'}>
        {kpis.map((item) => (
          <KpiCard item={item} key={item.key} />
        ))}
      </div>

      <div className="mobile-summary-row" aria-label="Live compact dashboard numbers">
        <section className="ops-panel compact-number-panel" aria-labelledby="mobile-revenue-title">
          <span className="panel-label">Revenue</span>
          <h3 id="mobile-revenue-title">{formatMoney(revenueTotal)}</h3>
          <p>
            {revenuePeak
              ? `7-day net total. Largest move: ${revenuePeak.day} (${revenuePeak.displayValue}).`
              : 'No revenue recorded in the live dashboard range.'}
          </p>
        </section>
        <section className="ops-panel compact-number-panel" aria-labelledby="mobile-plan-title">
          <span className="panel-label">Plan mix</span>
          <h3 id="mobile-plan-title">{planTotal.toLocaleString('en-US')}</h3>
          <p>
            {topPlan
              ? `Live subscriptions. Top plan: ${topPlan.label} (${topPlan.value.toLocaleString('en-US')}).`
              : 'No live plan distribution yet.'}
          </p>
        </section>
      </div>

      <div className="dashboard-chart-row">
        <section className="ops-panel chart-panel" aria-labelledby="revenue-title">
          <div className="panel-head">
            <div>
              <span className="panel-label">Revenue</span>
              <h3 id="revenue-title">Last 7 days</h3>
              <p>Live net revenue, in million VND.</p>
            </div>
            <span className="panel-pill">Live</span>
          </div>
          <RevenueBars items={revenueItems} />
        </section>

        <section className="ops-panel chart-panel" aria-labelledby="plans-title">
          <div className="panel-head">
            <div>
              <span className="panel-label">Membership mix</span>
              <h3 id="plans-title">Plan distribution</h3>
              <p>Live subscription count by plan.</p>
            </div>
          </div>
          <PlanDonut items={planItems} />
        </section>
      </div>

      <div className="dashboard-bottom-row">
        <section className="ops-panel members-panel" aria-labelledby="latest-members-title">
          <div className="panel-head">
            <div>
              <span className="panel-label">Members</span>
              <h3 id="latest-members-title">Latest registrations</h3>
            </div>
            <span className="panel-pill">{latestMembers.length} recent</span>
          </div>
          <MemberTable members={latestMembers} />
        </section>

        <section className="ops-panel schedule-panel" aria-labelledby="today-sessions-title">
          <div className="panel-head">
            <div>
              <span className="panel-label">Classes</span>
              <h3 id="today-sessions-title">Today schedule</h3>
            </div>
            <span className="panel-pill">{todaySessions.length} classes</span>
          </div>
          <ScheduleList sessions={todaySessions} />
        </section>
      </div>

      <div className="mobile-expand-row">
        <section className="ops-panel mobile-expand-panel" aria-labelledby="mobile-members-title">
          <button
            type="button"
            aria-expanded={expandedMobilePanel === 'members'}
            onClick={() => setExpandedMobilePanel((current) => (current === 'members' ? null : 'members'))}
          >
            <span>
              <span className="panel-label">Members</span>
              <strong id="mobile-members-title">Latest registrations</strong>
            </span>
            <em>{expandedMobilePanel === 'members' ? 'Close' : 'Open'}</em>
          </button>
          {expandedMobilePanel === 'members' ? <MemberTable members={latestMembers} /> : null}
        </section>

        <section className="ops-panel mobile-expand-panel" aria-labelledby="mobile-sessions-title">
          <button
            type="button"
            aria-expanded={expandedMobilePanel === 'sessions'}
            onClick={() => setExpandedMobilePanel((current) => (current === 'sessions' ? null : 'sessions'))}
          >
            <span>
              <span className="panel-label">Classes</span>
              <strong id="mobile-sessions-title">{todaySessions.length} classes today</strong>
            </span>
            <em>{expandedMobilePanel === 'sessions' ? 'Close' : 'Open'}</em>
          </button>
          {expandedMobilePanel === 'sessions' ? <ScheduleList sessions={todaySessions} /> : null}
        </section>
      </div>
    </div>
  )
}

export default DashboardHome
