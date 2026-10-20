import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { booksApi, guestApi } from '../lib/api';

export const useBooks = (guest = false) => {
  return useQuery({
    queryKey: guest ? ['guestBooks'] : ['books'],
    queryFn: guest ? guestApi.getBooks : booksApi.getBooks,
  });
};

export const useCreateBook = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { title: string; bg_color?: string }) => booksApi.createBook(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['books'] });
    },
  });
};

export const useUpdateBook = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ bookId, data }: { bookId: string; data: any }) =>
      booksApi.updateBook(bookId, data),
    onSuccess: (_data, { bookId }) => {
      queryClient.invalidateQueries({ queryKey: ['books'] });
      queryClient.invalidateQueries({ queryKey: ['book', bookId] });
    },
  });
};

export const useDeleteBook = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (bookId: string) => booksApi.deleteBook(bookId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['books'] });
    },
  });
};

export const useUploadPhotos = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ bookId, files, onUploadProgress }: { bookId: string; files: File[]; onUploadProgress?: (progress: number) => void }) =>
      booksApi.uploadPhotos(bookId, files, onUploadProgress),
    onSuccess: (_data, { bookId }) => {
      queryClient.invalidateQueries({ queryKey: ['books'] });
      queryClient.invalidateQueries({ queryKey: ['bookPhotos', bookId] });
      queryClient.invalidateQueries({ queryKey: ['book', bookId] });
    },
  });
};

export const useDeletePhoto = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ photoId }: { photoId: string; bookId: string }) =>
      booksApi.deletePhoto(photoId),
    onSuccess: (_data, { bookId }) => {
      queryClient.invalidateQueries({ queryKey: ['books'] });
      queryClient.invalidateQueries({ queryKey: ['bookPhotos', bookId] });
      queryClient.invalidateQueries({ queryKey: ['book', bookId] });
    },
  });
};

export const useBookDetail = (bookId: string | undefined, guest = false) => {
  return useQuery({
    queryKey: guest ? ['guestBook', bookId] : ['book', bookId],
    queryFn: () =>
      bookId
        ? guest
          ? guestApi.getBookDetail(bookId)
          : booksApi.getBookDetail(bookId)
        : Promise.reject(new Error('No bookId')),
    enabled: !!bookId,
  });
};

export const useBookPhotos = (bookId: string | undefined, guest = false) => {
  return useQuery({
    queryKey: guest ? ['guestBookPhotos', bookId] : ['bookPhotos', bookId],
    queryFn: () =>
      bookId
        ? guest
          ? guestApi.getBookPages(bookId)
          : booksApi.getBookPages(bookId)
        : Promise.reject(new Error('No bookId')),
    enabled: !!bookId,
  });
};
