-- RoamAdmin 初始化 SQL（示例：demo 模块）
-- 方言兼容：本脚本不含 MySQL 特有语法（反引号 / ON DUPLICATE KEY 等），
-- 后续由迁移工具统一管理 schema，并进入 MySQL8 / PostgreSQL 双库 CI。

CREATE TABLE IF NOT EXISTS ra_demo_item (
    id         BIGINT       NOT NULL COMMENT '雪花主键',
    name       VARCHAR(64)  NOT NULL COMMENT '名称',
    remark     VARCHAR(255) DEFAULT '' COMMENT '备注',
    created_at DATETIME     NOT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
) COMMENT='示例数据表';
