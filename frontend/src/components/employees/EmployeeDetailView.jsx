import { useCallback, useEffect, useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import PageHeader from '../PageHeader.jsx'
import StateBlock from '../StateBlock.jsx'
import { useAuth } from '../../context/AuthContext.jsx'
import { listBranches } from '../../lib/branchesApi.js'
import { getEmployee, updateEmployee } from '../../lib/employeesApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import EmployeeForm from './EmployeeForm.jsx'
import PasswordResetPanel from './PasswordResetPanel.jsx'
import {
  cleanEmployeePayload,
  compactId,
  employeeValuesFromEmployee,
  formatDateTime,
  isObjectId,
  roleLabel,
} from './employeeFormatters.js'

function EmployeeDetailView({ employeeId, navigate }) {
  const { accessToken } = useAuth()
  const [employeeState, setEmployeeState] = useState({ status: 'loading', data: null, error: null })
  const [values, setValues] = useState(employeeValuesFromEmployee(null))
  const [errors, setErrors] = useState({})
  const [branchesState, setBranchesState] = useState({ status: 'loading', data: [], error: null })
  const [mutation, setMutation] = useState({ status: 'idle', error: null, notice: '' })

  const loadEmployee = useCallback(async () => {
    if (!isObjectId(employeeId)) {
      setEmployeeState({
        status: 'error',
        data: null,
        error: { code: 'INVALID_ID', message: 'Employee ID must be a 24 character ObjectID.' },
      })
      return
    }

    setEmployeeState((current) => ({ ...current, status: 'loading', error: null }))

    try {
      const response = await getEmployee(accessToken, employeeId)
      setEmployeeState({ status: 'success', data: response.data, error: null })
      setValues(employeeValuesFromEmployee(response.data))
    } catch (error) {
      setEmployeeState({ status: 'error', data: null, error })
    }
  }, [accessToken, employeeId])

  useEffect(() => {
    loadEmployee()
  }, [loadEmployee])

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

  async function handleUpdate(payload) {
    const originalPayload = cleanEmployeePayload(employeeValuesFromEmployee(employeeState.data))

    if (JSON.stringify(payload) === JSON.stringify(originalPayload)) {
      setMutation({ status: 'success', error: null, notice: 'No employee changes to save.' })
      return
    }

    if (employeeState.data?.status === 'active' && payload.status === 'inactive' && !window.confirm('Deactivate this employee and revoke active refresh tokens?')) {
      return
    }

    setMutation({ status: 'submitting', error: null, notice: '' })

    try {
      await updateEmployee(accessToken, employeeId, payload)
      setMutation({ status: 'success', error: null, notice: 'Employee updated.' })
      await loadEmployee()
    } catch (error) {
      setMutation({ status: 'error', error, notice: '' })
    }
  }

  if (employeeState.status === 'loading') {
    return (
      <div className="module-page resource-workspace">
        <PageHeader eyebrow="Employees" title="Employee detail" description={`Loading ${compactId(employeeId)}.`} />
        <StateBlock tone="loading" title="Loading employee" message="Fetching staff account from the API." />
      </div>
    )
  }

  if (employeeState.status === 'error') {
    return (
      <div className="module-page resource-workspace">
        <PageHeader
          eyebrow="Employees"
          title={employeeState.error?.code === 'NOT_FOUND' ? 'Employee not found' : 'Employee lookup failed'}
          description="Return to employees or use a valid ObjectID."
          actions={<button className="btn-outline" type="button" onClick={() => navigate('/app/employees')}>Employees</button>}
        />
        <StateBlock tone={employeeState.error?.code === 'NOT_FOUND' ? 'notFound' : 'error'} title="Could not load employee" message={apiErrorText(employeeState.error)} />
      </div>
    )
  }

  const employee = employeeState.data

  return (
    <div className="module-page resource-workspace employees-workspace">
      <PageHeader
        eyebrow="Employees"
        title={employee.full_name}
        description={`Employee ID ${employee.id}`}
        actions={<button className="btn-outline" type="button" onClick={() => navigate('/app/employees')}>Employees</button>}
      />

      {branchesState.status === 'error' ? <div className="form-alert" role="alert">Branch options failed to load. Manual branch IDs are still available. {apiErrorText(branchesState.error)}</div> : null}

      <DataPanel title="Profile summary">
        <dl className="detail-grid">
          <div><dt>Staff code</dt><dd>{employee.employee_id}</dd></div>
          <div><dt>Email</dt><dd>{employee.email}</dd></div>
          <div><dt>Roles</dt><dd>{roleLabel(employee.role)}</dd></div>
          <div><dt>Status</dt><dd>{employee.status}</dd></div>
          <div><dt>Level</dt><dd>{employee.level || 'Not set'}</dd></div>
          <div><dt>Updated</dt><dd>{formatDateTime(employee.updated_at)}</dd></div>
        </dl>
      </DataPanel>

      <DataPanel title="Edit employee" description="Status inactive deactivates the account and backend revokes active refresh tokens.">
        <EmployeeForm
          values={values}
          setValues={setValues}
          errors={errors}
          setErrors={setErrors}
          branches={branchesState.data}
          onSubmit={handleUpdate}
          submitLabel="Save employee"
          submittingLabel="Saving"
          status={mutation.status}
        />
        {mutation.error ? <div className="form-alert" role="alert">{apiErrorText(mutation.error, 'Employee could not be updated.')}</div> : null}
        {mutation.notice ? <div className="form-success" role="status">{mutation.notice}</div> : null}
      </DataPanel>

      <PasswordResetPanel accessToken={accessToken} employeeId={employee.id} />
    </div>
  )
}

export default EmployeeDetailView
