from pydantic import BaseModel, Field
from typing import Optional

# --- CÁC CLASS THÀNH PHẦN (Phục vụ validate API) ---
class Location(BaseModel):
    lat: float
    lng: float

class Vehicle(BaseModel):
    type: str
    license_plate: str

# --- KHUÔN API GIAO TIẾP VỚI FRONTEND ---
class DriverCreate(BaseModel):
    """Dữ liệu Frontend gửi lên khi tạo Driver (POST)"""
    driver_code: str
    full_name: str  # Dùng full_name để khớp với DB
    phone: str
    vehicle: Vehicle
    status: str = Field(default="AVAILABLE", description="AVAILABLE, BUSY, OFFLINE")
    current_location: Optional[Location] = None
    warehouse_id: str = Field(..., description="Mã kho trực thuộc, ví dụ: HUB_MY_THO")

class DriverResponse(DriverCreate):
    """Dữ liệu Backend trả về cho Frontend (GET/Response)"""
    id: str = Field(alias="_id", description="Mã ObjectId do MongoDB tự sinh")
    created_at: str
    updated_at: Optional[str] = None
    created_by: Optional[str] = None
    model_config = {
        "populate_by_name": True,
        "from_attributes": True
    }