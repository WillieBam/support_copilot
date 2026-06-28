import type React from 'react'
import { Link } from 'react-router-dom'
import { useFirebaseTotpAuth } from '../service/auth/useFirebaseTotpAuth'
import { useLoginState } from '../service/auth/useLoginState'

type LoginPageProps = {
  auth: ReturnType<typeof useFirebaseTotpAuth>
}

export const LoginPage: React.FC<LoginPageProps> = ({ auth }) => {
  const state = useLoginState(auth)

  if (state.isSignedIn && !state.isEmailVerified) {
    return (
      <div className="min-h-screen bg-black text-white flex flex-col justify-center items-center font-sans px-4">
        <div className="w-full max-w-[440px] p-8 border border-neutral-800 bg-neutral-950/60 backdrop-blur-xl rounded-[20px] flex flex-col gap-6 shadow-2xl relative overflow-hidden">
          {/* Subtle decorative glow */}
          <div className="absolute -top-20 -left-20 w-48 h-48 bg-emerald-500/10 rounded-full blur-[100px]" />
          <div className="absolute -bottom-20 -right-20 w-48 h-48 bg-orange-500/10 rounded-full blur-[100px]" />

          <div className="flex flex-col gap-2 text-center relative z-10">
            <span className="text-[11px] font-bold tracking-[0.2em] text-emerald-500 uppercase">Support Copilot</span>
            <h1 className="text-2xl font-bold tracking-tight mt-1">Verify Email</h1>
            <p className="text-sm text-neutral-400 mt-2">
              We sent a verification email to <strong className="text-white">{state.email}</strong>.
            </p>
            <p className="text-xs text-neutral-500 mt-1">
              Please check your inbox and spam folder, then verify your email before continuing.
            </p>
          </div>

          <div className="flex flex-col gap-3 relative z-10">
            <button
              type="button"
              onClick={() => void state.checkVerificationStatus()}
              disabled={state.isBusy}
              className="w-full bg-emerald-600 hover:bg-emerald-500 text-white font-semibold py-3 rounded-[20px] transition duration-200 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed text-center text-sm"
            >
              {state.isBusy ? 'Checking...' : 'Check Status'}
            </button>

            <button
              type="button"
              onClick={() => void state.resendVerification()}
              disabled={state.isBusy}
              className="w-full bg-transparent border border-neutral-850 hover:bg-neutral-900 text-neutral-300 font-medium py-3 rounded-[20px] transition duration-200 cursor-pointer disabled:opacity-50 text-sm"
            >
              Resend Verification Email
            </button>

            <button
              type="button"
              onClick={() => void state.signOut()}
              disabled={state.isBusy}
              className="w-full bg-transparent border border-neutral-850 hover:bg-neutral-900 text-red-400 font-medium py-3 rounded-[20px] transition duration-200 cursor-pointer disabled:opacity-50 text-sm"
            >
              Sign Out
            </button>
          </div>

          {state.submitError && (
            <p className="text-red-500 text-xs text-center z-10">{state.submitError}</p>
          )}
          {state.authStatus && (
            <p className="text-neutral-400 text-xs text-center z-10">{state.authStatus}</p>
          )}
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-black text-white flex flex-col justify-center items-center font-sans px-4">
      <div className="w-full max-w-[440px] p-8 border border-neutral-800 bg-neutral-950/60 backdrop-blur-xl rounded-[20px] flex flex-col gap-6 shadow-2xl relative overflow-hidden">
        
        <div className="flex flex-col gap-1 text-center relative z-10">
          <h1 className="text-2xl tracking-tight mt-1">Welcome to Support Copilot</h1>
        </div>

        <form onSubmit={(e) => void state.handleSubmit(e)} className="flex flex-col gap-4 relative z-10">
          <div className="flex flex-col gap-1.5">
            <label className="text-xs font-semibold text-neutral-400 tracking-wide px-1">Email</label>
            <input
              type="email"
              value={state.email}
              onChange={(e) => state.handleEmailChange(e.target.value)}
              placeholder="name@example.com"
              disabled={state.isBusy}
              className="w-full bg-neutral-900 border border-neutral-800 text-white rounded-[20px] px-4 py-3 text-sm focus:outline-none focus:border-neutral-700 focus:ring-1 focus:ring-neutral-700 transition placeholder-neutral-500 disabled:opacity-50"
              required
            />
            {state.emailError && (
              <p className="text-red-500 text-xs mt-1 px-1">{state.emailError}</p>
            )}
          </div>

          <div className="flex flex-col gap-1.5">
            <label className="text-xs font-semibold text-neutral-400 tracking-wide px-1">Password</label>
            <input
              type="password"
              value={state.password}
              onChange={(e) => state.handlePasswordChange(e.target.value)}
              placeholder="••••••••"
              disabled={state.isBusy}
              className="w-full bg-neutral-900 border border-neutral-800 text-white rounded-[20px] px-4 py-3 text-sm focus:outline-none focus:border-neutral-700 focus:ring-1 focus:ring-neutral-700 transition placeholder-neutral-500 disabled:opacity-50"
              required
            />
            {state.passwordError && (
              <p className="text-red-500 text-xs mt-1 px-1">{state.passwordError}</p>
            )}
          </div>

          <div className="flex justify-end mt-2">
            <button
              type="submit"
              disabled={state.isBusy}
              className="bg-transparent border border-neutral-800 hover:bg-neutral-900 hover:border-neutral-700 text-white font-medium px-6 py-2.5 rounded-[20px] transition cursor-pointer disabled:opacity-50 text-sm"
            >
              {state.isBusy ? 'Signing In...' : 'Sign In'}
            </button>
          </div>
        </form>

        {state.submitError && (
          <p className="text-red-500 text-xs text-center z-10">{state.submitError}</p>
        )}
        {state.authStatus && !state.submitError && (
          <p className="text-neutral-400 text-xs text-center z-10">{state.authStatus}</p>
        )}

        <div className="text-center text-sm text-neutral-400 border-t border-neutral-900 pt-4 mt-2 relative z-10">
          First Time Here?{' '}
          <Link to="/register" className="text-white hover:underline font-medium">
            Create an Account
          </Link>
        </div>
      </div>
    </div>
  )
}

export default LoginPage