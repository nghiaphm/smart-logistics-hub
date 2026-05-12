from pydantic import BaseModel, Field
from typing import Optional

class AIEventInDB(BaseModel):
    """Cấu trúc log sự kiện do AI Service gửi về lưu trong MongoDB"""
    event_code: str = Field(..., description="Mã sự kiện duy nhất")
    
    # Dữ liệu AI trích xuất được
    license_plate: str = Field(..., description="Biển số xe nhận diện được")
    confidence_score: float = Field(..., description="Độ tự tin của model (0.0 -> 1.0)")
    
    # Ngữ cảnh sự kiện
    event_type: str = Field(..., description="INBOUND (Vào cổng) hoặc OUTBOUND (Ra cổng)")
    gate_id: str = Field(default="MAIN_GATE", description="Khu vực cổng camera")
    
    # Metadata
    timestamp: str = Field(..., description="Thời gian thực tế camera ghi nhận")
    created_at: str
    
    # Nếu đối chiếu được với database tài xế/chuyến xe thì gán vào đây
    matched_driver_id: Optional[str] = None
    matched_trip_id: Optional[str] = None