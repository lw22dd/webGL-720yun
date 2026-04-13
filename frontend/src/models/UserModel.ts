export interface User {
  id?: string | number;
  username: string;
  email: string;
  password?: string;
  phone?: string;
  nickname?: string;
  role_id: number;
  role?: {
    id: number;
    name: string;
  };
  is_super_admin?: boolean;
  status: number;
  class_id?: number;
  created_at?: string;
  updated_at?: string;
}

export interface CreateUserRequest {
  id?: number;
  username: string;
  password: string;
  email?: string;
  phone?: string;
  nickname?: string;
  role_id: number;
  class_id?: number;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
  user: User;
}

export interface RegisterRequest {
  id?: number;
  username: string;
  password: string;
  email?: string;
  phone?: string;
  nickname?: string;
  role_id: number;
  class_id?: number;
  teacher_ids?: number[];
}

export interface UpdateUserRequest {
  email?: string;
  phone?: string;
  nickname?: string;
  avatar?: string;
  status?: number;
  class_id?: number;
  teacher_ids?: number[];
}

export interface UserListRequest {
  page?: number;
  page_size?: number;
  username?: string;
  email?: string;
  role_id?: number;
  status?: number;
  class_id?: number;
  keyword?: string;
}

export interface UserListResponse {
  total: number;
  users: User[];
  page: number;
  page_size: number;
}

export interface DeleteUserBatchRequest {
  ids: number[];
}

export interface DeleteUserBatchResponse {
  message: string;
  success_count: number;
  failed_count?: number;
  failed_ids?: number[];
}

export interface BatchRegisterResult {
  index: number;
  username: string;
  student_id: string;
  status: string;
  message: string;
}

export interface BatchRegisterResponse {
  success_count: number;
  failed_count: number;
  results: BatchRegisterResult[];
  errors: Record<string, unknown>[];
}

export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

export interface ResetPasswordRequest {
  username: string;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface RefreshTokenResponse {
  access_token: string;
  token_type: string;
}

export interface ClassInfo {
  id: number;
  name: string;
  description: string;
  teacher_id: number;
  created_at: string;
  updated_at: string;
}

export interface StudentExcelData {
  index: number;
  name: string;
  student_id: string;
  email: string;
  phone: string;
  class_name: string;
  class_id: number;
  teacher_ids: number[];
}
