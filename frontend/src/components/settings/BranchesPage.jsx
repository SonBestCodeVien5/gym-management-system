import { useCallback, useEffect, useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import PageHeader from '../PageHeader.jsx'
import StateBlock from '../StateBlock.jsx'
import { useAuth } from '../../context/AuthContext.jsx'
import { createBranch, deleteBranch, listBranches } from '../../lib/branchesApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import BranchForm from './BranchForm.jsx'
import NearbyBranchesPanel from './NearbyBranchesPanel.jsx'
import { compactId, formatCoordinates } from './settingsFormatters.js'

function BranchesPage({ navigate }) {
  const { accessToken } = useAuth()
  const [branchesState, setBranchesState] = useState({ status: 'loading', data: [], error: null })
  const [values, setValues] = useState(BranchForm.initialValues)
  const [errors, setErrors] = useState({})
  const [mutation, setMutation] = useState({ status: 'idle', error: null, notice: '' })

  const loadBranches = useCallback(async () => {
    setBranchesState((current) => ({ ...current, status: 'loading', error: null }))

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
    setMutation({ status: 'submitting', error: null, notice: '' })

    try {
      const response = await createBranch(accessToken, payload)
      const branchId = response.data?.id
      setValues(BranchForm.initialValues)
      setMutation({ status: 'success', error: null, notice: 'Branch created.' })

      if (branchId) {
        navigate(`/app/settings/branches/${branchId}`, { replace: true })
        return
      }

      await loadBranches()
    } catch (error) {
      setMutation({ status: 'error', error, notice: '' })
    }
  }

  async function handleDelete(branch) {
    if (!window.confirm(`Delete branch ${branch.name}?`)) {
      return
    }

    setMutation({ status: 'submitting', error: null, notice: '' })

    try {
      await deleteBranch(accessToken, branch.id)
      setMutation({ status: 'success', error: null, notice: 'Branch deleted.' })
      await loadBranches()
    } catch (error) {
      setMutation({ status: 'error', error, notice: '' })
    }
  }

  return (
    <div className="module-page resource-workspace settings-workspace">
      <PageHeader
        eyebrow="Settings"
        title="Branches"
        description="Manage branch reference data and numeric nearby search."
        actions={<button className="btn-outline" type="button" onClick={loadBranches}>Refresh</button>}
      />

      <div className="module-page__grid">
        <DataPanel title="Create branch" description="GeoJSON coordinates use longitude then latitude.">
          <BranchForm
            values={values}
            setValues={setValues}
            errors={errors}
            setErrors={setErrors}
            onSubmit={handleCreate}
            submitLabel="Create branch"
            submittingLabel="Creating"
            status={mutation.status}
          />
          {mutation.error ? <div className="form-alert" role="alert">{apiErrorText(mutation.error, 'Branch could not be saved.')}</div> : null}
          {mutation.notice ? <div className="form-success" role="status">{mutation.notice}</div> : null}
        </DataPanel>

        <DataPanel title="Nearby search" description="Search by coordinates without adding a map dependency.">
          <NearbyBranchesPanel accessToken={accessToken} navigate={navigate} />
        </DataPanel>
      </div>

      <DataPanel title="Branch list">
        {branchesState.status === 'loading' ? <StateBlock tone="loading" title="Loading branches" message="Fetching branches from the API." /> : null}
        {branchesState.status === 'error' ? <StateBlock tone="error" title="Could not load branches" message={apiErrorText(branchesState.error)} /> : null}
        {branchesState.status === 'success' && !branchesState.data.length ? <StateBlock tone="empty" title="No branches yet" message="Create the first branch above." /> : null}
        {branchesState.status === 'success' && branchesState.data.length ? (
          <div className="resource-list">
            {branchesState.data.map((branch) => (
              <article className="resource-row" key={branch.id}>
                <div>
                  <strong>{branch.name}</strong>
                  <span>{branch.branch_code} · {branch.province}</span>
                  <small>{formatCoordinates(branch.location)}</small>
                </div>
                <div>
                  <span>{compactId(branch.id)}</span>
                  <button className="btn-outline" type="button" onClick={() => navigate(`/app/settings/branches/${branch.id}`)}>Open</button>
                  <button className="btn-outline" type="button" onClick={() => handleDelete(branch)}>Delete</button>
                </div>
              </article>
            ))}
          </div>
        ) : null}
      </DataPanel>
    </div>
  )
}

export default BranchesPage
