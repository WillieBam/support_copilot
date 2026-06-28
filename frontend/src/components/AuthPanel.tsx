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
  checkVerificationStatus: () => Promise<void>
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

  // --- RENDER CONDITION 1: User is already Authenticated with Firebase ---
  if (props.isSignedIn) {
    return (
      <section className="auth-box">
        <p className="auth-title">Firebase Auth Status</p>
        <p className="auth-status">{props.authStatus}</p>

        {!props.isEmailVerified ? (
          <div className="auth-grid auth-grid-totp">
            <p className="auth-help">
              Your email is not verified yet. We sent a verification email to <strong>{props.email}</strong>.
            </p>
            <p className="auth-help">
              Please check your inbox and verify your email, then click "Check status" below.
            </p>
            <div className="auth-actions">
              <button type="button" onClick={() => void props.checkVerificationStatus()} disabled={props.isBusy}>
                {props.isBusy ? 'Checking...' : 'Check status'}
              </button>
              <button type="button" onClick={() => void props.resendVerification()} disabled={props.isBusy}>
                Resend verification email
              </button>
              <button type="button" onClick={() => void props.signOut()} disabled={props.isBusy}>
                Sign out
              </button>
            </div>
          </div>
        ) : (
          <>
            {!props.hasTotpEnabled && !props.needsTotpEnrollment && (
              <div className="auth-grid auth-grid-totp">
                <p className="auth-help">
                  Email verified! You must set up a Time-based One-Time Password (TOTP) to access the workspace.
                </p>
                <div className="auth-actions">
                  <button type="button" onClick={() => void props.startTotpEnrollment()} disabled={props.isBusy}>
                    Set up TOTP
                  </button>
                  <button type="button" onClick={() => void props.signOut()} disabled={props.isBusy}>
                    Sign out
                  </button>
                </div>
              </div>
            )}

            {props.needsTotpEnrollment && (
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
                  required
                />
                <div className="auth-actions">
                  <button type="submit" disabled={props.isBusy || !props.enrollCode.trim()}>
                    {props.isBusy ? 'Enrolling...' : 'Complete TOTP setup'}
                  </button>
                  <button type="button" onClick={() => void props.signOut()} disabled={props.isBusy}>
                    Cancel & Sign out
                  </button>
                </div>
              </form>
            )}

            {props.hasTotpEnabled && (
              <div className="auth-grid auth-grid-totp">
                <p className="auth-help">TOTP is enabled. Redirecting to workspace...</p>
                <div className="auth-actions">
                  <button type="button" onClick={() => void props.signOut()} disabled={props.isBusy}>
                    Sign out
                  </button>
                </div>
              </div>
            )}
          </>
        )}

        {props.authError ? <p className="error">Auth Error: {props.authError}</p> : null}
      </section>
    )
  }

  // --- RENDER CONDITION 2: User is NOT Signed In (Show Login Forms) ---
  return (
    <section className="auth-box">
      <p className="auth-title">Firebase Auth</p>

      {/* Conditional Rendering: If TOTP Multi-factor validation challenge is active, swap out forms */}
      {props.needsTotpSignIn ? (
        <form className="auth-grid auth-grid-totp" onSubmit={handleVerifyTotp}>
          <p className="auth-help">Multi-Factor Authentication Required.</p>
          <label htmlFor="totp">Authenticator code</label>
          <input
            id="totp"
            type="text"
            inputMode="numeric"
            value={props.totpCode}
            onChange={(event) => props.setTotpCode(event.target.value)}
            placeholder="123456"
            required
          />
          <div className="auth-actions">
            <button type="submit" disabled={props.isBusy || !props.totpCode.trim()}>
              {props.isBusy ? 'Verifying...' : 'Verify TOTP'}
            </button>
          </div>
        </form>
      ) : (
        <>
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
            </div>
          </form>
        </>
      )}

      <p className="auth-status">{props.authStatus}</p>
      {props.authError ? <p className="error">Auth Error: {props.authError}</p> : null}
    </section>
  )
}

export default AuthPanel