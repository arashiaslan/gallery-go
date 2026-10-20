import { getApiClient } from './client';
import type { User, LoginRequest, RegisterRequest, AuthResponse } from '../types/auth';
import type { Book, Photo } from '../types';

const api = getApiClient();

// Backend wraps every success payload as { data: ... } and every error as
// { error: ... }. Unwrapping happens here only, never in components.
interface Envelope<T> {
  data: T;
}
interface PagedPhotos {
  photos: Photo[];
  count: number;
}
interface UploadResult {
  photos: Photo[];
  errors?: string[];
}

const unwrap = <T>(payload: Envelope<T>): T => payload.data;

export const authApi = {
  login: async (data: LoginRequest): Promise<User> => {
    const response = await api.post<AuthResponse>('/api/auth/login', data);
    return unwrap(response.data);
  },

  register: async (data: RegisterRequest): Promise<User> => {
    const response = await api.post<AuthResponse>('/api/auth/register', data);
    return unwrap(response.data);
  },

  logout: async (): Promise<void> => {
    await api.post('/api/auth/logout');
  },

  getMe: async (): Promise<User> => {
    const response = await api.get<AuthResponse>('/api/me');
    return unwrap(response.data);
  },
};

export const booksApi = {
  getBooks: async (): Promise<Book[]> => {
    const response = await api.get<Envelope<Book[]>>('/api/books');
    return unwrap(response.data) ?? [];
  },

  createBook: async (data: { title: string; bg_color?: string }): Promise<Book> => {
    const response = await api.post<Envelope<Book>>('/api/books', data);
    return unwrap(response.data);
  },

  updateBook: async (bookId: string, data: Partial<Book>): Promise<Book> => {
    const response = await api.patch<Envelope<Book>>(`/api/books/${bookId}`, data);
    return unwrap(response.data);
  },

  deleteBook: async (bookId: string): Promise<void> => {
    await api.delete(`/api/books/${bookId}`);
  },

  uploadPhotos: async (
    bookId: string,
    files: File[],
    onUploadProgress?: (progress: number) => void
  ): Promise<UploadResult> => {
    const formData = new FormData();
    files.forEach((file) => {
      formData.append('photos', file);
    });

    const response = await api.post<Envelope<UploadResult>>(`/api/books/${bookId}/photos`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        if (progressEvent.total && onUploadProgress) {
          const progress = (progressEvent.loaded / progressEvent.total) * 100;
          onUploadProgress(Math.round(progress));
        }
      },
    });
    return unwrap(response.data);
  },

  deletePhoto: async (photoId: string): Promise<void> => {
    await api.delete(`/api/photos/${photoId}`);
  },

  getBookDetail: async (bookId: string): Promise<Book> => {
    const response = await api.get<Envelope<Book>>(`/api/books/${bookId}`);
    return unwrap(response.data);
  },

  getBookPages: async (bookId: string): Promise<PagedPhotos> => {
    const response = await api.get<Envelope<PagedPhotos>>(`/api/books/${bookId}/pages`);
    return unwrap(response.data);
  },

  getPhotoUrl: (bookId: string, photoId: string, size: 'thumb' | 'medium' | 'original' = 'medium'): string => {
    return `/media/${bookId}/${photoId}/${size}`;
  },
};

// Public, read-only endpoints for the guest view (no session required).
export const guestApi = {
  getBooks: async (): Promise<Book[]> => {
    const response = await api.get<Envelope<Book[]>>('/api/guest/books');
    return unwrap(response.data) ?? [];
  },

  getBookDetail: async (bookId: string): Promise<Book> => {
    const response = await api.get<Envelope<Book>>(`/api/guest/books/${bookId}`);
    return unwrap(response.data);
  },

  getBookPages: async (bookId: string): Promise<PagedPhotos> => {
    const response = await api.get<Envelope<PagedPhotos>>(`/api/guest/books/${bookId}/pages`);
    return unwrap(response.data);
  },
};
