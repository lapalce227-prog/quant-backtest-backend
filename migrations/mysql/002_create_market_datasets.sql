CREATE TABLE IF NOT EXISTS `market_datasets` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `symbol` VARCHAR(64) NOT NULL,
    `timeframe` VARCHAR(32) NOT NULL,
    `source_name` VARCHAR(255) NOT NULL,
    `storage_path` VARCHAR(512) NOT NULL,
    `start_time` DATETIME(3) NOT NULL,
    `end_time` DATETIME(3) NOT NULL,
    `record_count` BIGINT NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    KEY `idx_market_datasets_symbol_timeframe` (`symbol`, `timeframe`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
