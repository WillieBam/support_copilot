import { useMemo, useState } from 'react'
import type { FormEvent } from 'react'
import './App.css'
type Role = 'user' | 'assistant'

type ChatMessage = {
  role: Role
  text: string
}

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'

function App() {
  const [token, setToken] = useState('')
  const [prompt, setPrompt] = useState('')
  const [isSending, setIsSending] = useState(false)
  const [error, setError] = useState('')
  const [messages, setMessages] = useState<ChatMessage[]>([
    {
      role: 'assistant',
      text: 'Hello. Ask me anything and I will respond using Gemini through your backend endpoint.',
    },
  ])

  const endpoint = useMemo(() => `${API_BASE_URL}/query/sc`, [])

  const submitPrompt = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const input = prompt.trim()
    if (!input || isSending) return

    setError('')
    setMessages((prev) => [...prev, { role: 'user', text: input }])
    setPrompt('')
    setIsSending(true)

    try {
      const response = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...(token.trim() ? { Authorization: token.trim() } : {}),
        },
        body: JSON.stringify({ input }),
      })

      const payload = (await response.json()) as { output?: string; error?: string }

      if (!response.ok) {
        throw new Error(payload.error || 'Request failed')
      }

      setMessages((prev) => [...prev, { role: 'assistant', text: payload.output ?? '' }])
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unable to contact backend'
      setError(message)
      setMessages((prev) => [
        ...prev,
        {
          role: 'assistant',
          text: `I could not process that request: ${message}`,
        },
      ])
    } finally {
      setIsSending(false)
    }
  }

  return (
    <div className="page-shell">
      <main className="chat-panel">
        <header className="chat-header">
          <p className="eyebrow">Support Copilot</p>
          <h1>Gemini Chat Console</h1>
          <p className="header-subtitle">Connected to {endpoint}</p>
        </header>

        <section className="auth-box">
          <label htmlFor="token">Authorization Token (optional in UI, required if backend auth is enabled)</label>
          <input
            id="token"
            type="text"
            value={token}
            onChange={(event) => setToken(event.target.value)}
            placeholder="Paste Bearer token"
          />
        </section>

        <section className="messages" aria-live="polite">
          {messages.map((message, index) => (
            <article key={`${message.role}-${index}`} className={`message ${message.role}`}>
              <span className="badge">{message.role === 'user' ? 'You' : 'Gemini'}</span>
              <p>{message.text}</p>
            </article>
          ))}
        </section>

        <form className="composer" onSubmit={submitPrompt}>
          <textarea
            value={prompt}
            onChange={(event) => setPrompt(event.target.value)}
            placeholder="Ask something..."
            rows={4}
          />
          <button type="submit" disabled={isSending || !prompt.trim()}>
            {isSending ? 'Sending...' : 'Send'}
          </button>
        </form>

        {error ? <p className="error">Error: {error}</p> : null}
      </main>
      <div className="bg-grid" aria-hidden="true">
        <div className="orb orb-a" />
        <div className="orb orb-b" />
      </div>
    </div>
  )
}

export default App
