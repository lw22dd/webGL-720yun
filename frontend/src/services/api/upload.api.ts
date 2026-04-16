import type {
  Result,
} from "@/models/result.model";
import type {
  InitUploadRequest,
  InitUploadResponse,
  ChunkUploadResponse,
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
        chunkMd5: string,
        signal?: AbortSignal
    ): Promise<Result<ChunkUploadResponse>> {
        const formData = new FormData();
        formData.append('upload_id', uploadId);
        formData.append('chunk_index', chunkIndex.toString());
        formData.append('chunk_hash', chunkMd5);
        formData.append('chunk_data', chunk);

        const controller = new AbortController();
        if (signal) {
            signal.addEventListener('abort', () => controller.abort());
        }

        try {
            return await Axios.post('/upload/chunk', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data'
                },
                timeout: 60000,
                signal: controller.signal as any
            });
        } catch (error: any) {
            if (error.name === 'CanceledError' || error.name === 'AbortError') {
                throw new Error('UPLOAD_CANCELLED');
            }
            throw error;
        }
    }

    public static async completeUpload(uploadId: string, fileHash: string): Promise<Result<CompleteUploadResponse>> {
        return await Axios.post('/upload/complete', {
            upload_id: uploadId,
            file_hash: fileHash
        });
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
