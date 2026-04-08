/**
 * 通用 API 响应结果封装
 * 所有后端 API 接口的响应数据都遵循此格式
 */
export type Result<T> = {
    /** 状态码，200 表示成功 */
    code: number
    /** 响应消息，成功时通常为空字符串 */
    msg: string
    /** 响应数据，类型由调用方指定 */
    data?: T
}
