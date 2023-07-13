PRAGMA foreign_keys = ON;


-- -----------------------------------------------------
-- Table `chat`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `chat` (
  `chat_id`          TEXT NOT NULL PRIMARY KEY,
  `telegram_chat_id` TEXT NOT NULL UNIQUE
);


-- -----------------------------------------------------
-- Table `event_type`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `event_type` (
  `event_type` TEXT NOT NULL PRIMARY KEY
);

INSERT INTO `event_type` (`event_type`)
VALUES
  ('birthday'),
  ('nameday'),
  ('wedding anniversary');


-- -----------------------------------------------------
-- Table `event`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `event` (
  `event_id`   TEXT NOT NULL PRIMARY KEY,
  `event_type` TEXT NOT NULL,
  `year`       INTEGER,
  `month`      INTEGER NOT NULL,
  `day`        INTEGER NOT NULL,
  FOREIGN KEY (`event_type`) REFERENCES `event_type` (`event_type`)
);


-- -----------------------------------------------------
-- Table `chat_event`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `chat_event` (
  `chat_id`  TEXT NOT NULL,
  `event_id` TEXT NOT NULL,
  FOREIGN KEY (`chat_id`) REFERENCES `chat` (`chat_id`),
  FOREIGN KEY (`event_id`) REFERENCES `event` (`event_id`)
);
CREATE UNIQUE INDEX IF NOT EXISTS `chat_id_event_id_UNIQUE`
ON `chat_event` (`chat_id`, `event_id`);


-- -----------------------------------------------------
-- Table `person`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `person` (
  `person_id` TEXT NOT NULL PRIMARY KEY
  `name`      TEXT NOT NULL,
  `surname`   TEXT
);


-- -----------------------------------------------------
-- Table `event_person`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `event_person` (
  `event_id`  TEXT NOT NULL,
  `person_id` TEXT NOT NULL,
  FOREIGN KEY (`event_id`) REFERENCES `event` (`event_id`),
  FOREIGN KEY (`person_id`) REFERENCES `person` (`person_id`)
);
CREATE UNIQUE INDEX IF NOT EXISTS `event_id_person_id_UNIQUE`
ON `event_person` (`event_id`, `person_id`);
