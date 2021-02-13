DROP DATABASE IF EXISTS `hello`;

CREATE DATABASE `hello`;

USE `hello`;

CREATE TABLE `names` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `names` VALUES (4,'This is the app running in 2 containers!');