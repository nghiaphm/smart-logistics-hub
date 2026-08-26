const ACCESS_TOKEN_KEY = "access_token";
const REFRESH_TOKEN_KEY = "refresh_token";
const OAUTH_STATE_KEY = "oauth_state";
const AUTH_COOKIE_NAME = "slh_access_token";

export type KeycloakTokens = {
  access_token: string;
  refresh_token?: string;
  id_token?: string;
  expires_in?: number;
};

function getKeycloakConfig() {
  const url = process.env.NEXT_PUBLIC_KEYCLOAK_URL;
  const realm = process.env.NEXT_PUBLIC_KEYCLOAK_REALM;
  const clientId = process.env.NEXT_PUBLIC_KEYCLOAK_CLIENT_ID;
  if (!url || !realm || !clientId) {
    throw new Error(
      "NEXT_PUBLIC_KEYCLOAK_URL, NEXT_PUBLIC_KEYCLOAK_REALM and NEXT_PUBLIC_KEYCLOAK_CLIENT_ID must be set"
    );
  }
  return { url, realm, clientId };
}

function getTokenEndpoint(): string {
  const { url, realm } = getKeycloakConfig();
  return `${url}/realms/${realm}/protocol/openid-connect/token`;
}

export function createAuthorizationUrl(redirectUri: string, state: string, isRegister = false): string {
  const { url, realm, clientId } = getKeycloakConfig();
  const params = new URLSearchParams({
    client_id: clientId,
    redirect_uri: redirectUri,
    response_type: "code",
    scope: "openid",
    state,
  });
  const endpoint = isRegister ? "reg" : "auth";
  return `${url}/realms/${realm}/protocol/openid-connect/${endpoint}?${params.toString()}`;
}

export function saveOAuthState(state: string): void {
  if (typeof window === "undefined") {
    return;
  }
  window.sessionStorage.setItem(OAUTH_STATE_KEY, state);
}

export function consumeOAuthState(): string | null {
  if (typeof window === "undefined") {
    return null;
  }
  const state = window.sessionStorage.getItem(OAUTH_STATE_KEY);
  window.sessionStorage.removeItem(OAUTH_STATE_KEY);
  return state;
}

export function getAccessToken(): string | null {
  if (typeof window === "undefined") {
    return null;
  }
  return window.localStorage.getItem(ACCESS_TOKEN_KEY);
}

export function getRefreshToken(): string | null {
  if (typeof window === "undefined") {
    return null;
  }
  return window.localStorage.getItem(REFRESH_TOKEN_KEY);
}

export function setTokens(tokens: KeycloakTokens): void {
  if (typeof window === "undefined") {
    return;
  }
  window.localStorage.setItem(ACCESS_TOKEN_KEY, tokens.access_token);
  if (tokens.refresh_token) {
    window.localStorage.setItem(REFRESH_TOKEN_KEY, tokens.refresh_token);
  }
  setAccessTokenCookie(tokens.access_token);
  notifyTokenChanged();
}

export function clearTokens(): void {
  if (typeof window === "undefined") {
    return;
  }
  window.localStorage.removeItem(ACCESS_TOKEN_KEY);
  window.localStorage.removeItem(REFRESH_TOKEN_KEY);
  clearAccessTokenCookie();
  notifyTokenChanged();
}

type TokenChangeListener = () => void;

const tokenChangeListeners = new Set<TokenChangeListener>();

export function subscribeTokenChanges(listener: TokenChangeListener): () => void {
  tokenChangeListeners.add(listener);
  return () => {
    tokenChangeListeners.delete(listener);
  };
}

function notifyTokenChanged() {
  tokenChangeListeners.forEach((listener) => listener());
}

function setAccessTokenCookie(token: string): void {
  document.cookie = `${AUTH_COOKIE_NAME}=${encodeURIComponent(token)}; path=/; samesite=lax`;
}

function clearAccessTokenCookie(): void {
  document.cookie = `${AUTH_COOKIE_NAME}=; path=/; max-age=0; samesite=lax`;
}

export function decodeJwt<T>(token: string): T | null {
  try {
    const [, payload] = token.split(".");
    if (!payload) {
      return null;
    }
    const base64 = payload.replace(/-/g, "+").replace(/_/g, "/");
    const json = decodeURIComponent(
      atob(base64)
        .split("")
        .map((char) => "%" + char.charCodeAt(0).toString(16).padStart(2, "0"))
        .join("")
    );
    return JSON.parse(json) as T;
  } catch {
    return null;
  }
}

export function getAccessTokenExpiry(): number | null {
  const token = getAccessToken();
  if (!token) {
    return null;
  }
  const payload = decodeJwt<{ exp?: number }>(token);
  if (!payload || typeof payload.exp !== "number") {
    return null;
  }
  return payload.exp * 1000;
}

export function isAccessTokenExpired(): boolean {
  const expiry = getAccessTokenExpiry();
  return expiry === null || Date.now() >= expiry;
}

export function isAccessTokenExpiringSoon(bufferMs = 60_000): boolean {
  const expiry = getAccessTokenExpiry();
  if (expiry === null) {
    return true;
  }
  return Date.now() + bufferMs >= expiry;
}

async function postToTokenEndpoint(body: URLSearchParams): Promise<KeycloakTokens> {
  const response = await fetch(getTokenEndpoint(), {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body,
  });
  if (!response.ok) {
    throw new Error("Keycloak token request failed");
  }
  return (await response.json()) as KeycloakTokens;
}

export async function refreshAccessToken(): Promise<string | null> {
  const refreshToken = getRefreshToken();
  if (!refreshToken) {
    return null;
  }

  const { clientId } = getKeycloakConfig();
  try {
    const tokens = await postToTokenEndpoint(
      new URLSearchParams({
        client_id: clientId,
        grant_type: "refresh_token",
        refresh_token: refreshToken,
      })
    );
    setTokens(tokens);
    return tokens.access_token ?? null;
  } catch {
    clearTokens();
    return null;
  }
}

export async function exchangeCodeForTokens(code: string, redirectUri: string): Promise<KeycloakTokens> {
  const { clientId } = getKeycloakConfig();
  return postToTokenEndpoint(
    new URLSearchParams({
      client_id: clientId,
      grant_type: "authorization_code",
      code,
      redirect_uri: redirectUri,
    })
  );
}

export async function ensureFreshAccessToken(): Promise<string | null> {
  const token = getAccessToken();
  if (!token) {
    return null;
  }
  if (!isAccessTokenExpiringSoon()) {
    return token;
  }
  return refreshAccessToken();
}
