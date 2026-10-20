import axios from 'axios';
import type { AxiosInstance } from 'axios';

let api: AxiosInstance | null = null;

export const getApiClient = (): AxiosInstance => {
  if (!api) {
    api = axios.create({
      baseURL: '/',
      withCredentials: true,
    });
  }
  return api;
};

// Backend errors are always { error: string }; older/other endpoints may use
// { message: string }. Read both, never guess per call site.
export const getErrorMessage = (err: unknown, fallback: string): string => {
  const data = (err as { response?: { data?: { error?: string; message?: string } } })
    ?.response?.data;
  return data?.error || data?.message || fallback;
};

export default getApiClient;
