/* Ideally a team with the same name should be able to exist over multiple tournaments */
ALTER TABLE teams DROP INDEX name, ADD UNIQUE(name, tournament);