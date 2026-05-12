from pydantic import BaseModel, Field
from typing import List, Optional

class TripInDB(BaseModel):
    """Cấu trúc dữ liệu chuyến xe lưu trữ vật lý trong MongoDB"""
    trip_code: str = Field(..., description="Mã chuyến xe (VD: TRIP-2026-001)")
    driver_id: Optional[str] = Field(None, description="Mã _id của tài xế phụ trách")
    
    # Danh sách các mã đơn hàng trong chuyến này (giữ đúng thứ tự giao hàng)
    order_ids: List[str] = Field(default=[], description="Danh sách các _id đơn hàng")
    
    vehicle_info: Optional[dict] = None
    
    status: str = Field(
        default="PLANNED", 
        description="Trạng thái: PLANNED, DISPATCHED, IN_PROGRESS, COMPLETED, CANCELLED"
    )
    
    # Metadata hành trình
    total_distance_km: float = 0.0
    estimated_duration_min: int = 0
    
    actual_start_at: Optional[str] = None
    actual_end_at: Optional[str] = None
    
    created_at: str
    updated_at: str
    created_by: Optional[str] = None