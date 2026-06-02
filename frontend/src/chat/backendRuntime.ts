import { useEffect, useRef } from 'react'
import { useLocalRuntime, type ChatModelAdapter, type ThreadMessage } from '@assistant-ui/react'

const DEFAULT_API_BASE_URL = 'http://localhost:8080'
const SYSTEM_INSTRUCTION =
  'You are Support Copilot. Be concise, accurate, and helpful. If the user asks about authentication, mention the Firebase TOTP flow when relevant.'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || DEFAULT_API_BASE_URL
const ENDPOINT = new URL('/query/chat', API_BASE_URL).toString()

function extractText(message: ThreadMessage): string {
  return message.content
    .map((part) => {
      if (part.type === 'text') return part.text
      return ''
    })
    .filter((part) => part.trim().length > 0)
    .join('\n')
}

function buildPrompt(messages: readonly ThreadMessage[]): string {
  const lines = [SYSTEM_INSTRUCTION, 'Conversation:']

  for (const message of messages) {
    const roleLabel = message.role === 'assistant' ? 'Assistant' : 'User'
    const content = extractText(message).trim()
    if (!content) continue
    lines.push(`${roleLabel}: ${content}`)
  }

  lines.push('Assistant:')
  return lines.join('\n\n')
}

export function useBackendRuntime(authToken: string) {
  const authTokenRef = useRef(authToken)
  useEffect(() => {
    authTokenRef.current = authToken
  }, [authToken])

  const chatModelRef = useRef<ChatModelAdapter>({
    async run({messages, abortSignal}) {
      const response = await fetch(ENDPOINT, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...(authTokenRef.current.trim() ? {Authorization: authTokenRef.current.trim()} : {})
        },
        body: JSON.stringify({input: buildPrompt(messages)}),
        signal: abortSignal,
      })

      let payload: {output?: string; error?: string} = {}
      try{
        payload = (await response.json() as {output?: string; error?: string})
      }catch {
        payload = {}
      }
      if (!response.ok) {
        throw new Error(payload.error || `Backend request failed (${response.status})`)
      }

      return {
        content: [{ type: 'text', text: payload.output ?? '' }],
        status: { type: 'complete', reason: 'stop' },
        metadata: {
          custom: {
            backendUrl: ENDPOINT,
          },
        },
      }
    }
  })


  const runtime = useLocalRuntime(chatModelRef.current)

  return { runtime, endpoint:ENDPOINT }
}