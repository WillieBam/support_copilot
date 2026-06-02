# Support Copilot

This project now includes:

- A Go backend API (`POST /query/sc`) that returns Ollama responses.
- A React chat frontend in `frontend/`.
- A Python FastMCP server in `mcp_server/` exposing an Ollama chat tool.

## Environment Variables

Copy `.env.example` to `.env.local` (or `.env`) and set:

- `OLLAMA_BASE_URL`: Ollama server URL, default `http://localhost:11434`.
- `OLLAMA_MODEL`: Ollama model name, default `llama3.1`.
- `AUTH_ENABLED`: set `true` to enforce Authorization header checks on `/query/sc`.
- `AUTH_TOTP_REQUIRED`: set `true` to require TOTP as the Firebase second factor.
- `FIREBASE_PROJECT_ID`: Firebase project ID used by the backend Admin SDK.

Frontend (`frontend/.env` or `frontend/.env.local`):

- `VITE_FIREBASE_API_KEY`
- `VITE_FIREBASE_AUTH_DOMAIN`
- `VITE_FIREBASE_PROJECT_ID`
- `VITE_FIREBASE_APP_ID`
- `VITE_API_BASE_URL` (optional, default `http://localhost:8080`)

## Run with Docker Compose

From workspace root:

```bash
docker compose --env-file .env.local up --build
```

Services:

- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`
- FastMCP (streamable HTTP): `http://localhost:9000/mcp`

## Backend Chat Endpoint

Endpoint: `POST /query/sc`

Headers:

- `Content-Type: application/json`
- `Authorization: <bearer-token>` (required when auth middleware is enabled)

Notes:

- The backend validates Firebase ID tokens with Firebase Admin SDK.
- When `AUTH_TOTP_REQUIRED=true`, the token's second factor must be TOTP.

Body:

```json
{
	"input": "How can I reset my password?"
}
```

Response:

```json
{
	"output": "...assistant response..."
}
```

## FastMCP Tool

The FastMCP server defines a tool:

- `chat_with_ollama(message: string, model?: string) -> string`

You can connect using any MCP client that supports streamable HTTP transport.
