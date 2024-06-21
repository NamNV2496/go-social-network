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

CREATE TABLE IF NOT exists `user_user` (
    id int AUTO_INCREMENT NOT NULL PRIMARY KEY,
    user_id varchar(10),
    follower varchar(10),
    created_at timestamp,
    updated_at timestamp);

INSERT INTO `user` (id, email, `name`, picture, user_id, `password`, created_at, updated_at) VALUES
    (1,'hello@gmail.com', "nguyễn văn Nam", "namnv.png", "namnv", "$2a$10$vlcVb/AN1KbHVEWYltm4/ORdeY7xVVfK0SAEMGUQfEJLEwQoM0DDi", "2014-01-06 18:36:00", "2014-01-06 18:36:00"),
    (2,'hello2@gmail.com', "nguyễn văn bảo", "banbq.png", "baobq", "$2a$10$ku6te0wdnR7q8VZ23MbmxOc0TpkohOUOTzRf8Icy3zhUZPIUqZQeK", "2014-01-06 18:36:00", "2014-01-06 18:36:00"),
    (3,'hello3@gmail.com', "nguyễn mỹ khánh", "knm.png", "knm", "$2a$10$zMDwOUVxu/e90vkgQ2yXMuj/dBS38uFSWYbaW3iuMu.wAr0rAz8Zq", "2014-01-06 18:36:00", "2014-01-06 18:36:00");

INSERT INTO `user_user` (id, user_id, follower, created_at, updated_at) VALUES
    (1,'namnv', 'baobq', "2014-01-06 18:36:00", "2014-01-06 18:36:00"),
    (2,'namnv', 'knm', "2014-01-06 18:36:00", "2014-01-06 18:36:00"),
    (3,'knm', 'namnv', "2014-01-06 18:36:00", "2014-01-06 18:36:00"),
    (4,'baobq', 'namnv', "2014-01-06 18:36:00", "2014-01-06 18:36:00");

CREATE INDEX idx_follower
ON `user_user` (follower);

COMMIT;

-- +migrate Down

