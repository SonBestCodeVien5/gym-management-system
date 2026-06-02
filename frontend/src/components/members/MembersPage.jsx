import DataPanel from '../DataPanel.jsx'
import PageHeader from '../PageHeader.jsx'
import StateBlock from '../StateBlock.jsx'
import MemberLookupPanel from './MemberLookupPanel.jsx'

function MembersPage({ navigate }) {
  return (
    <div className="module-page members-workspace">
      <PageHeader
        eyebrow="Members"
        title="Members"
        description="Create member profiles and open existing profiles by direct ID."
        actions={(
          <button className="btn-outline" type="button" onClick={() => navigate('/app/members/new')}>
            New member
          </button>
        )}
      />

      <div className="module-page__grid">
        <DataPanel
          title="Open member"
          description="Use the backend ObjectID from a created profile, subscription record, or API sample."
        >
          <MemberLookupPanel navigate={navigate} />
        </DataPanel>

        <DataPanel
          title="Create profile"
          description="Register a new member record before creating subscriptions in the later subscription workspace."
          action={(
            <button className="btn-outline" type="button" onClick={() => navigate('/app/members/new')}>
              Create
            </button>
          )}
        >
          <ul className="feature-list">
            <li>CCID and full name are required.</li>
            <li>Email, phone, gender, and level are optional.</li>
            <li>New members start pending until offline payment is confirmed.</li>
          </ul>
        </DataPanel>
      </div>

      <DataPanel title="Directory search">
        <StateBlock
          tone="planned"
          title="List and search need backend support"
          message="The current API does not expose GET /api/v1/members or CCID search, so this page uses direct ID lookup without fake directory rows."
        />
      </DataPanel>
    </div>
  )
}

export default MembersPage
