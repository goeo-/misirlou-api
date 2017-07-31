<?php

require_once __DIR__ . "/../../registration_validation.php";

function run_method($state) {
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

	$teamID = (int) @$_GET["id"];
	if (!$state->db->fetch("SELECT 1 FROM team_users WHERE team = ? AND user = ? AND attributes = 0 LIMIT 1", [$teamID, $uid])) {
		error_message("That team does not exist!", 404);
		return;
	}

	$tournID = $state->db->fetch("SELECT tournament FROM teams WHERE id = ? LIMIT 1", $teamID)["tournament"];

	if (!validate_registration($state, $tournID, $uid)) {
		return;
	}

	$state->db->execute("UPDATE team_users SET attributes = 1 WHERE user = ? AND team = ?", [$uid, $teamID]);

	echo json_encode(["ok" => true]);
	return;
}
