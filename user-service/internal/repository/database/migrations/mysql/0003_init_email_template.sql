-- +migrate Up
START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS `network`;

CREATE TABLE IF NOT exists `email_template` (
    id int AUTO_INCREMENT NOT NULL PRIMARY KEY,
    template_id varchar(50),
    template varchar(255),
    created_at timestamp,
    updated_at timestamp);

INSERT INTO `email_template` (id, template_id, template, created_at, updated_at) VALUES
(1,'otp_email', "<b>Hi {{.full_name}} this is your OTP: {{.otp}}, please don't publish this OTP to another people", "2014-01-06 18:36:00", "2014-01-06 18:36:00"),
(2,'register_email', "<b>Hi {{.full_name}} Thank for yout attendion to us. Click here {{.deep_link}} to done it.", "2014-01-06 18:36:00", "2014-01-06 18:36:00");
    
COMMIT;

-- +migrate Down

