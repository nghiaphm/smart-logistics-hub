from pydantic import BaseModel, Field
from typing import Optional

class WarehouseInDB(BaseModel):
    """Cấu trúc dữ liệu kho hàng lưu trữ vật lý trong MongoDB"""
    warehouse_code: str = Field(..., description="Mã kho (VD: WH-HCM-01)")
    name: str = Field(..., description="Tên kho hàng")
    address: str
    
    # Lưu tọa độ dưới dạng dictionary trong DB
    location: Optional[dict] = None
    
    contact_phone: Optional[str] = None
    manager_name: Optional[str] = None
    
    # Metadata quản trị
    created_at: str
    updated_at: str
    is_active: bool = True