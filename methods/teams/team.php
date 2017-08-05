<?php

function run_method($state) {
	$uid = $state->getSelfId();
	if (!$uid) {
		error_message("Missing or invalid session token.", 401);
		return;
	}

	$teamID = (int) @$_GET["id"];
	if (!$teamID) {
		error_message("Missing team ID", 422);
		return;
	}

	$members = $state->db->fetchAll("SELECT attributes, user FROM team_users WHERE team = ?", $teamID);

	array_walk($members, "member_walker");

	$team = $state->db->fetch(
"SELECT
	teams.id, teams.name, teams.captain, teams.created_at,
	t.id as tourn_id, t.name as tourn_name, t.description,
	t.mode, t.status, t.team_size, t.exclusivity_starts,
	t.exclusivity_ends, t.min_team_size
FROM teams
INNER JOIN tournaments t ON t.id = teams.tournament
WHERE teams.id = ?", [$teamID]);

	echo json_encode([
		"ok" => true,
		"team" => team_obj_creator($team, $members),
	], JSON_HEX_TAG);
}

function member_walker(&$member) {
	$member["attributes"] = (int) $member["attributes"];
	$member["user"] = (int) $member["user"];
}

function team_obj_creator($team, $members) {
	return [
		"id"         => (int) $team["id"],
		"name"       => $team["name"],
		"captain"    => (int) $team["captain"],
		"created_at" => to_3339($team["created_at"]),
		"members"    => $members,
		"tournament" => [
			"id"                 => (int) $team["tourn_id"],
			"name"               => $team["tourn_name"],
			"description"        => $team["description"],
			"mode"               => (int) $team["mode"],
			"status"             => (int) $team["status"],
			"team_size"          => (int) $team["team_size"],
			"exclusivity_starts" => $team["exclusivity_starts"],
			"exclusivity_ends"   => $team["exclusivity_ends"],
			"min_team_size"      => (int) $team["min_team_size"],
		],
	];
}
