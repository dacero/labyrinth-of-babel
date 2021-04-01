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
  
INSERT INTO `rooms` VALUES ('Labyrinth Wall'),('README');

INSERT INTO `cells` VALUES ('08b3c476-7f64-492c-be15-c9447bd67138','Rooms','The cells of the Labyrinth of Babel are grouped in rooms.\r\n\r\nA cell can only belong to one room.\r\n\r\nThe architect of the Labyrinth can decide how to organize those rooms, so that similar cells stay with each other.\r\n\r\nThis cell belongs to the README room.\r\n\r\nYou can enter a room to check all cells that are contained in it.',now(),now(),'README'),('16237ec7-e812-4d23-8650-4a070bd45018','Cells','Each cell represents an idea, a concept, or any other form of annotation.\r\n\r\nThey are represented as cards to incentivize brevity and conciseness, but there’s no explicit limit to them\r\n\r\nCells are written in [markdown](https://en.wikipedia.org/wiki/Markdown). They can contain links, images, videos or anything really that a web page accepts.',now(),now(),'README'),('2e79f5da-a402-4058-8a24-f28e0f38507f','README','Labyrinth of Babel is a space of connected ideas.\r\n\r\nEvery cell of this Labyrinth represents an idea.\r\n\r\nCells are connected with each other creating a labyrinth of passageways.\r\n\r\nThe Labyrinth of Babel grows as new ideas arrive and unexpected connections are found.',now(),now(),'README'),('4d25d822-5319-4279-b46a-22a753f746b0','Sources','Most cells in the Labyrinth of Babel stem from somewhere outside the Labyrinth: a book, an author, a lecture… the source section of a cell captures its origin.\r\n\r\nSome cells have more than one source.\r\n\r\nSome cells have no explicit source, like this one.\r\n\r\nThe Labyrinth allows you to see all cells that came out of the same source.',now(),now(),'README'),('aa431416-6e06-44d9-8bfd-f1a7a63cbf8c','Links','The cells in the Labyrinth of Babel can connect with each other.\r\n\r\nThese are the links below each cell. They are doors between ideas.\r\n\r\nWhen a door is created between cell A and cell B, that same door works in the opposite direction, from B to A. These are called backlinks.',now(),now(),'README'),('b3bf792e-372d-4386-95f2-c8310a2aa02b','References','The Labyrinth of Babel is inspired by the [zettlekasten](https://en.wikipedia.org/wiki/Zettelkasten) method of knowledge management.\r\n\r\nIt’s a simplified version of my beloved platform [Roam Research](https://roamresearch.com). The Labyrinth of Babel has only fundamental features and is more opinionated in terms of its structure. It allows for a more beautiful visualization.\r\n\r\nIt aspires to become a curated display of the inner world of ideas and knowledge we all develop over time.\r\n\r\nIt’s name and logo are inspired by the short story _Library of Babel by Jorge Luis Borges_.',now(),now(),'README'),('entry','Labyrinth Entry','![](http://deliris.net/thoughts/labyrinth/images/labyrinth.jpg)\r\n\r\n_“Like all men of the Labyrinth, in my younger days I travelled.”_\r\n\r\n- [All of the Labyrinth\'s rooms](/rooms)\r\n- [The Labyrinth’s Wall](/room/Labyrinth%20Wall)\r\n- [Create one new cell](/page/new_card.html)',now(),now(),'Labyrinth Wall');

INSERT INTO `cells_links` VALUES ('08b3c476-7f64-492c-be15-c9447bd67138','16237ec7-e812-4d23-8650-4a070bd45018'),('4d25d822-5319-4279-b46a-22a753f746b0','16237ec7-e812-4d23-8650-4a070bd45018'),('aa431416-6e06-44d9-8bfd-f1a7a63cbf8c','16237ec7-e812-4d23-8650-4a070bd45018'),('16237ec7-e812-4d23-8650-4a070bd45018','2e79f5da-a402-4058-8a24-f28e0f38507f'),('b3bf792e-372d-4386-95f2-c8310a2aa02b','2e79f5da-a402-4058-8a24-f28e0f38507f'),('2e79f5da-a402-4058-8a24-f28e0f38507f','entry');

