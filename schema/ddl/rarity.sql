CREATE TABLE IF NOT EXISTS `app`.`rarity` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'pk',
    `name` varchar(255) NOT NULL DEFAULT '' COMMENT 'Name of the rarity tier',
    `drop_probability_weight` int(11) DEFAULT 1 COMMENT 'Drop probability weight for random loot generation',
    `lvl_requirement` int(11) DEFAULT 1 COMMENT 'Minimum player level requirement to be eligible for this rarity loot',
    `item_type` varchar(50) DEFAULT '' COMMENT 'Type of item associated with this rarity (e.g., weapon, armor, consumable)',
    `enchantable` tinyint(1) DEFAULT 0 COMMENT 'Whether the item can be enchanted (1 for yes, 0 for no)',
    `max_bonus_scaling` int(11) DEFAULT 0 COMMENT 'Maximum bonus scaling with the rarity tier',
    `description` TEXT NOT NULL DEFAULT '' COMMENT 'Description of the rarity',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;