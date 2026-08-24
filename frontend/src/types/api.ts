export class ApiError extends Error {
  readonly status: number;
  readonly code: string | number;

  constructor(status: number, code: string | number, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.code = code;
  }
}

export type ApiResponse<T> = {
  items: T[];
  total: number;
  skip: number;
  limit: number;
};
