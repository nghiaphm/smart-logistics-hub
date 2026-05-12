from pydantic import BaseModel, Field
from typing import Optional

class InventoryInDB(BaseModel):
    """Cấu trúc dữ liệu tồn kho lưu trữ vật lý trong MongoDB"""
    product_id: str = Field(..., description="Mã định danh sản phẩm")
    warehouse_id: str = Field(default="MAIN_HUB", description="Mã định danh kho hàng")
    
    # 4 trạng thái cốt lõi để theo dõi chính xác dòng hàng
    available: int = Field(default=0, description="Số lượng sẵn sàng bán")
    reserved: int = Field(default=0, description="Số lượng đã được giữ (chờ đóng gói/giao)")
    damaged: int = Field(default=0, description="Số lượng hàng lỗi/hỏng chờ xử lý")
    hold: int = Field(default=0, description="Số lượng hàng đang bị tạm giữ")
    
    # Thông tin quản lý cập nhật
    created_at: str
    updated_at: str
    last_updated_by: Optional[str] = None