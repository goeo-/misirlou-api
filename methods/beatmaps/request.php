<?php

require_once __DIR__ . "/../../classes/Status.php";

function run_method($state)
{
	$uid = $state->getSelfID();
	if (!$uid) {
		error_message("Missing or invalid session token.", 401);
		return;
	}

	$tournID = (int) @$_GET["tourn_id"];
	if ($tournID === 0) {
		error_message("Missing tournament parameter.", 401);
		return;
	}

	$tourn = $state->db->fetch("SELECT max_beatmap_requests, status FROM tournaments WHERE id = ? LIMIT 1", [$tournID]);
	if (!$tourn) {
		error_message("Tournament does not exist.");
		return;
	}
	// TODO: temporary fix
	/*if ($tourn["max_beatmap_requests"] < 1 || $tourn["status"] != Status::Open) {
		error_message("Tournament does not accept beatmap requests.");
		return;
	}*/

	// Decode POST body and check tournament is set
	$maps = json_decode(file_get_contents('php://input'));
	if ($maps === false || !is_array($maps)) {
		error_message("Invalid JSON body", 422);
		return;
	}

	$sprintfString = "($tournID, $uid, %d, %d)";
	$inserts = [];
	$retValues = [];
	foreach ($maps as $map) {
		if (!isset($map->beatmap) || ((int) $map->beatmap) < 1)
			continue;
		$inserts[] = sprintf($sprintfString, (int) $map->beatmap, (int) $map->category);
		$retValues[] = [
			"beatmap"  => (int) $map->beatmap,
			"category" => (int) $map->category,
		];
	}
	$inserts   = array_slice(array_unique($inserts), 0, $tourn["max_beatmap_requests"]);
	// We need array_values because when array_unique detects a duplicate, it
	// converts the array to an assoc array, and thus in the response it is
	// returned as an object, not an array.
	$retValues = array_slice(array_values(array_unique($retValues, SORT_REGULAR)), 0, $tourn["max_beatmap_requests"]);

	$state->db->execute("DELETE FROM beatmap_requests WHERE user = ? AND tournament = ?", [$uid, $tournID]);

	if ($inserts > 0) {
		$state->db->execute(
			"INSERT INTO beatmap_requests(tournament, user, beatmap, category) VALUES " . implode(", ", $inserts)
		);
	}

	echo json_encode([
		"ok"       => true,
		"requests" => $retValues,
	]);
}
