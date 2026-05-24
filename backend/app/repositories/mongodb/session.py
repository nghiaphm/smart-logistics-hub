from app.repository.mongodb.database import db_instance

class MongoTransaction:
    """
    Sử dụng class này trong tầng Services để thực hiện các nghiệp vụ ACID.
    Ví dụ: Vừa trừ tồn kho, vừa cập nhật trạng thái đơn hàng.
    """
    def __init__(self):
        self.client = db_instance.client
        self.session = None

    async def __aenter__(self):
        self.session = await self.client.start_session()
        self.session.start_transaction()
        return self.session

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        if exc_type is not None:
            # Nếu có lỗi xảy ra (Exception) -> Hủy bỏ toàn bộ thay đổi
            await self.session.abort_transaction()
            print(f"Transaction Aborted due to: {exc_val}")
        else:
            # Nếu chạy mượt mà -> Lưu thay đổi vào đĩa cứng
            await self.session.commit_transaction()
        
        # Đóng phiên làm việc
        self.session.end_session()