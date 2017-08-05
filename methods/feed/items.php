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

	$can_post = $state->db->fetch("SELECT 1 FROM tournament_staff WHERE id = ? AND tournament = ? AND (privileges & 1) = 1", [$state->getSelfID(), $tournID]);

	array_walk($elements, "walker");

	echo json_encode([
		"ok" => true,
		"items" => $elements,
		"can_post" => (bool) $can_post,
	], JSON_HEX_TAG);
}

function walker(&$el)
{
	$el["id"] = (int) $el["id"];
	$el["author"] = (int) $el["author"];
	$el["created_at"] = to_3339($el["created_at"]);
}
