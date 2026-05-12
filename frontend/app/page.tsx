"use client";
import { signIn, signOut, useSession } from "next-auth/react";
import { useState } from "react";

export default function Home() {
  const { data: session } = useSession();
  const [orders, setOrders] = useState<any[]>([]);
  const [error, setError] = useState("");

  // Hàm gọi sang FastAPI với Thẻ An Ninh (Token)
  const fetchOrdersFromFastAPI = async () => {
    try {
      // Ép kiểu để lấy accessToken mà chúng ta đã cấu hình trong route.ts
      const token = (session as any)?.accessToken; 
      console.log("\n=== COPY TOKEN Ở DƯỚI ĐÂY ===");
      console.log(token);
      console.log("==============================\n");
      const res = await fetch("http://localhost:8000/api/v1/orders", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`, // Gắn thẻ an ninh vào đây!
        },
      });

      if (res.ok) {
        const data = await res.json();
        setOrders(data);
        setError("");
      } else {
        const errData = await res.json();
        setError(errData.detail || "Lỗi xác thực!");
      }
    } catch (err) {
      setError("Không thể kết nối tới Backend FastAPI");
    }
  };

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24 bg-slate-100">
      <div className="w-full max-w-2xl rounded-xl border bg-white p-8 shadow-xl">
        <h1 className="text-3xl font-bold mb-2 text-blue-900">Smart Logistics Hub</h1>
        <p className="text-gray-500 mb-8">Hệ thống Điều phối & Kho vận Thông minh</p>
        
        {session ? (
          <div className="space-y-6">
            <div className="p-4 bg-green-50 border border-green-200 rounded-lg">
              <p className="text-green-700 font-bold text-lg">
                🟢 Chào mừng, {session.user?.name || "Tài xế"}!
              </p>
              <p className="text-sm text-green-600">Email: {session.user?.email}</p>
            </div>

            {/* Nút gọi Backend */}
            <button 
              onClick={fetchOrdersFromFastAPI}
              className="px-6 py-2 bg-blue-600 text-white font-semibold rounded-lg hover:bg-blue-700 transition"
            >
              Tải danh sách đơn hàng từ AI Backend
            </button>

            {/* Vùng hiển thị lỗi (nếu có) */}
            {error && <p className="text-red-500 font-medium bg-red-50 p-3 rounded">{error}</p>}

            {/* Vùng hiển thị dữ liệu từ FastAPI */}
            {orders.length > 0 && (
              <div className="mt-4">
                <h3 className="font-bold text-gray-700 mb-3">📦 Danh sách chuyến hàng của bạn:</h3>
                <div className="grid gap-3">
                  {orders.map((order, index) => (
                    <div key={index} className="p-4 border rounded-lg flex justify-between items-center bg-gray-50">
                      <div>
                        {/* Dùng order_code (KH-9999) hoặc _id của MongoDB */}
                        <p className="font-bold text-blue-800">Mã: {order.order_code || order._id}</p>
                        
                        {/* Nếu chưa có tài xế thì hiện chữ "Chưa phân công" */}
                        <p className="text-sm text-gray-600">
                          Tài xế: <span className="font-medium text-gray-800">{order.assigned_driver_id || "Chưa phân công"}</span>
                        </p>
                      </div>
                      <span className="px-3 py-1 bg-yellow-100 text-yellow-800 rounded-full text-sm font-semibold">
                        {order.status}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            )}

            <hr className="my-6" />
            <button 
              onClick={() => signOut()}
              className="px-4 py-2 bg-gray-200 text-gray-700 font-semibold rounded-lg hover:bg-gray-300 transition"
            >
              Đăng xuất
            </button>
          </div>
        ) : (
          <div className="space-y-4">
            <div className="p-6 bg-blue-50 rounded-lg border border-blue-100">
              <p className="text-blue-800 mb-4">Vui lòng xuất trình thẻ an ninh (Đăng nhập) để vào trạm kiểm soát.</p>
              <button 
                className="w-full px-4 py-3 bg-slate-900 text-white font-bold rounded-lg hover:bg-slate-800 transition" 
                onClick={() => signIn("keycloak")}
              >
                Đăng nhập hệ thống (Keycloak)
              </button>
            </div>
          </div>
        )}
      </div>
    </main>
  );
}