DROP DATABASE IF EXISTS q_mate_db;
CREATE DATABASE q_mate_db;
USE q_mate_db;

-- todos
DROP TABLE IF EXISTS `todos`;
CREATE TABLE `todos`
(
    `id`         INT         NOT NULL AUTO_INCREMENT,
    `title`      VARCHAR(45) NOT NULL,
    `completed`  TINYINT(1)  NOT NULL DEFAULT 0,
    `userId`     INT         NOT NULL,
    `deleted`    TINYINT(1)  NOT NULL DEFAULT 0,
    `created_at` DATETIME    NOT NULL,
    PRIMARY KEY (`id`)
);
INSERT INTO `todos` (`id`, `title`, `completed`, `userId`, `deleted`, `created_at`)
VALUES ('1', 'テスト1', '0', '1', '0', '2022-01-01');
INSERT INTO `todos` (`id`, `title`, `completed`, `userId`, `deleted`, `created_at`)
VALUES ('2', 'テスト2', '0', '1', '0', '2022-01-01');
INSERT INTO `todos` (`id`, `title`, `completed`, `userId`, `deleted`, `created_at`)
VALUES ('3', 'テスト3', '0', '1', '0', '2022-01-01');

-- users
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `id`   INT         NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(45) NOT NULL,
    PRIMARY KEY (`id`)
);
INSERT INTO `users` (`id`, `name`)
VALUES ('1', '山田');

-- customers
DROP TABLE IF EXISTS `customers`;
CREATE TABLE `customers`
(
    `customer_id`     INT PRIMARY KEY,
    `cognito_user_id` VARCHAR(255),
    `name`            VARCHAR(255),
    `email`           VARCHAR(255),
    `gender`          INT,
    `birthday`        DATE
);

-- companies
DROP TABLE IF EXISTS `companies`;
CREATE TABLE `companies`
(
    `company_id`   INT PRIMARY KEY,
    `company_name` VARCHAR(255)
);

-- stores
DROP TABLE IF EXISTS `stores`;
CREATE TABLE `stores`
(
    `store_id`       INT PRIMARY KEY,
    `company_id`     INT,
    `store_name`     VARCHAR(255),
    `address`        VARCHAR(255),
    `phone_number`   VARCHAR(20),
    `business_hours` VARCHAR(255),
    FOREIGN KEY (`company_id`) REFERENCES `companies` (`company_id`)
);

-- favorite_stores
DROP TABLE IF EXISTS `favorite_stores`;
CREATE TABLE `favorite_stores`
(
    `customer_id`  INT,
    `store_id` INT,
    PRIMARY KEY (`customer_id`, `store_id`),
    FOREIGN KEY (`customer_id`) REFERENCES `customers` (`customer_id`),
    FOREIGN KEY (`store_id`) REFERENCES `stores` (`store_id`)
);

-- staff
DROP TABLE IF EXISTS `staff`;
CREATE TABLE `staff`
(
    `staff_id`        INT PRIMARY KEY,
    `name`            VARCHAR(255),
    `cognito_user_id` VARCHAR(255)
);

-- store_staff
DROP TABLE IF EXISTS `store_staff`;
CREATE TABLE `store_staff`
(
    `staff_id` INT,
    `store_id` INT,
    PRIMARY KEY (`staff_id`, `store_id`),
    FOREIGN KEY (`staff_id`) REFERENCES `staff` (`staff_id`),
    FOREIGN KEY (`store_id`) REFERENCES `stores` (`store_id`)
);

-- active_staff
DROP TABLE IF EXISTS `active_staff`;
CREATE TABLE `active_staff`
(
    `staff_id`             INT PRIMARY KEY,
    `store_id`             INT,
    `break_start_datetime` DATETIME,
    `break_end_datetime`   DATETIME,
    FOREIGN KEY (`staff_id`) REFERENCES `staff` (`staff_id`),
    FOREIGN KEY (`store_id`) REFERENCES `stores` (`store_id`)
);

-- reservations
DROP TABLE IF EXISTS `reservations`;
CREATE TABLE `reservations`
(
    `reservation_id`         INT PRIMARY KEY,
    `customer_id`            INT NOT NULL,
    `store_id`               INT NOT NULL,
    `staff_id`               INT,
    `reservation_number`     INT NOT NULL,
    `reserved_datetime`      DATETIME NOT NULL,
    `hold_start_datetime`    DATETIME,
    `service_start_datetime` DATETIME,
    `service_end_datetime`   DATETIME,
    `status`                 INT,
    `arrival_flag`           BOOLEAN NOT NULL DEFAULT 0,
    `cancel_flag`            BOOLEAN NOT NULL DEFAULT 0,
    `cancel_type`            INT,
    FOREIGN KEY (`customer_id`) REFERENCES `customers` (`customer_id`),
    FOREIGN KEY (`store_id`) REFERENCES `stores` (`store_id`),
    FOREIGN KEY (`staff_id`) REFERENCES `staff` (`staff_id`)
);

