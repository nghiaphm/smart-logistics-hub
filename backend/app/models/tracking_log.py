from pydantic import BaseModel, Field
from typing import Optional

class TrackingLogCreate(BaseModel):
    id: str
    order_code: str
    driver_id: str
    status_update: str
    gps_location: Optional[dict] = None
    timestamp: str

class TrackingLogResponse(TrackingLogCreate):
    id: str = Field(alias="_id")
    
    