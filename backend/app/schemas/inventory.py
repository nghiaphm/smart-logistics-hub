from pydantic import BaseModel, Field
from typing import Optional

class InventoryBase(BaseModel):
    product_id: str
    warehouse_id: str = "MAIN_HUB"
    available: int = 0
    reserved: int = 0
    damaged: int = 0
    hold: int = 0

class InventoryCreate(InventoryBase):
    """Sử dụng khi khởi tạo tồn kho lần đầu cho một sản phẩm tại một kho"""
    pass

class InventoryUpdate(BaseModel):
    """Dùng để cập nhật tăng/giảm hoặc thay đổi số lượng ở các trạng thái cụ thể"""
    available: Optional[int] = None
    reserved: Optional[int] = None
    damaged: Optional[int] = None
    hold: Optional[int] = None
    reason: Optional[str] = Field(None, description="Lý do điều chỉnh (Ví dụ: Trả hàng, Hỏng khi vận chuyển)")

class InventoryResponse(InventoryBase):
    """Dữ liệu trả về cho Frontend kèm theo mã định danh hệ thống"""
    id: str = Field(alias="_id")
    updated_at: str