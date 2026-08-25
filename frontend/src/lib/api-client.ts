import { ApiError } from "@/types/api";
import { ensureFreshAccessToken } from "@/lib/auth";

export { ApiError };

type BackendErrorBody = {
  error?: {
    code?: string | number;
    message?: string;
  };
};

function isBackendErrorBody(value: unknown): value is BackendErrorBody {
  if (typeof value !== "object" || value === null) {
    return false;
  }
  const error = (value as { error?: unknown }).error;
  return typeof error === "object" && error !== null;
}

async function parseError(response: Response): Promise<ApiError> {
  const status = response.status;
  let code: string | number = status;
  let message = `Request failed with status ${status}`;

  try {
    const body: unknown = await response.json();
    if (isBackendErrorBody(body)) {
      code = body.error?.code ?? status;
      message = body.error?.message ?? message;
    }
  } catch {
  }

  return new ApiError(status, code, message);
}

export async function apiClient<T>(path: string, options: RequestInit = {}): Promise<T> {
  const baseUrl = process.env.NEXT_PUBLIC_API_URL;
  if (!baseUrl) {
    throw new Error("NEXT_PUBLIC_API_URL is not set");
  }

  const headers = new Headers(options.headers);
  const token = await ensureFreshAccessToken();
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  const response = await fetch(`${baseUrl}${path}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    throw await parseError(response);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return (await response.json()) as T;
}
