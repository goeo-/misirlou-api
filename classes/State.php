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
}
