-- 数据库初始化脚本

-- 创建角色表
CREATE TABLE IF NOT EXISTS `roles` (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL UNIQUE,
  `description` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建权限表
CREATE TABLE IF NOT EXISTS `permissions` (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL UNIQUE,
  `description` varchar(255) DEFAULT NULL,
  `resource` varchar(50) DEFAULT NULL,
  `action` varchar(50) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建角色-权限关联表
CREATE TABLE IF NOT EXISTS `role_permissions` (
  `role_id` int UNSIGNED NOT NULL,
  `permission_id` int UNSIGNED NOT NULL,
  PRIMARY KEY (`role_id`, `permission_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_permission_id` (`permission_id`),
  CONSTRAINT `fk_role_permissions_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_role_permissions_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建班级表
CREATE TABLE IF NOT EXISTS `classes` (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL UNIQUE,
  `description` varchar(255) DEFAULT NULL,
  `teacher_id` int UNSIGNED NOT NULL,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_teacher_id` (`teacher_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL UNIQUE,
  `password` varchar(255) NOT NULL,
  `email` varchar(100) UNIQUE,
  `phone` varchar(20) UNIQUE,
  `nickname` varchar(50) DEFAULT NULL,
  `avatar` varchar(255) DEFAULT NULL,
  `role_id` int UNSIGNED NOT NULL,
  `status` tinyint DEFAULT 1,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `student_id` varchar(20) DEFAULT NULL,
  `class_id` int UNSIGNED DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_class_id` (`class_id`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_student_id` (`student_id`),
  CONSTRAINT `fk_users_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`),
  CONSTRAINT `fk_users_class` FOREIGN KEY (`class_id`) REFERENCES `classes` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建学生-教师关联表
CREATE TABLE IF NOT EXISTS `student_teachers` (
  `student_id` int UNSIGNED NOT NULL,
  `teacher_id` int UNSIGNED NOT NULL,
  PRIMARY KEY (`student_id`, `teacher_id`),
  KEY `idx_student_id` (`student_id`),
  KEY `idx_teacher_id` (`teacher_id`),
  CONSTRAINT `fk_student_teachers_student` FOREIGN KEY (`student_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_student_teachers_teacher` FOREIGN KEY (`teacher_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建用户会话表
CREATE TABLE IF NOT EXISTS `user_sessions` (
  `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` int UNSIGNED NOT NULL,
  `access_token` text NOT NULL,
  `refresh_token` text NOT NULL,
  `ip` varchar(45) DEFAULT NULL,
  `user_agent` varchar(255) DEFAULT NULL,
  `expires_at` datetime(3) NOT NULL,
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  CONSTRAINT `fk_user_sessions_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入初始角色数据
INSERT IGNORE INTO `roles` (`id`, `name`, `description`) VALUES
(1, 'admin', '管理员角色'),
(2, 'student', '学生角色'),
(3, 'teacher', '教师角色');

-- 插入初始权限数据
INSERT IGNORE INTO `permissions` (`id`, `name`, `description`, `resource`, `action`) VALUES
(1, 'user:create', '创建用户', 'user', 'create'),
(2, 'user:read', '读取用户信息', 'user', 'read'),
(3, 'user:update', '更新用户信息', 'user', 'update'),
(4, 'user:delete', '删除用户', 'user', 'delete'),
(5, 'class:create', '创建班级', 'class', 'create'),
(6, 'class:read', '读取班级信息', 'class', 'read'),
(7, 'class:update', '更新班级信息', 'class', 'update'),
(8, 'class:delete', '删除班级', 'class', 'delete');

-- 插入初始角色-权限关联数据
INSERT IGNORE INTO `role_permissions` (`role_id`, `permission_id`) VALUES
-- 管理员拥有所有权限
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7), (1, 8),
-- 教师权限
(3, 2), (3, 3), (3, 5), (3, 6), (3, 7),
-- 学生权限
(2, 2), (2, 3);
