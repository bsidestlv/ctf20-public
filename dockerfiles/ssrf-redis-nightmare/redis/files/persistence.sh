#!/bin/bash
while true; do
	echo "flushall" | redis-cli
	echo 'set flag "BsidesTLV2020{Y0u_4r3_A_Genius!!}"'| redis-cli
	sleep 10
done