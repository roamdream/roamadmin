package service

import (
	"context"

	"github.com/roamdream/roamadmin/kernel/ra"
)

// EnsureSchema 幂等建表（方言无关，无反引号 / ON DUPLICATE KEY；与 manifest/sql/iot.sql 保持一致）。
func EnsureSchema(ctx context.Context) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS ra_device_model (
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
)`,
		`CREATE TABLE IF NOT EXISTS ra_device_model_property (
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
)`,
		`CREATE TABLE IF NOT EXISTS ra_device (
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
)`,
		`CREATE TABLE IF NOT EXISTS ra_device_shadow (
    id         BIGINT   NOT NULL,
    device_id  BIGINT   NOT NULL,
    reported   TEXT,
    desired    TEXT,
    version    BIGINT   NOT NULL DEFAULT 0,
    updated_at DATETIME NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (device_id)
)`,
	}
	for _, s := range stmts {
		if _, err := ra.DB().DB().Exec(ctx, s); err != nil {
			return err
		}
	}
	return nil
}
