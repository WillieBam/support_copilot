import AuthPanel from './components/AuthPanel'
import { useFirebaseTotpAuth } from './auth/useFirebaseTotpAuth'
import { Thread } from './components/assistant-ui/thread'
import { AssistantRuntimeProvider } from '@assistant-ui/react'
import { useBackendRuntime } from './chat/backendRuntime'
import { Navigate, Route, Routes, useNavigate } from 'react-router-dom'
import { useEffect } from 'react'

type AuthState = ReturnType<typeof useFirebaseTotpAuth>

function LoadingScreen() {
  return (
    <div className="login-shell">
      <div className="login-card">
        <div className="eyebrow-shell">
          <p className="eyebrow">Support Copilot</p>
          <h1>Loading session</h1>
          <p className="header-subtitle">Restoring your Firebase login state.</p>
        </div>
      </div>
      <div className="bg-grid" aria-hidden="true">
        <div className="orb orb-a" />
        <div className="orb orb-b" />
      </div>
    </div>
  )
}

function LoginPage({ auth }: { auth: AuthState }) {
  const navigate = useNavigate()

  useEffect(() => {
    if (auth.isAuthReady && auth.isSignedIn) {
      if (auth.hasTotpEnabled) {
        navigate('/chat', { replace: true })
      }
    }
  }, [auth.isAuthReady, auth.isSignedIn, auth.hasTotpEnabled, navigate])

  if (!auth.isAuthReady) return <LoadingScreen />

  return (
    <div className="login-shell">
      <div className="login-card">
        <div className="eyebrow-shell">
          <p className="eyebrow">Support Copilot</p>
          <h1>Sign in to continue</h1>
          <p className="header-subtitle">
            Use Firebase email/password and TOTP to unlock the assistant workspace.
          </p>
        </div>

        <AuthPanel {...auth} />
      </div>

      <div className="bg-grid" aria-hidden="true">
        <div className="orb orb-a" />
        <div className="orb orb-b" />
      </div>
    </div>
  )
}

function ChatPage({ auth, runtime, endpoint }: { auth: AuthState; runtime: ReturnType<typeof useBackendRuntime>['runtime']; endpoint: string }) {
  const navigate = useNavigate()

  useEffect(() => {
    if (auth.isAuthReady) {
      if (!auth.isSignedIn) {
        navigate('/login', { replace: true })
      } else if (!auth.isEmailVerified) {
        navigate('/login', { replace: true })
      } else if (!auth.hasTotpEnabled) {
        navigate('/login', { replace: true })
      }
    }
  }, [auth.isAuthReady, auth.isSignedIn, auth.isEmailVerified, auth.hasTotpEnabled, navigate])

  if (!auth.isAuthReady) return <LoadingScreen />

  return (
    <div className="app-shell">
      <main className="chat-panel">
        <header className="chat-header">
          <div className="chat-header-copy">
            <p className="eyebrow">assistant-ui</p>
            <h1>Support Copilot</h1>
            <p className="header-subtitle">
              Chat is powered by assistant-ui and your Go backend at {endpoint}.
            </p>
          </div>

          <div className="chat-status-block">
            <div className="chat-status-pill">
              <span>Signed in</span>
              <span>TOTP enabled</span>
            </div>
            <button type="button" className="sign-out-link" onClick={() => void auth.signOut()} disabled={auth.isBusy}>
              Sign out
            </button>
          </div>
        </header>

        <p className="chat-intro">
          Ask a question and the assistant-ui runtime will call the backend with your Firebase bearer token.
        </p>

        <div className="chat-thread-shell">
          <AssistantRuntimeProvider runtime={runtime}>
            <Thread />
          </AssistantRuntimeProvider>
        </div>
      </main>

      <div className="bg-grid" aria-hidden="true">
        <div className="orb orb-a" />
        <div className="orb orb-b" />
      </div>
    </div>
  )
}

function App() {
  const auth = useFirebaseTotpAuth()
  const { runtime, endpoint } = useBackendRuntime(auth.token)

  return (
    <Routes>
      <Route path="/" element={<Navigate to={auth.isSignedIn ? '/chat' : '/login'} replace />} />
      <Route path="/login" element={<LoginPage auth={auth} />} />
      <Route path="/chat" element={<ChatPage auth={auth} runtime={runtime} endpoint={endpoint} />} />
      <Route path="*" element={<Navigate to={auth.isSignedIn ? '/chat' : '/login'} replace />} />
    </Routes>
  )
}

export default App
