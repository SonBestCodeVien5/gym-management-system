import { useState } from 'react'
import { API_BASE_URL } from '../lib/api.js'
import { useAuth } from '../context/AuthContext.jsx'
import BrandMark from './BrandMark.jsx'

function LoginView({ navigate }) {
  const { login, notice, clearAuthMessage } = useAuth()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [formError, setFormError] = useState(null)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [submitted, setSubmitted] = useState(false)

  const trimmedEmail = email.trim()
  const emailError = submitted && !trimmedEmail ? 'Email la bat buoc.' : ''
  const passwordError = submitted && !password
    ? 'Mat khau la bat buoc.'
    : ''

  async function handleSubmit(event) {
    event.preventDefault()
    setSubmitted(true)
    setFormError(null)
    clearAuthMessage()

    if (!trimmedEmail || !password) {
      return
    }

    setIsSubmitting(true)

    try {
      await login({ email: trimmedEmail, password })
      navigate('/app', { replace: true })
    } catch (error) {
      setFormError(error.message || 'Dang nhap khong thanh cong.')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <main className="login-page grid-bg">
      <section className="login-panel" aria-labelledby="login-title">
        <div className="login-brand">
          <BrandMark meta="Staff Portal" />
        </div>

        <div className="login-heading">
          <p className="section-eyebrow">Nhan vien</p>
          <h1 id="login-title">Dang nhap van hanh phong tap</h1>
          <p>Su dung tai khoan employee de truy cap dashboard noi bo.</p>
        </div>

        <form className="login-form" onSubmit={handleSubmit} noValidate>
          <div className="field-group">
            <label htmlFor="email">Email</label>
            <input
              id="email"
              name="email"
              type="email"
              autoComplete="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              aria-invalid={Boolean(emailError)}
              aria-describedby={emailError ? 'email-error' : undefined}
              disabled={isSubmitting}
            />
            {emailError ? <span id="email-error">{emailError}</span> : null}
          </div>

          <div className="field-group">
            <label htmlFor="password">Mat khau</label>
            <input
              id="password"
              name="password"
              type="password"
              autoComplete="current-password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              aria-invalid={Boolean(passwordError)}
              aria-describedby={passwordError ? 'password-error' : undefined}
              disabled={isSubmitting}
            />
            {passwordError ? <span id="password-error">{passwordError}</span> : null}
          </div>

          {(formError || notice) ? (
            <p className="form-alert" aria-live="polite">
              {formError || notice}
            </p>
          ) : null}

          <button className="btn-primary" type="submit" disabled={isSubmitting}>
            {isSubmitting ? 'Dang xu ly' : 'Dang nhap'}
          </button>
        </form>

        <footer className="login-meta">
          <span>API</span>
          <code>{API_BASE_URL}</code>
        </footer>
      </section>
    </main>
  )
}

export default LoginView
