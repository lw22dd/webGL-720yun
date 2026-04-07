import type { Result } from "@/models/Result";
import type { User } from "@/models/UserModel";
import Axios from "@/utils/axios";

export default class UserApi {
    public static async login(username: string, password: string): Promise<Result<any | null>> {
        console.log('login', username, password);
        return await Axios.post('/auth/login', { username, password });
    }

    public static async register(user: User): Promise<Result<User | null>> {
        return await Axios.post('/user/register', user);
    }
    public static async logout(): Promise<Result<boolean>> {
        return await Axios.post('/user/logout');
    }
    public static async test(str?: string ): Promise<Result<string | null>> {
        return await Axios.get('/user/test', { params: { str } });
    }

    /**
     * 批量注册学生
     * @param file Excel文件
     * @param onUploadProgress 上传进度回调
     * @returns 注册结果
     */
    public static async batchRegister(file: File, onUploadProgress?: (progress: number) => void): Promise<Result<any>> {
        const formData = new FormData();
        formData.append('file', file);

        return await Axios.post('/user/admin/batch-register', formData, {
            onUploadProgress: (progressEvent) => {
                if (progressEvent.total && onUploadProgress) {
                    const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
                    onUploadProgress(percentCompleted);
                }
            }
        });
    }

    /**
     * 获取用户列表
     * @param params 查询参数
     * @returns 用户列表
     */
    public static async getUserList(params?: {
        keyword?: string;
        status?: string;
        page?: number;
        page_size?: number;
    }): Promise<Result<any>> {
        return await Axios.get('/user/admin/list', { params });
    }

    /**
     * 获取用户详情
     * @param userId 用户ID
     * @returns 用户详情
     */
    public static async getUserDetail(userId: string): Promise<Result<any>> {
        return await Axios.get(`/user/admin/${userId}`);
    }

    /**
     * 删除用户
     * @param userId 用户ID
     * @returns 删除结果
     */
    public static async deleteUser(userId: string): Promise<Result<boolean>> {
        return await Axios.delete(`/user/admin/${userId}`);
    }

    /**
     * 更新用户
     * @param userId 用户ID
     * @param userInfo 用户信息
     * @returns 更新结果
     */
    public static async updateUser(userId: string, userInfo: Partial<User>): Promise<Result<boolean>> {
        return await Axios.put(`/user/admin/${userId}`, userInfo);
    }

    /**
     * 新增用户
     * @param userInfo 用户信息
     * @returns 新增结果
     */
    public static async addUser(userInfo: User): Promise<Result<string | null>> {
        return await Axios.post('/user/admin/create', userInfo);
    }
}