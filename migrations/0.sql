CREATE TABLE sessions(
	id CHAR(64) NOT NULL,
	user_id INT NOT NULL,
	access_token VARCHAR(64) NOT NULL,
	PRIMARY KEY(id)
);