# Misirlou API

>Three tomatoes are walking down the street - a poppa tomato, a momma tomato,
>and a little baby tomato. Baby tomato starts lagging behind. Poppa tomato gets
>angry, goes over to the baby tomato, and smooshes him... and says, 'ketchup!'

Misirlou is Ripple's system for managing tournaments. This is the API/backend,
and why write it in PHP, you may ask? Because I thought that the thing wouldn't
have grown and that writing a few PHP files would have been enough. What a fool!

To get started, simply copy config.sample.php to config.php, and edit the
values in `$config`. Everything should be pretty obvious and need no
explanation.

Done that, create/update the database schema by running from the command line
`migrate.php`. Finally, create an nginx config so that all requests are routed
to /router.php:

```nginx
server {
	autoindex on;
	listen 80;
	server_name quarterpounderwithcheese.org;
	root /home/howl/oc/misirlou-api;
	charset utf-8;
	location / {
		try_files $uri $uri/ /router.php;
	}
	include php;
}
```

(The configuration above is actually not recommended in production - you should
at least remove the autoindex.)
