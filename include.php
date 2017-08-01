<?php
require_once(__DIR__ . "/vendor/autoload.php");

require_once(__DIR__ . "/config.php");

// Set up sentry
if (!empty($config["sentry"]["dsn"])) {
	$sentry = new Raven_Client($config["sentry"]["dsn"]);
	$sentry->install();
}

require_once(__DIR__ . "/classes/DBPDO.php");

require_once(__DIR__ . "/classes/State.php");
require_once(__DIR__ . "/classes/RippleAPI.php");
require_once(__DIR__ . "/common_functions.php");

$db = new DBPDO();
$state = new State($db);

session_start();
