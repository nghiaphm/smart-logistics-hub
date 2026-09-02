-- DROP bảng `profiles` mồ côi (migration 000003).
-- Domain thật dùng `user_profiles` (migration 000005); `profiles` không còn
-- được tham chiếu ở bất kỳ đâu (grep backend + frontend = 0 kết quả ngoài
-- chính migration tạo bảng). Xem WN-012.
DROP TABLE IF EXISTS profiles;
