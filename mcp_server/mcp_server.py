from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import requests
from typing import Optional, List
import uvicorn

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

        # Find the user with matching first name
        matching_users = [user for user in users if user["firstName"].lower() == request.firstName.lower()]
        
        if not matching_users:
            raise HTTPException(status_code=404, detail=f"No user found with first name: {request.firstName}")
        
        # Return the last name with 100% confidence since we have exact match
        return GuessResponse(
            lastName=matching_users[0]["lastName"],
            confidence=1.0
        )

    except requests.RequestException as e:
        raise HTTPException(status_code=500, detail=f"Error connecting to user service: {str(e)}")

if __name__ == "__main__":
    uvicorn.run("mcp_server:app", host="0.0.0.0", port=8000, reload=True) 