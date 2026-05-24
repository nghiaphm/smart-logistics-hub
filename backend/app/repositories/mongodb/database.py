import urllib.parse
from motor.motor_asyncio import AsyncIOMotorClient
from app.core.config import settings
from app.helpers.logger_helper import logger # Giả sử bạn đã cấu hình logger này

class Database:
    client: AsyncIOMotorClient = None
    db = None

db_instance = Database()

async def connect_to_mongo():
    try:
        # 1. Mã hóa thông tin đăng nhập an toàn như cách của bạn
        encoded_user = urllib.parse.quote_plus(settings.MONGO_USER)
        encoded_pwd = urllib.parse.quote_plus(settings.MONGO_PWD)
        
        # 2. Lắp ráp chuỗi kết nối chuẩn
        conn_uri = f"mongodb://{encoded_user}:{encoded_pwd}@{settings.MONGO_HOST}:{settings.MONGO_PORT}/admin?authMechanism=DEFAULT&authSource={settings.MONGO_DB_NAME}&replicaSet=rs0&ssl=false"
        
        # 3. Kết nối bằng Motor (Bất đồng bộ)
        db_instance.client = AsyncIOMotorClient(conn_uri)
        db_instance.db = db_instance.client[settings.MONGO_DB_NAME]
        
        # 4. Kiểm tra và tự động tạo Collection (Chạy Async)
        required_collections = {"countries", "workspaces", "users"}
        
        # Dùng await để lấy danh sách collection mà không block server
        existing_collections = set(await db_instance.db.list_collection_names())
        collections_to_create = required_collections - existing_collections

        for collection_name in collections_to_create:
            # Lệnh tạo collection bất đồng bộ
            await db_instance.db.create_collection(collection_name)
            logger.info(f"Collection '{collection_name}' đã được tạo.")
            
        logger.info(f"Kết nối MongoDB ERP ({settings.MONGO_DB_NAME}) thành công!")
        
    except Exception as e:
        logger.error(f"Lỗi khi kết nối MongoDB: {e}")
        raise e

async def close_mongo_connection():
    if db_instance.client is not None:
        db_instance.client.close()
        logger.info("Đã ngắt kết nối MongoDB.")