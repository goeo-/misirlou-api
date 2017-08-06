<?php

function run_method($state)
{
	$tourn_id = (int) @$_GET["tourn_id"];

	$offset = max((int)@$_GET["p"], 0) * 50;

	$teams = $state->db->fetchAll("SELECT
	id, name, captain, created_at
FROM teams
LIMIT $offset, 50
WHERE tournament = ?", [$tourn_id]);

	array_walk($teams, "walker");

	echo json_encode([
		"ok" => true,
		"teams" => $teams,
	], JSON_HEX_TAG);
}

function walker(&$team) {
	$team = [
		"id" => (int) $team["id"],
		"name" => $team["name"],
		"captain" => (int) $team["captain"],
		"created_at" => to_3339($team["created_at"]),
	];
}
