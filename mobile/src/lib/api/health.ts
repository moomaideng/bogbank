import { apiFetch } from '@/lib/api/client';

export type HealthStatus = {
  status: string;
  error?: string;
};

export function getLiveness(): Promise<HealthStatus> {
  return apiFetch<HealthStatus>('/livez');
}
