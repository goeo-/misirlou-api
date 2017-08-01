<?php

$config = [
	"db" => [
		"host" => "",
		"user" => "",
		"pass" => "",
		"name" => "",
	],
	"oauth" => [
		"client_id"     => "",
		"client_secret" => "",
		"redirect_uri"  => "https://localhost/api/oauth_flow_finish"
	],
	"sentry" => [
		"dsn" => "",
	],
];

define("DATABASE_HOST", $config["db"]["host"]);
define("DATABASE_USER", $config["db"]["user"]);
define("DATABASE_PASS", $config["db"]["pass"]);
define("DATABASE_NAME", $config["db"]["name"]);
