from pydantic import BaseModel, Field
from typing import List, Optional

class OrderInDB(BaseModel):
    """Cấu trúc đơn hàng vật lý lưu trong MongoDB"""
    order_code: str
    sender: dict
    receiver: dict
    items: List[dict]
    
    # Các trạng thái sinh tử của đơn hàng
    status: str = Field(
        default="PENDING", 
        description="PENDING, RESERVED, PICKING, PACKING, SORTING, SHIPPING, COMPLETED"
    )
    
    # Tracking thời gian và người phụ trách
    timeline: dict = {}
    assigned_driver_id: Optional[str] = None
    
    # Metadata lưu trữ
    created_at: str
    updated_at: str
    created_by: Optional[str] = None