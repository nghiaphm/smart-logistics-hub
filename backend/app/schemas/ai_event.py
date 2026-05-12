from pydantic import BaseModel, Field
from typing import Optional

class AIEventCreate(BaseModel):
    """Payload mà AI Service (Colab/Microservice) sẽ bắn về Core Backend"""
    license_plate: str
    confidence_score: float = Field(..., ge=0.0, le=1.0) # Ép kiểu điểm AI từ 0-1
    event_type: str = Field(..., description="INBOUND, OUTBOUND")
    gate_id: str = "MAIN_GATE"
    timestamp: str

class AIEventResponse(AIEventCreate):
    """Trả về cho Dashboard Admin hiển thị log xe ra vào"""
    id: str = Field(alias="_id")
    event_code: str
    matched_driver_id: Optional[str] = None
    matched_trip_id: Optional[str] = None