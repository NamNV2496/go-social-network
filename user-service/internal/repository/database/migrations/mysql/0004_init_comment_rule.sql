-- +migrate Up
START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS `network`;

CREATE TABLE IF NOT exists `comment_rule` (
    id int AUTO_INCREMENT NOT NULL PRIMARY KEY,
    comment_text varchar(255),
    created_at timestamp,
    updated_at timestamp);

COMMIT;

-- +migrate Down

