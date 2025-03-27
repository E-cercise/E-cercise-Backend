
from fastapi import FastAPI
from pydantic import BaseModel
from typing import List, Optional
import uvicorn
import json

app = FastAPI()

with open("equipment_options.json") as f:
    EQUIPMENT_DATA = json.load(f)

class Preference(BaseModel):
    tag: Optional[str] = None
    max_price: Optional[float] = None
    min_weight: Optional[float] = None

class RecommendRequest(BaseModel):
    user_type: Optional[str] = None
    weight: Optional[float] = None
    height: Optional[float] = None
    age: Optional[int] = None
    gender: Optional[str] = None
    goal: Optional[str] = None
    experience: Optional[str] = None
    preferences: Optional[List[Preference]] = []

@app.post("/recommend")
async def recommend(req: RecommendRequest):
    results = []

    for option in EQUIPMENT_DATA:
        score = 0
        text = json.dumps(option).lower()

        # Age-based personalization
        if req.age:
            if req.age >= 55 and "joint-friendly" in text:
                score += 10
            if req.age < 25 and "multi-function" in text:
                score += 5

        # Gender-aware scoring
        if req.gender == "female":
            if "glutes" in text or "abs" in text or "core" in text:
                score += 6
            if option["weight"] <= 60:
                score += 4
        elif req.gender == "male":
            if "arms" in text or "chest" in text or "pull-up" in text:
                score += 6
            if option["weight"] >= 60:
                score += 4

        # Experience scoring
        if req.experience == "beginner" and option["weight"] <= 50:
            score += 10
        elif req.experience == "athlete" and option["weight"] >= 100:
            score += 10
        elif req.experience == "elderly":
            if "joint-friendly" in text or "low-impact" in text:
                score += 10
            if option["weight"] < 30:
                score += 5

        # Weight/height personalization
        if req.weight and req.weight < 60 and option["weight"] < 50:
            score += 5
        if req.height and req.height > 180 and "adjustable" in text:
            score += 3

        # Goal-based scoring
        goal = (req.goal or "").lower()
        if goal == "tone":
            if "resistance" in text or "bodyweight" in text:
                score += 8
            if "abs" in text or "core" in text:
                score += 5
            if option["weight"] <= 50:
                score += 5
        elif goal == "build-muscle":
            if option["weight"] >= 80:
                score += 10
            if "chest" in text or "arms" in text or "back" in text or "legs" in text:
                score += 6
            if "weighted" in text or "resistance" in text:
                score += 5
        elif goal == "rehab":
            if "joint-friendly" in text or "low-impact" in text:
                score += 10
            if option["weight"] < 40:
                score += 6
            if "stretching" in text or "mobility" in text:
                score += 5

        # Tag Preferences
        for pref in req.preferences or []:
            if pref.tag and pref.tag.lower() in text:
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