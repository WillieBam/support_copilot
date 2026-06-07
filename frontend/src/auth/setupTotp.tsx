import type React from 'react'
import { useFirebaseTotpAuth } from './useFirebaseTotpAuth'
import { useSetupTotpState } from './useSetupTotpState'
import TotpQrCode from '../components/TotpQrCode'

type SetupTotpProps = {
  auth: ReturnType<typeof useFirebaseTotpAuth>
}

export const SetupTotp: React.FC<SetupTotpProps> = ({ auth }) => {
  const state = useSetupTotpState(auth)

  return (
    <div className="min-h-screen bg-black text-white flex flex-col justify-center items-center font-sans px-4">
      <div className="w-full max-w-[440px] p-8 border border-neutral-800 bg-neutral-950/60 backdrop-blur-xl rounded-[20px] flex flex-col gap-6 shadow-2xl relative overflow-hidden">
        {/* Subtle decorative glow */}
        <div className="absolute -top-20 -left-20 w-48 h-48 bg-emerald-500/10 rounded-full blur-[100px]" />
        <div className="absolute -bottom-20 -right-20 w-48 h-48 bg-orange-500/10 rounded-full blur-[100px]" />

        <div className="flex flex-col gap-1 text-center relative z-10">
          <span className="text-[11px] font-bold tracking-[0.2em] text-emerald-500 uppercase">Support Copilot</span>
          <h1 className="text-2xl font-bold tracking-tight mt-1">Setup TOTP</h1>
        </div>

        {state.needsTotpEnrollment ? (
          <div className="flex flex-col gap-5 relative z-10 items-center w-full">
            <p className="text-sm text-neutral-400 text-center">
              Scan this QR code in your authenticator app (e.g. Google Authenticator or 1Password) to begin:
            </p>

            <div className="p-3 bg-white rounded-lg flex items-center justify-center">
              <TotpQrCode otpauthUrl={state.enrollOtpAuthUrl} />
            </div>

            <div className="w-full flex flex-col gap-2">
              <span className="text-xs font-semibold text-neutral-500 tracking-wide uppercase px-1">Or Setup Key Manually:</span>
              <input
                type="text"
                value={state.enrollOtpAuthUrl}
                readOnly
                className="w-full bg-neutral-900 border border-neutral-800 text-neutral-400 rounded-[20px] px-4 py-2.5 text-xs focus:outline-none"
                onClick={(e) => (e.target as HTMLInputElement).select()}
              />
            </div>

            <form onSubmit={(e) => void state.handleSubmit(e)} className="w-full flex flex-col gap-4">
              <div className="flex flex-col gap-1.5">
                <label htmlFor="totp-code" className="text-xs font-semibold text-neutral-400 tracking-wide uppercase px-1">
                  Enter TOTP
                </label>
                <div className="flex items-center gap-3">
                  <input
                    id="totp-code"
                    type="text"
                    inputMode="numeric"
                    pattern="[0-9]*"
                    maxLength={6}
                    value={state.enrollCode}
                    onChange={(e) => state.handleCodeChange(e.target.value)}
                    placeholder="123456"
                    disabled={state.isBusy}
                    className="flex-1 bg-neutral-900 border border-neutral-800 text-white rounded-[20px] px-4 py-3 text-sm focus:outline-none focus:border-neutral-700 focus:ring-1 focus:ring-neutral-700 transition placeholder-neutral-500 disabled:opacity-50"
                    required
                  />
                  <button
                    type="submit"
                    disabled={state.isBusy || !state.enrollCode.trim()}
                    className="bg-transparent border border-neutral-800 hover:bg-neutral-900 hover:border-neutral-700 text-white font-medium px-6 py-3 rounded-[20px] transition cursor-pointer disabled:opacity-50 text-sm whitespace-nowrap"
                  >
                    {state.isBusy ? 'Verifying...' : 'Submit'}
                  </button>
                </div>
                {state.codeError && (
                  <p className="text-red-500 text-xs mt-1 px-1">{state.codeError}</p>
                )}
              </div>
            </form>
          </div>
        ) : (
          <div className="flex flex-col gap-4 text-center py-4 relative z-10">
            <div className="animate-pulse text-sm text-neutral-400">Initializing TOTP generator...</div>
          </div>
        )}

        <div className="flex flex-col gap-3 relative z-10 w-full">
          {state.submitError && (
            <p className="text-red-500 text-xs text-center">{state.submitError}</p>
          )}
          {state.authStatus && !state.submitError && (
            <p className="text-neutral-400 text-xs text-center">{state.authStatus}</p>
          )}

          <button
            type="button"
            onClick={() => void state.handleCancel()}
            disabled={state.isBusy}
            className="w-full bg-transparent border border-neutral-850 hover:bg-neutral-900 text-neutral-400 hover:text-white font-medium py-3 rounded-[20px] transition duration-200 cursor-pointer disabled:opacity-50 text-sm mt-2"
          >
            Cancel & Sign Out
          </button>
        </div>
      </div>
    </div>
  )
}

export default SetupTotp