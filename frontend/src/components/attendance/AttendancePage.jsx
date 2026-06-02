import { useCallback, useEffect, useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import PageHeader from '../PageHeader.jsx'
import { useAuth } from '../../context/AuthContext.jsx'
import { listBranches } from '../../lib/branchesApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import CheckInPanel from './CheckInPanel.jsx'
import ReportMissedPanel from './ReportMissedPanel.jsx'
import MakeupPanel from './MakeupPanel.jsx'
import SubscriptionHistoryLookup from './SubscriptionHistoryLookup.jsx'

function AttendancePage({ navigate }) {
  const { accessToken } = useAuth()
  const [branchesState, setBranchesState] = useState({ status: 'loading', data: [], error: null })

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

  const branches = branchesState.data

  return (
    <div className="module-page resource-workspace attendance-workspace">
      <PageHeader
        eyebrow="Attendance"
        title="Attendance"
        description="Record free check-ins, report missed sessions, create makeup attendance, and open subscription history."
      />

      {branchesState.status === 'error' ? (
        <div className="form-alert" role="alert">Branch options failed to load. Manual branch ObjectID entry is still available. {apiErrorText(branchesState.error)}</div>
      ) : null}

      <DataPanel title="History lookup" description="Attendance history is subscription-scoped.">
        <SubscriptionHistoryLookup navigate={navigate} />
      </DataPanel>

      <div className="module-page__grid">
        <CheckInPanel accessToken={accessToken} branches={branches} />
        <ReportMissedPanel accessToken={accessToken} branches={branches} />
      </div>

      <MakeupPanel accessToken={accessToken} branches={branches} />
    </div>
  )
}

export default AttendancePage
