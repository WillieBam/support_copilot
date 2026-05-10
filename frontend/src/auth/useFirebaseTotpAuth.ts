import { useEffect, useMemo, useState } from 'react'
import { onIdTokenChanged, type MultiFactorInfo, type MultiFactorResolver, type TotpSecret } from 'firebase/auth'
import { firebaseAuth } from '../firebase'
import {
  beginTotpEnrollment,
  confirmTotpEnrollment,
  createAccount,
  hasTotpEnrollment,
  resolveTotpSignIn,
  sendVerificationEmail,
  signInWithPassword,
  signOutCurrentUser,
  toErrorMessage,
} from './authService'

export function useFirebaseTotpAuth() {
  const [token, setToken] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [totpCode, setTotpCode] = useState('')
  const [enrollCode, setEnrollCode] = useState('')
  const [isBusy, setIsBusy] = useState(false)
  const [totpResolver, setTotpResolver] = useState<MultiFactorResolver | null>(null)
  const [totpHint, setTotpHint] = useState<MultiFactorInfo | null>(null)
  const [enrollSecret, setEnrollSecret] = useState<TotpSecret | null>(null)
  const [enrollOtpAuthUrl, setEnrollOtpAuthUrl] = useState('')
  const [hasTotpEnabled, setHasTotpEnabled] = useState(false)
  const [isEmailVerified, setIsEmailVerified] = useState(false)
  const [authStatus, setAuthStatus] = useState('Not signed in')
  const [authError, setAuthError] = useState('')

  useEffect(() => {
    const unsubscribe = onIdTokenChanged(firebaseAuth, async (user) => {
      if (!user) {
        setToken('')
        setHasTotpEnabled(false)
        setIsEmailVerified(false)
        setAuthStatus('Not signed in')
        return
      }

      const idToken = await user.getIdToken()
      setToken(`Bearer ${idToken}`)
      setHasTotpEnabled(hasTotpEnrollment(user))
      setIsEmailVerified(user.emailVerified)
      setAuthStatus(`Signed in as ${user.email ?? user.uid}`)
    })

    return () => unsubscribe()
  }, [])

  const isSignedIn = useMemo(() => token !== '', [token])

  const signIn = async () => {
    if (isBusy) return

    setAuthError('')
    setIsBusy(true)
    try {
      const result = await signInWithPassword(firebaseAuth, email.trim(), password)
      if (result.type === 'totp-required') {
        setTotpResolver(result.resolver)
        setTotpHint(result.hint)
        setAuthStatus('TOTP required. Enter your authenticator code to continue.')
        return
      }

      setTotpResolver(null)
      setTotpHint(null)
      setTotpCode('')
      if (!firebaseAuth.currentUser?.emailVerified) {
        setAuthStatus('Signed in, but email is not verified yet. Verify email first, then set up TOTP.')
      }
    } catch (error) {
      setAuthError(toErrorMessage(error, 'Sign-in failed'))
      setAuthStatus('Sign-in failed')
    } finally {
      setIsBusy(false)
    }
  }

  const verifyTotpSignIn = async () => {
    if (isBusy || !totpResolver || !totpHint) return

    const code = totpCode.trim()
    if (!code) {
      setAuthError('TOTP code is required')
      return
    }

    setAuthError('')
    setIsBusy(true)
    try {
      await resolveTotpSignIn(totpResolver, totpHint.uid, code)
      setTotpResolver(null)
      setTotpHint(null)
      setTotpCode('')
      setAuthStatus('TOTP verified. Signed in.')
    } catch (error) {
      setAuthError(toErrorMessage(error, 'TOTP verification failed'))
      setAuthStatus('TOTP verification failed')
    } finally {
      setIsBusy(false)
    }
  }

  const register = async () => {
    if (isBusy) return

    setAuthError('')
    setIsBusy(true)
    try {
      const user = await createAccount(firebaseAuth, email.trim(), password)
      await sendVerificationEmail(user)
      await signOutCurrentUser(firebaseAuth)
      setEnrollSecret(null)
      setEnrollOtpAuthUrl('')
      setEnrollCode('')
      setAuthStatus('Account created. Verification email sent. Verify your email, then sign in and set up TOTP.')
      setTotpResolver(null)
      setTotpHint(null)
      setTotpCode('')
    } catch (error) {
      setAuthError(toErrorMessage(error, 'Account creation failed'))
      setAuthStatus('Account creation failed')
    } finally {
      setIsBusy(false)
    }
  }

  const enrollTotp = async () => {
    if (isBusy || !enrollSecret) return

    const user = firebaseAuth.currentUser
    if (!user) {
      setAuthError('No signed-in user for TOTP enrollment')
      return
    }

    const code = enrollCode.trim()
    if (!code) {
      setAuthError('Enrollment code is required')
      return
    }

    setAuthError('')
    setIsBusy(true)
    try {
      await confirmTotpEnrollment(user, enrollSecret, code)
      setEnrollSecret(null)
      setEnrollOtpAuthUrl('')
      setEnrollCode('')
      setAuthStatus('TOTP enrolled. Sign in again to verify TOTP during login.')
      await signOutCurrentUser(firebaseAuth)
    } catch (error) {
      setAuthError(toErrorMessage(error, 'TOTP enrollment failed'))
      setAuthStatus('TOTP enrollment failed')
    } finally {
      setIsBusy(false)
    }
  }

  const startTotpEnrollment = async () => {
    if (isBusy || enrollSecret) return

    const user = firebaseAuth.currentUser
    if (!user) {
      setAuthError('You must sign in first')
      return
    }
    if (!user.emailVerified) {
      setAuthError('Please verify your email first before enrolling TOTP')
      setAuthStatus('Email verification required before TOTP setup')
      return
    }
    if (hasTotpEnrollment(user)) {
      setAuthStatus('TOTP is already enrolled for this account.')
      return
    }

    setAuthError('')
    setIsBusy(true)
    try {
      const { secret, otpauthUrl } = await beginTotpEnrollment(user)
      setEnrollSecret(secret)
      setEnrollOtpAuthUrl(otpauthUrl)
      setEnrollCode('')
      setAuthStatus('Scan QR and enter code to complete TOTP setup.')
    } catch (error) {
      setAuthError(toErrorMessage(error, 'Unable to start TOTP enrollment'))
      setAuthStatus('Unable to start TOTP setup')
    } finally {
      setIsBusy(false)
    }
  }

  const resendVerification = async () => {
    if (isBusy) return
    const user = firebaseAuth.currentUser
    if (!user) {
      setAuthError('Sign in first to resend verification email')
      return
    }
    if (user.emailVerified) {
      setAuthStatus('Email is already verified.')
      return
    }

    setAuthError('')
    setIsBusy(true)
    try {
      await sendVerificationEmail(user)
      setAuthStatus('Verification email sent. Check your inbox and spam folder.')
    } catch (error) {
      setAuthError(toErrorMessage(error, 'Failed to resend verification email'))
    } finally {
      setIsBusy(false)
    }
  }

  const signOut = async () => {
    if (isBusy) return

    setAuthError('')
    setIsBusy(true)
    try {
      await signOutCurrentUser(firebaseAuth)
      setTotpResolver(null)
      setTotpHint(null)
      setTotpCode('')
      setEnrollSecret(null)
      setEnrollOtpAuthUrl('')
      setEnrollCode('')
      setHasTotpEnabled(false)
      setIsEmailVerified(false)
      setAuthStatus('Signed out')
    } catch (error) {
      setAuthError(toErrorMessage(error, 'Sign-out failed'))
    } finally {
      setIsBusy(false)
    }
  }

  return {
    token,
    email,
    password,
    setEmail,
    setPassword,
    totpCode,
    setTotpCode,
    enrollCode,
    setEnrollCode,
    enrollOtpAuthUrl,
    hasTotpEnabled,
    isEmailVerified,
    needsTotpSignIn: totpResolver !== null && totpHint !== null,
    needsTotpEnrollment: enrollSecret !== null,
    canStartTotpEnrollment: isSignedIn && isEmailVerified && !hasTotpEnabled && enrollSecret === null,
    isBusy,
    isSignedIn,
    authStatus,
    authError,
    signIn,
    verifyTotpSignIn,
    register,
    startTotpEnrollment,
    enrollTotp,
    resendVerification,
    signOut,
  }
}
