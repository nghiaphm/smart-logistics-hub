from pydantic import BaseModel, Field
from typing import Optional

# --- CÁC CLASS THÀNH PHẦN ---
class Location(BaseModel):
    lat: float
    lng: float

# --- KHUÔN API GIAO TIẾP VỚI FRONTEND ---
class WarehouseCreate(BaseModel):
    """Dữ liệu Frontend gửi lên khi tạo kho mới (POST)"""
    warehouse_code: str = Field(..., description="Mã kho duy nhất")
    name: str
    address: str
    location: Optional[Location] = None
    contact_phone: Optional[str] = None
    manager_name: Optional[str] = None

class WarehouseUpdate(BaseModel):
    """Dữ liệu dùng để cập nhật thông tin kho (PATCH)"""
    name: Optional[str] = None
    address: Optional[str] = None
    location: Optional[Location] = None
    contact_phone: Optional[str] = None
    manager_name: Optional[str] = None
    is_active: Optional[bool] = None

class WarehouseResponse(WarehouseCreate):
    """Dữ liệu Backend trả về cho Frontend (GET)"""
    id: str = Field(alias="_id")
    is_active: bool
    created_at: str
    updated_at: str