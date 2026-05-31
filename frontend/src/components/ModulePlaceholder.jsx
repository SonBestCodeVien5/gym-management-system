import { formatRoles } from '../lib/permissions.js'
import DataPanel from './DataPanel.jsx'
import PageHeader from './PageHeader.jsx'
import StateBlock from './StateBlock.jsx'

function ModulePlaceholder({ route, params = {} }) {
  const plannedApis = route.plannedApis || []
  const plannedScope = route.plannedScope || []
  const hasParams = Object.keys(params).length > 0

  return (
    <div className="module-page">
      <PageHeader
        eyebrow={route.group || 'Module'}
        title={route.title}
        description={route.description}
      />

      <DataPanel
        title="Module status"
        description={route.status === 'blocked' ? 'Waiting for later backend or product scope.' : 'Foundation route is ready; workflow implementation comes later.'}
      >
        <StateBlock
          tone={route.status === 'blocked' ? 'planned' : 'empty'}
          title={route.status === 'blocked' ? 'Blocked for later scope' : 'Placeholder ready'}
          message={`Allowed roles: ${formatRoles(route.roles)}.`}
          details={hasParams ? <p>Route params: {JSON.stringify(params)}</p> : null}
        />
      </DataPanel>

      <div className="module-page__grid">
        <DataPanel title="Planned workflow">
          {plannedScope.length ? (
            <ul className="feature-list">
              {plannedScope.map((item) => (
                <li key={item}>{item}</li>
              ))}
            </ul>
          ) : (
            <p className="panel-copy">Workflow details will be defined in its feature plan.</p>
          )}
        </DataPanel>

        <DataPanel title="Planned API surface">
          {plannedApis.length ? (
            <ul className="api-list">
              {plannedApis.map((endpoint) => (
                <li key={endpoint}>{endpoint}</li>
              ))}
            </ul>
          ) : (
            <p className="panel-copy">No business API planned for this placeholder yet.</p>
          )}
        </DataPanel>
      </div>
    </div>
  )
}

export default ModulePlaceholder
