DROP DATABASE IF EXISTS `labyrinth_of_babel_test`;

CREATE DATABASE `labyrinth_of_babel_test`;

USE `labyrinth_of_babel_test`;


CREATE TABLE `rooms` (
  `room` varchar(250) NOT NULL,
  PRIMARY KEY (`room`)
);

CREATE TABLE `sources` (
  `source` varchar(250) NOT NULL,
  PRIMARY KEY (`source`)
);

CREATE TABLE `cells` (
  `id` varchar(40) NOT NULL,
  `title` text NOT NULL,
  `body` longtext NOT NULL,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `room` varchar(250) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_cells_rooms1_idx` (`room`),
  CONSTRAINT `fk_cells_rooms1` FOREIGN KEY (`room`) REFERENCES `rooms` (`room`)
);

CREATE TABLE `cells_sources` (
  `cells_id` varchar(40) NOT NULL,
  `sources_source` varchar(250) NOT NULL,
  PRIMARY KEY (`cells_id`,`sources_source`),
  KEY `fk_cells_has_sources_sources1_idx` (`sources_source`),
  KEY `fk_cells_has_sources_cells1_idx` (`cells_id`),
  CONSTRAINT `fk_cells_has_sources_cells1` FOREIGN KEY (`cells_id`) REFERENCES `cells` (`id`),
  CONSTRAINT `fk_cells_has_sources_sources1` FOREIGN KEY (`sources_source`) REFERENCES `sources` (`source`)
);

CREATE TABLE `cells_links` (
  `cells_a` varchar(40) NOT NULL,
  `cells_b` varchar(40) NOT NULL,
  PRIMARY KEY (`cells_a`,`cells_b`),
  KEY `fk_cells_has_cells_cells2_idx` (`cells_b`),
  KEY `fk_cells_has_cells_cells1_idx` (`cells_a`),
  CONSTRAINT `fk_cells_has_cells_cells1` FOREIGN KEY (`cells_a`) REFERENCES `cells` (`id`),
  CONSTRAINT `fk_cells_has_cells_cells2` FOREIGN KEY (`cells_b`) REFERENCES `cells` (`id`)
);
  
INSERT INTO `rooms` VALUES ('This is a room');

INSERT INTO `sources` VALUES ('Analects'),('Confucius');

INSERT INTO `cells` VALUES ('417ecfe7-d2b4-4e43-afd4-dbf5f431d97d','Idea two','The second idea has a shorter body, but it\'s good enough.','2021-02-20 07:54:18','2021-02-20 07:54:18','This is a room'),('72aed05b-cb2d-4cad-bf70-05d8ae02a7bc','Idea one','Body of the first idea. Lengthy, useless, but interesting','2021-02-20 07:53:08','2021-02-20 07:53:08','This is a room'),('df38bd04-0ec4-41bf-9e53-d0eeb95a4939','','The third idea has no title, so that we can test what happens here','2021-02-20 07:55:36','2021-02-20 07:55:36','This is a room');

INSERT INTO `cells_sources` VALUES ('417ecfe7-d2b4-4e43-afd4-dbf5f431d97d','Confucius'),('72aed05b-cb2d-4cad-bf70-05d8ae02a7bc','Analects'),('72aed05b-cb2d-4cad-bf70-05d8ae02a7bc','Confucius'),('df38bd04-0ec4-41bf-9e53-d0eeb95a4939','Confucius');

INSERT INTO `cells_links` VALUES ('72aed05b-cb2d-4cad-bf70-05d8ae02a7bc','417ecfe7-d2b4-4e43-afd4-dbf5f431d97d'),('df38bd04-0ec4-41bf-9e53-d0eeb95a4939','72aed05b-cb2d-4cad-bf70-05d8ae02a7bc');
