<?php

require_once __DIR__ . "/../../classes/Status.php";

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

	$team = $state->db->fetch("SELECT id, tournament FROM teams WHERE id = ? AND captain = ? LIMIT 1", [$team, $uid]);
	if (!$team) {
		error_message("Team not found", 404);
		return;
	}

	$tournStatus = $state->db->fetch("SELECT status FROM tournaments WHERE id = ?", [$team["tournament"]])["status"];
	if ($tournStatus >= Status::RegClosedRequestsOpen) {
		error_message("If you want to disband a team after a tournament has begun, please contact an administrator of this tournament.", 403);
		return;
	}

	$state->db->execute("DELETE FROM teams WHERE id = ?", [$team["id"]]);

	echo json_encode(["ok" => true]);
}
