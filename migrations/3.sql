ALTER TABLE tournaments ADD team_size TINYINT NOT NULL;

CREATE TABLE teams(
	id INT NOT NULL AUTO_INCREMENT,
	name VARCHAR(25) NOT NULL,
	tournament INT NOT NULL,
	captain INT NOT NULL,
	/*
		max team size is 4 (int32s)
		1 + 4 * 4 (first byte is length)
		base64'd makes 27 bytes
	*/
	created_at DATETIME NOT NULL,
	PRIMARY KEY (id),
	UNIQUE (name),
	FOREIGN KEY (tournament) REFERENCES tournaments(id)
		ON UPDATE CASCADE
		ON DELETE CASCADE
);

CREATE TABLE team_users(
	team INT NOT NULL,
	attributes TINYINT NOT NULL,
	user INT NOT NULL,
	PRIMARY KEY(team, user),
	FOREIGN KEY(team) REFERENCES teams(id)
		ON DELETE CASCADE
		ON UPDATE CASCADE
);
