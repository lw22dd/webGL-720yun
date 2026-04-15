import type {
  Result,
} from "@/models/result.model";
import type {
  InitUploadRequest,
  InitUploadResponse,
  ChunkUploadResponse,
  CompleteUploadRequest,
  CompleteUploadResponse,
  UploadStatusResponse,
  FileInfo,
} from "@/models/upload.model";
import Axios from "@/utils/axios";

export default class UploadApi {
    public static async initUpload(data: InitUploadRequest): Promise<Result<InitUploadResponse>> {
        return await Axios.post('/upload/init', data);
    }

    public static async uploadChunk(
        uploadId: string,
        chunkIndex: number,
        chunk: Blob,
        onUploadProgress?: (progress: number) => void
    ): Promise<Result<ChunkUploadResponse>> {
        const formData = new FormData();
        formData.append('upload_id', uploadId);
        formData.append('chunk_index', chunkIndex.toString());
        formData.append('chunk', chunk);

        return await Axios.post('/upload/chunk', formData, {
            onUploadProgress: (progressEvent) => {
                if (progressEvent.total && onUploadProgress) {
                    const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
                    onUploadProgress(percentCompleted);
                }
            },
        });
    }

    public static async completeUpload(data: CompleteUploadRequest): Promise<Result<CompleteUploadResponse>> {
        return await Axios.post('/upload/complete', data);
    }

    public static async getUploadStatus(uploadId: string): Promise<Result<UploadStatusResponse>> {
        return await Axios.get(`/upload/status/${uploadId}`);
    }

    public static async cancelUpload(uploadId: string): Promise<Result<{ message: string }>> {
        return await Axios.delete(`/upload/${uploadId}`);
    }

    public static async getFileInfo(fileId: string): Promise<Result<FileInfo>> {
        return await Axios.get(`/upload/file/${fileId}`);
    }
}
