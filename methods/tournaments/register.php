<?php

require_once __DIR__ . "/../../registration_validation.php";
require_once __DIR__ . "/../../classes/Notify.php";

function run_method($state)
{
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

	// Decode POST body and check tournament is set
	$obj = json_decode(file_get_contents('php://input'));
	if (empty($obj->tournament)) {
		error_message("Missing tournament parameter", 422);
		return;
	}

	$tourn = validate_registration($state, $obj->tournament, $uid);
	if ($tourn === false) {
		return;
	}

	if ($tourn["team_size"] == 1) {
		$state->db->execute("INSERT INTO teams(name, tournament, captain, created_at) VALUES (?, ?, ?, NOW())", [
			$user->username,
			$obj->tournament,
			$uid,
		]);
		$id = $state->db->lastInsertId();
		$state->db->execute("INSERT INTO team_users(team, attributes, user) VALUES(?, 2, ?)", [$id, $uid]);
		echo json_encode([
			"ok" => true,
			"team_id" => $id,
		]);
		return;
	}

	// Oh boy...
	$members = @$obj->members;
	$name = @$obj->name;
	if (!$members) {
		error_message("For a tournament of size bigger than 1, at least one further member in the team is required.");
		return;
	}
	if (!$name) {
		error_message("Missing team name");
		return;
	}

	if ($state->db->fetch("SELECT 1 FROM teams WHERE tournament = ? AND name = ? LIMIT 1", [$obj->tournament, $name])) {
		error_message("A team with that name already exists!", 403);
	}

	array_walk($members, function(&$el) {
		$el = (int) $el;
	});
	$members = array_unique($members);
	$members = array_filter($members, function($value) use ($uid) {
		return $value > 0 && $value !== $uid;
	});

	if ((count($members)+1) < $tourn["min_team_size"] || (count($members)+1) > $tourn["team_size"]) {
		error_message("Number of members goes out of bonds of team size", 413);
		return;
	}

	$state->db->execute("INSERT INTO teams(name, tournament, captain, created_at) VALUES (?, ?, ?, NOW())", [
		remove_cchars($name),
		$obj->tournament,
		$uid,
	]);
	$id = $state->db->lastInsertId();

	$vals = ["($id, 2, $uid)"];
	// $members has been array_walked and all its elements are ints
	foreach ($members as $member) {
		$vals[] = "($id, 0, $member)";
	}

	$state->db->execute("INSERT INTO team_users(team, attributes, user) VALUES " . implode(", ", $vals));

	// Notify users of invites
	$sett = new NotifySettings();
	$sett->getUsers($state, $members);
	$sett->title = "You just got invited!";
	$sett->body = "Would you like to join " . $user->username . "'s team?";
	$sett->action = "https://tourn.ripple.moe/invites"; // TODO: HARDCODE
	$sett->icon = "https://tourn.ripple.moe/static/favicon.png"; // TODO: HARDCODE
	Notify($sett);

	echo json_encode([
		"ok" => true,
		"team_id" => $id,
	]);
}