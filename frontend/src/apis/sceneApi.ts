import type { Result } from "@/models/Result";
import type { SceneListRequest, SceneDetailResponse, SceneListResponse, BatchImportResponse } from "@/models/SceneModel";
import Axios from "@/utils/axios";

export default class SceneApi {
    public static async getSceneList(params?: SceneListRequest): Promise<Result<SceneListResponse>> {
        return await Axios.get('/resource/scenes', { params });
    }

    public static async getSceneDetail(id: number): Promise<Result<SceneDetailResponse>> {
        return await Axios.get(`/resource/scenes/${id}`);
    }

    public static async createScene(formData: FormData): Promise<Result<SceneDetailResponse>> {
        return await Axios.post('/resource/scenes', formData);
    }

    public static async updateScene(id: number, formData: FormData): Promise<Result<SceneDetailResponse>> {
        return await Axios.put(`/resource/scenes/${id}`, formData);
    }

    public static async deleteScene(id: number): Promise<Result<{ message: string }>> {
        return await Axios.delete(`/resource/scenes/${id}`);
    }

    public static async batchImportScene(spaceId: number, file: File, onUploadProgress?: (progress: number) => void): Promise<Result<BatchImportResponse>> {
        const formData = new FormData();
        formData.append('space_id', spaceId.toString());
        formData.append('file', file);

        return await Axios.post('/resource/scenes/batch-import', formData, {
            onUploadProgress: (progressEvent) => {
                if (progressEvent.total && onUploadProgress) {
                    const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
                    onUploadProgress(percentCompleted);
                }
            }
        });
    }
}