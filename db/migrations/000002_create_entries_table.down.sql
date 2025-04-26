ALTER TABLE `entries`
DROP INDEX `uq_entries_volume_id_and_key`;

ALTER TABLE `entries`
DROP FOREIGN KEY `fk_entries_volume_id`;

DROP TABLE IF EXISTS `entries`;
