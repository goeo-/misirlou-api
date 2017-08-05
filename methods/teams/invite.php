<?php

require_once __DIR__ . "/../../classes/Notify.php";

function run_method($state)
{
	$team = (int) @$_GET["team"];
	if ($team === 0) {
		error_message("Missing team ID", 422);
		return;
	}
	$target = (int)@$_GET["target"];
	if (!$target) {
		error_message("Missing target", 422);
		return;
	}
	$uid = $state->getSelfID();
	if (!$uid) {
		error_message("Missing or invalid session token.", 401);
		return;
	}

	$teamInfo = $state->db->fetch("SELECT tournament, captain, name FROM teams WHERE id = ?", $team);
	if (!$teamInfo || $uid != $teamInfo["captain"]) {
		error_message("Team does not exist.", 404);
		return;
	}

	// Get tournament team size, team's size and compare
	$size = $state->db->fetch("SELECT team_size FROM tournaments WHERE tournaments.id = ?", [$teamInfo["tournament"]]);
	$members = $state->db->fetch("SELECT COUNT(*) as members FROM team_users WHERE team = ?", [$team]);

	if ($members["members"] >= $size["team_size"]) {
		error_message("You can't add new team members.", 403);
		return;
	}

	// check user isn't already in team
	if ($state->db->fetch("SELECT 1 FROM team_users WHERE user = ? AND team = ?", [$target, $team])) {
		error_message("User is already in team", 409);
		return;
	}

	$state->db->execute("INSERT INTO team_users(team, attributes, user) VALUES (?, 0, ?)", [
		$team,
		$target,
	]);

	// Notify users of invites
	$sett = new NotifySettings();
	$sett->getUsers($state, [$target]);
	$sett->title = "You just got invited!";
	$sett->body = "Would you like to join " . $teamInfo["name"];
	$sett->action = "https://tourn.ripple.moe/invites"; // TODO: HARDCODE
	$sett->icon = "https://tourn.ripple.moe/static/favicon.png"; // TODO: HARDCODE
	Notify($sett);

	echo json_encode([
		"ok" => true,
		"new_user" => [
			"attributes" => 0,
			"user" => $target,
		],
	]);
}
