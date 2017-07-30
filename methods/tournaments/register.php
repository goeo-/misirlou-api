<?php

require_once __DIR__ . "/../../classes/Status.php";

function run_method($state)
{
	// We're at 25 lines and I already know this is gonna turn into a 150LOC
	// function, and I hate myself for that.

	$tok = $state->getAccessToken();
	if (!$tok) {
		error_message("Missing or invalid session token.", 401);
		return;
	}

	// Get own user from the Ripple API, and make sure they are not restricted.
	$user = RippleAPI::user("self", $tok);
	if ($user->code === 404) {
		error_message("Access token is invalid (you probably revoked the token).", 401);
		return;
	}
	if (($user->privileges & 3) != 3) {
		error_message("You don't have the required privileges to register for a tournament.", 403);
		return;
	}
	$uid = $user->id;

	// Decode POST body and check tournament is set
	$obj = json_decode(file_get_contents('php://input'));
	if (empty($obj->tournament)) {
		error_message("Missing tournament parameter", 422);
		return;
	}

	// Get information about our tournament
	$tourn = $state->db->fetch("SELECT status, created_at, team_size, min_team_size, exclusivity_starts, exclusivity_ends FROM tournaments WHERE id = ?", [$obj->tournament]);
	if (!$tourn || $tourn["status"] == Status::Organising) {
		error_message("Tournament does not exist.", 404);
		return;
	}

	$starts = $tourn["exclusivity_starts"];
	$ends   = $tourn["exclusivity_ends"];

	if ($tourn["status"] > Status::Open) {
		error_message("No more registrations are allowed.", 403);
		return;
	}

	// Check that user is not already in a team of this very tournament
	if ($state->db->fetch("SELECT 1 FROM teams
	INNER JOIN team_users ON teams.id = team_users.team
	WHERE teams.tournament = ? AND team_users.user = ? AND team_users.attributes > 0
	LIMIT 1", [$obj->tournament, $uid])) {
		error_message("You are already in another team in this tournament.", 403);
		return;
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
	LIMIT 1", [$uid, $obj->tournament, $starts, $ends, $starts, $starts])) {
		error_message("You can't join another tournament that overlaps with a tournament you're already in!", 403);
		return;
	}

	if ($tourn["team_size"] == 1) {
		$state->db->execute("INSERT INTO teams(name, tournament, captain, created_at) VALUES (?, ?, ?, NOW())", [
			$user->username,
			$obj->tournament,
			$uid,
		]);
		$id = $state->db->lastInsertId();
		$state->db->execute("INSERT INTO team_users(team, attributes, user) VALUES(?, 2, ?)", [$id, $uid]);
		echo json_encode([
			"ok" => true,
			"team_id" => $id,
		]);
		return;
	}

	// Oh boy...
	$members = @$obj->members;
	$name = @$obj->name;
	if (!$members) {
		error_message("For a tournament of size bigger than 1, at least one further member in the team is required.");
		return;
	}
	if (!$name) {
		error_message("Missing team name");
		return;
	}

	array_walk($members, function(&$el) {
		$el = (int) $el;
	});
	$members = array_unique($members);
	$members = array_filter($members, function($value) use ($uid) {
		return $value > 0 && $value !== $uid;
	});

	if ((count($members)+1) < $tourn["min_team_size"] || (count($members)+1) > $tourn["team_size"]) {
		error_message("Number of members goes out of bonds of team size", 413);
		return;
	}

	$state->db->execute("INSERT INTO teams(name, tournament, captain, created_at) VALUES (?, ?, ?, NOW())", [
		remove_cchars($user->username),
		$obj->tournament,
		$uid,
	]);
	$id = $state->db->lastInsertId();

	$vals = ["($id, 2, $uid)"];
	// $members has been array_walked and all its elements are ints
	foreach ($members as $member) {
		$vals[] = "($id, 0, $member)";
	}

	$state->db->execute("INSERT INTO team_users(team, attributes, user) VALUES " . implode(", ", $vals));
	echo json_encode([
		"ok" => true,
		"team_id" => $id,
	]);
}