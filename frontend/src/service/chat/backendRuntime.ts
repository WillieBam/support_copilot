import { useRef } from 'react'
import { useLocalRuntime, type ChatModelAdapter, type ThreadMessage } from '@assistant-ui/react'
import apiClient from '../apiClient'

// const DEFAULT_API_BASE_URL = 'http://localhost:8080'
const SYSTEM_INSTRUCTION =
  'You are Support Copilot. Be concise, accurate, and helpful. If the user asks about authentication, mention the Firebase TOTP flow when relevant.'

// const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || DEFAULT_API_BASE_URL
// const ENDPOINT = new URL('/query/chat', API_BASE_URL).toString()

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

export function useBackendRuntime() {
 const chatModelRef = useRef<ChatModelAdapter>({
  
    async run({messages, abortSignal}) {
      try{
        const response = await apiClient.post('/query/chat',
          {input: buildPrompt(messages)},
          {signal: abortSignal}
        );
        return {
          content: [{type: 'text', text: response.data.output ?? ''}],
          status: {type: 'complete', reason:'stop'},
        }
      } catch (error: any) {
        const errMsg = error.response?.data?.error || `Backend request failed`;
        throw new Error (errMsg);
      }
    }

  })
  return { runtime: useLocalRuntime(chatModelRef.current) }
}