// Design tokens khớp 1:1 với CSS variables trong `src/app/globals.css`.
// File này KHÔNG thay thế Tailwind/CSS variables — chỉ "xuất ra JS" cùng 1 nguồn,
// dùng khi cần giá trị màu trong JS (biểu đồ, canvas...), không dùng để style
// thay cho className thông thường.
//
// Nếu sửa giá trị trong globals.css → phải cập nhật đồng bộ ở đây.

export const RADIUS = "0.625rem"

/** Giá trị :root (light mode) trong globals.css */
export const COLORS = {
  background: "oklch(1 0 0)",
  foreground: "oklch(0.145 0 0)",
  card: "oklch(1 0 0)",
  cardForeground: "oklch(0.145 0 0)",
  popover: "oklch(1 0 0)",
  popoverForeground: "oklch(0.145 0 0)",
  primary: "oklch(0.205 0 0)",
  primaryForeground: "oklch(0.985 0 0)",
  secondary: "oklch(0.97 0 0)",
  secondaryForeground: "oklch(0.205 0 0)",
  muted: "oklch(0.97 0 0)",
  mutedForeground: "oklch(0.556 0 0)",
  accent: "oklch(0.97 0 0)",
  accentForeground: "oklch(0.205 0 0)",
  destructive: "oklch(0.577 0.245 27.325)",
  border: "oklch(0.922 0 0)",
  input: "oklch(0.922 0 0)",
  ring: "oklch(0.708 0 0)",
  chart: {
    1: "oklch(0.87 0 0)",
    2: "oklch(0.556 0 0)",
    3: "oklch(0.439 0 0)",
    4: "oklch(0.371 0 0)",
    5: "oklch(0.269 0 0)",
  },
  sidebar: {
    background: "oklch(0.985 0 0)",
    foreground: "oklch(0.145 0 0)",
    primary: "oklch(0.205 0 0)",
    primaryForeground: "oklch(0.985 0 0)",
    accent: "oklch(0.97 0 0)",
    accentForeground: "oklch(0.205 0 0)",
    border: "oklch(0.922 0 0)",
    ring: "oklch(0.708 0 0)",
  },
  admin: {
    sidebar: "oklch(0.145 0 0)",
    sidebarForeground: "oklch(0.985 0 0)",
    sidebarMuted: "oklch(0.708 0 0)",
    sidebarBorder: "oklch(1 0 0 / 12%)",
    accent: "oklch(0.666 0.179 58.318)",
    accentForeground: "oklch(0.205 0 0)",
  },
}

/** Giá trị `.dark` trong globals.css (override toàn bộ) */
export const DARK_COLORS = {
  background: "oklch(0.145 0 0)",
  foreground: "oklch(0.985 0 0)",
  card: "oklch(0.205 0 0)",
  cardForeground: "oklch(0.985 0 0)",
  popover: "oklch(0.205 0 0)",
  popoverForeground: "oklch(0.985 0 0)",
  primary: "oklch(0.922 0 0)",
  primaryForeground: "oklch(0.205 0 0)",
  secondary: "oklch(0.269 0 0)",
  secondaryForeground: "oklch(0.985 0 0)",
  muted: "oklch(0.269 0 0)",
  mutedForeground: "oklch(0.708 0 0)",
  accent: "oklch(0.269 0 0)",
  accentForeground: "oklch(0.985 0 0)",
  destructive: "oklch(0.704 0.191 22.216)",
  border: "oklch(1 0 0 / 10%)",
  input: "oklch(1 0 0 / 15%)",
  ring: "oklch(0.556 0 0)",
  chart: {
    1: "oklch(0.87 0 0)",
    2: "oklch(0.556 0 0)",
    3: "oklch(0.439 0 0)",
    4: "oklch(0.371 0 0)",
    5: "oklch(0.269 0 0)",
  },
  sidebar: {
    background: "oklch(0.205 0 0)",
    foreground: "oklch(0.985 0 0)",
    primary: "oklch(0.488 0.243 264.376)",
    primaryForeground: "oklch(0.985 0 0)",
    accent: "oklch(0.269 0 0)",
    accentForeground: "oklch(0.985 0 0)",
    border: "oklch(1 0 0 / 10%)",
    ring: "oklch(0.556 0 0)",
  },
  admin: {
    sidebar: "oklch(0.205 0 0)",
    sidebarForeground: "oklch(0.985 0 0)",
    sidebarMuted: "oklch(0.708 0 0)",
    sidebarBorder: "oklch(1 0 0 / 12%)",
    accent: "oklch(0.769 0.188 70.08)",
    accentForeground: "oklch(0.205 0 0)",
  },
}
