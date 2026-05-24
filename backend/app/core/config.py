import os
from pydantic_settings import BaseSettings, SettingsConfigDict

# Đọc biến môi trường gốc của OS để biết đang chạy ở đâu (mặc định là development)
# Bạn có thể truyền biến này khi chạy lệnh: APP_ENV=production uvicorn main:app
APP_ENV = os.getenv("APP_ENV", "development")

class Settings(BaseSettings):
    PROJECT_NAME: str = "Smart Logistics Hub API"
    VERSION: str = "1.0.0"
    ENVIRONMENT: str = APP_ENV
    
    # --- CẤU HÌNH MONGODB ---
    MONGO_URI: str
    MONGO_DB_NAME: str
    
    # --- CẤU HÌNH GIAO TIẾP (URL) ---
    FRONTEND_URL: str
    AI_SERVICE_URL: str
    
    # --- CẤU HÌNH KEYCLOAK (BẢO MẬT) ---
    KEYCLOAK_SERVER_URL: str
    KEYCLOAK_REALM: str
    KEYCLOAK_CLIENT_ID: str
    
    # Chỉ định file .env linh hoạt dựa theo APP_ENV
    model_config = SettingsConfigDict(
        env_file=f".env.{APP_ENV}", 
        env_file_encoding="utf-8",
        extra="ignore"  # Bỏ qua các biến thừa trong file .env nếu không khai báo ở trên
    )

# Khởi tạo instance duy nhất để import ở các file khác
settings = Settings()