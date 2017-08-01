<?php

function run_method($state)
{
	$parts = ["t.status != 0"];
	$params = [];

	if (@$_GET["id"]) {
		$parts[] = "t.id = ?";
		$params[] = $_GET["id"];
	}

	$offset = max((int)@$_GET["p"], 0) * 50;

	$query = "SELECT
	t.id, t.name, t.description, t.mode, t.status,
	t.status_data, t.created_at, t.updated_at,
	t.team_size, t.min_team_size, t.exclusivity_starts,
	t.exclusivity_ends,	t.max_beatmap_requests
FROM tournaments t
";
	$query .= build_where($parts);
	$query .= " GROUP BY t.id ORDER BY t.updated_at DESC LIMIT $offset, 50";

	$results = $state->db->fetchAll($query, $params);
	array_walk($results, 'walker', $state);

	echo json_encode([
		"ok" => true,
		"tournaments" => $results,
	], JSON_HEX_TAG);
}

function walker(&$el, $key, $state)
{
	$el["id"]                   = (int) $el["id"];
	$el["mode"]                 = (int) $el["mode"];
	$el["status"]               = (int) $el["status"];
	$el["team_size"]            = (int) $el["team_size"];
	$el["min_team_size"]        = (int) $el["min_team_size"];
	$el["max_beatmap_requests"] = (int) $el["max_beatmap_requests"];

	if (empty($el["status_data"])) {
		$el["status_data"] = null;
	} else {
		$data = json_decode($el["status_data"]);
		$el["status_data"] = $data;
	}

	$el["created_at"] = to_3339($el["created_at"]);
	$el["updated_at"] = to_3339($el["updated_at"]);

	$myTeam = $state->db->fetch("SELECT tu.team AS my_team, teams.name as my_team_name
FROM team_users tu
INNER JOIN teams ON teams.id = tu.team
WHERE teams.tournament = ? AND tu.attributes > 0
LIMIT 1", [$el["id"]]);

	if ($myTeam) {
		$el["my_team"] = (int) $myTeam["my_team"];
		$el["my_team_name"] = $myTeam["my_team_name"];
	}
}
