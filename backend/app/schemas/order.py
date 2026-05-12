from pydantic import BaseModel, Field
from typing import List, Optional

# --- CÁC CLASS THÀNH PHẦN (Phục vụ validate API) ---
class Location(BaseModel):
    lat: float
    lng: float

class Sender(BaseModel):
    name: str
    phone: str
    address: str
    location: Optional[Location] = None

class Receiver(BaseModel):
    name: str
    phone: str
    address: str
    location: Optional[Location] = None

class Item(BaseModel):
    product_id: str
    quantity: int
    price: float = 0.0

class OrderTimeline(BaseModel):
    created_at: Optional[str] = None
    reserved_at: Optional[str] = None  # Đã trừ tồn kho available, cộng vào reserved
    packed_at: Optional[str] = None    # Đã đóng gói xong
    shipped_at: Optional[str] = None   # Giao cho tài xế
    delivered_at: Optional[str] = None # Khách nhận thành công

# --- KHUÔN API GIAO TIẾP VỚI FRONTEND ---
class OrderCreate(BaseModel):
    """Dữ liệu Frontend gửi lên khi tạo Order (Input)"""
    order_code: str
    sender: Sender
    receiver: Receiver
    items: List[Item]
    
    status: str = Field(
        default="PENDING", 
        description="PENDING, RESERVED, PICKING, PACKING, SORTING, SHIPPING, COMPLETED"
    )
    timeline: Optional[OrderTimeline] = None
    assigned_driver_id: Optional[str] = None

class OrderResponse(OrderCreate):
    """Dữ liệu Backend trả về cho Frontend (Output)"""
    id: str = Field(alias="_id", description="Mã ObjectId do MongoDB tự sinh")
    created_by: Optional[str] = None
    created_at: Optional[str] = None
    updated_at: Optional[str] = None