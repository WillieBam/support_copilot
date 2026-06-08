import { useEffect } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { useFirebaseTotpAuth } from '../service/auth/useFirebaseTotpAuth'

type AuthState = ReturnType<typeof useFirebaseTotpAuth>

export function useAppRouter(auth: AuthState) {
  const navigate = useNavigate()
  const location = useLocation()

  useEffect(() => {
    if (!auth.isAuthReady) return

    const path = location.pathname

    if (!auth.isSignedIn) {
      if (auth.needsTotpSignIn) {
        if (path !== '/totp') {
          navigate('/totp', { replace: true })
        }
      } else {
        if (path !== '/login' && path !== '/register') {
          navigate('/login', { replace: true })
        }
      }
    } else {
      if (!auth.isEmailVerified) {
        if (path !== '/login') {
          navigate('/login', { replace: true })
        }
      } else if (!auth.hasTotpEnabled) {
        if (path !== '/setup-totp') {
          navigate('/login', { replace: true })
        }
      } else {
        if (
          path === '/login' ||
          path === '/register' ||
          path === '/setup-totp' ||
          path === '/totp' ||
          path === '/'
        ) {
          navigate('/chat', { replace: true })
        }
      }
    }
  }, [
    auth.isAuthReady,
    auth.isSignedIn,
    auth.isEmailVerified,
    auth.hasTotpEnabled,
    auth.needsTotpSignIn,
    location.pathname,
    navigate,
  ])
}
