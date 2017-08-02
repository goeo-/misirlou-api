<?php

function run_method($state) {
	$team = (int) @$_GET["team"];
	if ($team === 0) {
		error_message("Missing team ID", 422);
		return;
	}
	$newName = remove_cchars(@$_GET["name"]);
	if (!$newName) {
		error_message("Missing team name", 422);
		return;
	}
	$uid = $state->getSelfID();
	if (!$uid) {
		error_message("Missing or invalid session token.", 401);
		return;
	}

	$team = $state->db->fetch("SELECT t.id, tourn.team_size FROM teams t INNER JOIN tournaments tourn ON tourn.id = t.tournament WHERE t.id = ? AND t.captain = ? LIMIT 1", [$team, $uid]);
	if (!$team) {
		error_message("Team not found", 404);
		return;
	}

	if ($team["team_size"] == 1) {
		error_message("You can't change your username in a single-player tournament.");
		return;
	}

	$state->db->execute("UPDATE teams SET name = ? WHERE id = ?", [$newName, $team["id"]]);

	echo json_encode([
		"ok"       => true,
		"new_name" => $newName,
	], JSON_HEX_TAG);
}