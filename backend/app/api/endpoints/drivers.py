# app/api/endpoints/drivers.py

from backend.app.api.deps import get_current_driver_from_db
from fastapi import APIRouter, Depends, HTTPException
from pydantic import BaseModel
from datetime import datetime, timezone
from bson import ObjectId
from bson.errors import InvalidId

from ...schemas.driver import DriverCreate, DriverResponse
from ...db.database import drivers_collection
# Import thêm RequireRole ở đây
from ...core.security import verify_token, RequireRole 

router = APIRouter()

# --- CHỨC NĂNG DÀNH CHO ADMIN ---

@router.post("", dependencies=[Depends(RequireRole(["admin"]))])
async def create_driver(driver: DriverCreate, user: dict = Depends(verify_token)):
    """Chỉ Admin mới được tạo tài xế mới"""
    driver_dict = driver.model_dump(by_alias=True)
    driver_dict["status"] = "AVAILABLE"
    driver_dict["current_location"] = None
    driver_dict["created_at"] = datetime.now(timezone.utc).isoformat()
    driver_dict["created_by"] = user.get("preferred_username")
    
    await drivers_collection.insert_one(driver_dict)
    return {"message": "Driver created successfully", "driver_code": driver_dict["driver_code"]}

@router.get("", response_model=list[DriverResponse], dependencies=[Depends(RequireRole(["admin"]))])
async def get_drivers():
    """Chỉ Admin mới được xem danh sách tất cả tài xế"""
    drivers = []
    async for doc in drivers_collection.find({}):
        doc["_id"] = str(doc["_id"])
        drivers.append(doc)
    return drivers

@router.delete("/{driver_id}", dependencies=[Depends(RequireRole(["admin"]))])
async def delete_driver(driver_id: str):
    """Chỉ Admin mới có quyền xóa tài xế"""
    try:
        obj_id = ObjectId(driver_id)
    except InvalidId:
        raise HTTPException(status_code=400, detail="Định dạng ID không hợp lệ")

    result = await drivers_collection.delete_one({"_id": obj_id})
    if result.deleted_count == 1:
        return {"message": f"Đã xóa thành công tài xế có ID: {driver_id}"}
    raise HTTPException(status_code=404, detail="Không tìm thấy tài xế này!")


# --- CHỨC NĂNG DÀNH RIÊNG CHO DRIVER ---

@router.get("/me")
async def get_my_driver_profile(full_data: dict = Depends(get_current_driver_from_db)):
    # Bây giờ dữ liệu này chắc chắn đã tồn tại trong DB
    db_profile = full_data["profile"]
    identity = full_data["identity"]
    
    return {
        "status": "success",
        "driver_info": {
            "name": identity.get("name"),
            "hub": identity.get("HUB_MY_THO"),
            "db_status": db_profile.get("status"),
            "phone_in_db": db_profile.get("phone"),
            "license_plate": db_profile.get("vehicle", {}).get("license_plate")
        }
    }