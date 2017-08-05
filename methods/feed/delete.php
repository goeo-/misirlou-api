<?php

function run_method($state)
{
	$post_id = (int) @$_GET["post_id"];
	$tourn = $state->db->fetch("SELECT tournament FROM feed_items WHERE id = ?", [$post_id]);
	if (!$tourn) {
		error_message("does not exist", 404);
		return;
	}
	if (!$state->db->fetch("SELECT 1 FROM tournament_staff WHERE id = ? AND tournament = ? AND (privileges & 1) = 1", [$state->getSelfId(), $tourn["tournament"]])) {
		error_message("You are not allowed to delete the post.", 403);
		return;
	}

	$state->db->execute("DELETE FROM feed_items WHERE id = ?", [$post_id]);
	echo json_encode([
		"ok" => true,
		"id" => $post_id,
	]);
}
