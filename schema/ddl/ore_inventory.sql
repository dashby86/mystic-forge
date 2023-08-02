CREATE TABLE IF NOT EXISTS `app`.`ore_inventory` (
    `player_id` int(11) unsigned NOT NULL,
    `quantity` int(11) unsigned NOT NULL DEFAULT 0,
    PRIMARY KEY (`player_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;