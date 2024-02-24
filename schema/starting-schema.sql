CREATE TABLE IF NOT EXISTS `player` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'pk',
    `name` varchar(255) NOT NULL DEFAULT '',
    `player_level` int(11) unsigned NOT NULL DEFAULT 1,
    `player_exp` int(22) unsigned NOT NULL DEFAULT 0,
    `forge_level` int(11) unsigned NOT NULL DEFAULT 1,
    `forge_exp` int(22) unsigned NOT NULL DEFAULT 0,
    `is_active` TINYINT(1) DEFAULT 1,
    `created` int(11) NOT NULL DEFAULT 0,
    `updated` int(11) NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

INSERT IGNORE INTO `app`.`player` (`name`) VALUES ('dashpy');

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
    `level` smallint(4) unsigned NOT NULL DEFAULT 0,
    `rarity_id` int(11) unsigned NOT NULL,
    `hp` int(11) unsigned NOT NULL DEFAULT 0,
    `attack` int(11) unsigned NOT NULL DEFAULT 0,
    `defense` int(11) unsigned NOT NULL DEFAULT 0,
    `speed` int(11) unsigned NOT NULL DEFAULT 0,
    `crit` smallint(4) unsigned NOT NULL DEFAULT 0,
    `dodge` smallint(4) unsigned NOT NULL DEFAULT 0,
    `block` smallint(4) unsigned NOT NULL DEFAULT 0,
    CONSTRAINT UC_player_slot UNIQUE (player_id, gear_slot_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `ore` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'pk',
    `name` varchar(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `ore_inventory` (
    `player_id` int(11) unsigned NOT NULL,
    `quantity` int(11) unsigned NOT NULL DEFAULT 0,
    PRIMARY KEY (`player_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `player_level` (
    `player_id` int(11) unsigned NOT NULL,
    `level` int(11) unsigned NOT NULL DEFAULT 1,
    PRIMARY KEY (`player_id`)
)

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

INSERT IGNORE INTO `app`.`ore` (`name`) VALUES ('Copper');
INSERT IGNORE INTO `app`.`ore` (`name`) VALUES ('Iron');
INSERT IGNORE INTO `app`.`ore` (`name`) VALUES ('Nickle');
INSERT IGNORE INTO `app`.`ore` (`name`) VALUES ('Tin');
INSERT IGNORE INTO `app`.`ore` (`name`) VALUES ('Silver');
INSERT IGNORE INTO `app`.`ore` (`name`) VALUES ('Gold');
INSERT IGNORE INTO `app`.`ore` (`name`) VALUES ('Platinum');
INSERT IGNORE INTO `app`.`ore` (`name`) VALUES ('Cobalt');
INSERT IGNORE INTO `app`.`ore` (`name`) VALUES ('Mithril');