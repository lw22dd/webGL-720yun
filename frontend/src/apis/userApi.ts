import type {
  User,
  CreateUserRequest,
  LoginResponse,
  RegisterRequest,
  UpdateUserRequest,
  UserListRequest,
  UserListResponse,
  DeleteUserBatchRequest,
  DeleteUserBatchResponse,
  BatchRegisterResponse,
  ChangePasswordRequest,
  RefreshTokenRequest,
  RefreshTokenResponse,
} from "@/models/UserModel";
import type { Result } from "@/models/Result";
import Axios from "@/utils/axios";

export default class UserApi {
    public static async login(username: string, password: string): Promise<Result<LoginResponse | null>> {
        return await Axios.post('/auth/login', { username, password });
    }

    public static async register(user: RegisterRequest): Promise<Result<User | null>> {
        return await Axios.post('/user/register', user);
    }

    public static async logout(): Promise<Result<{ message: string }>> {
        return await Axios.post('/user/logout');
    }

    public static async refreshToken(refreshToken: string): Promise<Result<RefreshTokenResponse>> {
        return await Axios.post('/auth/refresh', { refresh_token: refreshToken } as RefreshTokenRequest);
    }

    public static async getProfile(): Promise<Result<User | null>> {
        return await Axios.get('/user/profile');
    }

    public static async updateProfile(userInfo: UpdateUserRequest): Promise<Result<User | null>> {
        return await Axios.put('/user/profile', userInfo);
    }

    public static async changePassword(data: ChangePasswordRequest): Promise<Result<{ message: string }>> {
        return await Axios.put('/user/change-password', data);
    }

    public static async test(str?: string): Promise<Result<string | null>> {
        return await Axios.get('/user/test', { params: { str } });
    }

    public static async batchRegister(
        file: File,
        onUploadProgress?: (progress: number) => void
    ): Promise<Result<BatchRegisterResponse>> {
        const formData = new FormData();
        formData.append('file', file);

        return await Axios.post('/user/admin/batch-register', formData, {
            onUploadProgress: (progressEvent) => {
                if (progressEvent.total && onUploadProgress) {
                    const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
                    onUploadProgress(percentCompleted);
                }
            },
        });
    }

    public static async getUserList(params?: UserListRequest): Promise<Result<UserListResponse>> {
        return await Axios.get('/user/admin/list', { params });
    }

    public static async getUserDetail(userId: number | string): Promise<Result<User | null>> {
        return await Axios.get(`/user/admin/${userId}`);
    }

    public static async deleteUser(userId: number | string): Promise<Result<{ message: string }>> {
        return await Axios.delete(`/user/admin/${userId}`);
    }

    public static async deleteUserBatch(request: DeleteUserBatchRequest): Promise<Result<DeleteUserBatchResponse>> {
        return await Axios.request({
            method: 'DELETE',
            url: '/user/admin/batch',
            data: request,
        });
    }

    public static async updateUser(userId: number | string, userInfo: UpdateUserRequest): Promise<Result<User | null>> {
        return await Axios.put(`/user/admin/${userId}`, userInfo);
    }

    public static async addUser(userInfo: CreateUserRequest): Promise<Result<User | null>> {
        return await Axios.post('/user/admin/create', userInfo);
    }
}
