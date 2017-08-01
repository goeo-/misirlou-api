<?php

function run_method($state)
{
	$tournID = (int) @$_GET["tourn_id"];
	if (!$tournID) {
		error_message("Missing tournament ID.", 422);
		return;
	}

	$offset = max((int)@$_GET["p"], 0) * 50;

	$elements = $state->db->fetchAll("SELECT id, content, created_at, author FROM feed_items WHERE tournament = ?
ORDER BY created_at DESC LIMIT $offset, 50", [$tournID]);

	array_walk($elements, "walker");

	echo json_encode([
		"ok" => true,
		"items" => $elements,
	], JSON_HEX_TAG);
}

function walker(&$el)
{
	$el["id"] = (int) $el["id"];
	$el["author"] = (int) $el["author"];
	$el["created_at"] = to_3339($el["created_at"]);
}
