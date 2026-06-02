import { useCallback, useEffect, useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import PageHeader from '../PageHeader.jsx'
import StateBlock from '../StateBlock.jsx'
import { useAuth } from '../../context/AuthContext.jsx'
import { listBranches } from '../../lib/branchesApi.js'
import { listEmployees } from '../../lib/employeesApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import EmployeeFilters from './EmployeeFilters.jsx'
import { compactId, formatDateTime, roleLabel } from './employeeFormatters.js'

function EmployeesPage({ navigate }) {
  const { accessToken } = useAuth()
  const [filters, setFilters] = useState({})
  const [employeesState, setEmployeesState] = useState({ status: 'loading', data: [], error: null })
  const [branchesState, setBranchesState] = useState({ status: 'loading', data: [], error: null })

  const loadEmployees = useCallback(async () => {
    setEmployeesState((current) => ({ ...current, status: 'loading', error: null }))

    try {
      const response = await listEmployees(accessToken, filters)
      setEmployeesState({ status: 'success', data: response.data || [], error: null })
    } catch (error) {
      setEmployeesState({ status: 'error', data: [], error })
    }
  }, [accessToken, filters])

  useEffect(() => {
    loadEmployees()
  }, [loadEmployees])

  useEffect(() => {
    async function loadBranches() {
      try {
        const response = await listBranches(accessToken)
        setBranchesState({ status: 'success', data: response.data || [], error: null })
      } catch (error) {
        setBranchesState({ status: 'error', data: [], error })
      }
    }

    loadBranches()
  }, [accessToken])

  return (
    <div className="module-page resource-workspace employees-workspace">
      <PageHeader
        eyebrow="Employees"
        title="Employees"
        description="Admin-only staff account management, profile updates, deactivation, and password reset."
        actions={<button className="btn-outline" type="button" onClick={() => navigate('/app/employees/new')}>New employee</button>}
      />

      {branchesState.status === 'error' ? <div className="form-alert" role="alert">Branch options failed to load. {apiErrorText(branchesState.error)}</div> : null}

      <DataPanel title="Filters">
        <EmployeeFilters branches={branchesState.data} onApply={setFilters} />
      </DataPanel>

      <DataPanel title="Staff list" action={<button className="btn-outline" type="button" onClick={loadEmployees}>Refresh</button>}>
        {employeesState.status === 'loading' ? <StateBlock tone="loading" title="Loading employees" message="Fetching staff accounts from the API." /> : null}
        {employeesState.status === 'error' ? <StateBlock tone="error" title="Could not load employees" message={apiErrorText(employeesState.error)} /> : null}
        {employeesState.status === 'success' && !employeesState.data.length ? <StateBlock tone="empty" title="No employees found" message="Adjust filters or create a staff account." /> : null}
        {employeesState.status === 'success' && employeesState.data.length ? (
          <div className="resource-list">
            {employeesState.data.map((employee) => (
              <article className="resource-row" key={employee.id}>
                <div>
                  <strong>{employee.full_name}</strong>
                  <span>{employee.employee_id} · {roleLabel(employee.role)} · {employee.status}</span>
                  <small>{employee.email} · updated {formatDateTime(employee.updated_at)}</small>
                </div>
                <div>
                  <span>{compactId(employee.id)}</span>
                  <button className="btn-outline" type="button" onClick={() => navigate(`/app/employees/${employee.id}`)}>Open</button>
                </div>
              </article>
            ))}
          </div>
        ) : null}
      </DataPanel>
    </div>
  )
}

export default EmployeesPage
