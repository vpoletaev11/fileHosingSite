CREATE DATABASE IF NOT EXISTS fileHostingSite;

CREATE TABLE IF NOT EXISTS fileHostingSite.users (
	PRIMARY KEY(username),
	username VARCHAR(20) NOT NULL,
	password VARCHAR(60) NOT NULL
);

CREATE TABLE IF NOT EXISTS fileHostingSite.sessions (
	username VARCHAR(20) NOT NULL,
	cookie VARCHAR(60) DEFAULT NULL
);
expiration DATETIME NOT NULL


CREATE TABLE IF NOT EXISTS fileHostingSite.files (
	filename VARCHAR(10) NOT NULL,
	PRIMARY KEY(filename),
	label VARCHAR(50) NOT NULL,
	description VARCHAR(500) NOT NULL,
	owner VARCHAR(20) NOT NULL,
	category VARCHAR(20) NOT NULL,
	upload_date DATE NOT NULL,
	rating INT NOT NULL
);


