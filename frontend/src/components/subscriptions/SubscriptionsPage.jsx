import DataPanel from '../DataPanel.jsx'
import PageHeader from '../PageHeader.jsx'
import StateBlock from '../StateBlock.jsx'
import SubscriptionLookupPanel from './SubscriptionLookupPanel.jsx'

function SubscriptionsPage({ navigate }) {
  return (
    <div className="module-page resource-workspace subscriptions-workspace">
      <PageHeader
        eyebrow="Subscriptions"
        title="Subscriptions"
        description="Create pending subscriptions and manage lifecycle actions by direct subscription ID."
        actions={<button className="btn-outline" type="button" onClick={() => navigate('/app/subscriptions/new')}>New subscription</button>}
      />

      <div className="module-page__grid">
        <DataPanel title="Open subscription" description="Use the backend ObjectID from member subscription history or a created subscription response.">
          <SubscriptionLookupPanel navigate={navigate} />
        </DataPanel>

        <DataPanel title="Create pending subscription" description="Course and branch options load from the live reference APIs.">
          <ul className="feature-list">
            <li>Member selection remains direct ObjectID because there is no member search endpoint.</li>
            <li>Payment activation remains in the member detail offline-payment panel.</li>
            <li>Lifecycle actions are available from the detail route.</li>
          </ul>
          <button className="btn-outline" type="button" onClick={() => navigate('/app/subscriptions/new')}>Create</button>
        </DataPanel>
      </div>

      <DataPanel title="Subscription directory">
        <StateBlock
          tone="planned"
          title="Global list not exposed"
          message="The backend does not expose GET /api/v1/subscriptions, so this workspace avoids fake directory rows and uses direct lookup."
        />
      </DataPanel>
    </div>
  )
}

export default SubscriptionsPage
