import { useCallback, useEffect, useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import PageHeader from '../PageHeader.jsx'
import StateBlock from '../StateBlock.jsx'
import { useAuth } from '../../context/AuthContext.jsx'
import { listBranches } from '../../lib/branchesApi.js'
import { listSessions } from '../../lib/sessionsApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import SessionFilters from './SessionFilters.jsx'
import SessionList from './SessionList.jsx'

function SessionsPage({ navigate }) {
  const { accessToken } = useAuth()
  const [filters, setFilters] = useState({})
  const [sessionsState, setSessionsState] = useState({ status: 'loading', data: [], error: null })
  const [branchesState, setBranchesState] = useState({ status: 'loading', data: [], error: null })

  const loadSessions = useCallback(async () => {
    setSessionsState((current) => ({ ...current, status: 'loading', error: null }))

    try {
      const response = await listSessions(accessToken, filters)
      setSessionsState({ status: 'success', data: response.data || [], error: null })
    } catch (error) {
      setSessionsState({ status: 'error', data: [], error })
    }
  }, [accessToken, filters])

  useEffect(() => {
    loadSessions()
  }, [loadSessions])

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
    <div className="module-page resource-workspace sessions-workspace">
      <PageHeader
        eyebrow="Sessions"
        title="Sessions"
        description="Filter session schedules, create classes, enroll subscriptions, and record session check-ins."
        actions={<button className="btn-outline" type="button" onClick={() => navigate('/app/sessions/new')}>New session</button>}
      />

      {branchesState.status === 'error' ? <div className="form-alert" role="alert">Branch options failed to load. {apiErrorText(branchesState.error)}</div> : null}

      <DataPanel title="Filters">
        <SessionFilters branches={branchesState.data} onApply={setFilters} />
      </DataPanel>

      <DataPanel title="Session list" action={<button className="btn-outline" type="button" onClick={loadSessions}>Refresh</button>}>
        {sessionsState.status === 'loading' ? <StateBlock tone="loading" title="Loading sessions" message="Fetching sessions from the API." /> : null}
        {sessionsState.status === 'error' ? <StateBlock tone="error" title="Could not load sessions" message={apiErrorText(sessionsState.error)} /> : null}
        {sessionsState.status === 'success' && !sessionsState.data.length ? <StateBlock tone="empty" title="No sessions found" message="Adjust filters or create a new session." /> : null}
        {sessionsState.status === 'success' && sessionsState.data.length ? <SessionList sessions={sessionsState.data} navigate={navigate} /> : null}
      </DataPanel>
    </div>
  )
}

export default SessionsPage
