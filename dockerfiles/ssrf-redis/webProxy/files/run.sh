#!/bin/bash
service nginx reload
service php7.0-fpm start

service nginx start

while true; do
	service php7.0-fpm restart
	sleep 10
done
