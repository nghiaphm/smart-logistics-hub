from pydantic import BaseModel, Field
from typing import Optional

class BillingInDB(BaseModel):
    """Cấu trúc dữ liệu hóa đơn/thanh toán lưu trữ vật lý trong MongoDB"""
    billing_code: str = Field(..., description="Mã hóa đơn duy nhất (VD: INV-2026-001)")
    order_code: str = Field(..., description="Mã đơn hàng tham chiếu")
    
    amount_total: float = Field(..., description="Tổng số tiền cần thanh toán")
    currency: str = Field(default="VND")
    
    payment_method: str = Field(
        default="COD", 
        description="Phương thức thanh toán: COD, VNPAY, BANK_TRANSFER"
    )
    payment_status: str = Field(
        default="UNPAID", 
        description="Trạng thái: UNPAID, PENDING, PAID, FAILED, REFUNDED"
    )
    
    # Lưu mã giao dịch đối soát từ VNPay hoặc Ngân hàng trả về
    transaction_id: Optional[str] = None 
    
    # Thông tin người trả tiền (lưu dạng dict)
    payer_info: dict 
    
    # Metadata thời gian
    created_at: str
    updated_at: str
    paid_at: Optional[str] = None
    created_by: Optional[str] = None