from pydantic import BaseModel, Field
from typing import List, Optional

# --- CÁC CLASS THÀNH PHẦN (Phục vụ validate API) ---
class Location(BaseModel):
    lat: float
    lng: float

class TripStop(BaseModel):
    order_code: str
    stop_type: str  # PICKUP hoặc DROP_OFF
    location: Location
    address: str
    status: str = Field(default="PENDING") # PENDING, ARRIVED, SUCCESS, FAILED
    estimated_arrival_time: Optional[str] = None
    actual_arrival_time: Optional[str] = None

# --- KHUÔN API GIAO TIẾP VỚI FRONTEND ---
class TripCreate(BaseModel):
    """Dữ liệu dùng khi tạo hoặc lập kế hoạch chuyến xe (POST)"""
    trip_code: str
    driver_code: str
    vehicle_license_plate: Optional[str] = None
    stops: List[TripStop]
    status: str = "DRAFT"
    total_distance_km: Optional[float] = 0.0

class TripUpdate(BaseModel):
    """Dùng khi điều phối viên hoặc tài xế cập nhật trạng thái chuyến đi"""
    status: Optional[str] = None
    stops: Optional[List[TripStop]] = None
    started_at: Optional[str] = None
    completed_at: Optional[str] = None

class TripResponse(TripCreate):
    """Dữ liệu Backend trả về cho Frontend (GET)"""
    id: str = Field(alias="_id")
    created_at: str
    updated_at: Optional[str] = None
    started_at: Optional[str] = None
    completed_at: Optional[str] = None
    created_by: Optional[str] = None