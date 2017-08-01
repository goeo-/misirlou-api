<?php

function run_method($state)
{
	$uid = $state->getSelfID();
	if (!$uid) {
		error_message("Missing or invalid session token.", 401);
		return;
	}

	$requests = $state->db->fetchAll(
		"SELECT beatmap, category FROM beatmap_requests WHERE user = ? AND tournament = ?",
		[$uid, (int) @$_GET["tourn_id"]]
	);

	array_walk($requests, function(&$request) {
		$request["beatmap"] = (int) $request["beatmap"];
		$request["category"] = (int) $request["category"];
	});

	echo json_encode([
		"ok"       => true,
		"requests" => $requests,
	]);
}
