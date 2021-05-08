-- enable foregin keys
PRAGMA foreign_keys = ON;

CREATE TABLE `boards` (
    `uuid` TEXT PRIMARY KEY,
    `url` TEXT UNIQUE,
    `descriptor` TEXT DEFAULT 'empty'
);


CREATE TABLE `threads`(
    `id` UNSIGNED BIG INT PRIMARY KEY,
    `date` UNSIGNED BIG INT,
    `file` TEXT,
    `responses` UNSIGNED INT,
    `images` UNSIGNED INT,
    `teaser` TEXT,
    `images-url` TEXT,
    `board` TEXT NOT NULL,
    CONSTRAINT `board_fk` FOREIGN KEY (`board`) REFERENCES `boards`(`uuid`)
);

-- insert example on sqlite3
-- INSERT OR IGNORE INTO `threads` (`id`, `date`, `file`, `responses`, `images`, `teaser`, `images-url`) VALUES (853825027, 1620353054, "1614140109332m.jpg", 16, 15, "ITT: we post images with feature. Bonus points for other feature.", "1620353054097")

