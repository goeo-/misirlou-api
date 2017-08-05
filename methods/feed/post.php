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

	notify_tournament_members($state, $i->tournament, $i->content);

	echo json_encode([
		"ok" => true,
		"item" => [
			"id" => $state->db->lastInsertId(),
			"content" => remove_cchars($i->content),
			"author" => $uid,
			"created_at" => (new \DateTime())->format(\DateTime::RFC3339),
		],
	], JSON_HEX_TAG);
}

require_once __DIR__ . "/../../classes/Notify.php";

function notify_tournament_members($state, $id, $content) {
	$users_raw = $state->db->fetchAll("SELECT user FROM team_users WHERE team IN (SELECT id FROM teams WHERE tournament = ?) AND attributes != 0", [$id]);

	$users = [];
	foreach ($users_raw as $user) {
		$users[] = $user["user"];
	}

	// Notify users of invites
	$sett = new NotifySettings();
	$sett->getUsers($state, $users);
	$sett->title = "New feed post!";
	$sett->body = substr(remove_cchars($content), 0, 20);
	$sett->action = "https://tourn.ripple.moe/feed/$id"; // TODO: HARDCODE
	$sett->icon = "https://tourn.ripple.moe/static/favicon.png"; // TODO: HARDCODE
	Notify($sett);
}
