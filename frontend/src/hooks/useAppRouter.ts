import { useEffect } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { useFirebaseTotpAuth } from '../service/auth/useFirebaseTotpAuth'

type AuthState = ReturnType<typeof useFirebaseTotpAuth>

export function useAppRouter(auth: AuthState) {
  const navigate = useNavigate()
  const location = useLocation()

  useEffect(() => {
    if (!auth.isAuthReady) return

    console.log("isSignedin", auth.isSignedIn)
    if(auth.needsTotpSignIn){
      if(location.pathname !== '/totp'){
        navigate('/totp');
        }
        return;
    }

    if(!auth.isSignedIn) {
      if(location.pathname !== '/login' && location.pathname !== '/register') {
        navigate('/login');
      }
      return;
    }

    console.log("hasTotpEnabled", auth.hasTotpEnabled)
    if(!auth.hasTotpEnabled) {
      if(location.pathname !== '/setup-totp'){
        navigate('/setup-totp');
      }
      return;
    }

    if(auth.isSignedIn && auth.hasTotpEnabled){
      if(location.pathname === '/login' || location.pathname === '/register' ||
        location.pathname == '/setup-totp' || location.pathname == '/totp'
      ){
        navigate('/chat');
      }
    }
  }, 
    [
    auth.isAuthReady,
    auth.isSignedIn,
    auth.isEmailVerified,
    auth.hasTotpEnabled,
    auth.needsTotpSignIn,
    location.pathname,
    navigate,
  ]
  )
}
