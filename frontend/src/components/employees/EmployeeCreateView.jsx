import { useCallback, useEffect, useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import PageHeader from '../PageHeader.jsx'
import { useAuth } from '../../context/AuthContext.jsx'
import { listBranches } from '../../lib/branchesApi.js'
import { createEmployee } from '../../lib/employeesApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import EmployeeForm, { EMPTY_EMPLOYEE_VALUES } from './EmployeeForm.jsx'

function EmployeeCreateView({ navigate }) {
  const { accessToken } = useAuth()
  const [values, setValues] = useState(EMPTY_EMPLOYEE_VALUES)
  const [errors, setErrors] = useState({})
  const [branchesState, setBranchesState] = useState({ status: 'loading', data: [], error: null })
  const [submitState, setSubmitState] = useState({ status: 'idle', error: null })

  const loadBranches = useCallback(async () => {
    try {
      const response = await listBranches(accessToken)
      setBranchesState({ status: 'success', data: response.data || [], error: null })
    } catch (error) {
      setBranchesState({ status: 'error', data: [], error })
    }
  }, [accessToken])

  useEffect(() => {
    loadBranches()
  }, [loadBranches])

  async function handleCreate(payload) {
    setSubmitState({ status: 'submitting', error: null })

    try {
      const response = await createEmployee(accessToken, payload)
      const employeeId = response.data?.id
      setValues(EMPTY_EMPLOYEE_VALUES)

      if (employeeId) {
        navigate(`/app/employees/${employeeId}`, { replace: true })
        return
      }

      setSubmitState({ status: 'success', error: null })
    } catch (error) {
      setSubmitState({ status: 'error', error })
    }
  }

  return (
    <div className="module-page resource-workspace employees-workspace">
      <PageHeader
        eyebrow="Employees"
        title="New employee"
        description="Create an admin-managed staff account with an initial password."
        actions={<button className="btn-outline" type="button" onClick={() => navigate('/app/employees')}>Employees</button>}
      />

      {branchesState.status === 'error' ? <div className="form-alert" role="alert">Branch options failed to load. Manual branch IDs are still available. {apiErrorText(branchesState.error)}</div> : null}

      <DataPanel title="Account details">
        <EmployeeForm
          values={values}
          setValues={setValues}
          errors={errors}
          setErrors={setErrors}
          branches={branchesState.data}
          onSubmit={handleCreate}
          submitLabel="Create employee"
          submittingLabel="Creating"
          status={submitState.status}
          requirePassword
          onCancel={() => navigate('/app/employees')}
        />
        {submitState.error ? <div className="form-alert" role="alert">{apiErrorText(submitState.error, 'Employee could not be created.')}</div> : null}
        {submitState.status === 'success' ? <div className="form-success" role="status">Employee created, but response did not include an ID.</div> : null}
      </DataPanel>
    </div>
  )
}

export default EmployeeCreateView
