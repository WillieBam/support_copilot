import os
import logging
import uuid
from dotenv import load_dotenv
from typing import Optional

from fastmcp import FastMCP
from google import genai

load_dotenv()

MCP_HOST = os.getenv("MCP_HOST")
MCP_PORT = int(os.getenv("MCP_PORT", "9000"))
MCP_PATH = os.getenv("MCP_PATH")
DEFAULT_MODEL = os.getenv("GEMINI_MODEL")
GEMINI_API_KEY = os.getenv("GEMINI_API_KEY")

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


def _build_client() -> genai.Client:
    if not GEMINI_API_KEY:
        raise ValueError("GEMINI_API_KEY is required")

    return genai.Client(api_key=GEMINI_API_KEY)


@mcp.tool(description="Send a user message to Gemini and return the model response")
def chat_with_gemini(message: str, model: Optional[str] = None) -> str:
    if not message or not message.strip():
        raise ValueError("message is required")

    client = _build_client()
    selected_model = model or DEFAULT_MODEL

    response = client.models.generate_content(
        model=selected_model,
        contents=message,
    )

    _logger.info("tool=chat_with_gemini model=%s input_chars=%d", selected_model, len(message))

    if not response.text:
        raise ValueError("Gemini returned an empty response")

    return response.text


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
