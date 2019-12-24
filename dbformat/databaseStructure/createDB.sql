CREATE DATABASE IF NOT EXISTS fileHostingSite;

USE fileHostingSite;

CREATE TABLE IF NOT EXISTS users (
	PRIMARY KEY(username),
	username VARCHAR(20) NOT NULL,
	password VARCHAR(60) NOT NULL,
	timezone VARCHAR(40) NOT NULL,
	rating INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS sessions (
	PRIMARY KEY(cookie, username),
	username VARCHAR(20) NOT NULL,
	cookie VARCHAR(60) NOT NULL,
	expires DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS files (
	PRIMARY KEY(id),
	id INT NOT NULL AUTO_INCREMENT,
	label VARCHAR(50) NOT NULL,
	filesizeBytes INT NOT NULL,
	description VARCHAR(500) NOT NULL,
	owner VARCHAR(20) NOT NULL,
	category VARCHAR(20) NOT NULL,
	uploadDate DATETIME NOT NULL,
	rating INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS filesRating (
	PRIMARY KEY(fileID, voter),
	fileID INT NOT NULL,
	voter VARCHAR(20) NOT NULL,
	rating SMALLINT
);

