import type { Result } from "@/models/result.model";
import type { HotspotListRequest, CreateHotspotRequest, UpdateHotspotRequest, HotspotDetailResponse, HotspotListResponse } from "@/models/hotspot.model";
import Axios from "@/utils/axios";

export default class HotspotApi {
    public static async getHotspotList(params?: HotspotListRequest): Promise<Result<HotspotListResponse>> {
        return await Axios.get('/resource/hotspots', { params });
    }

    public static async getHotspotDetail(id: number): Promise<Result<HotspotDetailResponse>> {
        return await Axios.get(`/resource/hotspots/${id}`);
    }

    public static async createHotspot(data: CreateHotspotRequest): Promise<Result<HotspotDetailResponse>> {
        return await Axios.post('/resource/hotspots', data);
    }

    public static async updateHotspot(id: number, data: UpdateHotspotRequest): Promise<Result<HotspotDetailResponse>> {
        return await Axios.put(`/resource/hotspots/${id}`, data);
    }

    public static async deleteHotspot(id: number): Promise<Result<{ message: string }>> {
        return await Axios.delete(`/resource/hotspots/${id}`);
    }
}