import { useFirebaseTotpAuth } from './service/auth/useFirebaseTotpAuth'
import { Thread } from './components/assistant-ui/thread'
import { AssistantRuntimeProvider } from '@assistant-ui/react'
import { useBackendRuntime } from './service/chat/backendRuntime'
import { Navigate, Route, Routes } from 'react-router-dom'
import { LoginPage } from './pages/loginPage'
import { RegisterPage } from './pages/registerPage'
import { SetupTotp } from './pages/setupTotp'
import { TotpPage } from './pages/totpPage'
import { useAppRouter } from './hooks/useAppRouter'
import { useWorkspaceState } from './hooks/useWorkspaceState'
import { Brain, FileText, Users, LogOut, PanelLeftClose, PanelLeftOpen, Sun, Moon } from 'lucide-react'
import { useTheme } from './hooks/useTheme'

type AuthState = ReturnType<typeof useFirebaseTotpAuth>

function LoadingScreen() {
  return (
    <div className="bg-transparent text-foreground min-h-screen flex items-center justify-center transition-colors duration-350">
      <div className="p-8 border border-border bg-card/60 backdrop-blur-xl rounded-[20px] w-full max-w-[440px]">
        <div className="text-center">
          <p className="text-emerald-500 uppercase tracking-widest text-[11px] font-bold">Support Copilot</p>
          <h1 className="text-2xl font-bold tracking-tight mt-1 text-foreground">Loading session</h1>
          <p className="text-muted-foreground text-sm mt-2">Restoring your Firebase login state.</p>
        </div>
      </div>
    </div>
  )
}

function GlobalHeader({ auth }: { auth: AuthState }) {
  const { theme, toggleTheme } = useTheme()

  return (
    <header className="flex w-full h-[73px] items-center justify-between px-6 bg-card border-b border-border shrink-0 transition-colors duration-350 z-20">
      <div className="flex items-center gap-3">
        <Brain className="w-6 h-6 text-emerald-500" />
        <span className="text-foreground font-bold text-lg tracking-tight">Support Copilot</span>
      </div>
      <div className="flex items-center gap-4">
        <button
          onClick={toggleTheme}
          className="flex items-center justify-center w-9 h-9 bg-transparent border border-border rounded-[20px] text-muted-foreground hover:text-foreground hover:bg-muted transition-colors cursor-pointer"
          title={`Switch to ${theme === 'light' ? 'dark' : 'light'} mode`}
        >
          {theme === 'light' ? <Moon className="w-4.5 h-4.5" /> : <Sun className="w-4.5 h-4.5" />}
        </button>
        {auth.isSignedIn && (
          <>
            <button className="flex items-center gap-2 bg-transparent border border-border rounded-[20px] px-4 py-1.5 text-foreground hover:bg-muted transition-colors text-sm cursor-pointer">
              <FileText className="w-4 h-4 text-muted-foreground" /> Manage Instruction
            </button>
            <button className="flex items-center gap-2 bg-transparent border border-border rounded-[20px] px-4 py-1.5 text-foreground hover:bg-muted transition-colors text-sm cursor-pointer">
              <Users className="w-4 h-4 text-muted-foreground" /> Manage Member
            </button>
            <button
              onClick={() => void auth.signOut()}
              disabled={auth.isBusy}
              className="flex items-center gap-2 bg-transparent border border-border rounded-[20px] px-4 py-1.5 text-red-500 hover:bg-muted transition-colors text-sm ml-2 disabled:opacity-50 cursor-pointer"
            >
              <LogOut className="w-4 h-4" /> Logout
            </button>
          </>
        )}
      </div>
    </header>
  )
}

function ChatWorkspace({ auth, runtime }: { auth: AuthState; runtime: ReturnType<typeof useBackendRuntime>['runtime'] }) {
  const { isSidebarOpen, toggleSidebar } = useWorkspaceState()
  
  const email = auth.userEmail;
  const initial = email ? email.charAt(0).toUpperCase() : 'U';
  const displayName = email ? email.split('@')[0] : '';

  return (
    <div className="flex w-full flex-1 overflow-hidden bg-transparent text-foreground relative">
      {/* Left Panel Drawer */}
      <aside
        className={`flex flex-col border-r border-border bg-card/40 backdrop-blur-md transition-all duration-300 ease-in-out ${
          isSidebarOpen ? 'w-[300px]' : 'w-0 border-r-0 overflow-hidden opacity-0'
        }`}
      >
        <div className="p-6 border-b border-border flex items-center gap-4 shrink-0 min-w-[300px]">
          <div className="w-10 h-10 rounded-full bg-muted border border-border flex items-center justify-center text-foreground font-bold text-lg shadow-inner">
            {initial}
          </div>
          <div className="flex flex-col overflow-hidden">
            <span className="text-foreground font-medium truncate">{displayName}</span>
            <span className="text-muted-foreground text-xs truncate">{email}</span>
          </div>
        </div>
        <div className="flex-1 p-6 text-muted-foreground text-sm italic min-w-[300px]">
          WIP: Chat lists WIP.
        </div>
      </aside>

      {/* Right Panel Main Workspace */}
      <main className="flex-1 relative overflow-hidden flex flex-col p-6 items-center">
        <button
          onClick={toggleSidebar}
          className="absolute top-6 left-6 z-20 flex items-center justify-center w-10 h-10 bg-card/60 border border-border rounded-[20px] text-muted-foreground hover:text-foreground hover:bg-muted transition-colors cursor-pointer"
        >
          {isSidebarOpen ? <PanelLeftClose className="w-5 h-5" /> : <PanelLeftOpen className="w-5 h-5" />}
        </button>

        <div className="w-full max-w-[800px] h-full mx-auto border border-border bg-card/60 backdrop-blur-xl rounded-[20px] flex flex-col shadow-2xl relative overflow-hidden transition-colors duration-350">
          <div className="absolute -top-40 -left-40 w-96 h-96 bg-emerald-500/5 rounded-[20px] blur-[120px] pointer-events-none" />
          <div className="absolute -bottom-40 -right-40 w-96 h-96 bg-orange-500/5 rounded-[20px] blur-[120px] pointer-events-none" />

          <div className="flex-1 flex bg-transparent flex-col pt-14 relative z-10 w-full overflow-hidden">
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
    <div className="flex flex-col min-h-screen bg-transparent text-foreground w-full overflow-hidden transition-colors duration-350">
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
