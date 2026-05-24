from bson import ObjectId
from app.repository.mongodb.database import db_instance
from datetime import datetime, timezone
from typing import Optional, Dict, Any

class BaseRepository:
    def __init__(self, collection_name: str):
        self.collection_name = collection_name

    @property
    def collection(self):
        """Lấy collection hiện tại từ instance database"""
        return db_instance.db[self.collection_name]

    def _get_current_time(self) -> str:
        return datetime.now(timezone.utc).isoformat()

    async def get_by_id(self, document_id: str) -> Optional[Dict[str, Any]]:
        """Lấy một bản ghi theo _id"""
        document = await self.collection.find_one({"_id": ObjectId(document_id)})
        if document:
            document["_id"] = str(document["_id"])  # Ép kiểu ObjectId thành string cho Pydantic
        return document

    async def insert(self, document_data: dict, session=None) -> str:
        """Tạo mới một bản ghi"""
        document_data["created_at"] = self._get_current_time()
        document_data["updated_at"] = document_data["created_at"]
        
        result = await self.collection.insert_one(document_data, session=session)
        return str(result.inserted_id)

    async def update(self, document_id: str, update_data: dict, session=None) -> bool:
        """Cập nhật bản ghi có sẵn"""
        update_data["updated_at"] = self._get_current_time()
        
        result = await self.collection.update_one(
            {"_id": ObjectId(document_id)},
            {"$set": update_data},
            session=session
        )
        return result.modified_count > 0

    async def delete(self, document_id: str, session=None) -> bool:
        """Xóa một bản ghi"""
        result = await self.collection.delete_one({"_id": ObjectId(document_id)}, session=session)
        return result.deleted_count > 0