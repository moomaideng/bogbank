const BASE_URL = process.env.EXPO_PUBLIC_API_URL;

if (!BASE_URL) {
  throw new Error('EXPO_PUBLIC_API_URL is not set. Copy .env.example to .env.');
}

export class ApiError extends Error {
  constructor(
    public readonly status: number,
    message: string,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

export async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${BASE_URL}${path}`, init);
  if (!response.ok) {
    throw new ApiError(response.status, `${init?.method ?? 'GET'} ${path} failed`);
  }
  return response.json() as Promise<T>;
}
