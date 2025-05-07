CREATE TABLE IF NOT EXISTS `entries` (
  `id` CHAR(36) NOT NULL COMMENT "ID",
  `account_id` CHAR(36) NOT NULL COMMENT "アカウントID",
  `volume_id` CHAR(36) NOT NULL COMMENT "ボリュームID",
  `key` VARCHAR(255) NOT NULL COMMENT "キー",
  `size` BIGINT UNSIGNED NOT NULL COMMENT "サイズ",
  `type` VARCHAR(255) NOT NULL COMMENT "タイプ",
  `created_at` DATETIME (6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT "作成日時",
  `updated_at` DATETIME (6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT "更新日時",
  PRIMARY KEY (`id`),
  UNIQUE `uq_entries_volume_id_and_key` (`volume_id`, `key`),
  CONSTRAINT `fk_entries_volume_id` FOREIGN KEY (`volume_id`) REFERENCES `volumes` (`id`)
);
