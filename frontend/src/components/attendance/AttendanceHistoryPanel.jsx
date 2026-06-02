import { useCallback, useEffect, useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import StateBlock from '../StateBlock.jsx'
import { listSubscriptionAttendance } from '../../lib/attendanceApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import {
  attendanceStatusLabel,
  compactId,
  formatDateTime,
  isObjectId,
} from './attendanceFormatters.js'

function AttendanceHistoryPanel({ accessToken, subscriptionId, refreshKey = 0 }) {
  const [state, setState] = useState({ status: 'idle', data: [], error: null })

  const loadHistory = useCallback(async () => {
    if (!subscriptionId) {
      setState({ status: 'idle', data: [], error: null })
      return
    }

    if (!isObjectId(subscriptionId)) {
      setState({
        status: 'error',
        data: [],
        error: { code: 'INVALID_ID', message: 'Subscription ID must be a 24 character ObjectID.' },
      })
      return
    }

    setState({ status: 'loading', data: [], error: null })

    try {
      const response = await listSubscriptionAttendance(accessToken, subscriptionId)
      setState({ status: 'success', data: response.data || [], error: null })
    } catch (error) {
      setState({ status: 'error', data: [], error })
    }
  }, [accessToken, subscriptionId])

  useEffect(() => {
    loadHistory()
  }, [loadHistory, refreshKey])

  return (
    <DataPanel title="Attendance history" action={<button className="btn-outline" type="button" onClick={loadHistory} disabled={!subscriptionId}>Refresh</button>}>
      {state.status === 'idle' ? <StateBlock tone="empty" title="Choose a subscription" message="Open a subscription attendance route to load history." /> : null}
      {state.status === 'loading' ? <StateBlock tone="loading" title="Loading history" message="Fetching attendance records." /> : null}
      {state.status === 'error' ? <StateBlock tone="error" title="Could not load history" message={apiErrorText(state.error)} /> : null}
      {state.status === 'success' && !state.data.length ? <StateBlock tone="empty" title="No attendance records" message="No records have been created for this subscription." /> : null}
      {state.status === 'success' && state.data.length ? (
        <div className="resource-list">
          {state.data.map((record) => (
            <article className="resource-row" key={record.id}>
              <div>
                <strong>{attendanceStatusLabel(record.status)}</strong>
                <span>{formatDateTime(record.date)} · Branch {compactId(record.branch_id)}</span>
                {record.is_makeup_for ? <small>Makeup for {formatDateTime(record.is_makeup_for)}</small> : null}
              </div>
              <div><span>{compactId(record.id)}</span></div>
            </article>
          ))}
        </div>
      ) : null}
    </DataPanel>
  )
}

export default AttendanceHistoryPanel
