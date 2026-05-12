from pydantic import BaseModel, Field
from typing import Optional

class ProductInDB(BaseModel):
    """Cấu trúc dữ liệu sản phẩm lưu trữ vật lý trong MongoDB"""
    sku: str = Field(..., description="Mã vạch/SKU duy nhất của sản phẩm")
    name: str
    category: str
    price: float
    weight_gram: float
    
    # Lưu dưới dạng dictionary (JSON) trong DB
    dimensions: Optional[dict] = None
    
    # Metadata quản trị
    created_at: str
    updated_at: str
    created_by: Optional[str] = None