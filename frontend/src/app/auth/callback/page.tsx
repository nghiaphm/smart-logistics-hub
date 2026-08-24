import { Suspense } from "react";
import CallbackClient from "./callback-client";

export default function Page() {
  return (
    <Suspense fallback={<p>Đang xử lý đăng nhập...</p>}>
      <CallbackClient />
    </Suspense>
  );
}
