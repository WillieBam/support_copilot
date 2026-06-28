import { useFirebaseTotpAuth } from './service/auth/useFirebaseTotpAuth'
import { Thread } from './components/assistant-ui/thread'
import { AssistantRuntimeProvider } from '@assistant-ui/react'
import { useBackendRuntime } from './service/chat/backendRuntime'
import { Navigate, Route, Routes, useNavigate } from 'react-router-dom'
import { LoginPage } from './pages/loginPage'
import { RegisterPage } from './pages/registerPage'
import { SetupTotp } from './pages/setupTotp'
import { TotpPage } from './pages/totpPage'
import { useAppRouter } from './hooks/useAppRouter'
import { useWorkspaceState } from './hooks/useWorkspaceState'
import { Brain, FileText, Users, LogOut, PanelLeftClose, PanelLeftOpen } from 'lucide-react'

type AuthState = ReturnType<typeof useFirebaseTotpAuth>

function LoadingScreen() {
  return (
    <div className="bg-black text-white min-h-screen flex items-center justify-center">
      <div className="p-8 border border-neutral-800 bg-neutral-950/60 backdrop-blur-xl rounded-[20px] w-full max-w-[440px]">
        <div className="text-center">
          <p className="text-emerald-500 uppercase tracking-widest text-[11px] font-bold">Support Copilot</p>
          <h1 className="text-2xl font-bold tracking-tight mt-1 text-white">Loading session</h1>
          <p className="text-neutral-400 text-sm mt-2">Restoring your Firebase login state.</p>
        </div>
      </div>
    </div>
  )
}

function GlobalHeader({ auth }: { auth: AuthState }) {
  return (
    <header className="flex w-full h-[73px] items-center justify-between px-6 bg-black border-b border-neutral-800 shrink-0">
      <div className="flex items-center gap-3">
        <Brain className="w-6 h-6 text-emerald-500" />
        <span className="text-white font-bold text-lg tracking-tight">Support Copilot</span>
      </div>
      <div className="flex items-center gap-4">
        {auth.isSignedIn ? (
          <>
            <div className="relative">
              <select className="appearance-none bg-neutral-900 border border-neutral-800 rounded-[20px] text-white px-4 py-1.5 pr-8 outline-none text-sm cursor-pointer hover:bg-neutral-800 transition-colors">
                <option>Personal Workspace</option>
              </select>
              <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-neutral-400">
                <svg className="fill-current h-4 w-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"><path d="M9.293 12.95l.707.707L15.657 8l-1.414-1.414L10 10.828 5.757 6.586 4.343 8z"/></svg>
              </div>
            </div>
            <button className="flex items-center gap-2 bg-transparent border border-neutral-800 rounded-[20px] px-4 py-1.5 text-white hover:bg-neutral-900 transition-colors text-sm">
              <FileText className="w-4 h-4 text-neutral-400" /> Manage Instruction
            </button>
            <button className="flex items-center gap-2 bg-transparent border border-neutral-800 rounded-[20px] px-4 py-1.5 text-white hover:bg-neutral-900 transition-colors text-sm">
              <Users className="w-4 h-4 text-neutral-400" /> Manage Member
            </button>
            <button
              onClick={() => void auth.signOut()}
              disabled={auth.isBusy}
              className="flex items-center gap-2 bg-transparent border border-neutral-800 rounded-[20px] px-4 py-1.5 text-red-400 hover:bg-neutral-900 transition-colors text-sm ml-2 disabled:opacity-50"
            >
              <LogOut className="w-4 h-4" /> Logout
            </button>
          </>
        ) : (
          <button>
          </button>
        )}
      </div>
    </header>
  )
}

