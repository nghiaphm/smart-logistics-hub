import requests
from jose import jwt, JWTError
from fastapi import Depends, HTTPException, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials

KEYCLOAK_URL = "http://localhost:8180"
REALM = "web-app-project"
JWKS_URL = f"{KEYCLOAK_URL}/realms/{REALM}/protocol/openid-connect/certs"

security = HTTPBearer()
_jwks_cache = None

def get_jwks():
    global _jwks_cache
    if _jwks_cache is None:
        try:
            _jwks_cache = requests.get(JWKS_URL, timeout=5).json()
        except (requests.RequestException, ValueError) as e:
            raise HTTPException(status_code=500, detail="Could not fetch JWKS from Keycloak")
    return _jwks_cache

def verify_token(credentials: HTTPAuthorizationCredentials = Depends(security)):
    token = credentials.credentials
    try:
        jwks = get_jwks()
        payload = jwt.decode(token, jwks, algorithms=["RS256"], 
                             issuer=f"{KEYCLOAK_URL}/realms/{REALM}", 
                             options={"verify_aud": False})
        return payload
    except JWTError:
        raise HTTPException(status_code=401, detail="Token invalid or expired")

# Viết thêm một Lớp Kiểm Tra Quyền (Role Checker)
class RequireRole:
    def __init__(self, allowed_roles: list[str]):
        self.allowed_roles = allowed_roles

    def __call__(self, user: dict = Depends(verify_token)):
        # Trích xuất danh sách huy hiệu từ Token của Keycloak
        realm_access = user.get("realm_access", {})
        user_roles = realm_access.get("roles", [])

        # Kiểm tra xem user có sở hữu ít nhất 1 quyền được phép không
        for role in self.allowed_roles:
            if role in user_roles:
                return user # Cho qua chốt kiểm soát
                
        # Nếu quét hết mà không thấy quyền, thẳng tay đuổi ra!
        raise HTTPException(
            status_code=403, 
            detail="403 Forbidden: Bạn không có quyền Admin để thực hiện hành động này!"
        )