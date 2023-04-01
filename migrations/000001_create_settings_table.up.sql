-- CREATE TABLE "settings" -------------------------------------
CREATE TABLE `settings`(
                           `id` BigInt( 255 ) AUTO_INCREMENT NOT NULL,
                           `username` VarChar( 255 ) NOT NULL,
                           `chat_id` BigInt( 255 ) NOT NULL,
                           `linkace_token` VarChar( 255 ) NOT NULL,
                           `easylist_token` VarChar( 255 ) NOT NULL,
                           `easywords_token` VarChar( 255 ) NOT NULL,
                           `mode` TinyInt( 255 ) NOT NULL DEFAULT 0,
                           `context` Text NOT NULL,
                           PRIMARY KEY (id),
                           CONSTRAINT `unique_id` UNIQUE( `id` ),
                           CONSTRAINT `unique_username` UNIQUE( `username` ),
                           CONSTRAINT `unique_chat_id` UNIQUE( `chat_id` ),
                           CONSTRAINT `unique_linkace_token` UNIQUE( `linkace_token` ),
                           CONSTRAINT `unique_easylist_token` UNIQUE( `easylist_token` ),
                           CONSTRAINT `unique_easywords_token` UNIQUE( `easywords_token` ) )
    ENGINE = InnoDB;-- -------------------------------------------------------------
