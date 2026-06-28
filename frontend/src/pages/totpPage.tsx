import type React from 'react'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useFirebaseTotpAuth } from '../service/auth/useFirebaseTotpAuth'
import { useTotpPageState } from '../service/auth/useTotpPageState'

type TotpPageProps = {
  auth: ReturnType<typeof useFirebaseTotpAuth>
}

export const TotpPage: React.FC<TotpPageProps> = ({ auth }) => {
  const state = useTotpPageState(auth)
  const [success, setSuccess] = useState(false)
  const navigate = useNavigate()

  console.log("hello")

  const handleFormSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    await state.handleSubmit(e)

    // If verification succeeded, auth.authError is cleared and token is present
    if (!auth.authError && auth.isSignedIn) {
      setSuccess(true)
      setTimeout(() => {
        navigate('/chat', { replace: true })
      }, 1200)
    }
  }

  return (
    <div className="min-h-screen bg-black text-white flex flex-col justify-center items-center font-sans px-4">
      <div className="w-full max-w-[440px] p-8 border border-neutral-800 bg-neutral-950/60 backdrop-blur-xl rounded-[20px] flex flex-col gap-6 shadow-2xl relative overflow-hidden">
        {/* Subtle decorative glow */}
        <div className="absolute -top-20 -left-20 w-48 h-48 bg-emerald-500/10 rounded-full blur-[100px]" />
        <div className="absolute -bottom-20 -right-20 w-48 h-48 bg-orange-500/10 rounded-full blur-[100px]" />

        <div className="flex flex-col gap-1 text-center relative z-10">
          <span className="text-[11px] font-bold tracking-[0.2em] text-emerald-500 uppercase">Support Copilot</span>
          <h1 className="text-2xl font-bold tracking-tight mt-1">Enter Totp</h1>
        </div>

        {success ? (
          <div className="bg-emerald-950/40 border border-emerald-500/30 text-emerald-400 p-4 rounded-[20px] text-center text-sm font-semibold animate-fade-in relative z-10">
            ✓ login successfully
          </div>
        ) : (
          <form onSubmit={handleFormSubmit} className="flex flex-col gap-4 relative z-10">
            <div className="flex flex-col gap-1.5">
              <label htmlFor="totp-input" className="text-xs font-semibold text-neutral-400 tracking-wide uppercase px-1">
                Enter your 6-digit TOTP code
              </label>
              <div className="flex items-center gap-3">
                <input
                  id="totp-input"
                  type="text"
                  inputMode="numeric"
                  pattern="[0-9]*"
                  maxLength={6}
                  value={state.totpCode}
                  onChange={(e) => state.handleCodeChange(e.target.value)}
                  placeholder="123456"
                  disabled={state.isBusy}
                  className="flex-1 bg-neutral-900 border border-neutral-800 text-white rounded-[20px] px-4 py-3 text-sm focus:outline-none focus:border-neutral-700 focus:ring-1 focus:ring-neutral-700 transition placeholder-neutral-500 disabled:opacity-50"
                  required
                />
                <button
                  type="submit"
                  disabled={state.isBusy || !state.totpCode.trim()}
                  className="bg-transparent border border-neutral-800 hover:bg-neutral-900 hover:border-neutral-700 text-white font-medium px-6 py-3 rounded-[20px] transition cursor-pointer disabled:opacity-50 text-sm whitespace-nowrap"
                >
                  {state.isBusy ? 'Verifying...' : 'Verify'}
                </button>
              </div>
              {state.codeError && (
                <p className="text-red-500 text-xs mt-1 px-1">{state.codeError}</p>
              )}
              {state.submitError && (
                <p className="text-red-500 text-xs mt-1 px-1">{state.submitError}</p>
              )}
            </div>
          </form>
        )}

        <div className="flex flex-col gap-3 relative z-10 w-full mt-2">
          {state.authStatus && !state.submitError && !success && (
            <p className="text-neutral-400 text-xs text-center">{state.authStatus}</p>
          )}

          <button
            type="button"
            onClick={() => void state.handleCancel()}
            disabled={state.isBusy}
            className="w-full bg-transparent border border-neutral-850 hover:bg-neutral-900 text-neutral-400 hover:text-white font-medium py-3 rounded-[20px] transition duration-200 cursor-pointer disabled:opacity-50 text-sm"
          >
            Cancel & Sign Out
          </button>
        </div>
      </div>
    </div>
  )
}

export default TotpPage