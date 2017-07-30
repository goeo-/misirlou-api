<?php

// State contains the information about the application and is passed to every
// method function.
class State {
	/**
	 * SQL database from which to retrieve data.
	 *
	 * @var DBPDO
	 */
	public $db;

	/**
	 * @param DBPDO $db
	 */
	public function __construct($db)
	{
		$this->db = $db;
	}

	public function getAccessToken()
	{
		$at = $this->db->fetch("SELECT access_token FROM sessions WHERE id = ?", [hash("sha256", @$_SERVER["HTTP_AUTHORIZATION"])]);
		return @$at["access_token"];
	}

	public function getSelfID()
	{
		$at = $this->db->fetch("SELECT user_id FROM sessions WHERE id = ?", [hash("sha256", @$_SERVER["HTTP_AUTHORIZATION"])]);
		if (!$at) return 0;
		return @$at["user_id"];
	}
}
