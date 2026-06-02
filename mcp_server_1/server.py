import os
import logging
import uuid
import json
from dotenv import load_dotenv
from typing import Optional
from urllib import error as urllib_error
from urllib import request as urllib_request

from fastmcp import FastMCP

load_dotenv()

MCP_HOST = os.getenv("MCP_HOST")
MCP_PORT = int(os.getenv("MCP_PORT", "9000"))
MCP_PATH = os.getenv("MCP_PATH")
OLLAMA_BASE_URL = os.getenv("OLLAMA_BASE_URL", "http://localhost:11434").rstrip("/")
OLLAMA_MODEL = os.getenv("OLLAMA_MODEL", "llama3.2")

mcp = FastMCP("support-copilot-mcp")

_logger = logging.getLogger("mcp-tools")
if not _logger.handlers:
    _logger.setLevel(logging.INFO)
    formatter = logging.Formatter("%(asctime)s %(levelname)s %(message)s")

    console_handler = logging.StreamHandler()
    console_handler.setFormatter(formatter)
    _logger.addHandler(console_handler)

    file_handler = logging.FileHandler("mcp_tool_calls.log", encoding="utf-8")
    file_handler.setFormatter(formatter)
    _logger.addHandler(file_handler)


def _chat_with_ollama(message: str, model: Optional[str] = None) -> str:
    selected_model = model or OLLAMA_MODEL
    if not selected_model:
        raise ValueError("OLLAMA_MODEL is required")

    payload = json.dumps(
        {
            "model": selected_model,
            "stream": False,
            "messages": [{"role": "user", "content": message}],
        }
    ).encode("utf-8")

    req = urllib_request.Request(
        f"{OLLAMA_BASE_URL}/api/chat",
        data=payload,
        headers={"Content-Type": "application/json"},
        method="POST",
    )

    try:
        with urllib_request.urlopen(req, timeout=30) as response:
            response_payload = json.loads(response.read().decode("utf-8"))
    except urllib_error.HTTPError as exc:
        raise ValueError(f"Ollama API error ({exc.code}): {exc.read().decode('utf-8')}") from exc

    content = response_payload.get("message", {}).get("content", "").strip()
    if not content:
        content = response_payload.get("response", "").strip()

    if not content:
        raise ValueError("Ollama returned an empty response")

    return content


@mcp.tool(description="Send a user message to Ollama and return the model response")
def chat_with_ollama(message: str, model: Optional[str] = None) -> str:
    if not message or not message.strip():
        raise ValueError("message is required")

    response_text = _chat_with_ollama(message, model)
    _logger.info("tool=chat_with_ollama model=%s input_chars=%d", model or OLLAMA_MODEL, len(message))
    return response_text


@mcp.tool(description="Add two numbers and return the sum")
def add_numbers(a: float, b: float) -> float:
    result = a + b
    _logger.info("tool=add_numbers a=%s b=%s result=%s", a, b, result)
    return result


@mcp.tool(description="Add two numbers and return a verifiable proof token from the MCP server")
def add_numbers_with_proof(a: float, b: float) -> dict:
    result = a + b
    proof = str(uuid.uuid4())
    _logger.info("tool=add_numbers_with_proof a=%s b=%s result=%s proof=%s", a, b, result, proof)
    return {
        "tool": "add_numbers_with_proof",
        "result": result,
        "proof": proof,
    }


if __name__ == "__main__":
    mcp.run(transport="streamable-http", host=MCP_HOST, port=MCP_PORT, path=MCP_PATH)
