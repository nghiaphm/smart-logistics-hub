from pydantic import BaseModel, Field
from typing import Optional, Literal

class VehicleInfo(BaseModel):
    type: str # eg: "VAN", "TRUCK", "MOTORBIKE", ...
    license_plate: str
    capacity_max_kg: float
class CurrentCoordinates(BaseModel):
    # Dùng Literal để "ép" MongoDB hiểu đây là định dạng GeoJSON chuẩn
    type: Literal["Point"] = "Point"
    # Array chứa 2 số : [Kinh độ(lng), Vĩ độ(lat)]
    coordinates: list[float]

# Khuôn đầu vào (khi thêm mới / đăng ký tài xế)
class DriverCreate(BaseModel):
    # Id trùng với Id do keycloak sinh ra để dễ đồng bộ bảo mật
    id: str = Field(alias="_id")
    full_name : str
    phone : str
    vehicle: VehicleInfo
class DriverResponse(DriverCreate):
    status: str = "AVAILABLE" # "AVAILABLE", "ON_DUTY", "OFFLINE"
    current_location: Optional[CurrentCoordinates] = None