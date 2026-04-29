from fastapi import APIRouter, Depends
from datetime import datetime
from ...models.order import OrderCreate, OrderResponse
from ...db.database import orders_collection
from ...core.security import verify_token

router = APIRouter()

@router.post("/")
async def create_order(order: OrderCreate, user: dict = Depends(verify_token)):
    order_dict = order.model_dump()
    order_dict["status"] = "PENDING"
    # order_dict["created_at"] = datetime.utcnow().isoformat()
    order_dict["created_by"] = user.get("preferred_username")
    
    # Tạo object timeline mặc định
    order_dict["timeline"] = {
        "created_at": datetime.utcnow().isoformat(),
        "picked_up_at": None,
        "delivered_at": None
    }

    await orders_collection.insert_one(order_dict)
    return {"message": "Order created successfully"}

# Khai báo response_model để FastAPI tự động chuyển đổi ObjectId thành string và validate dữ liệu trả về
@router.get("/", response_model=list[OrderResponse])
async def get_orders(user: dict = Depends(verify_token)):
    orders = []
    async for doc in orders_collection.find({}):
        # Pydantic OrderResponse có alias="_id" nên nó sẽ tự nhận doc["_id"],
        # nhưng MongoDB trả về ObjectId, ta cần ép về string trước.
        doc["_id"] = str(doc["_id"])
        orders.append(doc)
    return orders