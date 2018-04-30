CREATE TABLE tournaments(
	id BIGINT NOT NULL,
	name VARCHAR(60) NOT NULL,
	description VARCHAR(250) NOT NULL,
	mode TINYINT(1) NOT NULL,
	status TINYINT(2) NOT NULL,
	status_data VARCHAR(800),
	updated_at DATETIME NOT NULL,
	PRIMARY KEY(id)
);
