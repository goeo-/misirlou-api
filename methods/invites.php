<?php

function run_method($state)
{
	$offset = max((int)@$_GET["p"], 0) * 50;

	$query = "SELECT
	t.id as tourn_id, t.name as tourn_name,
	t.description, t.mode, t.status,
	t.exclusivity_starts, t.exclusivity_ends,
	teams.id, teams.name, teams.captain
FROM tournaments t
LEFT JOIN teams ON teams.tournament = t.id
LEFT JOIN team_users tu ON tu.team = teams.id
WHERE tu.user = ? AND tu.attributes = 0 AND t.status = 1
ORDER BY t.exclusivity_starts ASC LIMIT $offset, 50";


	$results = $state->db->fetchAll($query, [$state->getSelfID()]);

	array_walk($results, 'walker');

	echo json_encode([
		"ok" => true,
		"invites" => $results,
	], JSON_HEX_TAG);
}

function walker(&$el)
{
	$el = [
		"id"         => (int) $el["id"],
		"name"       => $el["name"],
		"captain"    => (int) $el["captain"],
		"tournament" => [
			"id"       => (int) $el["tourn_id"],
			"name"     => $el["tourn_name"],
			"mode"     => (int) $el["mode"],
			"status"   =>  (int) $el["status"],
			"description" => $el["description"],
			"exclusivity_starts" =>   $el["exclusivity_starts"],
			"exclusivity_ends"   =>  $el["exclusivity_ends"],
		],
	];
}

