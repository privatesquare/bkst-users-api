CREATE TABLE users (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `first_name` varchar(45) NOT NULL,
    `last_name` varchar(45) NOT NULL,
    `email` varchar(45) NOT NULL,
    `password` varchar(255) NOT NULL,
    `status` varchar(45) NOT NULL,
    `date_created` datetime NOT NULL,
    `date_updated` datetime NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `email_UNIQUE` (`email`)
);