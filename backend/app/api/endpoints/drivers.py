from fastapi import APIRouter, Depends, HTTPException # Chuyển HTTPException sang đây
from pydantic import BaseModel
from datetime import datetime, timezone
from bson import ObjectId # Cần thiết để xử lý ID của MongoDB
from bson.errors import InvalidId # Xử lý lỗi nếu ID không đúng định dạng

from ...models.driver import DriverCreate, DriverResponse
from ...db.database import drivers_collection
from ...core.security import verify_token

router = APIRouter()

# Đã bỏ dấu "/" đi để tránh lỗi 307 Redirect
@router.post("")
async def create_driver(driver: DriverCreate, user: dict = Depends(verify_token)):
    driver_dict = driver.model_dump(by_alias=True)
    driver_dict["status"] = "AVAILABLE"
    driver_dict["current_location"] = None
    
    # Dùng datetime có timezone chuẩn
    driver_dict["created_at"] = datetime.now(timezone.utc).isoformat()
    driver_dict["created_by"] = user.get("preferred_username")
    
    await drivers_collection.insert_one(driver_dict)
    return {"message": "Driver created successfully", "driver_code": driver_dict["driver_code"]}
    

# Đã bỏ dấu "/" đi
@router.get("", response_model=list[DriverResponse])
async def get_drivers(user: dict = Depends(verify_token)):
    drivers = []
    async for doc in drivers_collection.find({}):
        doc["_id"] = str(doc["_id"])
        drivers.append(doc)
    return drivers


class LocationUpdate(BaseModel):
    lat: float
    lng: float


@router.delete("/{driver_id}")
async def delete_driver(driver_id: str, user: dict = Depends(verify_token)):
    try:
        # Ép kiểu String thành ObjectId để MongoDB có thể hiểu được
        obj_id = ObjectId(driver_id)
    except InvalidId:
        raise HTTPException(status_code=400, detail="Định dạng ID không hợp lệ")

    result = await drivers_collection.delete_one({"_id": obj_id})
    
    if result.deleted_count == 1:
        return {"message": f"Đã xóa thành công tài xế có ID: {driver_id}"}
    else:
        raise HTTPException(status_code=404, detail="Không tìm thấy tài xế này!")