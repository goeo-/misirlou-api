<?php

function run_method($state) {
	echo "SOMETHING FANCIER THAN THIS NEXT TIME PINKY SWEAR<br><br>";
	$tournament = (int) @$_GET["tournament"];
	$beatmaps = $state->db->fetchAll("SELECT beatmap, GROUP_CONCAT(category SEPARATOR ',') as category, COUNT(*) as popularity FROM beatmap_requests WHERE tournament = ? GROUP BY beatmap ORDER BY popularity DESC", [$tournament]);

	$categories = [
		[
			"NoMod",
			"Hidden",
			"HardRock",
			"DoubleTime",
			"FreeMod",
			"Tiebreaker",
		],
		[
			"NoMod",
			"Hidden",
			"HardRock",
			"DoubleTime",
			"FreeMod",
			"Tiebreaker",
		],
		[
			"NoMod",
			"Hidden",
			"HardRock",
			"DoubleTime",
			"Tiebreaker",
		],
		[
			"NoMod",
			"FreeMod",
			"Tiebreaker",
		],
	];

	echo "<ul>";
	foreach ($beatmaps as $b) {
		$cats = explode(",", $b["category"]);
		foreach ($cats as $k => $cat) {
			$cats[$k] = $categories[$tournament - 1][$cat];
		}
		echo "<li><b>$b[popularity] times</b> - <a href='https://osu.ppy.sh/b/$b[beatmap]'>https://osu.ppy.sh/b/$b[beatmap]</a> | CATEGORIES: " . implode(", ", $cats) . "</li>";
	}
	echo "</ul>";
}