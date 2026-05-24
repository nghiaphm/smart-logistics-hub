# app/api/deps.py
from fastapi import Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer
from app.core.security import RequireRole, decode_token
from ...db.database import drivers_collection


reusable_oauth2 = OAuth2PasswordBearer(tokenUrl="http://localhost:8080/realms/web-app-project/protocol/openid-connect/token")

async def get_current_user(token: str = Depends(reusable_oauth2)):
    try:
        payload = decode_token(token)
        return payload
    except Exception:
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Token không hợp lệ")

    # 1. Lấy username từ Token (đã qua kiểm tra của Keycloak)
    username = user.get("preferred_username")
    
    # 2. Truy vấn vào MongoDB để tìm tài xế có mã tương ứng
    db_driver = await drivers_collection.find_one({"driver_code": username})
    
    # 3. Nếu không tìm thấy trong DB, chặn lại ngay dù Token hợp lệ
    if not db_driver:
        raise HTTPException(
            status_code=404, 
            detail=f"Tài xế {username} có Token hợp lệ nhưng chưa được kích hoạt trong hệ thống dữ liệu nội bộ!"
        )
    
    # 4. Trả về kết hợp cả dữ liệu Token và dữ liệu DB
    return {
        "identity": user,      # Dữ liệu từ Keycloak
        "profile": db_driver   # Dữ liệu từ MongoDB
    }

# Chốt chặn dành riêng cho Driver
def check_driver_role(current_user: dict = Depends(get_current_user)):
    roles = current_user.get("roles", [])
    if "driver" not in roles:
        raise HTTPException(status_code=403, detail="Bạn không có quyền tài xế")
    return current_user