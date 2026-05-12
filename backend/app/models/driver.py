from pydantic import BaseModel, Field
from typing import Optional

class DriverInDB(BaseModel):
    """Cấu trúc dữ liệu tài xế lưu trữ vật lý trong MongoDB"""
    driver_code: str
    full_name: str
    phone: str
    
    # Lưu dưới dạng dictionary trong DB
    vehicle: dict 
    
    # Trạng thái để phục vụ điều phối
    status: str = Field(
        default="AVAILABLE", 
        description="Trạng thái: AVAILABLE, BUSY, OFFLINE"
    )
    current_location: Optional[dict] = None
    
    # Metadata quản trị
    created_at: str
    updated_at: str
    created_by: Optional[str] = None