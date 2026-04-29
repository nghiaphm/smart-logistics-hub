from http.client import HTTPException

from pydantic import BaseModel

from fastapi import APIRouter, Depends
from datetime import datetime
from ...models.driver import DriverCreate, DriverResponse
from ...db.database import drivers_collection
from ...core.security import verify_token

router = APIRouter()

@router.post("/")
async def create_driver(driver: DriverCreate, user: dict = Depends(verify_token)):

    driver_dict = driver.model_dump(by_alias=True) # Sử dụng model_dump() để có alias="_id" tự động
    driver_dict["status"] = "AVAILABLE"

    driver_dict["current_location"] = None
    driver_dict["created_at"] = datetime.utcnow().isoformat()
    driver_dict["created_by"] = user.get("preferred_username")
    
    await drivers_collection.insert_one(driver_dict)
    return {"message": "Driver created successfully"}
    
@router.get("/", response_model=list[DriverResponse])
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
    result = await drivers_collection.delete_one({"_id": driver_id})
    if result.deleted_count == 1:
        return {"message": "Driver deleted successfully"}
    else:
        raise HTTPException(status_code=404, detail="Driver not found")

