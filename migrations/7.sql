CREATE TABLE feed_items(
	id BIGINT NOT NULL,
	tournament BIGINT NOT NULL,
	content TEXT NOT NULL,
	author INT,
	PRIMARY KEY(id),
	FOREIGN KEY(tournament) REFERENCES tournaments(id)
		ON DELETE CASCADE
		ON UPDATE CASCADE
);
