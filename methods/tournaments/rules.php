<?php

function run_method($state)
{
	$id = (int) @$_GET["id"];
	if (!$id) {
		error_message("Missing ID");
		return;
	}

	$result = $state->db->fetch("SELECT rules FROM tournament_rules WHERE id = ? LIMIT 1", [$id]);

	echo json_encode([
		"ok" => true,
		"rules" => $result["rules"],
	], JSON_HEX_TAG);
}
