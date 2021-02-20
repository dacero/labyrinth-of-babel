DROP DATABASE IF EXISTS `labyrinth_of_babel`;

CREATE DATABASE `labyrinth_of_babel`;

USE `labyrinth_of_babel`;

CREATE TABLE IF NOT EXISTS `labyrinth_of_babel`.`cells` (
  `id` VARCHAR(32) NOT NULL,
  `title` TEXT(250) NULL,
  `body` LONGTEXT NOT NULL,
  `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`));

CREATE TABLE IF NOT EXISTS `labyrinth_of_babel`.`topics` (
  `topic` VARCHAR(250) NOT NULL,
  PRIMARY KEY (`topic`));

CREATE TABLE IF NOT EXISTS `labyrinth_of_babel`.`sources` (
  `source` VARCHAR(250) NOT NULL,
  PRIMARY KEY (`source`));

CREATE TABLE IF NOT EXISTS `labyrinth_of_babel`.`cells_topics` (
  `topics_topic` VARCHAR(250) NOT NULL,
  `cells_id` VARCHAR(32) NOT NULL,
  PRIMARY KEY (`topics_topic`, `cells_id`),
  INDEX `fk_topics_has_cells_cells1_idx` (`cells_id` ASC) VISIBLE,
  INDEX `fk_topics_has_cells_topics_idx` (`topics_topic` ASC) VISIBLE,
  CONSTRAINT `fk_topics_has_cells_topics`
	FOREIGN KEY (`topics_topic`)
	REFERENCES `labyrinth_of_babel`.`topics` (`topic`)
	ON DELETE NO ACTION
	ON UPDATE NO ACTION,
  CONSTRAINT `fk_topics_has_cells_cells1`
	FOREIGN KEY (`cells_id`)
	REFERENCES `labyrinth_of_babel`.`cells` (`id`)
	ON DELETE NO ACTION
	ON UPDATE NO ACTION);

CREATE TABLE IF NOT EXISTS `labyrinth_of_babel`.`cells_sources` (
  `cells_id` VARCHAR(32) NOT NULL,
  `sources_source` VARCHAR(250) NOT NULL,
  PRIMARY KEY (`cells_id`, `sources_source`),
  INDEX `fk_cells_has_sources_sources1_idx` (`sources_source` ASC) VISIBLE,
  INDEX `fk_cells_has_sources_cells1_idx` (`cells_id` ASC) VISIBLE,
  CONSTRAINT `fk_cells_has_sources_cells1`
	FOREIGN KEY (`cells_id`)
	REFERENCES `labyrinth_of_babel`.`cells` (`id`)
	ON DELETE NO ACTION
	ON UPDATE NO ACTION,
  CONSTRAINT `fk_cells_has_sources_sources1`
	FOREIGN KEY (`sources_source`)
	REFERENCES `labyrinth_of_babel`.`sources` (`source`)
	ON DELETE NO ACTION
	ON UPDATE NO ACTION);

CREATE TABLE IF NOT EXISTS `labyrinth_of_babel`.`cells_links` (
  `cells_a` VARCHAR(32) NOT NULL,
  `cells_b` VARCHAR(32) NOT NULL,
  PRIMARY KEY (`cells_a`, `cells_b`),
  INDEX `fk_cells_has_cells_cells2_idx` (`cells_b` ASC) VISIBLE,
  INDEX `fk_cells_has_cells_cells1_idx` (`cells_a` ASC) VISIBLE,
  CONSTRAINT `fk_cells_has_cells_cells1`
	FOREIGN KEY (`cells_a`)
	REFERENCES `labyrinth_of_babel`.`cells` (`id`)
	ON DELETE NO ACTION
	ON UPDATE NO ACTION,
  CONSTRAINT `fk_cells_has_cells_cells2`
	FOREIGN KEY (`cells_b`)
	REFERENCES `labyrinth_of_babel`.`cells` (`id`)
	ON DELETE NO ACTION
	ON UPDATE NO ACTION);
	
INSERT INTO `sources` VALUES ('A single source');

INSERT INTO `topics` VALUES ('First topic');

INSERT INTO `cells` VALUES ('2213f29185094571a4750dbb24f225ec','The second idea','This is the body of the second idea, which is, obviously related to the first','2021-02-14 17:22:49','2021-02-14 17:22:49'),('b2020ced60d743c99464b90d8d2f3440','','This is the body of the first idea and if the length is over 60 characters this should be cut down into pieces and make it ok to not show there','2021-02-14 17:21:43','2021-02-14 17:21:43');

INSERT INTO `cells_links` VALUES ('2213f29185094571a4750dbb24f225ec','b2020ced60d743c99464b90d8d2f3440');

INSERT INTO `cells_sources` VALUES ('2213f29185094571a4750dbb24f225ec','A single source'),('b2020ced60d743c99464b90d8d2f3440','A single source');

INSERT INTO `cells_topics` VALUES ('First topic','2213f29185094571a4750dbb24f225ec'),('First topic','b2020ced60d743c99464b90d8d2f3440');
