import type { Result } from "@/models/Result";
import type { SpaceListRequest, SpaceDetailResponse, SpaceListResponse } from "@/models/SpaceModel";
import Axios from "@/utils/axios";

export default class SpaceApi {
    public static async getSpaceList(params?: SpaceListRequest): Promise<Result<SpaceListResponse>> {
        return await Axios.get('/resource/spaces', { params });
    }

    public static async getSpaceDetail(slug: string): Promise<Result<SpaceDetailResponse>> {
        return await Axios.get(`/resource/spaces/${slug}`);
    }

    public static async createSpace(formData: FormData): Promise<Result<SpaceDetailResponse>> {
        return await Axios.post('/resource/spaces', formData);
    }

    public static async updateSpace(id: number, formData: FormData): Promise<Result<SpaceDetailResponse>> {
        return await Axios.put(`/resource/spaces/${id}`, formData);
    }

    public static async deleteSpace(id: number): Promise<Result<{ message: string }>> {
        return await Axios.delete(`/resource/spaces/${id}`);
    }
}