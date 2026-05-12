from pydantic import BaseModel, Field
from typing import Optional

# --- CÁC CLASS THÀNH PHẦN (Phục vụ validate API) ---
class Location(BaseModel):
    lat: float
    lng: float

# --- KHUÔN API GIAO TIẾP VỚI FRONTEND/APP TÀI XẾ ---
class TrackingLogCreate(BaseModel):
    """Khuôn dùng cho API POST (Khi App tài xế bắn tọa độ lên Server)"""
    order_code: str = Field(..., description="Mã đơn hàng (VD: KH-9999)")
    driver_code: str = Field(..., description="Mã tài xế (VD: TX-007)")
    status_update: str = Field(..., description="Trạng thái cập nhật (VD: PICKING_UP, DELIVERING, COMPLETED, DELAYED)")
    gps_location: Optional[Location] = Field(default=None, description="Tọa độ GPS hiện tại của tài xế")
    note: Optional[str] = Field(default=None, description="Ghi chú thêm (VD: Tắc đường, Khách không nghe máy,...)")

class TrackingLogResponse(TrackingLogCreate):
    """Khuôn dùng cho API GET (Khi App khách hàng muốn xem lịch sử di chuyển)"""
    id: str = Field(alias="_id", description="Mã ObjectId do MongoDB tự sinh")
    timestamp: str = Field(..., description="Thời gian server ghi nhận log này (ISO Format)")