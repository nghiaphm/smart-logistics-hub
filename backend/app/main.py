from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

# Import cả orders và drivers từ thư mục endpoints
from .api.endpoints import orders, drivers 

app = FastAPI(title="Smart Logistics Hub")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],
    allow_methods=["*"],
    allow_headers=["*"],
)

# Thêm route gốc để kiểm tra server và khắc phục lỗi 404
@app.get("/")
async def root():
    return {
        "message": "Welcome to Smart Logistics Hub API!",
        "docs": "Truy cập /docs để xem tài liệu API"
    }

# Gắn Router của Order vào hệ thống
app.include_router(orders.router, prefix="/api/v1/orders", tags=["Orders"])

# Gắn thêm Router của Driver vào hệ thống
app.include_router(drivers.router, prefix="/api/v1/drivers", tags=["Drivers"])