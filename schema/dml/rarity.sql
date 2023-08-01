INSERT IGNORE INTO `app`.`rarity` (`name`, `drop_probability_weight`, `lvl_requirement`, `item_type`, `enchantable`, `max_bonus_scaling`, `description`) VALUES
    ('Nothing', 28.90, 1, 'all', 0, 0, 'Literally nothing'),
    ('Junk', 20.23, 1, 'all', 0, 0, 'Remnants of lesser adventurers, offering little value and often abandoned on the wayside.'),
    ('Common', 14.45, 1, 'all', 0, 0, 'Reliable tools for every adventurer, standing strong on the journey ahead.'),
    ('Uncommon', 11.56, 1, 'all', 0, 0, 'Whispered about in hushed tones, possessing hidden potential waiting to be unlocked.'),
    ('Rare', 8.67, 10, 'all', 1, 10, 'Treasures sought by many, bestowing unique abilities and a testament to the adventurous spirit.'),
    ('Epic', 5.78, 25, 'all', 1, 25, 'Tales woven with each swing, granting powers capable of rewriting legends.'),
    ('Masterwork', 4.34, 40, 'all', 1, 40, 'Crafted with utmost skill, becoming the envy of artisans and the pride of champions.'),
    ('Legendary', 2.89, 60, 'all', 1, 60, 'The stuff of myth and wonder, carrying the weight of destinies and the glory of ages past.'),
    ('Mythic', 2.02, 80, 'all', 1, 75, 'Born of ancient lore, holding the secrets to untold power and the mysteries of forgotten realms.'),
    ('Transcendant', 1.16, 100, 'all', 1, 100, 'Beyond the grasp of ordinary mortals, bestowing godlike might and the ability to shape worlds.');
