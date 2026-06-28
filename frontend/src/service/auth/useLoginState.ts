import { useState } from 'react'
import type { useFirebaseTotpAuth } from './useFirebaseTotpAuth'

export const useLoginState = (auth: ReturnType<typeof useFirebaseTotpAuth>) => {
  const [emailError, setEmailError] = useState('')
  const [passwordError, setPasswordError] = useState('')

  const handleEmailChange = (val: string) => {
    auth.setEmail(val)
    if (emailError) setEmailError('')
  }

  const handlePasswordChange = (val: string) => {
    auth.setPassword(val)
    if (passwordError) setPasswordError('')
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setEmailError('')
    setPasswordError('')

    let hasError = false
    if (!auth.email.trim()) {
      setEmailError('Email is required')
      hasError = true
    } else if (!/\S+@\S+\.\S+/.test(auth.email)) {
      setEmailError('Invalid email format')
      hasError = true
    }

    if (!auth.password) {
      setPasswordError('Password is required')
      hasError = true
    } else if (auth.password.length < 6) {
      setPasswordError('Password must be at least 6 characters')
      hasError = true
    }

    if (hasError) return

    await auth.signIn()
  }
  return {
    email: auth.email,
    password: auth.password,
    emailError,
    passwordError,
    submitError: auth.authError,
    authStatus: auth.authStatus,
    isBusy: auth.isBusy,
    isSignedIn: auth.isSignedIn,
    isEmailVerified: auth.isEmailVerified,
    handleEmailChange,
    handlePasswordChange,
    handleSubmit,
    checkVerificationStatus: auth.checkVerificationStatus,
    resendVerification: auth.resendVerification,
    signOut: auth.signOut,
  }
}
