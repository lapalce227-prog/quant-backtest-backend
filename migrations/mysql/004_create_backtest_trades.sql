CREATE TABLE IF NOT EXISTS `backtest_trades` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `backtest_run_id` BIGINT UNSIGNED NOT NULL,
    `symbol` VARCHAR(64) NOT NULL,
    `side` VARCHAR(16) NOT NULL,
    `entry_time` DATETIME(3) NOT NULL,
    `entry_price` DECIMAL(20,8) NOT NULL,
    `exit_time` DATETIME(3) NOT NULL,
    `exit_price` DECIMAL(20,8) NOT NULL,
    `quantity` DECIMAL(20,8) NOT NULL,
    `pnl` DECIMAL(20,8) NOT NULL,
    `return_rate` DECIMAL(12,6) NOT NULL,
    `commission` DECIMAL(20,8) NOT NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    KEY `idx_backtest_trades_backtest_run_id` (`backtest_run_id`),
    CONSTRAINT `fk_backtest_trades_backtest_run_id`
        FOREIGN KEY (`backtest_run_id`) REFERENCES `backtest_runs` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
