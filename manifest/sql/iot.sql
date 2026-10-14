-- RoamAdmin IoT（P3）模块初始化 SQL：物模型 / 设备 / 设备影子
-- 方言兼容：不含反引号 / ON DUPLICATE KEY 等 MySQL 特有语法；
-- 主键统一雪花 bigint，租户相关表带 tenant_id 列，表名前缀统一 ra_。
-- 注：device 模块启动时会幂等建表（EnsureSchema），此文件作为部署 / 审阅基线，内容保持一致。

-- 物模型 / 产品表（租户隔离；product_key 全局唯一，作为 MQTT 上行主题段）
CREATE TABLE IF NOT EXISTS ra_device_model (
    id          BIGINT       NOT NULL,
    tenant_id   BIGINT       NOT NULL,
    name        VARCHAR(64)  NOT NULL,
    product_key VARCHAR(64)  NOT NULL,
    description VARCHAR(255) DEFAULT '',
    status      SMALLINT     NOT NULL DEFAULT 1,
    created_at  DATETIME     NOT NULL,
    updated_at  DATETIME     NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (product_key)
);

-- 物模型属性表（同一物模型内 identifier 唯一）
CREATE TABLE IF NOT EXISTS ra_device_model_property (
    id         BIGINT      NOT NULL,
    model_id   BIGINT      NOT NULL,
    identifier VARCHAR(64) NOT NULL,
    name       VARCHAR(64) NOT NULL,
    data_type  VARCHAR(16) NOT NULL,
    unit       VARCHAR(32) DEFAULT '',
    access     VARCHAR(8)  DEFAULT 'r',
    sort       INT         DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE (model_id, identifier)
);

-- 设备表（租户隔离；product_key + device_key 全局唯一，作为 MQTT 上行主题段）
CREATE TABLE IF NOT EXISTS ra_device (
    id          BIGINT      NOT NULL,
    tenant_id   BIGINT      NOT NULL,
    model_id    BIGINT      NOT NULL,
    product_key VARCHAR(64) NOT NULL,
    device_key  VARCHAR(64) NOT NULL,
    name        VARCHAR(64) NOT NULL,
    secret      VARCHAR(64) NOT NULL,
    status      SMALLINT    NOT NULL DEFAULT 1,
    created_at  DATETIME    NOT NULL,
    updated_at  DATETIME    NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (product_key, device_key)
);

-- 设备影子表（device_id 唯一；reported / desired 为 JSON 文本）
CREATE TABLE IF NOT EXISTS ra_device_shadow (
    id         BIGINT   NOT NULL,
    device_id  BIGINT   NOT NULL,
    reported   TEXT,
    desired    TEXT,
    version    BIGINT   NOT NULL DEFAULT 0,
    updated_at DATETIME NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (device_id)
);
