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
Guesses the last name of a user given their first name.

Request body:
```json
{
    "firstName": "Anjana"
}
```

Response:
```json
{
    "lastName": "Paulose",
    "confidence": 1.0
}
```

## Error Handling

- 404: User not found
- 500: Error connecting to user service

## Testing

You can test the API using curl:
```bash
curl -X POST http://localhost:8000/guess \
     -H "Content-Type: application/json" \
     -d '{"firstName": "Anjana"}'
``` 