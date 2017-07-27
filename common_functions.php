<?php

function get_oauth_provider()
{
	global $config;
	return new \League\OAuth2\Client\Provider\GenericProvider([
		'clientId'                => $config["oauth"]["client_id"],
		'clientSecret'            => $config["oauth"]["client_secret"],
		'redirectUri'             => $config["oauth"]["redirect_uri"],
		'urlAuthorize'            => 'https://ripple.moe/oauth/authorize',
		'urlAccessToken'          => 'https://ripple.moe/oauth/token',
		'urlResourceOwnerDetails' => 'https://ripple.moe/api/v1/ping',
		'scopeSeparator'          => ' ',
		'scopes'                  => [],
	]);
}

/**
 * outputs an error in JSON format.
 *
 * @param string $message error message to write
 * @param int $code HTTP response code, defaults to 400
 * @return void
 */
function error_message($message, $code = 400)
{
	http_response_code($code);
	echo json_encode([
		"ok" => false,
		"message" => $message,
	]);
}

