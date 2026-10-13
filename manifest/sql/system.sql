-- RoamAdmin system 模块初始化 SQL（系统管理：用户/角色/权限/菜单/字典/租户/操作日志）
-- 方言兼容：不含反引号 / ON DUPLICATE KEY 等 MySQL 特有语法；
-- 主键统一雪花 bigint，租户相关表带 tenant_id 列，表名前缀统一 ra_。

-- 租户表（租户注册表本身，无 tenant_id）
CREATE TABLE IF NOT EXISTS ra_tenant (
    id         BIGINT       NOT NULL COMMENT '雪花主键',
    name       VARCHAR(64)  NOT NULL COMMENT '租户名称',
    code       VARCHAR(64)  NOT NULL COMMENT '租户编码',
    status     SMALLINT     NOT NULL DEFAULT 1 COMMENT '状态（1启用 0禁用）',
    remark     VARCHAR(255) DEFAULT '' COMMENT '备注',
    created_at DATETIME     NOT NULL COMMENT '创建时间',
    updated_at DATETIME     NOT NULL COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE (code)
);

-- 用户表（租户隔离）
CREATE TABLE IF NOT EXISTS ra_user (
    id         BIGINT       NOT NULL COMMENT '雪花主键',
    tenant_id  BIGINT       NOT NULL COMMENT '租户 ID',
    username   VARCHAR(64)  NOT NULL COMMENT '用户名（全局唯一，用于登录）',
    password   VARCHAR(100) NOT NULL COMMENT '密码（bcrypt）',
    nickname   VARCHAR(64)  NOT NULL COMMENT '昵称',
    email      VARCHAR(128) DEFAULT '' COMMENT '邮箱',
    phone      VARCHAR(32)  DEFAULT '' COMMENT '手机号',
    avatar     VARCHAR(255) DEFAULT '' COMMENT '头像',
    status     SMALLINT     NOT NULL DEFAULT 1 COMMENT '状态（1启用 0禁用）',
    remark     VARCHAR(255) DEFAULT '' COMMENT '备注',
    created_at DATETIME     NOT NULL COMMENT '创建时间',
    updated_at DATETIME     NOT NULL COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE (username)
);

-- 角色表（租户隔离）
CREATE TABLE IF NOT EXISTS ra_role (
    id         BIGINT       NOT NULL COMMENT '雪花主键',
    tenant_id  BIGINT       NOT NULL COMMENT '租户 ID',
    name       VARCHAR(64)  NOT NULL COMMENT '角色名称',
    code       VARCHAR(64)  NOT NULL COMMENT '角色编码',
    status     SMALLINT     NOT NULL DEFAULT 1 COMMENT '状态（1启用 0禁用）',
    remark     VARCHAR(255) DEFAULT '' COMMENT '备注',
    created_at DATETIME     NOT NULL COMMENT '创建时间',
    updated_at DATETIME     NOT NULL COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE (tenant_id, code)
);

-- 用户-角色关联表（租户隔离）
CREATE TABLE IF NOT EXISTS ra_user_role (
    id        BIGINT NOT NULL COMMENT '雪花主键',
    user_id   BIGINT NOT NULL COMMENT '用户 ID',
    role_id   BIGINT NOT NULL COMMENT '角色 ID',
    tenant_id BIGINT NOT NULL COMMENT '租户 ID',
    PRIMARY KEY (id),
    UNIQUE (user_id, role_id)
);

-- 菜单/权限表（租户隔离，树形）
CREATE TABLE IF NOT EXISTS ra_menu (
    id         BIGINT       NOT NULL COMMENT '雪花主键',
    tenant_id  BIGINT       NOT NULL COMMENT '租户 ID',
    parent_id  BIGINT       NOT NULL DEFAULT 0 COMMENT '父菜单 ID（0 为根）',
    name       VARCHAR(64)  NOT NULL COMMENT '菜单名称',
    type       VARCHAR(16)  NOT NULL DEFAULT 'menu' COMMENT '类型（dir/menu/button）',
    perm       VARCHAR(128) DEFAULT '' COMMENT '权限码',
    path       VARCHAR(255) DEFAULT '' COMMENT '路由地址',
    component  VARCHAR(255) DEFAULT '' COMMENT '组件路径',
    icon       VARCHAR(64)  DEFAULT '' COMMENT '图标',
    sort       INT          NOT NULL DEFAULT 0 COMMENT '排序',
    status     SMALLINT     NOT NULL DEFAULT 1 COMMENT '状态（1启用 0禁用）',
    created_at DATETIME     NOT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
);

-- 角色-菜单关联表（租户隔离）
CREATE TABLE IF NOT EXISTS ra_role_menu (
    id        BIGINT NOT NULL COMMENT '雪花主键',
    role_id   BIGINT NOT NULL COMMENT '角色 ID',
    menu_id   BIGINT NOT NULL COMMENT '菜单 ID',
    tenant_id BIGINT NOT NULL COMMENT '租户 ID',
    PRIMARY KEY (id),
    UNIQUE (role_id, menu_id)
);

