-- Run DDL scripts to DROP and CREATE tables
source ddl/gear_slot.sql
source ddl/player.sql
source ddl/player_gear.sql
source ddl/rarity.sql

-- Wrap DML scripts in a transaction to ensure atomicity
BEGIN;

-- Insert data using DML scripts
source dml/gear_slot.sql
source dml/player.sql
source dml/rarity.sql

-- Commit the transaction
COMMIT;