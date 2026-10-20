export interface User {
  id: string;
  email: string;
  role: 'admin' | 'member';
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  passwordConfirm: string;
  invitationToken: string;
}

export interface AuthResponse {
  data: User;
}

export interface ErrorResponse {
  message: string;
}
