-- 智汇选课系统数据库表创建语句 (优化版V2)
-- 生成时间：2025-10-12
-- ----------------------------
-- 权限表 (permissions)
-- ----------------------------
CREATE TABLE `permissions` (
    `permission_id` INT NOT NULL AUTO_INCREMENT COMMENT '权限ID',
    `permission_name` VARCHAR(50) NOT NULL COMMENT '权限名称',
    `description` VARCHAR(500) COMMENT '权限描述',
    PRIMARY KEY (`permission_id`),
    UNIQUE KEY `uk_permission_name` (`permission_name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '权限表';

-- ----------------------------
-- 角色表 (roles)
-- ----------------------------
CREATE TABLE `roles` (
    `role_id` INT NOT NULL AUTO_INCREMENT COMMENT '角色ID',
    `role_name` VARCHAR(50) NOT NULL COMMENT '角色名称',
    `description` VARCHAR(500) COMMENT '角色描述',
    PRIMARY KEY (`role_id`),
    UNIQUE KEY `uk_role_name` (`role_name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '角色表';

-- ----------------------------
-- 角色-权限关联表 (role_permissions)
-- ----------------------------
CREATE TABLE `role_permissions` (
    `role_id` INT NOT NULL COMMENT '角色ID',
    `permission_id` INT NOT NULL COMMENT '权限ID',
    PRIMARY KEY (`role_id`, `permission_id`),
    CONSTRAINT `fk_rp_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`role_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_rp_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`permission_id`) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '角色-权限关联表';

-- ----------------------------
-- 学院表 (colleges)
-- ----------------------------
CREATE TABLE `colleges` (
    `college_id` INT NOT NULL AUTO_INCREMENT COMMENT '学院ID',
    `college_name` VARCHAR(50) NOT NULL COMMENT '学院名称',
    `school` VARCHAR(50) NOT NULL COMMENT '所属学校',
    PRIMARY KEY (`college_id`),
    UNIQUE KEY `uk_college_name` (`college_name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '学院表';

-- ----------------------------
-- 专业表 (majors)
-- ----------------------------
CREATE TABLE `majors` (
    `major_id` INT NOT NULL AUTO_INCREMENT COMMENT '专业ID',
    `major_name` VARCHAR(50) NOT NULL COMMENT '专业名称',
    `college_id` INT NOT NULL COMMENT '学院ID',
    PRIMARY KEY (`major_id`),
    KEY `fk_major_college` (`college_id`),
    CONSTRAINT `fk_major_college` FOREIGN KEY (`college_id`) REFERENCES `colleges` (`college_id`),
    UNIQUE KEY `uk_major_college` (`major_name`, `college_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '专业表';

-- ----------------------------
-- 教师表 (teachers)
-- ----------------------------
CREATE TABLE `teachers` (
    `teacher_id` INT NOT NULL AUTO_INCREMENT COMMENT '教师ID',
    `name` VARCHAR(50) NOT NULL COMMENT '教师姓名',
    `college_id` INT NOT NULL COMMENT '学院ID',
    `introduction` VARCHAR(1000) COMMENT '教师简介',
    PRIMARY KEY (`teacher_id`),
    KEY `fk_teacher_college` (`college_id`),
    CONSTRAINT `fk_teacher_college` FOREIGN KEY (`college_id`) REFERENCES `colleges` (`college_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '教师表';

-- ----------------------------
-- 用户表 (users) - 时间字段改为TIMESTAMP
-- ----------------------------
CREATE TABLE `users` (
    `user_id` INT NOT NULL AUTO_INCREMENT COMMENT '用户唯一标识',
    `username` VARCHAR(16) NOT NULL COMMENT '用户名',
    `password_hash` VARCHAR(255) NOT NULL COMMENT '加密密码',
    `email` VARCHAR(100) NOT NULL COMMENT '邮箱',
    `college_id` INT COMMENT '所属学院ID',
    `major_id` INT COMMENT '所属专业ID',
    `avatar_url` VARCHAR(255) COMMENT '头像URL',
    `reputation_score` INT NOT NULL DEFAULT 80 COMMENT '信誉分',
    `role_id` INT NOT NULL COMMENT '角色ID',
    `status` ENUM('active','unactive', 'locked', 'banned') NOT NULL DEFAULT 'unactive' COMMENT '账户状态',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`user_id`),
    UNIQUE KEY `uk_username` (`username`),
    UNIQUE KEY `uk_email` (`email`),
    KEY `fk_user_role` (`role_id`),
    KEY `fk_user_college` (`college_id`),
    KEY `fk_user_major` (`major_id`),
    CONSTRAINT `fk_user_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`role_id`),
    CONSTRAINT `fk_user_college` FOREIGN KEY (`college_id`) REFERENCES `colleges` (`college_id`),
    CONSTRAINT `fk_user_major` FOREIGN KEY (`major_id`) REFERENCES `majors` (`major_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '用户表';

-- ----------------------------
-- 课程表 (courses) - 时间字段改为TIMESTAMP
-- ----------------------------
CREATE TABLE `courses` (
    `course_id` INT NOT NULL AUTO_INCREMENT COMMENT '课程ID',
    `course_name` VARCHAR(100) NOT NULL COMMENT '课程名称',
    `teacher_id` INT NOT NULL COMMENT '教师ID',
    `credit` DECIMAL(2, 1) NOT NULL COMMENT '学分',
    `major_id` INT NOT NULL COMMENT '专业ID',
    `grade` VARCHAR(20) NOT NULL COMMENT '年级',
    `description` VARCHAR(1000) COMMENT '课程简介',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`course_id`),
    KEY `fk_course_teacher` (`teacher_id`),
    KEY `fk_course_major` (`major_id`),
    KEY `idx_course_name` (`name`),
    CONSTRAINT `fk_course_teacher` FOREIGN KEY (`teacher_id`) REFERENCES `teachers` (`teacher_id`),
    CONSTRAINT `fk_course_major` FOREIGN KEY (`major_id`) REFERENCES `majors` (`major_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '课程表';

-- ----------------------------
-- 标签表 (tags)
-- ----------------------------
CREATE TABLE `tags` (
    `tag_id` INT NOT NULL AUTO_INCREMENT COMMENT '标签ID',
    `tag_name` VARCHAR(50) NOT NULL COMMENT '标签名称',
    PRIMARY KEY (`tag_id`),
    UNIQUE KEY `uk_tag_name` (`tag_name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '标签表';

-- ----------------------------
-- 资源表 (resources) - 时间字段改为TIMESTAMP
-- ----------------------------
CREATE TABLE `resources` (
    `resource_id` INT NOT NULL AUTO_INCREMENT COMMENT '资源ID',
    `resource_name` VARCHAR(255) NOT NULL COMMENT '资源标题',
    `description` VARCHAR(500) COMMENT '资源简介',
    `resource_url` VARCHAR(255) NOT NULL COMMENT '文件URL',
    `type` ENUM('pdf', 'docx', 'pptx', 'zip') NOT NULL COMMENT '文件类型',
    `size` INT NOT NULL COMMENT '文件大小(字节)',
    `uploader_id` INT NOT NULL COMMENT '上传者ID',
    `course_id` INT NOT NULL COMMENT '课程ID',
    `download_count` INT DEFAULT 0 COMMENT '下载次数',
    `average_rating` DECIMAL(2, 1) DEFAULT 0 COMMENT '平均评分',
    `rating_count` INT DEFAULT 0 COMMENT '评分人数',
    `status` ENUM('normal', 'low_quality', 'pending_review') DEFAULT 'pending_review' COMMENT '状态',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`resource_id`),
    KEY `fk_resource_uploader` (`uploader_id`),
    KEY `fk_resource_course` (`course_id`),
    KEY `idx_resource_status` (`status`),
    CONSTRAINT `fk_resource_uploader` FOREIGN KEY (`uploader_id`) REFERENCES `users` (`user_id`),
    CONSTRAINT `fk_resource_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '资源表';

-- ----------------------------
-- 资源标签关联表 (resource_tags)
-- ----------------------------
CREATE TABLE `resource_tags` (
    `resource_id` INT NOT NULL COMMENT '资源ID',
    `tag_id` INT NOT NULL COMMENT '标签ID',
    PRIMARY KEY (`resource_id`, `tag_id`),
    KEY `fk_rt_tag` (`tag_id`),
    CONSTRAINT `fk_rt_resource` FOREIGN KEY (`resource_id`) REFERENCES `resources` (`resource_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_rt_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`tag_id`) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '资源标签关联表';

-- ----------------------------
-- 课程评分表 (course_ratings) - 时间字段改为TIMESTAMP
-- ----------------------------
CREATE TABLE `course_ratings` (
    `rating_id` INT NOT NULL AUTO_INCREMENT COMMENT '评分ID',
    `user_id` INT NOT NULL COMMENT '用户ID',
    `course_id` INT NOT NULL COMMENT '课程ID',
    `recommendation` DECIMAL(2, 1) NOT NULL COMMENT '综合推荐度(0-5)',
    `difficulty` INT NOT NULL COMMENT '课程难度(1-5)',
    `workload` INT NOT NULL COMMENT '作业压力(1-5)',
    `usefulness` INT NOT NULL COMMENT '知识实用性(1-5)',
    `is_visible` BOOLEAN DEFAULT TRUE COMMENT '是否显示',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`rating_id`),
    KEY `fk_cr_user` (`user_id`),
    KEY `fk_cr_course` (`course_id`),
    CONSTRAINT `fk_cr_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`),
    CONSTRAINT `fk_cr_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE,
    UNIQUE KEY `uk_user_course` (`user_id`, `course_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '课程评分表';

-- ----------------------------
-- 课程评论表 (course_comments) - 时间字段改为TIMESTAMP
-- ----------------------------
CREATE TABLE `course_comments` (
    `comment_id` INT NOT NULL AUTO_INCREMENT COMMENT '评论ID',
    `user_id` INT NOT NULL COMMENT '用户ID',
    `course_id` INT NOT NULL COMMENT '课程ID',
    `content` TEXT NOT NULL COMMENT '评论内容',
    `parent_id` INT COMMENT '父评论ID',
    `likes` INT DEFAULT 0 COMMENT '点赞数',
    `is_visible` BOOLEAN DEFAULT TRUE COMMENT '是否显示',
    `status` ENUM('normal', 'deleted_by_user', 'deleted_by_admin') DEFAULT 'normal' COMMENT '评论状态',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`comment_id`),
    KEY `fk_cc_user` (`user_id`),
    KEY `fk_cc_course` (`course_id`),
    KEY `fk_cc_parent` (`parent_id`),
    CONSTRAINT `fk_cc_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`),
    CONSTRAINT `fk_cc_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_cc_parent` FOREIGN KEY (`parent_id`) REFERENCES `course_comments` (`comment_id`) ON DELETE NO ACTION
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '课程评论表';

-- ----------------------------
-- 资源评分表 (resource_ratings) - 时间字段改为TIMESTAMP
-- ----------------------------
CREATE TABLE `resource_ratings` (
    `rating_id` INT NOT NULL AUTO_INCREMENT COMMENT '评分ID',
    `user_id` INT NOT NULL COMMENT '用户ID',
    `resource_id` INT NOT NULL COMMENT '资源ID',
    `recommendation` DECIMAL(2, 1) NOT NULL COMMENT '综合推荐度(0-5)',
    `is_visible` BOOLEAN DEFAULT TRUE COMMENT '是否显示',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`rating_id`),
    KEY `fk_rr_user` (`user_id`),
    KEY `fk_rr_resource` (`resource_id`),
    CONSTRAINT `fk_rr_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`),
    CONSTRAINT `fk_rr_resource` FOREIGN KEY (`resource_id`) REFERENCES `resources` (`resource_id`) ON DELETE CASCADE,
    UNIQUE KEY `uk_user_resource` (`user_id`, `resource_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '资源评分表';

-- ----------------------------
-- 资源评论表 (resource_comments) - 时间字段改为TIMESTAMP
-- ----------------------------
CREATE TABLE `resource_comments` (
    `comment_id` INT NOT NULL AUTO_INCREMENT COMMENT '评论ID',
    `user_id` INT NOT NULL COMMENT '用户ID',
    `resource_id` INT NOT NULL COMMENT '资源ID',
    `content` TEXT NOT NULL COMMENT '评论内容',
    `parent_id` INT COMMENT '父评论ID',
    `likes` INT DEFAULT 0 COMMENT '点赞数',
    `is_visible` BOOLEAN DEFAULT TRUE COMMENT '是否显示',
    `status` ENUM('normal', 'deleted_by_user', 'deleted_by_admin') DEFAULT 'normal' COMMENT '评论状态',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`comment_id`),
    KEY `fk_rc_user` (`user_id`),
    KEY `fk_rc_resource` (`resource_id`),
    KEY `fk_rc_parent` (`parent_id`),
    CONSTRAINT `fk_rc_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`),
    CONSTRAINT `fk_rc_resource` FOREIGN KEY (`resource_id`) REFERENCES `resources` (`resource_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_rc_parent` FOREIGN KEY (`parent_id`) REFERENCES `resource_comments` (`comment_id`) ON DELETE NO ACTION
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '资源评论表';

-- ----------------------------
-- 信誉分记录表 (reputation_records) - 时间字段改为TIMESTAMP
-- ----------------------------
CREATE TABLE `reputation_records` (
    `record_id` INT NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    `user_id` INT NOT NULL COMMENT '用户ID',
    `change_score` INT NOT NULL COMMENT '分数变化',
    `reason` TEXT NOT NULL COMMENT '变化原因',
    `related_id` INT COMMENT '关联对象ID',
    `related_type` ENUM('resource', 'comment', 'rating') COMMENT '关联对象类型',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录时间',
    PRIMARY KEY (`record_id`),
    KEY `fk_reprec_user` (`user_id`),
    KEY `idx_reputation_time` (`created_at`),
    CONSTRAINT `fk_reprec_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '信誉分记录表';

-- ----------------------------
-- 收藏表 (favorites) - 时间字段改为TIMESTAMP
-- ----------------------------
CREATE TABLE `favorites` (
    `favorite_id` INT NOT NULL AUTO_INCREMENT COMMENT '收藏ID',
    `user_id` INT NOT NULL COMMENT '用户ID',
    `target_id` INT NOT NULL COMMENT '对象ID',
    `target_type` ENUM('course', 'resource', 'teacher') NOT NULL COMMENT '对象类型',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '收藏时间',
    PRIMARY KEY (`favorite_id`),
    KEY `fk_favorite_user` (`user_id`),
    KEY `idx_favorite_type` (`user_id`, `target_type`),
    CONSTRAINT `fk_favorite_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE,
    UNIQUE KEY `uk_user_target` (`user_id`, `target_id`, `target_type`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '收藏表';

-- ----------------------------
-- 审核表 (reviews) - 时间字段改为TIMESTAMP
-- ----------------------------
CREATE TABLE `reviews` (
    `review_id` INT NOT NULL AUTO_INCREMENT COMMENT '审核ID',
    `target_id` INT NOT NULL COMMENT '对象ID',
    `target_type` ENUM(
        'resource',
        'course_rating',
        'resource_rating',
        'comment'
    ) NOT NULL COMMENT '对象类型',
    `reason` VARCHAR(500) NOT NULL COMMENT '审核原因',
    `status` ENUM('pending', 'approved', 'rejected') DEFAULT 'pending' COMMENT '状态',
    `priority` INT DEFAULT 3 COMMENT '优先级(1-5)',
    `reviewer_id` INT COMMENT '审核员ID',
    `reviewed_at` TIMESTAMP COMMENT '审核时间',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`review_id`),
    KEY `fk_review_reviewer` (`reviewer_id`),
    KEY `idx_review_status` (`status`),
    KEY `idx_review_priority` (`priority`),
    CONSTRAINT `fk_review_reviewer` FOREIGN KEY (`reviewer_id`) REFERENCES `users` (`user_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '审核表';

-- ----------------------------
-- 物品表 (items)
-- ----------------------------
CREATE TABLE `items` (
    `item_id` INT NOT NULL AUTO_INCREMENT COMMENT '物品ID',
    `name` VARCHAR(100) NOT NULL COMMENT '物品名称',
    `type` ENUM(
        'avatar_frame',
        'background',
        'medal',
        'nickname_color'
    ) NOT NULL COMMENT '物品类型',
    `price` INT NOT NULL COMMENT '所需积分',
    `description` VARCHAR(500) COMMENT '物品描述',
    `image_url` VARCHAR(255) COMMENT '预览图URL',
    PRIMARY KEY (`item_id`),
    UNIQUE KEY `uk_item_name` (`name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '物品表';

-- ----------------------------
-- 用户物品表 (user_items) - 时间字段改为TIMESTAMP
-- ----------------------------
CREATE TABLE `user_items` (
    `user_item_id` INT NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    `user_id` INT NOT NULL COMMENT '用户ID',
    `item_id` INT NOT NULL COMMENT '物品ID',
    `is_used` BOOLEAN DEFAULT FALSE COMMENT '是否启用',
    `obtained_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '获取时间',
    PRIMARY KEY (`user_item_id`),
    KEY `fk_ui_user` (`user_id`),
    KEY `fk_ui_item` (`item_id`),
    CONSTRAINT `fk_ui_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_ui_item` FOREIGN KEY (`item_id`) REFERENCES `items` (`item_id`) ON DELETE CASCADE,
    UNIQUE KEY `uk_user_item` (`user_id`, `item_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '用户物品表';
