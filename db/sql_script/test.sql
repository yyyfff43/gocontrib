CREATE SCHEMA `test` ;
use test;
CREATE TABLE `book` (
    `book_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `book_name` varchar(45) NOT NULL,
    `author` varchar(45) NOT NULL,
    `size` int(11) NOT NULL,
    PRIMARY KEY (`book_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;
CREATE TABLE `hachi` (
     `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
     `test_a` varchar(45) NOT NULL,
     `test_b` varchar(45) DEFAULT NULL,
     PRIMARY KEY (`id`),
     UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;

CREATE SCHEMA `testa` ;
use testa;
CREATE TABLE `book` (
    `book_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `book_name` varchar(45) NOT NULL,
    `author` varchar(45) NOT NULL,
    `size` int(11) NOT NULL,
    PRIMARY KEY (`book_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;
CREATE TABLE `hachi` (
     `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
     `test_a` varchar(45) NOT NULL,
     `test_b` varchar(45) DEFAULT NULL,
     PRIMARY KEY (`id`),
     UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;

CREATE SCHEMA `testb` ;
use testb;
CREATE TABLE `book` (
    `book_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `book_name` varchar(45) NOT NULL,
    `author` varchar(45) NOT NULL,
    `size` int(11) NOT NULL,
    PRIMARY KEY (`book_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;
CREATE TABLE `hachi` (
     `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
     `test_a` varchar(45) NOT NULL,
     `test_b` varchar(45) DEFAULT NULL,
     PRIMARY KEY (`id`),
     UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;
