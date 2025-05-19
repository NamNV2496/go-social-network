-- +migrate Up
START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS `network`;

CREATE TABLE IF NOT exists `notification` (
    id int AUTO_INCREMENT NOT NULL PRIMARY KEY,
    title varchar(255),
    description varchar(255),
    template varchar(255),
    image varchar(255),
    application varchar(100),
    visible boolean,
    link varchar(255),
    created_at timestamp,
    updated_at timestamp);

COMMIT;

-- +migrate Down
