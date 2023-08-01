CREATE TABLE IF NOT EXISTS `app`.`player` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'pk',
    `name` varchar(255) NOT NULL DEFAULT '',
    `is_active` TINYINT(1) DEFAULT 1,
    `created` int(11) NOT NULL DEFAULT 0,
    `updated` int(11) NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
