# Support Copilot

This project now includes:

- A Go backend API (`POST /query/sc`) that returns Gemini responses.
- A React chat frontend in `frontend/`.
- A Python FastMCP server in `mcp_server/` exposing a Gemini chat tool.

## Environment Variables

Copy `.env.example` to `.env.local` (or `.env`) and set:

- `GEMINI_API_KEY`: your Google GenAI API key (required).
- `GEMINI_MODEL`: Gemini model name, default `gemini-2.0-flash`.
- `AUTH_ENABLED`: set `true` to enforce Authorization header checks on `/query/sc`.

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

Body:

```json
{
	"input": "How can I reset my password?"
}
```

Response:

```json
{
	"output": "...Gemini response..."
}
```

## FastMCP Tool

The FastMCP server defines a tool:

- `chat_with_gemini(message: string, model?: string) -> string`

You can connect using any MCP client that supports streamable HTTP transport.
