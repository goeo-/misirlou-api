<?php

function run_method($state) {
	$id = $state->getSelfID();
	if (!$id) {
		error_message("Missing or invalid session token.", 401);
		return;
	}
	$teamID = @$_GET["id"];
	if (!$teamID) {
		error_message("Missing parameter id.", 422);
		return;
	}

	$state->db->execute("DELETE FROM team_users WHERE user = ? AND team = ? AND attributes = 0", [$id, $teamID]);

	echo json_encode(["ok" => true]);
}
