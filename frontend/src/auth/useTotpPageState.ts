import { useState } from 'react'
import type { useFirebaseTotpAuth } from './useFirebaseTotpAuth'

export const useTotpPageState = (auth: ReturnType<typeof useFirebaseTotpAuth>) => {
  const [codeError, setCodeError] = useState('')

  const handleCodeChange = (val: string) => {
    auth.setTotpCode(val)
    if (codeError) setCodeError('')
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setCodeError('')

    const code = auth.totpCode.trim()
    if (!code) {
      setCodeError('Verification code is required')
      return
    }
    if (code.length !== 6 || !/^\d+$/.test(code)) {
      setCodeError('Code must be a 6-digit number')
      return
    }

    await auth.verifyTotpSignIn()
  }

  return {
    totpCode: auth.totpCode,
    codeError,
    submitError: auth.authError,
    authStatus: auth.authStatus,
    isBusy: auth.isBusy,
    handleCodeChange,
    handleSubmit,
    handleCancel: auth.signOut,
  }
}
