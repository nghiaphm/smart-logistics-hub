from motor.motor_asyncio import AsyncIOMotorClient

# sau này nên đưa vào file .env
MONGO_URL = "mongodb://admin:password@localhost:27017/"
client = AsyncIOMotorClient(MONGO_URL)

# Trỏ vào database có tên là smart_logistics
db = client.smart_logistics

# Khai báo các collection
orders_collection = db.orders
drivers_collection = db.drivers