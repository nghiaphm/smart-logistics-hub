from pydantic import BaseModel, Field
from typing import List, Optional

# --- CÁC CLASS THÀNH PHẦN (Phục vụ validate API) ---
class InboundItem(BaseModel):
    product_id: str
    expected_qty: int = Field(..., description="Số lượng dự kiến nhập từ nhà cung cấp")
    received_qty: int = Field(default=0, description="Số lượng thực tế nhận tại cửa kho")
    rejected_qty: int = Field(default=0, description="Số lượng rớt QC (sẽ cộng vào trạng thái Damaged)")
    qc_passed: bool = Field(default=False, description="Đánh dấu đã hoàn tất khâu kiểm định")

# --- KHUÔN API GIAO TIẾP VỚI FRONTEND ---
class InboundCreate(BaseModel):
    """Dữ liệu Frontend gửi lên khi tạo Phiếu nhập kho (POST)"""
    receipt_code: str = Field(..., description="Mã phiếu nhập (VD: INB-20260507)")
    supplier_name: str = Field(..., description="Tên nhà cung cấp")
    items: List[InboundItem]
    
    status: str = Field(
        default="PENDING", 
        description="Trạng thái luồng: PENDING, RECEIVING, QC_CHECKING, COMPLETED"
    )

class InboundUpdate(BaseModel):
    """Dữ liệu dùng khi nhân viên kho cập nhật số lượng QC hoặc thay đổi trạng thái"""
    items: Optional[List[InboundItem]] = None
    status: Optional[str] = None

class InboundResponse(InboundCreate):
    """Dữ liệu Backend trả về cho Frontend (GET/Response)"""
    id: str = Field(alias="_id", description="Mã ObjectId do MongoDB tự sinh")
    created_at: str
    updated_at: Optional[str] = None
    completed_at: Optional[str] = None
    created_by: Optional[str] = None