-- 既存のテーブル削除
DROP DATABASE IF EXISTS q_pal_db;
CREATE DATABASE q_pal_db;
USE q_pal_db;


-- 既存のテーブル定義を削除
DROP TABLE IF EXISTS `todos`, `users`, `customers`, `companies`, `stores`, `favorite_stores`, `staffs`, `store_staff`, `active_staff`, `reservations`, `reservation_menus`, `visit_history`, `menus`, `menu_sets`, `menu_set_details`, `notifications`;


-- todos
CREATE TABLE `todos`
(
    `id`              INT         NOT NULL AUTO_INCREMENT,
    `title`           VARCHAR(45) NOT NULL,
    `completed`       TINYINT(1) NOT NULL DEFAULT 0,
    `userId`          INT         NOT NULL,
    `deleted`         TINYINT(1) NOT NULL DEFAULT 0,
    `created_at`      DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`      INT                  DEFAULT 0,
    `created_by_type` VARCHAR(45) NOT NULL DEFAULT 'system',
    `updated_by`      INT                  DEFAULT 0,
    `updated_by_type` VARCHAR(45) NOT NULL DEFAULT 'system',
    PRIMARY KEY (`id`)
);

-- users
CREATE TABLE `users`
(
    `id`              INT         NOT NULL AUTO_INCREMENT,
    `name`            VARCHAR(45) NOT NULL,
    `created_at`      DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`      INT                  DEFAULT 0,
    `created_by_type` VARCHAR(45) NOT NULL DEFAULT 'system',
    `updated_by`      INT                  DEFAULT 0,
    `updated_by_type` VARCHAR(45) NOT NULL DEFAULT 'system',
    PRIMARY KEY (`id`)
);

-- サンプルデータ投入
INSERT INTO `todos` (`id`, `title`, `completed`, `userId`, `deleted`, `created_at`, `updated_at`, `created_by`, `created_by_type`, `updated_by`, `updated_by_type`)
VALUES ('1', 'テスト1', '0', '1', '0', '2022-01-01', '2022-01-01', 0, 'system', 0, 'system'),
       ('2', 'テスト2', '0', '1', '0', '2022-01-01', '2022-01-01', 0, 'system', 0, 'system'),
       ('3', 'テスト3', '0', '1', '0', '2022-01-01', '2022-01-01', 0, 'system', 0, 'system');


SET
FOREIGN_KEY_CHECKS=0;


-- customers
CREATE TABLE `customers`
(
    `customer_id`     INT PRIMARY KEY,
    `cognito_user_id` VARCHAR(255),
    `name`            VARCHAR(255),
    `email`           VARCHAR(255),
    `gender`          INT,
    `birthday`        DATE,
    `created_at`             DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`             DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`             INT DEFAULT 0,
    `created_by_type`        VARCHAR(255) NOT NULL DEFAULT 'system',
    `updated_by`             INT DEFAULT 0,
    `updated_by_type`        VARCHAR(255) NOT NULL DEFAULT 'system'
);

-- companies
CREATE TABLE `companies`
(
    `company_id`      INT PRIMARY KEY,
    `company_name`    VARCHAR(255),
    `created_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`      INT                   DEFAULT 0,
    `created_by_type` VARCHAR(255) NOT NULL DEFAULT 'system',
    `updated_by`      INT                   DEFAULT 0,
    `updated_by_type` VARCHAR(255) NOT NULL DEFAULT 'system'
);

-- stores
CREATE TABLE `stores`
(
    `store_id`        INT PRIMARY KEY,
    `company_id`      INT,
    `store_name`      VARCHAR(255),
    `address`         VARCHAR(255),
    `phone_number`    VARCHAR(20),
    `business_hours`  VARCHAR(255),
    `created_at`      DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`      INT                  DEFAULT 0,
    `created_by_type` VARCHAR(45) NOT NULL DEFAULT 'system',
    `updated_by`      INT                  DEFAULT 0,
    `updated_by_type` VARCHAR(45) NOT NULL DEFAULT 'system',
    FOREIGN KEY (`company_id`) REFERENCES `companies` (`company_id`)
);

-- favorite_stores
CREATE TABLE `favorite_stores`
(
    `customer_id`     INT,
    `store_id`        INT,
    `created_at`      DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`      INT                  DEFAULT 0,
    `created_by_type` VARCHAR(45) NOT NULL DEFAULT 'system',
    `updated_by`      INT                  DEFAULT 0,
    `updated_by_type` VARCHAR(45) NOT NULL DEFAULT 'system',
    PRIMARY KEY (`customer_id`, `store_id`),
    FOREIGN KEY (`customer_id`) REFERENCES `customers` (`customer_id`),
    FOREIGN KEY (`store_id`) REFERENCES `stores` (`store_id`)
);

