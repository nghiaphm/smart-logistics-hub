from fastapi import APIRouter, Depends, HTTPException
from datetime import datetime, timezone
from pydantic import BaseModel
from bson import ObjectId
from bson.errors import InvalidId

from ...schemas.order import OrderCreate, OrderResponse
from ...db.database import orders_collection, drivers_collection # Cần thêm drivers_collection để kiểm tra tài xế
from ...core.security import verify_token, RequireRole # Import RequireRole để phân quyền

router = APIRouter()

# Schema nhận dữ liệu phân công
class AssignDriverRequest(BaseModel):
    driver_code: str

# 1. API Tạo đơn hàng - Dành cho mọi người dùng đã đăng nhập
@router.post("")
async def create_order(order: OrderCreate, user: dict = Depends(verify_token)):
    order_dict = order.model_dump()
    order_dict["status"] = "PENDING"
    
    now_iso = datetime.now(timezone.utc).isoformat()
    order_dict["created_by"] = user.get("preferred_username")
    
    order_dict["timeline"] = {
        "created_at": now_iso,
        "picked_up_at": None,
        "delivered_at": None
    }

    order_dict["spx_order_code"] = f"SPXVN{int(datetime.now().timestamp())}"

    await orders_collection.insert_one(order_dict)
    return {"message": "Order created successfully", "spx_code": order_dict["spx_order_code"]}


# 2. API Lấy danh sách đơn hàng - Dành cho mọi người dùng đã đăng nhập
@router.get("", response_model=list[OrderResponse])
async def get_orders(user: dict = Depends(verify_token)):
    orders = []
    async for doc in orders_collection.find({}):
        doc["_id"] = str(doc["_id"])
        orders.append(doc)
    return orders


# 3. API Phân công tài xế - CHỈ DÀNH CHO ADMIN
@router.patch("/{order_id}/assign-driver", dependencies=[Depends(RequireRole(["admin"]))])
async def assign_driver(order_id: str, request: AssignDriverRequest):
    try:
        obj_order_id = ObjectId(order_id)
    except InvalidId:
        raise HTTPException(status_code=400, detail="Mã đơn hàng không hợp lệ")

    # Kiểm tra sự tồn tại và trạng thái của tài xế
    driver = await drivers_collection.find_one({"driver_code": request.driver_code})
    if not driver:
        raise HTTPException(status_code=404, detail=f"Không tìm thấy tài xế {request.driver_code}")
    
    if driver.get("status") != "AVAILABLE":
        raise HTTPException(status_code=400, detail="Tài xế đang bận hoặc không online")

    # Cập nhật đơn hàng
    result = await orders_collection.update_one(
        {"_id": obj_order_id},
        {"$set": {
            "assigned_driver_id": request.driver_code,
            "status": "PICKING_UP"
        }}
    )

    if result.modified_count == 0:
        raise HTTPException(status_code=404, detail="Không tìm thấy đơn hàng")

    # Cập nhật trạng thái tài xế sang BUSY
    await drivers_collection.update_one(
        {"driver_code": request.driver_code},
        {"$set": {"status": "BUSY"}}
    )

    return {"message": "Phân công tài xế thành công!"}


# 4. API Dọn rác dữ liệu - CHỈ DÀNH CHO ADMIN
@router.delete("/clear-old-data", dependencies=[Depends(RequireRole(["admin"]))])
async def clear_database():
    """API này giúp xóa sạch toàn bộ đơn hàng cũ bị lỗi cấu trúc"""
    result = await orders_collection.delete_many({})
    return {"message": f"Đã xóa sạch {result.deleted_count} đơn hàng cũ khỏi MongoDB!"}