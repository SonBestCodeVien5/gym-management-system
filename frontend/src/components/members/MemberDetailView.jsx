import { useCallback, useEffect, useState } from 'react'
import DataPanel from '../DataPanel.jsx'
import PageHeader from '../PageHeader.jsx'
import StateBlock from '../StateBlock.jsx'
import { useAuth } from '../../context/AuthContext.jsx'
import { getMember, listMemberSubscriptions } from '../../lib/membersApi.js'
import MemberLookupPanel from './MemberLookupPanel.jsx'
import MemberProfilePanel from './MemberProfilePanel.jsx'
import MemberSubscriptionsPanel from './MemberSubscriptionsPanel.jsx'
import OfflinePaymentPanel from './OfflinePaymentPanel.jsx'
import { compactId, isObjectId } from './memberFormatters.js'

function MemberDetailView({ memberId, navigate }) {
  const { accessToken } = useAuth()
  const [memberState, setMemberState] = useState({ status: 'loading', data: null, error: null })
  const [subscriptionsState, setSubscriptionsState] = useState({ status: 'loading', data: [], error: null })
  const [selectedSubscriptionId, setSelectedSubscriptionId] = useState('')
  const [activationNotice, setActivationNotice] = useState('')

  const loadMember = useCallback(async ({ background = false } = {}) => {
    if (!isObjectId(memberId)) {
      setMemberState({
        status: 'error',
        data: null,
        error: { code: 'INVALID_ID', message: 'Member ID must be a 24 character hex ObjectID.' },
      })
      setSubscriptionsState({ status: 'idle', data: [], error: null })
      return
    }

    if (!background) {
      setMemberState((current) => ({ ...current, status: 'loading', error: null }))
      setSelectedSubscriptionId('')
      setActivationNotice('')
    }

    setSubscriptionsState((current) => ({ ...current, status: 'loading', error: null }))

    try {
      const memberResponse = await getMember(accessToken, memberId)
      setMemberState({ status: 'success', data: memberResponse.data, error: null })
    } catch (error) {
      if (background) {
        setSubscriptionsState({ status: 'error', data: [], error })
        setActivationNotice('Offline payment was submitted, but the latest member data could not be refreshed.')
        return
      }

      setMemberState({ status: 'error', data: null, error })
      setSubscriptionsState({ status: 'idle', data: [], error: null })
      return
    }

    try {
      const subscriptionsResponse = await listMemberSubscriptions(accessToken, memberId)
      setSubscriptionsState({
        status: 'success',
        data: subscriptionsResponse.data || [],
        error: null,
      })
    } catch (error) {
      setSubscriptionsState({ status: 'error', data: [], error })
    }
  }, [accessToken, memberId])

  const handleActivated = useCallback(async () => {
    setActivationNotice('Offline payment confirmed. Member and subscriptions were refreshed.')
    await loadMember({ background: true })
  }, [loadMember])

  useEffect(() => {
    loadMember()
  }, [loadMember])

  const member = memberState.data
  const subscriptions = subscriptionsState.data || []

  if (memberState.status === 'loading') {
    return (
      <div className="module-page members-workspace">
        <PageHeader eyebrow="Members" title="Member detail" description={`Loading ${compactId(memberId)}.`} />
        <StateBlock tone="loading" title="Loading member" message="Fetching member profile from the API." />
      </div>
    )
  }

  if (memberState.status === 'error') {
    const isNotFound = memberState.error?.code === 'NOT_FOUND'

    return (
      <div className="module-page members-workspace">
        <PageHeader
          eyebrow="Members"
          title={isNotFound ? 'Member not found' : 'Member lookup failed'}
          description="Use a valid member ObjectID or return to the member command center."
          actions={(
            <button className="btn-outline" type="button" onClick={() => navigate('/app/members')}>
              Members
            </button>
          )}
        />

        <DataPanel title="Lookup">
          <StateBlock
            tone={isNotFound ? 'notFound' : 'error'}
            title={isNotFound ? 'No member for this ID' : 'Could not load member'}
            message={memberState.error?.message || 'Member profile could not be loaded.'}
            details={<MemberLookupPanel navigate={navigate} initialValue={memberId} />}
          />
        </DataPanel>
      </div>
    )
  }

  return (
    <div className="module-page members-workspace">
      <PageHeader
        eyebrow="Members"
        title={member.full_name}
        description={`Member ID ${member.id}`}
        actions={(
          <button className="btn-outline" type="button" onClick={() => navigate('/app/members')}>
            Members
          </button>
        )}
      />

      {activationNotice ? (
        <div className="form-success member-activation-notice" role="status" aria-live="polite">
          {activationNotice}
        </div>
      ) : null}

      <div className="member-detail-layout">
        <MemberProfilePanel member={member} />
        <OfflinePaymentPanel
          accessToken={accessToken}
          memberId={member.id}
          subscriptions={subscriptions}
          selectedSubscriptionId={selectedSubscriptionId}
          onSelectSubscription={setSelectedSubscriptionId}
          onActivated={handleActivated}
        />
      </div>

      <MemberSubscriptionsPanel
        subscriptions={subscriptions}
        status={subscriptionsState.status}
        error={subscriptionsState.error}
        selectedSubscriptionId={selectedSubscriptionId}
        onSelectSubscription={setSelectedSubscriptionId}
      />
    </div>
  )
}

export default MemberDetailView
