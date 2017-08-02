<?php

// Life is short
// Filled with stuff
// I don't know what for
// I ain't had enough
// I learned all I know
// At about the age of nine
// But I can better myself
// If I could only find...
// SOME new kind of kick

function run_method($state) {
	$team = (int) @$_GET["team"];
	if ($team === 0) {
		error_message("Missing team ID", 422);
		return;
	}
	$target = (int) @$_GET["target"];
	if (!$target) {
		error_message("Missing target user", 422);
		return;
	}
	$uid = $state->getSelfID();
	if (!$uid) {
		error_message("Missing or invalid session token.", 401);
		return;
	}

	$info = $state->db->fetch("SELECT tu.attributes, teams.captain
FROM team_users tu
INNER JOIN teams ON tu.team = teams.id
WHERE tu.user = ? AND tu.team = ?", [$target, $team]);

	if (!$info) {
		error_message("Either that team does not exist or the user is not enrolled in it.", 404);
		return;
	}

	if ($info["captain"] == $uid) {
		// KICK
		if ($target == $uid) {
			error_message("You're a captain of this team, you can't kick yourself! You may only disband the team if you do not wish to take part in the tournament anymore.", 403);
			return;
		}
		$state->db->execute("DELETE FROM team_users WHERE user = ? AND team = ?", [$target, $team]);
		echo json_encode(["ok" => true]);
		return;
	} elseif ($target == $uid) {
		// LEAVE
		$state->db->execute("DELETE FROM team_users WHERE user = ? AND team = ?", [$target, $team]);
		echo json_encode(["ok" => true]);
		return;
	} else {
		error_message("You are not the captain of this team, and you are not trying to leave it.", 403);
		return;
	}
}
