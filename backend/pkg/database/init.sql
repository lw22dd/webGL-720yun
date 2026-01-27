-- WebGL-720yun 数据库初始化脚本
-- 使用 CREATE TABLE IF NOT EXISTS 和 INSERT IGNORE 确保幂等执行

-- 创建角色表
CREATE TABLE IF NOT EXISTS `sys_roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_roles_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建权限表
CREATE TABLE IF NOT EXISTS `sys_permissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `resource` varchar(50) DEFAULT NULL,
  `action` varchar(50) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_permissions_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建角色-权限关联表
CREATE TABLE IF NOT EXISTS `sys_role_permissions` (
  `role_id` bigint unsigned NOT NULL,
  `permission_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`role_id`,`permission_id`),
  KEY `idx_role_permissions_permission_id` (`permission_id`),
  CONSTRAINT `fk_role_permissions_permission_id` FOREIGN KEY (`permission_id`) REFERENCES `sys_permissions` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_role_permissions_role_id` FOREIGN KEY (`role_id`) REFERENCES `sys_roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建用户表
CREATE TABLE IF NOT EXISTS `sys_users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(100) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `nickname` varchar(50) DEFAULT NULL,
  `avatar` varchar(255) DEFAULT NULL,
  `role_id` bigint unsigned NOT NULL,
  `is_super_admin` tinyint(1) DEFAULT 0,
  `status` int DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `student_id` varchar(20) DEFAULT NULL,
  `class_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_username` (`username`),
  UNIQUE KEY `idx_users_email` (`email`),

  UNIQUE KEY `idx_users_student_id` (`student_id`),
  KEY `idx_users_role_id` (`role_id`),
  KEY `idx_users_status` (`status`),
  KEY `idx_users_deleted_at` (`deleted_at`),
  KEY `idx_users_class_id` (`class_id`),
  CONSTRAINT `fk_users_role_id` FOREIGN KEY (`role_id`) REFERENCES `sys_roles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建班级表
CREATE TABLE IF NOT EXISTS `sys_classes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `teacher_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_classes_name` (`name`),
  KEY `idx_classes_teacher_id` (`teacher_id`),
  CONSTRAINT `fk_classes_teacher_id` FOREIGN KEY (`teacher_id`) REFERENCES `sys_users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建学生-教师关联表
CREATE TABLE IF NOT EXISTS `sys_student_teachers` (
  `student_id` bigint unsigned NOT NULL,
  `teacher_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`student_id`,`teacher_id`),
  KEY `idx_student_teachers_teacher_id` (`teacher_id`),
  CONSTRAINT `fk_student_teachers_student_id` FOREIGN KEY (`student_id`) REFERENCES `sys_users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_student_teachers_teacher_id` FOREIGN KEY (`teacher_id`) REFERENCES `sys_users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建用户会话表
CREATE TABLE IF NOT EXISTS `sys_user_sessions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `access_token` text NOT NULL,
  `refresh_token` text NOT NULL,
  `ip` varchar(50) DEFAULT NULL,
  `user_agent` text DEFAULT NULL,
  `expires_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_sessions_user_id` (`user_id`),
  KEY `idx_user_sessions_expires_at` (`expires_at`),
  CONSTRAINT `fk_user_sessions_user_id` FOREIGN KEY (`user_id`) REFERENCES `sys_users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入基础角色数据
INSERT IGNORE INTO `sys_roles` (`id`, `name`, `description`, `created_at`, `updated_at`) VALUES
(1, 'admin', '管理员', NOW(), NOW()),
(2, 'teacher', '教师', NOW(), NOW()),
(3, 'student', '学生', NOW(), NOW());

-- 插入基础权限数据
INSERT IGNORE INTO `sys_permissions` (`id`, `name`, `description`, `resource`, `action`, `created_at`, `updated_at`) VALUES
(1, 'user:create', '创建用户', 'user', 'create', NOW(), NOW()),
(2, 'user:read', '读取用户', 'user', 'read', NOW(), NOW()),
(3, 'user:update', '更新用户', 'user', 'update', NOW(), NOW()),
(4, 'user:delete', '删除用户', 'user', 'delete', NOW(), NOW()),
(5, 'class:create', '创建班级', 'class', 'create', NOW(), NOW()),
(6, 'class:read', '读取班级', 'class', 'read', NOW(), NOW()),
(7, 'class:update', '更新班级', 'class', 'update', NOW(), NOW()),
(8, 'class:delete', '删除班级', 'class', 'delete', NOW(), NOW());

-- 关联角色与权限
INSERT IGNORE INTO `sys_role_permissions` (`role_id`, `permission_id`) VALUES
(1, 1),
(1, 2),
(1, 3),
(1, 4),
(1, 5),
(1, 6),
(1, 7),
(1, 8),
(2, 2),
(2, 3),
(2, 5),
(2, 6),
(2, 7),
(2, 8),
(3, 2);

-- 插入或更新默认管理员用户 admin
-- 密码：admin123（已加密）
INSERT INTO `sys_users` (`id`, `username`, `password`, `email`, `phone`, `nickname`, `role_id`, `is_super_admin`, `status`, `created_at`, `updated_at`) 
VALUES (1, 'admin', '$2a$10$0QfJWCtOYeMEPr4JBfmLK.nhudaQnVHsSZcjgS4x.4YFAsqpB7SDe', 'admin@example.com', '13800138000', '管理员', 1, 1, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE 
password = VALUES(password), 
email = VALUES(email), 
phone = VALUES(phone), 
nickname = VALUES(nickname), 
role_id = VALUES(role_id), 
is_super_admin = VALUES(is_super_admin),
status = VALUES(status), 
updated_at = NOW();

-- 索引已在表创建时定义，无需重复创建
