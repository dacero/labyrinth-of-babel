CREATE DATABASE `labyrinth_of_babel`;

USE `labyrinth_of_babel`;


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
  
INSERT INTO `rooms` VALUES ('Labyrinth Wall');

INSERT INTO `cells` VALUES ('entry','Labyrinth Entry','Enter the labyrinth','2021-02-20 07:54:18','2021-02-20 07:54:18','Labyrinth Wall');