-- staffs
CREATE TABLE `staffs`
(
    `staff_id`        INT PRIMARY KEY,
    `name`            VARCHAR(255),
    `cognito_user_id` VARCHAR(255),
    `created_at`             DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`             DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`             INT DEFAULT 0,
    `created_by_type`        VARCHAR(255) NOT NULL DEFAULT 'system',
    `updated_by`             INT DEFAULT 0,
    `updated_by_type`        VARCHAR(255) NOT NULL DEFAULT 'system'
);

-- store_staffs
CREATE TABLE `store_staffs`
(
    `staff_id`        INT,
    `store_id`        INT,
    `created_at`      DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`      INT                  DEFAULT 0,
    `created_by_type` VARCHAR(45) NOT NULL DEFAULT 'system',
    `updated_by`      INT                  DEFAULT 0,
    `updated_by_type` VARCHAR(45) NOT NULL DEFAULT 'system',
    PRIMARY KEY (`staff_id`, `store_id`),
    FOREIGN KEY (`staff_id`) REFERENCES `staffs` (`staff_id`),
    FOREIGN KEY (`store_id`) REFERENCES `stores` (`store_id`)
);

-- active_staffs
CREATE TABLE `active_staffs`
(
    `staff_id`             INT PRIMARY KEY,
    `store_id`             INT          NOT NULL,
    `break_start_datetime` DATETIME,
    `break_end_datetime`   DATETIME,
    `created_at`           DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`           DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`           INT                   DEFAULT 0,
    `created_by_type`      VARCHAR(255) NOT NULL DEFAULT 'system',
    `updated_by`           INT                   DEFAULT 0,
    `updated_by_type`      VARCHAR(255) NOT NULL DEFAULT 'system',
    FOREIGN KEY (`staff_id`) REFERENCES `staffs` (`staff_id`),
    FOREIGN KEY (`store_id`) REFERENCES `stores` (`store_id`)
);

-- reservations
CREATE TABLE `reservations`
(
    `reservation_id`         INT PRIMARY KEY AUTO_INCREMENT,
    `customer_id`            INT          NOT NULL,
    `store_id`               INT          NOT NULL,
    `staff_id`               INT,
    `reservation_number`     INT          NOT NULL,
    `reserved_datetime`      DATETIME     NOT NULL,
    `hold_start_datetime`    DATETIME,
    `service_start_datetime` DATETIME,
    `service_end_datetime`   DATETIME,
    `status`                 INT          NOT NULL DEFAULT 0,
    `arrival_flag`           BOOLEAN      NOT NULL DEFAULT 0,
    `cancel_type`            INT,
    `created_at`             DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`             DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`             INT                   DEFAULT 0,
    `created_by_type`        VARCHAR(255) NOT NULL DEFAULT 'system',
    `updated_by`             INT                   DEFAULT 0,
    `updated_by_type`        VARCHAR(255) NOT NULL DEFAULT 'system',
    FOREIGN KEY (`customer_id`) REFERENCES `customers` (`customer_id`),
    FOREIGN KEY (`store_id`) REFERENCES `stores` (`store_id`),
    FOREIGN KEY (`staff_id`) REFERENCES `staffs` (`staff_id`)
);

-- reservation_menus
CREATE TABLE `reservation_menus`
(
    `reservation_id`  INT,
    `store_id`        INT,
    `menu_id`         INT,
    `created_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`      INT                   DEFAULT 0,
    `created_by_type` VARCHAR(255) NOT NULL DEFAULT 'system',
    `updated_by`      INT                   DEFAULT 0,
    `updated_by_type` VARCHAR(255) NOT NULL DEFAULT 'system',
    PRIMARY KEY (`reservation_id`, `store_id`, `menu_id`),
    FOREIGN KEY (`store_id`) REFERENCES `stores` (`store_id`),
    FOREIGN KEY (`reservation_id`) REFERENCES `reservations` (`reservation_id`)
);

-- visit_history
CREATE TABLE `visit_history`
(
    `visit_history_id` INT PRIMARY KEY,
    `reservation_id`   INT,
    `menu_name`        VARCHAR(255),
    `price`            DECIMAL(10, 2),
    `created_at`       DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`       DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`       INT                   DEFAULT 0,
    `created_by_type`  VARCHAR(255) NOT NULL DEFAULT 'system',
    `updated_by`       INT                   DEFAULT 0,
    `updated_by_type`  VARCHAR(255) NOT NULL DEFAULT 'system',
    FOREIGN KEY (`reservation_id`) REFERENCES `reservations` (`reservation_id`)
);

-- menus
CREATE TABLE `menus`
(
    `store_id`        INT,
    `menu_id`         INT,
    `menu_name`       VARCHAR(255),
    `price`           DECIMAL(10, 2),
    `time`            INT          NOT NULL DEFAULT 0,
    `created_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`      INT                   DEFAULT 0,
    `created_by_type` VARCHAR(255) NOT NULL DEFAULT 'system',
    `updated_by`      INT                   DEFAULT 0,
    `updated_by_type` VARCHAR(255) NOT NULL DEFAULT 'system',
    PRIMARY KEY (`store_id`, `menu_id`),
    FOREIGN KEY (`store_id`) REFERENCES `stores` (`store_id`)
);

-- menu_sets
CREATE TABLE `menu_sets`
(
    `set_id`          INT,
    `store_id`        INT,
    `set_name`        VARCHAR(255),
    `set_price`       DECIMAL(10, 2),
    `created_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`      INT                   DEFAULT 0,
    `created_by_type` VARCHAR(255) NOT NULL DEFAULT 'system',
    `updated_by`      INT                   DEFAULT 0,
    `updated_by_type` VARCHAR(255) NOT NULL DEFAULT 'system',
    PRIMARY KEY (`set_id`),
    FOREIGN KEY (`store_id`) REFERENCES `stores` (`store_id`)
);

-- menu_set_details
CREATE TABLE `menu_set_details`
(
    `set_id`   INT,
    `store_id` INT,
    `menu_id`  INT,
    `created_at`             DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`             DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`             INT DEFAULT 0,
    `created_by_type`        VARCHAR(255) NOT NULL DEFAULT 'system',
    `updated_by`             INT DEFAULT 0,
    `updated_by_type`        VARCHAR(255) NOT NULL DEFAULT 'system',
    PRIMARY KEY ( `set_id`, `store_id`, `menu_id`),
    FOREIGN KEY (`set_id`) REFERENCES `menu_sets` (`set_id`),
    FOREIGN KEY (`store_id`) REFERENCES `stores` (`store_id`)
);

-- notifications
CREATE TABLE `notifications`
(
    `notification_id`      INT PRIMARY KEY,
    `customer_id`          INT,
    `reservation_id`       INT,
    `notification_type`    INT,
    `notification_content` TEXT,
    `notification_status`  INT,
    `created_at`           DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`           DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by`           INT                   DEFAULT 0,
    `created_by_type`      VARCHAR(255) NOT NULL DEFAULT 'system',
    `updated_by`           INT                   DEFAULT 0,
    `updated_by_type`      VARCHAR(255) NOT NULL DEFAULT 'system',
    FOREIGN KEY (`customer_id`) REFERENCES `customers` (`customer_id`),
    FOREIGN KEY (`reservation_id`) REFERENCES `reservations` (`reservation_id`)
);


-- companiesテーブルにレコードを2つ挿入
INSERT INTO `companies` (`company_id`, `company_name`)
VALUES (1, 'Company A'),
       (2, 'Company B');

-- storesテーブルにレコードを2つ挿入
INSERT INTO `stores` (`store_id`, `company_id`, `store_name`, `address`, `phone_number`, `business_hours`)
VALUES (1, 1, 'store_ A1', 'Address A1', '123-456-7890', '9:00-18:00'),
       (2, 2, 'store_ B1', 'Address B1', '098-765-4321', '10:00-19:00');

-- customersテーブルにレコードを5つ挿入
INSERT INTO `customers` (`customer_id`, `cognito_user_id`, `name`, `email`, `gender`, `birthday`)
VALUES (0, null, 'guest', null, null, null),
       (1, 'cognito1', 'Customer A', 'customerA@example.com', 1, '1990-01-01'),
       (2, 'cognito2', 'Customer B', 'customerB@example.com', 2, '1991-02-02'),
       (3, 'cognito3', 'Customer C', 'customerC@example.com', 1, '1992-03-03'),
       (4, 'cognito4', 'Customer D', 'customerD@example.com', 2, '1993-04-04'),
       (5, 'cognito5', 'Customer E', 'customerE@example.com', 1, '1994-05-05'),
       (6, 'cognito6', 'Customer F', 'customerF@example.com', 2, '1995-06-06'),
       (7, 'cognito7', 'Customer G', 'customerG@example.com', 1, '1996-07-07'),
       (8, 'cognito8', 'Customer H', 'customerH@example.com', 2, '1997-08-08'),
       (9, 'cognito9', 'Customer I', 'customerI@example.com', 1, '1998-09-09'),
       (10, 'cognito10', 'Customer J', 'customerJ@example.com', 2, '1999-10-10');


-- reservationsテーブルにレコードを5つ挿入
INSERT INTO `reservations` (`reservation_id`, `customer_id`, `store_id`, `staff_id`, `reservation_number`, `reserved_datetime`, `service_start_datetime`, `service_end_datetime`, `status`, `arrival_flag`, `cancel_type`)
VALUES (1, 1, 1, null, 101, '2023-01-01 10:00:00', '2023-01-01 10:00:00', '2023-01-01 11:00:00', 2, false, null),
       (2, 2, 1, null, 102, '2023-01-02 10:00:00', '2023-01-02 10:00:00', '2023-01-02 11:00:00', 2, false, null),

       (3, 3, 2, null, 101, '2023-01-03 10:00:00', '2023-01-03 10:00:00', '2023-01-03 10:10:00', 2, false, null),
       (4, 4, 2, null, 102, '2023-01-03 10:00:00', '2023-01-03 10:00:00', '2023-01-03 10:11:00', 2, false, null),
       (5, 5, 2, null, 103, '2023-01-03 10:00:00', '2023-01-03 10:00:00', '2023-01-03 10:12:00', 2, false, null),

       (6, 3, 2, null, 104, '2023-02-03 10:00:00', '2023-02-03 10:00:00', '2023-02-03 10:12:00', 2, false, null),
       (7, 4, 2, null, 105, '2023-02-03 10:00:00', '2023-02-03 10:00:00', '2023-02-03 10:13:00', 2, false, null),
       (8, 5, 2, null, 106, '2023-02-03 10:00:00', '2023-02-03 10:00:00', '2023-02-03 10:14:00', 2, false, null),

       (9, 3, 2, null, 107, '2023-03-03 10:00:00', '2023-03-03 10:00:00', '2023-03-03 10:14:00', 2, false, null),
       (10, 4, 2, null, 108, '2023-03-03 10:00:00', '2023-03-03 10:00:00', '2023-03-03 10:15:00', 2, false, null),
       (11, 5, 2, null, 109, '2023-03-03 10:00:00', '2023-03-03 10:00:00', '2023-03-03 10:16:00', 2, false, null),

       (12, 3, 2, 3, 101, now(), now(), null, 1, false, null),
       (13, 4, 2, 4, 102, now(), now(), null, 1, false, null),
       (14, 5, 2, null, 103, now(), null, null, 0, false, null),
       (15, 6, 2, null, 106, now(), null, null, 0, false, null),
       (16, 7, 2, null, 107, now(), null, null, 0, false, null),
       (17, 8, 2, null, 108, now(), null, null, 0, false, null),
       (18, 9, 2, null, 109, now(), null, null, 0, false, null),
       (19, 10, 2, null, 110, now(), null, null, 0, false, null);


INSERT INTO `menus` (`store_id`, `menu_id`, `menu_name`, `price`, `time`)
VALUES (1, 1, 'Menu 1', 1000, 15),
       (2, 1, 'Menu 2', 1500, 15);


INSERT INTO `reservation_menus` (`store_id`, `reservation_id`, `menu_id`)
VALUES (2, 3, 1),
       (2, 4, 1),
       (2, 5, 1),
       (2, 6, 1),
       (2, 7, 1),
       (2, 8, 1),
       (2, 9, 1),
       (2, 10, 1),
       (2, 11, 1),
       (2, 12, 1),
       (2, 13, 1),
       (2, 14, 1),
       (2, 15, 1),
       (2, 16, 1),
       (2, 17, 1),
       (2, 18, 1),
       (2, 19, 1);

-- staffsテーブルへのデータ投入
INSERT INTO `staffs` (`staff_id`, `name`, `cognito_user_id`)
VALUES (1, 'Staff A', 'cognitoA'),
       (2, 'Staff B', 'cognitoB'),
       (3, 'Staff C', 'cognitoC'),
       (4, 'Staff D', 'cognitoD');

-- store_staffテーブルへのデータ投入
INSERT INTO `store_staffs` (`staff_id`, `store_id`)
VALUES (1, 1),
       (2, 1),
       (3, 2),
       (4, 2);


INSERT INTO `active_staffs` (`staff_id`, `store_id`, `break_start_datetime`, `break_end_datetime`)
VALUES (1, 1, null, null),
       (2, 1, null, null),
       (3, 2, null, null),
       (4, 2, null, null);

