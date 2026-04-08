/**
 * 用户相关的数据模型
 * 对应后端 User 模型
 */
export interface User {
    /** 用户ID */
    id?: string
    /** 用户名（登录名） */
    username: string
    /** 邮箱 */
    email: string
    /** 密码（仅在创建/更新时使用，不会从接口返回） */
    password?: string
    /** 手机号 */
    phone?: string
    /** 昵称 */
    nickname?: string
    /** 角色ID：1-超级管理员，2-教师，3-学生 */
    role_id: number
    /** 角色详情 */
    role?: {
        id: number
        name: string
    }
    /** 状态：1-启用，0-禁用 */
    status: number
    /** 创建时间 */
    created_at?: string
    /** 更新时间 */
    updated_at?: string
}

/**
 * 角色ID常量
 */
export const RoleId = {
    SUPER_ADMIN: 1,
    TEACHER: 2,
    STUDENT: 3
} as const

export type RoleIdValue = typeof RoleId[keyof typeof RoleId]
