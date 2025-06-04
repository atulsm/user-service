# Name Guessing MCP Server

This is a Model Control Protocol (MCP) server that connects to the user service API and provides name guessing functionality.

## Setup

1. Create a virtual environment:
```bash
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
```

2. Install dependencies:
```bash
pip install -r requirements.txt
```

## Running the Server

1. Make sure the user service API is running at http://localhost:8080
2. Start the MCP server:
```bash
python mcp_server.py
```

The server will start on http://localhost:8000

## API Endpoints

### POST /guess
Guesses the last name of a user given their first name. Supports partial matches with confidence scores.

Request body:
```json
{
    "firstName": "Atul"
}
```

Response:
```json
{
    "lastName": "Soman",
    "confidence": 1.0
}
```

## Error Handling

- 404: User not found
- 500: Error connecting to user service

## Cursor MCP Configuration

To use this server with Cursor, add the following configuration to your `~/.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "name-guessing-mcp-server": {
      "command": "curl",
      "args": [
        "http://localhost:8000/guess",
        "-X",
        "POST",
        "-H",
        "Content-Type: application/json",
        "-d",
        "{\"firstName\": \"{{input}}\"}"
      ],
      "responseMapping": {
        "lastName": "{{response.lastName}}",
        "confidence": "{{response.confidence}}"
      }
    }
  }
}
```

This configuration:
1. Sets up a name-guessing MCP server
2. Uses curl to make POST requests to the server
3. Maps the response fields to lastName and confidence
4. Supports partial name matching with confidence scores

## Testing

You can test the API using curl:
```bash
curl -X POST http://localhost:8000/guess \
     -H "Content-Type: application/json" \
     -d '{"firstName": "Anjana"}'
```

## Features

- Exact name matching
- Partial name matching with similarity scores
- Case-insensitive matching
- Confidence scores based on string similarity
- Integration with Cursor MCP 