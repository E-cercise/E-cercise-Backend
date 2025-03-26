from fastapi import FastAPI, Request
from pydantic import BaseModel
from typing import List, Optional
import uvicorn
import json

app = FastAPI()

# Load mock data (normally you'd load from a DB)
with open("equipment_options.json") as f:
    EQUIPMENT_DATA = json.load(f)

class RecommendRequest(BaseModel):
    user_type: Optional[str] = None
    weight: Optional[float] = None
    height: Optional[float] = None
    goal: Optional[str] = None   # <--- this stays!
    preferences: Optional[List[Preference]] = []

class RecommendRequest(BaseModel):
    user_type: str
    preferences: Optional[List[Preference]] = []

@app.post("/recommend")
async def recommend(req: RecommendRequest):
    results = []
    for option in EQUIPMENT_DATA:
        score = 0
        if "weight" in user_data and user_data["weight"] < 60 and option["weight"] < 50:
            score += 5
        if "height" in user_data and "adjustable" in json.dumps(option).lower():
            score += 3
        if req.user_type == "beginner" and option["weight"] <= 50:
            score += 10
        if req.user_type == "athlete" and option["weight"] >= 100:
            score += 10
        for pref in req.preferences:
            if pref.tag and pref.tag in json.dumps(option).lower():
                score += 5
            if pref.max_price and option["price"] <= pref.max_price:
                score += 5
            if pref.min_weight and option["weight"] >= pref.min_weight:
                score += 5
        if score > 0:
            option["score"] = score
            results.append(option)
    results.sort(key=lambda x: x["score"], reverse=True)
    return results[:10]
