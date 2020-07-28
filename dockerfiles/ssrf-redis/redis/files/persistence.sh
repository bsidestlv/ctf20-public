#!/bin/bash
while true; do
	find /var/www/html/ -type f -not -name 'index.php' -print0 | xargs -0 rm -fr
	cp /root/index.php /var/www/html/
	echo "flushall" | redis-cli
	sleep 10
done
