from pydantic import BaseModel, Field
from typing import List, Optional

class InboundInDB(BaseModel):
    """Cấu trúc dữ liệu phiếu nhập kho lưu trữ trong MongoDB"""
    receipt_code: str
    supplier_name: str
    
    # Lưu dưới dạng danh sách các dictionary trong DB
    items: List[dict] 
    
    # Trạng thái để phục vụ luồng vận hành WMS
    status: str = Field(
        default="PENDING", 
        description="Trạng thái: PENDING, RECEIVING, QC_CHECKING, COMPLETED"
    )
    
    # Metadata theo dõi thời gian và người xử lý
    created_at: str
    updated_at: str
    completed_at: Optional[str] = None
    created_by: Optional[str] = None