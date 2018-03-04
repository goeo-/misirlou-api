<?php

class RippleAPI {
	const BASE_URL = "https://api.ripple.moe/api/v1";

	/**
	 * @param int|string $id
	 * @param string $token
	 * @return null|string
	 */
	public static function user($id, $token) {
		try {
			return static::request("/users", "GET", $token, [
				"id" => $id,
			]);
		} catch (RippleAPIErrorNotFound $e) {
			return null;
		}
		return null;
	}

	/**
	 * @param string $path
	 * @param string $method
	 * @param string $token
	 * @param Array $params
	 */
	public static function request($path, $method, $token, $params = []) {
		// build URL
		$url = static::BASE_URL . $path;

		// create curl and set necessary headers and options
		$c = new \Curl\Curl();
		$c->setUserAgent("RippleAPIPHPClient/1.0 MisirlouAPI/1.0");
		$c->setHeader("Authorization", "Bearer " . $token);
		$c->setOpt(CURLOPT_ENCODING, "gzip");
		$c->setOpt(CURLOPT_RETURNTRANSFER, true);

		// do the actual request
		if ($method === "GET") {
			$c->get($url, $params);
		} elseif ($method == "POST") {
			$c->post($url, $params);
		} else {
			throw new RippleAPIClassError("Method not supported");
		}

		// error handling
		if ($c->error) {
			throw new RippleAPICurlError("", $c->error_code);
		}

		$resp = json_decode($c->response);

		if ($resp->code >= 400) {
			if ($resp->code === 404)
				throw new RippleAPIErrorNotFound($resp->message, $resp->code);
			throw new RippleAPIErrorResponse($resp->message, $resp->code);
		}

		return $resp;
	}
}

class RippleAPIErrorResponse extends Exception {}

class RippleAPIErrorNotFound extends RippleAPIErrorResponse {}

class RippleAPIClassError extends Exception {}

class RippleAPICurlError extends Exception {}
