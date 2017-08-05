<?php

function run_method($state)
{
	$rows = $state->db->execute("UPDATE sessions SET fcm_token = ? WHERE id = ?", [$_GET["fcm_token"], hash("sha256", @$_SERVER["HTTP_AUTHORIZATION"])])->rowCount();
	if ($rows === 0) {
		error_message("Session not found.", 404);
		return;
	}
	echo json_encode([
		"ok" => true,
	]);
}
