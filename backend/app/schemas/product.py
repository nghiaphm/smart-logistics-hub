from pydantic import BaseModel, Field
from typing import Optional

# --- CÁC CLASS THÀNH PHẦN (Phục vụ validate API) ---
class Dimensions(BaseModel):
    length: float = 0.0
    width: float = 0.0
    height: float = 0.0

# --- KHUÔN API GIAO TIẾP VỚI FRONTEND ---
class ProductCreate(BaseModel):
    """Dữ liệu Frontend gửi lên hoặc dữ liệu map từ file Excel (POST/Import)"""
    sku: str = Field(..., description="Mã vạch/SKU duy nhất của sản phẩm (VD: KEYCHRON-Q1)")
    name: str = Field(..., description="Tên sản phẩm")
    category: str = Field(..., description="Ngành hàng (VD: Electronics)")
    price: float = Field(..., description="Giá bán")
    weight_gram: float = Field(..., description="Trọng lượng (gram) để tính phí ship")
    dimensions: Optional[Dimensions] = None

class ProductUpdate(BaseModel):
    """Dữ liệu dùng để cập nhật thông tin sản phẩm (PUT/PATCH)"""
    name: Optional[str] = None
    category: Optional[str] = None
    price: Optional[float] = None
    weight_gram: Optional[float] = None
    dimensions: Optional[Dimensions] = None

class ProductResponse(ProductCreate):
    """Dữ liệu Backend trả về cho Frontend (GET/Response)"""
    id: str = Field(alias="_id", description="Mã ObjectId do MongoDB tự sinh")
    created_at: str
    updated_at: Optional[str] = None
    created_by: Optional[str] = None