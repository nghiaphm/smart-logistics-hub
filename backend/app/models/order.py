from pydantic import BaseModel, Field
from typing import Optional, List
from enum import Enum

# 1. ENUMS (GIỚI HẠN LỰA CHỌN GIỐNG DROPDOWN TRÊN UI)

class ServiceType(str, Enum):
    STANDARD = "STANDARD"
    EXPRESS = "EXPRESS"
    
class PickupType(str, Enum):
    PICKUP = "PICKUP" # Lấy hàng tận nơi
    DROPOFF = "DROPOFF" # Gửi tại điểm

class FeePayer(str, Enum):
    SENDER = "SENDER" # Người gửi trả phí
    RECEIVER = "RECEIVER" # Người nhận trả phí

# 2. Các khuông thành phần (Nested Models)

class Location(BaseModel):
    lat: float
    lng: float

class AddressInfo(BaseModel):
    name: str
    phone: str
    address: str #Số nhà, đường
    province: str # Tỉnh/Thành phố
    district: str # Quận/Huyện
    ward: str # Phường/Xã
    postal_code: Optional[str] = None
    location: Optional[Location] = None

# Định nghĩa 1 món hàng trong hộp
class ProductItem(BaseModel):
    name: str
    price: float = Field(ge = 0) # Giá phải >= 0
    weight_kg: float = Field(ge = 0) # Trọng lượng phải >= 0
    quantity: int = Field(default = 1, ge = 0) # Số lượng phải >= 0, mặc định là 1 nếu không cung cấp

class PackageInfo(BaseModel):
    items: List[ProductItem] # Tương đương danh sách sản phẩm trên UI
    total_weight_kg: float
    length_cm: Optional[float] = None
    width_cm: Optional[float] = None
    height_cm: Optional[float] = None
    declared_value: float = 0.0 # Giá trị bưu gửi (để bảo hiểm)

class ShippingInfo(BaseModel):
    service_type: ServiceType = ServiceType.STANDARD
    pickup_type: PickupType = PickupType.PICKUP
    pickup_time: Optional[str] = None # Ví dụ: "28-04-2026 10h-18h"

class PaymentInfo(BaseModel):
    fee_payer: FeePayer = FeePayer.SENDER
    cod_amount: float = 0.0   # Tiền thu hộ (COD)


class TimelineInfo(BaseModel):
    created_at: str
    picked_up_at: Optional[str] = None
    delivered_at: Optional[str] = None

# Request
class OrderCreate(BaseModel):
    order_code: str
    package: PackageInfo
    sender: AddressInfo
    receiver: AddressInfo

# Response
class OrderResponse(OrderCreate):
    id: str = Field(alias="_id")
    status: str
    assigned_driver_id: Optional[str] = None
    timeline: TimelineInfo
    created_by: str
