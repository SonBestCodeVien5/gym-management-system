import { useCallback, useEffect, useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import PageHeader from '../PageHeader.jsx'
import StateBlock from '../StateBlock.jsx'
import { useAuth } from '../../context/AuthContext.jsx'
import { deleteBranch, getBranch, updateBranch } from '../../lib/branchesApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import BranchForm, { branchValuesFromBranch } from './BranchForm.jsx'
import { compactId, formatCoordinates, formatText, isObjectId } from './settingsFormatters.js'

function BranchDetailView({ branchId, navigate }) {
  const { accessToken } = useAuth()
  const [branchState, setBranchState] = useState({ status: 'loading', data: null, error: null })
  const [values, setValues] = useState(BranchForm.initialValues)
  const [errors, setErrors] = useState({})
  const [mutation, setMutation] = useState({ status: 'idle', error: null, notice: '' })

  const loadBranch = useCallback(async () => {
    if (!isObjectId(branchId)) {
      setBranchState({
        status: 'error',
        data: null,
        error: { code: 'INVALID_ID', message: 'Branch ID must be a 24 character ObjectID.' },
      })
      return
    }

    setBranchState((current) => ({ ...current, status: 'loading', error: null }))

    try {
      const response = await getBranch(accessToken, branchId)
      setBranchState({ status: 'success', data: response.data, error: null })
      setValues(branchValuesFromBranch(response.data))
    } catch (error) {
      setBranchState({ status: 'error', data: null, error })
    }
  }, [accessToken, branchId])

  useEffect(() => {
    loadBranch()
  }, [loadBranch])

  async function handleUpdate(payload) {
    setMutation({ status: 'submitting', error: null, notice: '' })

    try {
      await updateBranch(accessToken, branchId, payload)
      setMutation({ status: 'success', error: null, notice: 'Branch updated.' })
      await loadBranch()
    } catch (error) {
      setMutation({ status: 'error', error, notice: '' })
    }
  }

  async function handleDelete() {
    if (!window.confirm('Delete this branch?')) {
      return
    }

    setMutation({ status: 'submitting', error: null, notice: '' })

    try {
      await deleteBranch(accessToken, branchId)
      navigate('/app/settings/branches', { replace: true })
    } catch (error) {
      setMutation({ status: 'error', error, notice: '' })
    }
  }

  if (branchState.status === 'loading') {
    return (
      <div className="module-page resource-workspace">
        <PageHeader eyebrow="Settings" title="Branch detail" description={`Loading ${compactId(branchId)}.`} />
        <StateBlock tone="loading" title="Loading branch" message="Fetching branch detail from the API." />
      </div>
    )
  }

  if (branchState.status === 'error') {
    return (
      <div className="module-page resource-workspace">
        <PageHeader
          eyebrow="Settings"
          title={branchState.error?.code === 'NOT_FOUND' ? 'Branch not found' : 'Branch lookup failed'}
          description="Return to branches or use a valid ObjectID."
          actions={<button className="btn-outline" type="button" onClick={() => navigate('/app/settings/branches')}>Branches</button>}
        />
        <StateBlock tone={branchState.error?.code === 'NOT_FOUND' ? 'notFound' : 'error'} title="Could not load branch" message={apiErrorText(branchState.error)} />
      </div>
    )
  }

  const branch = branchState.data

  return (
    <div className="module-page resource-workspace">
      <PageHeader
        eyebrow="Settings"
        title={branch.name}
        description={`Branch ID ${branch.id}`}
        actions={<button className="btn-outline" type="button" onClick={() => navigate('/app/settings/branches')}>Branches</button>}
      />

      <div className="module-page__grid">
        <DataPanel title="Summary">
          <dl className="detail-grid">
            <div><dt>Code</dt><dd>{branch.branch_code}</dd></div>
            <div><dt>Province</dt><dd>{branch.province}</dd></div>
            <div><dt>Coordinates</dt><dd>{formatCoordinates(branch.location)}</dd></div>
            <div><dt>Manager</dt><dd>{formatText(branch.manager_id)}</dd></div>
          </dl>
          <p className="panel-copy">{branch.address}</p>
        </DataPanel>

        <DataPanel title="Danger zone">
          <p className="panel-copy">Delete is permanent and may be rejected by backend reference rules.</p>
          <button className="btn-outline" type="button" onClick={handleDelete} disabled={mutation.status === 'submitting'}>
            Delete branch
          </button>
        </DataPanel>
      </div>

      <DataPanel title="Edit branch">
        <BranchForm
          values={values}
          setValues={setValues}
          errors={errors}
          setErrors={setErrors}
          onSubmit={handleUpdate}
          submitLabel="Save branch"
          submittingLabel="Saving"
          status={mutation.status}
        />
        {mutation.error ? <div className="form-alert" role="alert">{apiErrorText(mutation.error, 'Branch could not be updated.')}</div> : null}
        {mutation.notice ? <div className="form-success" role="status">{mutation.notice}</div> : null}
      </DataPanel>
    </div>
  )
}

export default BranchDetailView
