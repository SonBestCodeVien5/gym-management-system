import { useCallback, useEffect, useState } from 'react'
import PageHeader from '../PageHeader.jsx'
import StateBlock from '../StateBlock.jsx'
import { useAuth } from '../../context/AuthContext.jsx'
import { getSubscription } from '../../lib/subscriptionsApi.js'
import { apiErrorText } from '../../lib/featureHelpers.js'
import RefundPanel from './RefundPanel.jsx'
import SubscriptionLifecyclePanel from './SubscriptionLifecyclePanel.jsx'
import SubscriptionLookupPanel from './SubscriptionLookupPanel.jsx'
import SubscriptionSummaryPanel from './SubscriptionSummaryPanel.jsx'
import { compactId, isObjectId } from './subscriptionFormatters.js'

function SubscriptionDetailView({ subscriptionId, navigate }) {
  const { accessToken } = useAuth()
  const [subscriptionState, setSubscriptionState] = useState({ status: 'loading', data: null, error: null })

  const loadSubscription = useCallback(async ({ background = false } = {}) => {
    if (!isObjectId(subscriptionId)) {
      setSubscriptionState({
        status: 'error',
        data: null,
        error: { code: 'INVALID_ID', message: 'Subscription ID must be a 24 character ObjectID.' },
      })
      return
    }

    if (!background) {
      setSubscriptionState((current) => ({ ...current, status: 'loading', error: null }))
    }

    try {
      const response = await getSubscription(accessToken, subscriptionId)
      setSubscriptionState({ status: 'success', data: response.data, error: null })
    } catch (error) {
      if (background) {
        setSubscriptionState((current) => ({ ...current, error }))
        return
      }

      setSubscriptionState({ status: 'error', data: null, error })
    }
  }, [accessToken, subscriptionId])

  useEffect(() => {
    loadSubscription()
  }, [loadSubscription])

  if (subscriptionState.status === 'loading') {
    return (
      <div className="module-page resource-workspace">
        <PageHeader eyebrow="Subscriptions" title="Subscription detail" description={`Loading ${compactId(subscriptionId)}.`} />
        <StateBlock tone="loading" title="Loading subscription" message="Fetching subscription detail from the API." />
      </div>
    )
  }

  if (subscriptionState.status === 'error') {
    return (
      <div className="module-page resource-workspace">
        <PageHeader
          eyebrow="Subscriptions"
          title={subscriptionState.error?.code === 'NOT_FOUND' ? 'Subscription not found' : 'Subscription lookup failed'}
          description="Use a valid subscription ObjectID or return to the command center."
          actions={<button className="btn-outline" type="button" onClick={() => navigate('/app/subscriptions')}>Subscriptions</button>}
        />
        <StateBlock
          tone={subscriptionState.error?.code === 'NOT_FOUND' ? 'notFound' : 'error'}
          title="Could not load subscription"
          message={apiErrorText(subscriptionState.error)}
          details={<SubscriptionLookupPanel navigate={navigate} initialValue={subscriptionId} />}
        />
      </div>
    )
  }

  const subscription = subscriptionState.data

  return (
    <div className="module-page resource-workspace subscriptions-workspace">
      <PageHeader
        eyebrow="Subscriptions"
        title="Subscription detail"
        description={`Subscription ID ${subscription.id}`}
        actions={<button className="btn-outline" type="button" onClick={() => navigate('/app/subscriptions')}>Subscriptions</button>}
      />

      <SubscriptionSummaryPanel subscription={subscription} navigate={navigate} />

      <div className="module-page__grid">
        <SubscriptionLifecyclePanel accessToken={accessToken} subscription={subscription} onChanged={() => loadSubscription({ background: true })} />
        <RefundPanel accessToken={accessToken} subscription={subscription} onChanged={() => loadSubscription({ background: true })} />
      </div>
    </div>
  )
}

export default SubscriptionDetailView
