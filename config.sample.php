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
	"store_tokens" => "https://tourn.ripple.moe/store_tokens",
	"fcm_token" => "", // https://console.firebase.google.com/project/SOMETHING/settings/cloudmessaging/ - Server Key
];

define("DATABASE_HOST", $config["db"]["host"]);
define("DATABASE_USER", $config["db"]["user"]);
define("DATABASE_PASS", $config["db"]["pass"]);
define("DATABASE_NAME", $config["db"]["name"]);
