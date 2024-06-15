-- +migrate Up
START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS `user`;

CREATE TABLE IF NOT exists `user` (
    id int AUTO_INCREMENT NOT NULL PRIMARY KEY,
    email varchar(50),
    `name` varchar(50),
    picture varchar(50),
    user_id varchar(10) UNIQUE,
    `password` varchar(100),
    created_at timestamp,
    updated_at timestamp);

INSERT INTO `user` (id, email, `name`, picture, user_id, `password`, created_at, updated_at) VALUES
    (1,'hello@gmail.com', "nguyễn văn Nam", "namnv.png", "namnv", "$2a$10$vlcVb/AN1KbHVEWYltm4/ORdeY7xVVfK0SAEMGUQfEJLEwQoM0DDi", "2014-01-06 18:36:00", "2014-01-06 18:36:00"),
    (2,'hello2@gmail.com', "nguyễn văn bảo", "banbq.png", "banbq", "$2a$10$ku6te0wdnR7q8VZ23MbmxOc0TpkohOUOTzRf8Icy3zhUZPIUqZQeK", "2014-01-06 18:36:00", "2014-01-06 18:36:00");

COMMIT;

-- +migrate Down

