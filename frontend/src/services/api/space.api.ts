import type { Result } from "@/models/result.model";
import type { SpaceListRequest, SpaceDetailResponse, SpaceListResponse, CreateSpaceRequest, UpdateSpaceRequest } from "@/models/space.model";
import type { GraphDataResponse } from "@/models/scene.model";
import Axios from "@/utils/axios";

export default class SpaceApi {
    public static async getSpaceList(params?: SpaceListRequest): Promise<Result<SpaceListResponse>> {
        return await Axios.get('/resource/spaces', { params });
    }

    public static async getSpaceDetail(id: number): Promise<Result<SpaceDetailResponse>> {
        return await Axios.get(`/resource/spaces/${id}`);
    }

    public static async createSpace(data: CreateSpaceRequest, coverFile?: File): Promise<Result<SpaceDetailResponse>> {
        if (coverFile) {
            const formData = new FormData();
            Object.entries(data).forEach(([key, value]) => {
                if (value !== undefined && value !== null) {
                    formData.append(key, String(value));
                }
            });
            formData.append('cover', coverFile);
            return await Axios.post('/resource/spaces', formData);
        }
        return await Axios.post('/resource/spaces', data);
    }

    public static async updateSpace(id: number, data: UpdateSpaceRequest, coverFile?: File): Promise<Result<SpaceDetailResponse>> {
        if (coverFile) {
            const formData = new FormData();
            Object.entries(data).forEach(([key, value]) => {
                if (value !== undefined && value !== null) {
                    formData.append(key, String(value));
                }
            });
            formData.append('cover', coverFile);
            return await Axios.put(`/resource/spaces/${id}`, formData);
        }
        return await Axios.put(`/resource/spaces/${id}`, data);
    }

    public static async deleteSpace(id: number): Promise<Result<{ message: string }>> {
        return await Axios.delete(`/resource/spaces/${id}`);
    }

    public static async deleteSpaceBatch(ids: number[]): Promise<Result<{ message: string }>> {
        return await Axios.delete('/resource/spaces/batch', { data: { ids } });
    }

    public static async getSpaceGraph(id: number): Promise<Result<GraphDataResponse>> {
        return await Axios.get(`/resource/spaces/${id}/graph`);
    }
}
