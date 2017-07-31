<?php

function run_method($state)
{
	$parts = ["t.status != 0"];
	$params = [$state->getSelfID()];

	if (@$_GET["id"]) {
		$parts[] = "t.id = ?";
		$params[] = $_GET["id"];
	}

	$offset = min((int)@$_GET["p"], 0) * 50;

	$query = "SELECT
	t.id, t.name, t.description, t.mode, t.status,
	t.status_data, t.created_at, t.updated_at,
	t.team_size, t.min_team_size, t.exclusivity_starts,
	t.exclusivity_ends,	tu.team AS my_team, teams.name as my_team_name
FROM tournaments t
LEFT JOIN teams ON teams.tournament = t.id
LEFT JOIN team_users tu ON tu.team = teams.id AND tu.user = ? AND tu.attributes != 0
 ";
	$query .= build_where($parts);
	$query .= " ORDER BY t.updated_at DESC LIMIT $offset, 50";

	$results = $state->db->fetchAll($query, $params);
	array_walk($results, 'walker');

	echo json_encode([
		"ok" => true,
		"tournaments" => $results,
	], JSON_HEX_TAG);
}

function walker(&$el)
{
	$el["id"]     = (int) $el["id"];
	$el["mode"]   = (int) $el["mode"];
	$el["status"] = (int) $el["status"];
	$el["team_size"] = (int) $el["team_size"];
	$el["min_team_size"] = (int) $el["min_team_size"];

	if (empty($el["status_data"])) {
		$el["status_data"] = null;
	} else {
		$data = json_decode($el["status_data"]);
		$el["status_data"] = $data;
	}

	$el["created_at"] = to_3339($el["created_at"]);
	$el["updated_at"] = to_3339($el["updated_at"]);

	if ($el["my_team"] !== null) {
		$el["my_team"] = (int) $el["my_team"];
	}
}