-- reservation_menus
DROP TABLE IF EXISTS `reservation_menus`;
CREATE TABLE `reservation_menus`
(
    `reservation_id` INT,
    `menu_id`        INT,
    PRIMARY KEY (`reservation_id`, `menu_id`),
    FOREIGN KEY (`reservation_id`) REFERENCES `reservations` (`reservation_id`)
);

-- visit_history
DROP TABLE IF EXISTS `visit_history`;
CREATE TABLE `visit_history`
(
    `visit_history_id` INT PRIMARY KEY,
    `reservation_id`   INT,
    `menu_name`        VARCHAR(255),
    `price`            DECIMAL(10, 2),
    FOREIGN KEY (`reservation_id`) REFERENCES `reservations` (`reservation_id`)
);

-- menus
DROP TABLE IF EXISTS `menus`;
CREATE TABLE `menus`
(
    `menu_id`   INT PRIMARY KEY,
    `store_id`  INT,
    `menu_name` VARCHAR(255),
    `price`     DECIMAL(10, 2),
    `time`      INT,
    FOREIGN KEY (`store_id`) REFERENCES `stores` (`store_id`)
);

-- menu_sets
DROP TABLE IF EXISTS `menu_sets`;
CREATE TABLE `menu_sets`
(
    `set_id`    INT PRIMARY KEY,
    `set_name`  VARCHAR(255),
    `set_price` DECIMAL(10, 2)
);

-- menu_set_details
DROP TABLE IF EXISTS `menu_set_details`;
CREATE TABLE `menu_set_details`
(
    `set_id`  INT,
    `menu_id` INT,
    PRIMARY KEY (`set_id`, `menu_id`),
    FOREIGN KEY (`set_id`) REFERENCES `menu_sets` (`set_id`),
    FOREIGN KEY (`menu_id`) REFERENCES `menus` (`menu_id`)
);

-- menu_set_details
DROP TABLE IF EXISTS `notifications`;
CREATE TABLE `notifications`
(
    `notification_id`      INT PRIMARY KEY,
    `customer_id`              INT,
    `reservation_id`       INT,
    `notification_type`    INT,
    `notification_content` TEXT,
    `notification_status`  INT,
    FOREIGN KEY (`customer_id`) REFERENCES `customers` (`customer_id`),
    FOREIGN KEY (`reservation_id`) REFERENCES `reservations` (`reservation_id`)
);



-- companiesテーブルにレコードを2つ挿入
INSERT INTO `companies` (`company_id`, `company_name`)
VALUES
    (1, 'Company A'),
    (2, 'Company B');

-- storesテーブルにレコードを2つ挿入
INSERT INTO `stores` (`store_id`, `company_id`, `store_name`, `address`, `phone_number`, `business_hours`)
VALUES
    (1, 1, 'store_ A1', 'Address A1', '123-456-7890', '9:00-18:00'),
    (2, 2, 'store_ B1', 'Address B1', '098-765-4321', '10:00-19:00');

-- customersテーブルにレコードを5つ挿入
INSERT INTO `customers` (`customer_id`, `cognito_user_id`, `name`, `email`, `gender`, `birthday`)
VALUES
    (1, 'cognito1', 'Customer A', 'customerA@example.com', 1, '1990-01-01'),
    (2, 'cognito2', 'Customer B', 'customerB@example.com', 2, '1991-02-02'),
    (3, 'cognito3', 'Customer C', 'customerC@example.com', 1, '1992-03-03'),
    (4, 'cognito4', 'Customer D', 'customerD@example.com', 2, '1993-04-04'),
    (5, 'cognito5', 'Customer E', 'customerE@example.com', 1, '1994-05-05');


-- reservationsテーブルにレコードを5つ挿入
-- ここでは、user_id, store_id, staff_idは適当な値を設定しています。実際には、存在するCustomer, stores, staffのIDを設定してください。
INSERT INTO `reservations` (`reservation_id`, `customer_id`, `store_id`, `staff_id`, `reservation_number`, `reserved_datetime`, `status`, `arrival_flag`, `cancel_type`)
VALUES
    (1, 1, 1, null, 101, now(), 1, false, null),
    (2, 2, 1, null, 102, now(), 1, false, null),
    (3, 3, 2, null, 103, now(), 1, false, null),
    (4, 4, 2, null, 104, now(), 1, false, null),
    (5, 5, 2, null, 105, now(), 1, false, null);

