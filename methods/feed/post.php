<?php

function run_method($state)
{
	$i = json_decode(file_get_contents('php://input'));
	if (!@$i->tournament) {
		error_message("Missing tournament parameter", 422);
		return;
	}
	$i->tournament = (int) $i->tournament;
	if (!@$i->content) {
		error_message("Missing content", 422);
		return;
	}

	$uid = $state->getSelfID();
	if (!$state->db->fetch("SELECT 1 FROM tournament_staff WHERE id = ? AND tournament = ? AND (privileges & 1) = 1", [$uid, $i->tournament])) {
		error_message("You are not allowed to post on this tournament feed.", 403);
		return;
	}

	$state->db->execute("INSERT INTO feed_items(content, created_at, author, tournament) VALUES (?, NOW(), ?, ?)",
	[$i->content, $uid, $i->tournament]);

	echo json_encode([
		"ok" => true,
		"item" => [
			"id" => $state->db->lastInsertId(),
			"content" => $i->content,
			"author" => $uid,
			"created_at" => (new \DateTime())->format(\DateTime::RFC3339),
		],
	], JSON_HEX_TAG);
}
