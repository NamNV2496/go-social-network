-- +migrate Up
START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS `network`;

CREATE TABLE IF NOT exists `post` (
    id int AUTO_INCREMENT NOT NULL PRIMARY KEY,
    user_id varchar(10),
    content_text varchar(255),
    images varchar(255),
    tags varchar(100),
    visible boolean,
    created_at timestamp,
    updated_at timestamp);

COMMIT;

-- +migrate Down

