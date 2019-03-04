<?php

$start = microtime(true);

if (php_sapi_name() != "cli")
	exit("Execute from command line");

require_once __DIR__ . "/include.php";

if (!$db->table_exists("db_version")) {
	$db->execute("CREATE TABLE db_version(version INT NOT NULL)");
	$db->execute("INSERT INTO db_version(version) VALUES (?)", [-1]);
}

$version = $db->fetch("SELECT version FROM db_version")["version"];

for (;;) {
	$version++;
	if (!file_exists("migrations/$version.sql")) {
		echo "Migrated in " . ((microtime(true)) - $start) . " seconds\n";
		return;
	}
	echo "Running $version.sql... ";
	$db->execute(file_get_contents("migrations/$version.sql"));
	echo "done.\n";
	$db->execute("UPDATE db_version SET version = ?", [$version]);
}

