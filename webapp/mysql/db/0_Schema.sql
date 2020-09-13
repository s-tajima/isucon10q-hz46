DROP DATABASE IF EXISTS isuumo;
CREATE DATABASE isuumo;

DROP TABLE IF EXISTS isuumo.estate;
DROP TABLE IF EXISTS isuumo.chair;

CREATE TABLE isuumo.estate (
	`id` int(11) NOT NULL,
	`name` varchar(64) NOT NULL,
	`description` varchar(4096) NOT NULL,
	`thumbnail` varchar(128) NOT NULL,
	`address` varchar(128) NOT NULL,
	`latitude` double NOT NULL,
	`longitude` double NOT NULL,
	`rent` int(11) NOT NULL,
	`door_height` int(11) NOT NULL,
	`door_width` int(11) NOT NULL,
	`features` varchar(64) NOT NULL,
	`popularity` int(11) NOT NULL,
	`geo` GEOMETRY AS (ST_GeomFromText(CONCAT('POINT(', latitude,' ', longitude, ' )'))) STORED NOT NULL,
	PRIMARY KEY (`id`),
	INDEX i1 (popularity, id),
	INDEX i2 (rent, id),
	INDEX i3 (door_width, door_height),
	INDEX i4 (latitude, longitude, popularity),
	SPATIAL INDEX(geo)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;

CREATE TABLE isuumo.chair (
	`id` int(11) NOT NULL,
	`name` varchar(64) NOT NULL,
	`description` varchar(4096) NOT NULL,
	`thumbnail` varchar(128) NOT NULL,
	`price` int(11) NOT NULL,
	`height` int(11) NOT NULL,
	`width` int(11) NOT NULL,
	`depth` int(11) NOT NULL,
	`color` varchar(64) NOT NULL,
	`features` varchar(64) NOT NULL,
	`kind` varchar(64) NOT NULL,
	`popularity` int(11) NOT NULL,
	`stock` int(11) NOT NULL,
	PRIMARY KEY (`id`),
	INDEX i1 (id, stock),
	INDEX i2 (price, id),
	INDEX i3 (kind, stock),
	INDEX i4 (height, stock),
	INDEX i5 (color, stock)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


