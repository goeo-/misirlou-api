<?php

function run_method($state)
{
	// Check whether the state is invalid (mitigate CSRF)
	if (check_invalid_state()) {
		$_SESSION = [];
		error_message("Invalid state");
		return;
	}

	global $config;

	try {
		$provider = get_oauth_provider();
		// retrieve access token from the ripple server
		$accessToken = $provider->getAccessToken("authorization_code", [
			"code" => $_GET["code"],
		]);
		$token = $accessToken->getToken();
		$user = RippleAPI::user("self", $token);
		if ($user === null) {
			error_message("user doesn't exist", 404);
			return;
		}
		$session_token = create_session($user->id, $token, $state->db);
		header("Location: " . $config["store_tokens"] . "?session=" . $session_token . "&access=" . $token);
	} catch (\League\OAuth2\Client\Provider\Exception\IdentityProviderException $e) {
		error_message($e->getMessage());
	}
}

// returns true when the state is invalid
function check_invalid_state()
{
	return empty($_GET["state"]) || @$_SESSION["oauth_state"] !== $_GET["state"];
}

/**
 * @param int $user_id
 * @param string $access_token
 * @param DBPDO $db
 * @return string
 */
function create_session($user_id, $access_token, $db)
{
	$token = random_str(32);
	$db->execute("INSERT INTO sessions(id, user_id, access_token) VALUES (?, ?, ?)", [
		hash("sha256", $token),
		$user_id,
		$access_token,
	]);
	return $token;
}

/**
 * Generate a random string, using a cryptographically secure
 * pseudorandom number generator (random_int)
 * https://stackoverflow.com/a/31107425/5328069
 *
 * @param int $length      How many characters do we want?
 * @param string $keyspace A string of all possible characters
 *                         to select from
 * @return string
 */
function random_str($length, $keyspace = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ')
{
    $str = '';
    $max = mb_strlen($keyspace, '8bit') - 1;
    for ($i = 0; $i < $length; ++$i) {
        $str .= $keyspace[random_int(0, $max)];
    }
    return $str;
}
