<?php

function run_method($state) {
	$team = (int) @$_GET["team"];
	if ($team === 0) {
		error_message("Missing team ID", 422);
		return;
	}
	$uid = $state->getSelfID();
	if (!$uid) {
		error_message("Missing or invalid session token.", 401);
		return;
	}

	$team = $state->db->fetch("SELECT id FROM teams WHERE id = ? AND captain = ? LIMIT 1", [$team, $uid]);
	if (!$team) {
		error_message("Team not found", 404);
		return;
	}

	$state->db->execute("DELETE FROM teams WHERE id = ?", [$team["id"]]);

	echo json_encode(["ok" => true]);
}
