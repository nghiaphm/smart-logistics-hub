export const SYSTEM_ADMIN_ROLE = "system_admin";

export type JwtUser = {
  sub?: string;
  preferred_username?: string;
  exp?: number;
  realm_access?: { roles?: string[] };
  resource_access?: Record<string, { roles?: string[] }>;
  [key: string]: unknown;
};

export function getUserId(user: JwtUser | null): string | undefined {
  return user?.sub ?? user?.preferred_username;
}

export function isSystemAdmin(user: JwtUser | null): boolean {
  if (!user) {
    return false;
  }
  const roles = new Set<string>();
  for (const role of user.realm_access?.roles ?? []) {
    roles.add(role);
  }
  for (const client of Object.values(user.resource_access ?? {})) {
    for (const role of client?.roles ?? []) {
      roles.add(role);
    }
  }
  return roles.has(SYSTEM_ADMIN_ROLE);
}
