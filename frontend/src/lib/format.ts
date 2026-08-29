/**
 * Format datetime string (ISO) sang định dạng vi-VN: dd/MM/yyyy HH:mm.
 * Trả về "—" khi rỗng hoặc không parse được.
 */
export function formatDateTime(value?: string): string {
  if (!value) return "—"
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return "—"
  return new Intl.DateTimeFormat("vi-VN", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date)
}
