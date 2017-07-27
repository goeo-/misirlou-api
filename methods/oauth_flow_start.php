<?php

// This method starts the OAuth 2 flow to login with Ripple.

function run_method($state)
{
	$provider = get_oauth_provider();

	$url = $provider->getAuthorizationUrl();

	$_SESSION["oauth_state"] = $provider->getState();

	header("Location: " . $url);
}
