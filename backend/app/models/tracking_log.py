from pydantic import BaseModel, Field
from typing import Optional

class TrackingLogInDB(BaseModel):
    """Cấu trúc Log lưu vết hành trình vật lý trong MongoDB (Dữ liệu thô cho Big Data)"""
    order_code: str = Field(..., description="Mã đơn hàng")
    driver_code: str = Field(..., description="Mã tài xế")
    status_update: str = Field(..., description="Trạng thái cập nhật tại thời điểm log")
    
    # Lưu dưới dạng dictionary để tối ưu tốc độ ghi của MongoDB
    gps_location: Optional[dict] = None
    
    note: Optional[str] = None
    
    # Timestamp là trường bắt buộc ở DB để phục vụ truy vấn Time-series và vẽ biểu đồ
    timestamp: str