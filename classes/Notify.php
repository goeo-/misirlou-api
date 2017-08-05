<?php

class NotifySettings {
	public $users = [];
	// https://firebase.google.com/docs/cloud-messaging/http-server-ref#notification-payload-support
	// table 2c
	public $title;
	public $body;
	public $icon;
	public $action;

	public function encodeNotification($page) {
		$offset = $page * 1000;
		$users = array_slice($this->users, $offset, 1000);
		if (count($users) === 0)
			return false;
		return json_encode([
			"registration_ids" => $this->users,
			"notification" => [
				"title" => $this->title,
				"body" => $this->body,
				"icon" => $this->icon,
				"click_action" => $this->action,
			],
		], JSON_HEX_TAG);
	}

	public function getUsers($state, $rippleIDs) {
		$results = $state->db->fetchAll("SELECT fcm_token FROM sessions WHERE user_id IN (" . implode(", ", $rippleIDs) . ")");
		if (!$results) {
			return false;
		}
		$users = [];
		foreach ($results as $res) {
			if ($res["fcm_token"] === "")
				continue;
			$users[] = $res["fcm_token"];
		}
		$this->users = array_values(array_unique($users));
		return true;
	}
}

// Basically a very basic client to send notifications using FCM.
function Notify(NotifySettings $settings) {
	global $config;

	$token = @$config["fcm_token"];
	if (!$token) {
		return false;
	}

	$c = new \Curl\Curl();
	// yeah hardcoded whatever
	$c->setUserAgent("PHP/7.1.7");
	$c->setHeader("Authorization", "key=" . $token);
	$c->setHeader("Content-Type", "application/json");
	$c->setOpt(CURLOPT_ENCODING, "gzip");
	$c->setOpt(CURLOPT_RETURNTRANSFER, true);
	$c->setOpt(CURLOPT_URL, "https://fcm.googleapis.com/fcm/send");
	$c->setOpt(CURLOPT_POST, true);
	$c->setOpt(CURLOPT_VERBOSE, true);
	$verbose = fopen('/home/howl/meme.txt', 'w');
	$c->setOpt(CURLOPT_STDERR, $verbose);

	$i = 0;
	for (;;) {
		$notif = $settings->encodeNotification($i);
		if ($notif === false)
			return;
		// we can't use $c->post because it sets www-data as post
		$c->setOpt(CURLOPT_POSTFIELDS, $notif);

		$c->_exec();

		// error handling
		if ($c->error) {
			throw new NotifyException("", $c->error_code);
		}
		$i++;
	}
}

class NotifyException extends Exception {}