function ChatWorkspace({ auth, runtime }: { auth: AuthState; runtime: ReturnType<typeof useBackendRuntime>['runtime'] }) {
  const { isSidebarOpen, toggleSidebar } = useWorkspaceState()
  
  // Extract user info
  const initial = auth.email?.charAt(0).toUpperCase() || 'U'
  const email = auth.email || 'user@example.com'
  const displayName = email.split('@')[0] // Fallback for name

  return (
    <div className="flex w-full flex-1 overflow-hidden bg-black text-white relative">
      {/* Left Panel Drawer */}
      <aside
        className={`flex flex-col border-r border-neutral-800 bg-black transition-all duration-300 ease-in-out ${
          isSidebarOpen ? 'w-[300px]' : 'w-[0px] border-r-0 overflow-hidden opacity-0'
        }`}
      >
        <div className="p-6 border-b border-neutral-800 flex items-center gap-4 shrink-0 min-w-[300px]">
          <div className="w-10 h-10 rounded-full bg-neutral-800/80 border border-neutral-700 flex items-center justify-center text-white font-bold text-lg shadow-inner">
            {initial}
          </div>
          <div className="flex flex-col overflow-hidden">
            <span className="text-white font-medium truncate">{displayName}</span>
            <span className="text-neutral-400 text-xs truncate">{email}</span>
          </div>
        </div>
        <div className="flex-1 p-6 text-neutral-500 text-sm italic min-w-[300px]">
          WIP:Chat lists WIP.
        </div>
      </aside>

      {/* Right Panel Main Workspace */}
      <main className="flex-1 relative overflow-hidden flex flex-col p-6 items-center">
        <button
          onClick={toggleSidebar}
          className="absolute top-6 left-6 z-20 flex items-center justify-center w-10 h-10 bg-neutral-900/60 border border-neutral-800 rounded-[20px] text-neutral-400 hover:text-white hover:bg-neutral-800 transition-colors"
        >
          {isSidebarOpen ? <PanelLeftClose className="w-5 h-5" /> : <PanelLeftOpen className="w-5 h-5" />}
        </button>

        <div className="w-full max-w-[800px] h-full mx-auto border border-neutral-800 bg-neutral-950/60 backdrop-blur-xl rounded-[20px] flex flex-col shadow-2xl relative overflow-hidden">
          {/* Subtle decorative glow */}
          <div className="absolute -top-40 -left-40 w-96 h-96 bg-emerald-500/5 rounded-[20px] blur-[120px] pointer-events-none" />
          <div className="absolute -bottom-40 -right-40 w-96 h-96 bg-orange-500/5 rounded-[20px] blur-[120px] pointer-events-none" />

          <div className="flex-1 flex flex-col pt-14 relative z-10 w-full overflow-hidden">
            <AssistantRuntimeProvider runtime={runtime}>
              <Thread />
            </AssistantRuntimeProvider>
          </div>
        </div>
      </main>
    </div>
  )
}

function App() {
  const auth = useFirebaseTotpAuth()
  const { runtime } = useBackendRuntime()
  
  // App routing logic has been decoupled into this hook
  useAppRouter(auth)

  if (!auth.isAuthReady) return <LoadingScreen />

  return (
    <div className="flex flex-col min-h-screen bg-black w-full overflow-hidden">
      <GlobalHeader auth={auth} />
      <Routes>
        <Route path="/" element={<Navigate to={auth.isSignedIn ? '/chat' : '/login'} replace />} />
        
        {/* Auth routes centered over black background */}
        <Route path="/login" element={<div className="flex-1 flex items-center justify-center"><LoginPage auth={auth} /></div>} />
        <Route path="/register" element={<div className="flex-1 flex items-center justify-center"><RegisterPage auth={auth} /></div>} />
        <Route path="/setup-totp" element={<div className="flex-1 flex items-center justify-center"><SetupTotp auth={auth} /></div>} />
        <Route path="/totp" element={<div className="flex-1 flex items-center justify-center"><TotpPage auth={auth} /></div>} />
        
        {/* Chat Workspace */}
        <Route path="/chat" element={<ChatWorkspace auth={auth} runtime={runtime} />} />
        
        <Route path="*" element={<Navigate to={auth.isSignedIn ? '/chat' : '/login'} replace />} />
      </Routes>
    </div>
  )
}

export default App
