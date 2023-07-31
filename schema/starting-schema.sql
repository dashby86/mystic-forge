CREATE TABLE IF NOT EXISTS `player` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'pk',
    `name` varchar(255) NOT NULL DEFAULT '',
    `is_active` TINYINT(1) DEFAULT 1,
    `created` int(11) NOT NULL DEFAULT 0,
    `updated` int(11) NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `rarity` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'pk',
    `name` varchar(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `gear_slot` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'pk',
    `name` varchar(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `player_gear` (
    `player_id` int(11) unsigned NOT NULL,
    `gear_slot_id` int(11) unsigned NOT NULL,
    `rarity_id` int(11) unsigned NOT NULL,
    `hp` int(11) unsigned NOT NULL DEFAULT 0,
    `attack` int(11) unsigned NOT NULL DEFAULT 0,
    `defense` int(11) unsigned NOT NULL DEFAULT 0,
    `speed` decimal(5, 2) DEFAULT 0.00,
    `crit` decimal(5, 2) DEFAULT 0.00,
    `dodge` decimal(5, 2) DEFAULT 0.00,
    `block` decimal(5, 2) DEFAULT 0.00,
    CONSTRAINT UC_player_slot UNIQUE (player_id, gear_slot_id).
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT IGNORE INTO `app`.`gear_slot` (`name`) VALUES ('Helmet');
INSERT IGNORE INTO `app`.`gear_slot` (`name`) VALUES ('Pauldrons');
INSERT IGNORE INTO `app`.`gear_slot` (`name`) VALUES ('Gloves');
INSERT IGNORE INTO `app`.`gear_slot` (`name`) VALUES ('Boots');
INSERT IGNORE INTO `app`.`gear_slot` (`name`) VALUES ('Greaves');
INSERT IGNORE INTO `app`.`gear_slot` (`name`) VALUES ('Ring');
INSERT IGNORE INTO `app`.`gear_slot` (`name`) VALUES ('Necklace');
INSERT IGNORE INTO `app`.`gear_slot` (`name`) VALUES ('Weapon');
INSERT IGNORE INTO `app`.`gear_slot` (`name`) VALUES ('Chest');

INSERT IGNORE INTO `app`.`rarity` (`name`) VALUES ('Common');
INSERT IGNORE INTO `app`.`rarity` (`name`) VALUES ('Uncommon');
INSERT IGNORE INTO `app`.`rarity` (`name`) VALUES ('Rare');
INSERT IGNORE INTO `app`.`rarity` (`name`) VALUES ('Epic');