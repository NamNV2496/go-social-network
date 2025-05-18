-- +migrate Up
START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS `network`;

CREATE TABLE IF NOT exists `comment` (
    id int AUTO_INCREMENT NOT NULL PRIMARY KEY,
    post_id int,
    user_id varchar(10),
    comment_text varchar(255),
    comment_level int,
    comment_parent int,
    images varchar(255),
    tags varchar(100),
    created_at timestamp,
    updated_at timestamp);

CREATE TABLE IF NOT exists `like` (
    id int AUTO_INCREMENT NOT NULL PRIMARY KEY,
    post_id int,
    user_id varchar(10),
    `like` boolean,
    created_at timestamp,
    updated_at timestamp);

CREATE TABLE IF NOT exists `like_count` (
    id int AUTO_INCREMENT NOT NULL PRIMARY KEY,
    post_id int,
    total_like int,
    created_at timestamp,
    updated_at timestamp);

CREATE INDEX idx_comment_postId
ON `comment` (post_id);

CREATE INDEX idx_like_postId
ON `like` (post_id);

COMMIT;

-- +migrate Down

