import type { FormEvent } from 'react'
import TotpQrCode from './TotpQrCode'

type AuthPanelProps = {
  email: string
  password: string
  setEmail: (value: string) => void
  setPassword: (value: string) => void
  totpCode: string
  setTotpCode: (value: string) => void
  enrollCode: string
  setEnrollCode: (value: string) => void
  enrollOtpAuthUrl: string
  hasTotpEnabled: boolean
  isEmailVerified: boolean
  needsTotpSignIn: boolean
  needsTotpEnrollment: boolean
  canStartTotpEnrollment: boolean
  isBusy: boolean
  isSignedIn: boolean
  authStatus: string
  authError: string
  signIn: () => Promise<void>
  verifyTotpSignIn: () => Promise<void>
  register: () => Promise<void>
  startTotpEnrollment: () => Promise<void>
  enrollTotp: () => Promise<void>
  resendVerification: () => Promise<void>
  signOut: () => Promise<void>
}

function AuthPanel(props: AuthPanelProps) {
  const handleSignIn = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    await props.signIn()
  }

  const handleRegister = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    await props.register()
  }

  const handleVerifyTotp = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    await props.verifyTotpSignIn()
  }

  const handleEnrollTotp = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    await props.enrollTotp()
  }

  return (
    <section className="auth-box">
      <p className="auth-title">Firebase Auth</p>

      <form className="auth-grid" onSubmit={handleSignIn}>
        <label htmlFor="email">Email</label>
        <input
          id="email"
          type="email"
          value={props.email}
          onChange={(event) => props.setEmail(event.target.value)}
          placeholder="name@example.com"
          autoComplete="email"
          required
        />

        <label htmlFor="password">Password</label>
        <input
          id="password"
          type="password"
          value={props.password}
          onChange={(event) => props.setPassword(event.target.value)}
          placeholder="Your password"
          autoComplete="current-password"
          required
        />

        <div className="auth-actions">
          <button type="submit" disabled={props.isBusy || !props.email.trim() || !props.password}>
            {props.isBusy ? 'Working...' : 'Sign in'}
          </button>
        </div>
      </form>

      <form className="auth-grid" onSubmit={handleRegister}>
        <div className="auth-actions">
          <button type="submit" disabled={props.isBusy || !props.email.trim() || !props.password}>
            {props.isBusy ? 'Working...' : 'Create account'}
          </button>
          {props.isSignedIn && !props.isEmailVerified ? (
            <button type="button" onClick={() => void props.resendVerification()} disabled={props.isBusy}>
              Resend verification email
            </button>
          ) : null}
          {props.canStartTotpEnrollment ? (
            <button type="button" onClick={() => void props.startTotpEnrollment()} disabled={props.isBusy}>
              Set up TOTP
            </button>
          ) : null}
          {props.isSignedIn ? (
            <button type="button" onClick={() => void props.signOut()} disabled={props.isBusy}>
              Sign out
            </button>
          ) : null}
        </div>
      </form>

      {props.isSignedIn && props.hasTotpEnabled ? (
        <p className="auth-help">TOTP is already enabled for this account.</p>
      ) : null}

      {props.needsTotpSignIn ? (
        <form className="auth-grid auth-grid-totp" onSubmit={handleVerifyTotp}>
          <label htmlFor="totp">Authenticator code</label>
          <input
            id="totp"
            type="text"
            inputMode="numeric"
            value={props.totpCode}
            onChange={(event) => props.setTotpCode(event.target.value)}
            placeholder="123456"
          />
          <button type="submit" disabled={props.isBusy || !props.totpCode.trim()}>
            {props.isBusy ? 'Verifying...' : 'Verify TOTP'}
          </button>
        </form>
      ) : null}

      {props.needsTotpEnrollment ? (
        <form className="auth-grid auth-grid-totp" onSubmit={handleEnrollTotp}>
          <p className="auth-help">Scan this QR code in your authenticator app:</p>
          <TotpQrCode otpauthUrl={props.enrollOtpAuthUrl} />
          <p className="auth-help">Or copy setup URL manually:</p>
          <input type="text" value={props.enrollOtpAuthUrl} readOnly />
          <label htmlFor="enroll-totp">Enter code from app</label>
          <input
            id="enroll-totp"
            type="text"
            inputMode="numeric"
            value={props.enrollCode}
            onChange={(event) => props.setEnrollCode(event.target.value)}
            placeholder="123456"
          />
          <button type="submit" disabled={props.isBusy || !props.enrollCode.trim()}>
            {props.isBusy ? 'Enrolling...' : 'Complete TOTP setup'}
          </button>
        </form>
      ) : null}

      <p className="auth-status">{props.authStatus}</p>
      {props.authError ? <p className="error">Auth Error: {props.authError}</p> : null}
    </section>
  )
}

export default AuthPanel
