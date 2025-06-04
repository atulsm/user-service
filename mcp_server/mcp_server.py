from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import requests
from typing import Optional, List
import uvicorn
from difflib import SequenceMatcher

app = FastAPI(title="Name Guessing MCP Server")

class User(BaseModel):
    id: str
    email: str
    firstName: str
    lastName: str
    createdAt: str
    phoneNumber: Optional[str] = None

class GuessRequest(BaseModel):
    firstName: str

class GuessResponse(BaseModel):
    lastName: str
    confidence: float

def calculate_similarity(str1: str, str2: str) -> float:
    """Calculate similarity ratio between two strings."""
    return SequenceMatcher(None, str1.lower(), str2.lower()).ratio()

@app.get("/")
async def root():
    return {"message": "Name Guessing MCP Server is running"}

@app.post("/guess", response_model=GuessResponse)
async def guess_last_name(request: GuessRequest):
    try:
        # Fetch users from the API
        response = requests.get("http://localhost:8080/api/users")
        response.raise_for_status()
        users = response.json()

        # Find users with partial matches
        matches = []
        for user in users:
            similarity = calculate_similarity(request.firstName, user["firstName"])
            if similarity > 0.5:  # Threshold for considering a match
                matches.append((user, similarity))

        if not matches:
            raise HTTPException(status_code=404, detail=f"No user found with first name similar to: {request.firstName}")
        
        # Sort matches by similarity score
        matches.sort(key=lambda x: x[1], reverse=True)
        
        # Return the best match
        best_match, confidence = matches[0]
        return GuessResponse(
            lastName=best_match["lastName"],
            confidence=confidence
        )

    except requests.RequestException as e:
        raise HTTPException(status_code=500, detail=f"Error connecting to user service: {str(e)}")

if __name__ == "__main__":
    uvicorn.run("mcp_server:app", host="0.0.0.0", port=8000, reload=True) 