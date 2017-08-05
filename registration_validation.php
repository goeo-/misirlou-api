<?php

require_once __DIR__ . "/classes/Status.php";

function validate_registration($state, $tournID, $uid, $minStatus = 1) {
	// Get information about our tournament
	$tourn = $state->db->fetch("SELECT status, created_at, team_size, min_team_size, exclusivity_starts, exclusivity_ends FROM tournaments WHERE id = ?", [$tournID]);
	if (!$tourn || $tourn["status"] == Status::Organising) {
		error_message("Tournament does not exist.", 404);
		return false;
	}

	$starts = $tourn["exclusivity_starts"];
	$ends   = $tourn["exclusivity_ends"];

	if ($tourn["status"] > $minStatus) {
		error_message("No more registrations are allowed.", 403);
		return false;
	}

	// Check that user is not already in a team of this very tournament
	if ($state->db->fetch("SELECT 1 FROM teams
	INNER JOIN team_users ON teams.id = team_users.team
	WHERE teams.tournament = ? AND team_users.user = ? AND team_users.attributes > 0
	LIMIT 1", [$tournID, $uid])) {
		error_message("You are already in another team in this tournament.", 403);
		return false;
	}

	// Check that user is not already in any other team of a tournament
	// currently running
	// Explain exclusivity thing:
	// 1. if there is something that begins while our tournament is in session
	// 2. or if there is something that ends after our tournament is started
	//    but has started before our tournament starts, thus including both the
	//    case of our tournament being in a subsection of a bigger one and
	//    another tournament ending after ours starts
	if ($state->db->fetch("SELECT 1 FROM team_users
	INNER JOIN teams ON teams.id = team_users.team
	INNER JOIN tournaments ON teams.tournament = tournaments.id
	WHERE
		team_users.user = ? AND team_users.attributes > 0 AND
		tournaments.id != ? AND
		((tournaments.exclusivity_starts >= ? AND tournaments.exclusivity_starts <= ?) OR
		 (tournaments.exclusivity_ends >= ? AND tournaments.exclusivity_starts <= ?))
	LIMIT 1", [$uid, $tournID, $starts, $ends, $starts, $starts])) {
		error_message("You can't join another tournament that overlaps with a tournament you're already in!", 403);
		return false;
	}

	return $tourn;
}