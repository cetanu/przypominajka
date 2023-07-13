PRAGMA foreign_keys = ON;


-- -----------------------------------------------------
-- Table `chat`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `chat` (
  `chat_id` TEXT NOT NULL PRIMARY KEY
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
  `event_id` TEXT NOT NULL PRIMARY KEY,
  `event_type` TEXT NOT NULL ,
  FOREIGN KEY (`event_type`) REFERENCES `event_type` (`event_type`)
);
