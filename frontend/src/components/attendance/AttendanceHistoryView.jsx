import { useCallback, useEffect, useState } from 'react'
import PageHeader from '../PageHeader.jsx'
import { useAuth } from '../../context/AuthContext.jsx'
import { listBranches } from '../../lib/branchesApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import AttendanceHistoryPanel from './AttendanceHistoryPanel.jsx'
import CheckInPanel from './CheckInPanel.jsx'
import MakeupPanel from './MakeupPanel.jsx'
import ReportMissedPanel from './ReportMissedPanel.jsx'

function AttendanceHistoryView({ subscriptionId, navigate }) {
  const { accessToken } = useAuth()
  const [branchesState, setBranchesState] = useState({ status: 'loading', data: [], error: null })
  const [refreshKey, setRefreshKey] = useState(0)

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

  const refreshHistory = () => setRefreshKey((current) => current + 1)
  const branches = branchesState.data

  return (
    <div className="module-page resource-workspace attendance-workspace">
      <PageHeader
        eyebrow="Attendance"
        title="Subscription attendance"
        description={`History and quick actions for subscription ${subscriptionId}.`}
        actions={<button className="btn-outline" type="button" onClick={() => navigate(`/app/subscriptions/${subscriptionId}`)}>Subscription</button>}
      />

      {branchesState.status === 'error' ? (
        <div className="form-alert" role="alert">Branch options failed to load. Manual branch ObjectID entry is still available. {apiErrorText(branchesState.error)}</div>
      ) : null}

      <AttendanceHistoryPanel accessToken={accessToken} subscriptionId={subscriptionId} refreshKey={refreshKey} />

      <div className="module-page__grid">
        <CheckInPanel accessToken={accessToken} branches={branches} subscriptionId={subscriptionId} onSuccess={refreshHistory} />
        <ReportMissedPanel accessToken={accessToken} branches={branches} subscriptionId={subscriptionId} onSuccess={refreshHistory} />
      </div>

      <MakeupPanel accessToken={accessToken} branches={branches} subscriptionId={subscriptionId} onSuccess={refreshHistory} />
    </div>
  )
}

export default AttendanceHistoryView
