CREATE TABLE IF NOT EXISTS `app`.`player_gear` (
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
    CONSTRAINT UC_player_slot UNIQUE (player_id, gear_slot_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
