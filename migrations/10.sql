CREATE TABLE tournament_staff(
	id INT NOT NULL,
	tournament INT NOT NULL,
	privileges INT NOT NULL,
	PRIMARY KEY(id, tournament)
);