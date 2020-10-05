#!/bin/sh
sh clear.sh
echo "Starting $1 instances of image '$2'";
for i in `seq 1 $1`; do sudo docker run -dit --name "u$i" $2 /bin/bash; done
