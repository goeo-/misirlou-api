CREATE TABLE tournament_rules(
	id INT NOT NULL,
	rules TEXT NOT NULL,
	FOREIGN KEY (id) REFERENCES tournaments (id)
		ON DELETE CASCADE
		ON UPDATE CASCADE
);