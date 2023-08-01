-- Run DDL scripts to DROP and CREATE tables
@ddl/gear_slot.sql
@ddl/player.sql
@ddl/player_gear.sql
@ddl/rarity.sql

-- Wrap DML scripts in a transaction to ensure atomicity
BEGIN;

-- Insert data using DML scripts
@dml/gear_slot.sql
@dml/player.sql
@dml/player_gear.sql
@dml/rarity.sql

-- Commit the transaction
COMMIT;