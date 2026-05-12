from pydantic import BaseModel, Field
from typing import Optional

# --- CÁC CLASS THÀNH PHẦN ---
class PayerInfo(BaseModel):
    name: str
    phone: str
    email: Optional[str] = None

# --- KHUÔN API GIAO TIẾP VỚI FRONTEND/PAYMENT GATEWAY ---
class BillingCreate(BaseModel):
    """Dữ liệu Frontend gửi lên khi tạo hóa đơn mới (POST)"""
    billing_code: str = Field(..., description="Mã hóa đơn (VD: INV-2026-001)")
    order_code: str = Field(..., description="Mã đơn hàng liên kết")
    
    amount_total: float = Field(..., description="Tổng tiền thanh toán")
    currency: str = "VND"
    
    payment_method: str = Field(default="COD", description="COD, VNPAY, BANK_TRANSFER")
    payment_status: str = Field(default="UNPAID", description="UNPAID, PENDING, PAID, FAILED")
    
    payer_info: PayerInfo

class BillingUpdate(BaseModel):
    """Dữ liệu dùng khi nhận Webhook/IPN từ cổng thanh toán để cập nhật trạng thái (PATCH)"""
    payment_status: Optional[str] = None
    transaction_id: Optional[str] = Field(None, description="Mã giao dịch do VNPay/Ngân hàng cấp")
    paid_at: Optional[str] = None

class BillingResponse(BillingCreate):
    """Dữ liệu trả về cho Frontend hiển thị chi tiết hóa đơn (GET)"""
    id: str = Field(alias="_id")
    transaction_id: Optional[str] = None
    created_at: str
    updated_at: str
    paid_at: Optional[str] = None
    created_by: Optional[str] = None