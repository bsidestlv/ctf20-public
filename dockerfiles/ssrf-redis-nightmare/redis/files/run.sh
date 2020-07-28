#!/bin/bash
/root/persistence.sh &
service nginx start
cd /home/admin/
su admin -c redis-server
