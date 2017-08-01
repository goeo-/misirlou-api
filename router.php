<?php

require_once __DIR__ . "/include.php";

$request_uri = $_SERVER["REQUEST_URI"];

if (stripos($request_uri, "/api") === 0)
	$request_uri = substr($request_uri, 4);

if (strpos($request_uri, "?") !== FALSE) {
	$qs = parse_url($request_uri, PHP_URL_QUERY);
	parse_str($qs, $_GET);
	$request_uri = substr($request_uri, 0, strpos($request_uri, "?"));
}

$method_path = realpath(__DIR__ . "/methods" . $request_uri . ".php");

// make sure the path is valid
if (strpos($method_path, __DIR__ . "/methods/") === FALSE) {
	error_message("not found", 404);
	return;
}

try {
	// include method
	require_once $method_path;

	run_method($state);
} catch (Exception $e) {
	if (isset($sentry)) {
		$sentry->captureException($e);
	}
	error_message("server-side error", 500);
}
