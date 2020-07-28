#!/bin/bash
/root/persistence.sh &
service nginx start
service php7.0-fpm start
cd /home/admin/
su admin -c redis-server
