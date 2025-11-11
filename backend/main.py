from fastapi import FastAPI, Depends, HTTPException
from fastapi.responses import StreamingResponse
from fastapi.middleware.cors import CORSMiddleware
from sqlalchemy.orm import Session
from sqlalchemy import func
import csv
import io

from database import SessionLocal, Vote, get_db
from schemas import VoteCreate, ReportItem

app = FastAPI()

# Enable CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.post("/api/vote", status_code=201)
def submit_vote(vote: VoteCreate, db: Session = Depends(get_db)):
    """Submit a vote for a character"""
    if not vote.characterId:
        raise HTTPException(status_code=400, detail="characterId is required")
    
    db_vote = Vote(
        character_id=str(vote.characterId),
        name=vote.name,
        image_url=vote.image_url
    )
    db.add(db_vote)
    db.commit()
    db.refresh(db_vote)
    return {"id": db_vote.id}

@app.get("/api/report")
def get_report(db: Session = Depends(get_db)):
    """Get aggregated vote counts per character"""
    results = db.query(
        Vote.character_id,
        Vote.name,
        Vote.image_url,
        func.count(Vote.id).label("votes")
    ).group_by(Vote.character_id).order_by(func.count(Vote.id).desc()).all()
    
    return [
        ReportItem(
            characterId=r[0],
            name=r[1],
            imageUrl=r[2],
            votes=r[3]
        ).dict()
        for r in results
    ]

@app.get("/api/votes")
def get_all_votes(db: Session = Depends(get_db)):
    """Get all raw votes"""
    votes = db.query(Vote).order_by(Vote.created_at.desc()).all()
    return votes

@app.get("/api/export")
def export_csv(db: Session = Depends(get_db)):
    """Export aggregated report as CSV"""
    results = db.query(
        Vote.character_id,
        Vote.name,
        func.count(Vote.id).label("votes")
    ).group_by(Vote.character_id).order_by(func.count(Vote.id).desc()).all()
    
    output = io.StringIO()
    writer = csv.writer(output)
    writer.writerow(["characterId", "name", "votes"])
    for r in results:
        writer.writerow([r[0], r[1] or "", r[2]])
    
    output.seek(0)
    return StreamingResponse(
        iter([output.getvalue()]),
        media_type="text/csv",
        headers={"Content-Disposition": "attachment; filename=report.csv"}
    )

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=3000)
