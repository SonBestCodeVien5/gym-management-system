import { ROLE_OPTIONS } from './employeeFormatters.js'

function RoleSelector({ value, onChange, error }) {
  function toggleRole(role) {
    if (value.includes(role)) {
      onChange(value.filter((item) => item !== role))
      return
    }

    onChange([...value, role])
  }

  return (
    <fieldset className="role-selector">
      <legend>Roles</legend>
      <div className="role-selector__options">
        {ROLE_OPTIONS.map((role) => (
          <label key={role}>
            <input
              type="checkbox"
              checked={value.includes(role)}
              onChange={() => toggleRole(role)}
            />
            <span>{role}</span>
          </label>
        ))}
      </div>
      {error ? <span>{error}</span> : null}
    </fieldset>
  )
}

export default RoleSelector
