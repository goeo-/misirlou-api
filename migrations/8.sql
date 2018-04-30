ALTER TABLE tournament_rules ADD PRIMARY KEY(id);

CREATE TABLE beatmap_requests(
	id BIGINT NOT NULL,
	tournament BIGINT NOT NULL,
	user INT NOT NULL,
	beatmap INT NOT NULL,
	category TINYINT(2) NOT NULL,
	PRIMARY KEY(id),
	FOREIGN KEY(tournament) REFERENCES tournaments(id)
		ON DELETE CASCADE
		ON UPDATE CASCADE
);

ALTER TABLE tournaments ADD max_beatmap_requests INT NOT NULL;
