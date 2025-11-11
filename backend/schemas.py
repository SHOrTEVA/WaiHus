from pydantic import BaseModel
from typing import Optional

class VoteCreate(BaseModel):
    characterId: str
    name: Optional[str] = None
    image_url: Optional[str] = None

class VoteResponse(BaseModel):
    id: int
    character_id: str
    name: Optional[str]
    image_url: Optional[str]
    
    class Config:
        from_attributes = True

class ReportItem(BaseModel):
    characterId: str
    name: Optional[str]
    imageUrl: Optional[str]
    votes: int