-- 字典类型表（租户隔离）
CREATE TABLE IF NOT EXISTS ra_dict_type (
    id         BIGINT       NOT NULL COMMENT '雪花主键',
    tenant_id  BIGINT       NOT NULL COMMENT '租户 ID',
    name       VARCHAR(64)  NOT NULL COMMENT '字典名称',
    type       VARCHAR(64)  NOT NULL COMMENT '字典类型编码',
    status     SMALLINT     NOT NULL DEFAULT 1 COMMENT '状态（1启用 0禁用）',
    remark     VARCHAR(255) DEFAULT '' COMMENT '备注',
    created_at DATETIME     NOT NULL COMMENT '创建时间',
    PRIMARY KEY (id),
    UNIQUE (tenant_id, type)
);

-- 字典数据表（租户隔离）
CREATE TABLE IF NOT EXISTS ra_dict_data (
    id         BIGINT       NOT NULL COMMENT '雪花主键',
    tenant_id  BIGINT       NOT NULL COMMENT '租户 ID',
    dict_type  VARCHAR(64)  NOT NULL COMMENT '字典类型编码',
    label      VARCHAR(64)  NOT NULL COMMENT '标签',
    value      VARCHAR(64)  NOT NULL COMMENT '值',
    sort       INT          NOT NULL DEFAULT 0 COMMENT '排序',
    status     SMALLINT     NOT NULL DEFAULT 1 COMMENT '状态（1启用 0禁用）',
    remark     VARCHAR(255) DEFAULT '' COMMENT '备注',
    created_at DATETIME     NOT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
);

-- 操作日志表（租户隔离）
CREATE TABLE IF NOT EXISTS ra_oper_log (
    id         BIGINT       NOT NULL COMMENT '雪花主键',
    tenant_id  BIGINT       NOT NULL COMMENT '租户 ID',
    user_id    BIGINT       NOT NULL DEFAULT 0 COMMENT '用户 ID',
    username   VARCHAR(64)  DEFAULT '' COMMENT '用户名',
    module     VARCHAR(64)  DEFAULT '' COMMENT '模块',
    action     VARCHAR(64)  DEFAULT '' COMMENT '操作',
    method     VARCHAR(16)  DEFAULT '' COMMENT '请求方法',
    path       VARCHAR(255) DEFAULT '' COMMENT '请求路径',
    ip         VARCHAR(64)  DEFAULT '' COMMENT '客户端 IP',
    status     SMALLINT     NOT NULL DEFAULT 0 COMMENT '状态（0成功 1失败）',
    err_msg    VARCHAR(512) DEFAULT '' COMMENT '错误信息',
    created_at DATETIME     NOT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
);

-- 初始化数据：默认租户 / 超级管理员（用户名 admin，密码 admin123，bcrypt）
INSERT INTO ra_tenant (id, name, code, status, remark, created_at, updated_at)
VALUES (1, 'Default Tenant', 'default', 1, 'Built-in default tenant', NOW(), NOW());

INSERT INTO ra_user (id, tenant_id, username, password, nickname, status, remark, created_at, updated_at)
VALUES (1, 1, 'admin', '$2a$10$RSN7r8q8dTLFvAP8XuA7wuu5cE4DQ5TioHhFQ/LINsmi5Cra5nU3m', 'Super Admin', 1, 'Built-in super admin', NOW(), NOW());

INSERT INTO ra_role (id, tenant_id, name, code, status, remark, created_at, updated_at)
VALUES (1, 1, 'Super Admin', 'super', 1, 'Built-in super admin role', NOW(), NOW());

INSERT INTO ra_user_role (id, user_id, role_id, tenant_id) VALUES (1, 1, 1, 1);

-- 初始化菜单（租户 1）与超级管理员角色的菜单关联（权限码 RBAC）
INSERT INTO ra_menu (id, tenant_id, parent_id, name, type, perm, path, component, icon, sort, status, created_at)
VALUES
    (1, 1, 0, 'System', 'dir', '', '', '', 'setting', 0, 1, NOW()),
    (2, 1, 1, 'User Management', 'menu', 'system:user:list', '/system/user', '', 'user', 1, 1, NOW()),
    (3, 1, 1, 'Role Management', 'menu', 'system:role:list', '/system/role', '', 'team', 2, 1, NOW()),
    (4, 1, 1, 'Menu Management', 'menu', 'system:menu:list', '/system/menu', '', 'menu', 3, 1, NOW()),
    (5, 1, 1, 'Dict Management', 'menu', 'system:dict:list', '/system/dict', '', 'book', 4, 1, NOW()),
    (6, 1, 1, 'Tenant Management', 'menu', 'system:tenant:list', '/system/tenant', '', 'apartment', 5, 1, NOW()),
    (7, 1, 1, 'Operation Logs', 'menu', 'system:operlog:list', '/system/operlog', '', 'file', 6, 1, NOW()),
    (8, 1, 2, 'Create User', 'button', 'system:user:create', '', '', '', 1, 1, NOW()),
    (9, 1, 2, 'Edit User', 'button', 'system:user:update', '', '', '', 2, 1, NOW()),
    (10, 1, 2, 'Delete User', 'button', 'system:user:delete', '', '', '', 3, 1, NOW());

INSERT INTO ra_role_menu (id, role_id, menu_id, tenant_id)
VALUES
    (1, 1, 1, 1),
    (2, 1, 2, 1),
    (3, 1, 3, 1),
    (4, 1, 4, 1),
    (5, 1, 5, 1),
    (6, 1, 6, 1),
    (7, 1, 7, 1),
    (8, 1, 8, 1),
    (9, 1, 9, 1),
    (10, 1, 10, 1);
