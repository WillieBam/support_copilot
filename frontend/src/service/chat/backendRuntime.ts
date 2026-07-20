import { useRef } from "react";
import {
  useLocalRuntime,
  type ChatModelAdapter,
  type ThreadMessage,
  type ChatModelRunOptions,
  type ChatModelRunResult,
} from "@assistant-ui/react";
// import apiClient from "../apiClient";
import { firebaseAuth } from "@/firebase";
import { exchangeToken } from "../auth/authService";

// const DEFAULT_API_BASE_URL = 'http://localhost:8080'
function extractText(message: ThreadMessage): string {
  return message.content
    .map((part) => {
      if (part.type === "text") return part.text;
      return "";
    })
    .filter((part) => part.trim().length > 0)
    .join("\n");
}

function buildPrompt(messages: readonly ThreadMessage[]): string {
  const lastUserMessage = [...messages].reverse().find((m) => m.role === "user");
  if (!lastUserMessage) return "";
  return extractText(lastUserMessage).trim();
}

export function useBackendRuntime() {
  const chatModelRef = useRef<ChatModelAdapter>({
    async *run({
      messages,
      abortSignal,
    }: ChatModelRunOptions): AsyncGenerator<ChatModelRunResult, void, unknown> {
      const API_BASE_URL =
        import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";
      const ENDPOINT = `${API_BASE_URL}/query/chat`;
      const fetchOptions: RequestInit = {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ input: buildPrompt(messages) }),
        signal: abortSignal,
        credentials: "include", // <-- This is the fetch equivalent of withCredentials: true
      };

      try {
        let response = await fetch(ENDPOINT, fetchOptions);
        const user = firebaseAuth.currentUser;
        if (response.status == 401) {
          if (!user) {
            throw new Error("Unauthorized: No active user session");
          }
          try {
            await exchangeToken(user);
            response = await fetch(ENDPOINT, fetchOptions);
          } catch (refreshErr: any) {
            if (refreshErr.message !== "mfa_required") {
              await firebaseAuth.signOut().catch(() => {});
              window.location.href = "/login";
            }
            throw refreshErr;
          }
        }
        if (!response.ok) {
          throw new Error(`HTTP error! status ${response.status}`);
        }
        if (!response.body) throw new Error("No response body");
        const stream = response.body;
        if (!stream) throw new Error("Missing response body");
        const reader = stream.getReader();
        const decoder = new TextDecoder();

        let currentReasoning = "";
        let fullText = "";

        while (true) {
          const { done, value } = await reader.read();
          if (done) break;

          const rawChunk = decoder.decode(value, { stream: true });
          const events = rawChunk.split("\n\n");

          for (const event of events) {
            if (!event.startsWith("data: ")) continue;

            const jsonString = event.replace("data: ", "").trim();
            if (!jsonString) continue;

            try {
              const parsed = JSON.parse(jsonString);
              if (parsed.type === "text") {
                fullText += parsed.content;
              } else if (parsed.type === "reasoning") {
                currentReasoning += parsed.content + "\n";
              }

              const contentParts: any[] = [];
              if (currentReasoning.trim()) {
                contentParts.push({
                  type:"reasoning",
                  text:currentReasoning.trim(),
                });
              }
              if(fullText.trim() || contentParts.length === 0){
                contentParts.push({
                  type: "text",
                  text: fullText,
                });
              }
              yield {
                content: contentParts,
              };
            } catch (e) {
              console.error("Error parsing chunk:", jsonString, e);
            }
          }
        }
        const finalContentParts: any[] = [];
        if(currentReasoning.trim()){
          finalContentParts.push({type:"reasoning", text:currentReasoning.trim()});
        }
        
        if (fullText.trim() || finalContentParts.length === 0){
          finalContentParts.push({type: "text", text: fullText});
        }

        yield {
          content: finalContentParts,
          status: { type: "complete", reason: "stop" } as const,
          // metadata: { custom: { reasoningText: currentReasoning } },
        };
        return;
        
      } catch (error: any) {
        if (error.name === "AbortError") {
          console.log("Stream aborted by user");
          return;        
        }
        throw new Error(error.message || "Backend request failed");
      }
    },
  });
  return { runtime: useLocalRuntime(chatModelRef.current) };
}
