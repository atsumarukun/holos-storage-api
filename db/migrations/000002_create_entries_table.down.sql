ALTER TABLE `entries`
DROP FOREIGN KEY `fk_entries_volume_id`;

ALTER TABLE `entries`
DROP INDEX `uq_entries_volume_id_and_key`;

DROP TABLE IF EXISTS `entries`;
