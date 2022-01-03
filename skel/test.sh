#!/bin/bash
#
# Dacă ați configurat avahi pentru mdns, this should work...
#
HOSTNAME=tema2.local

ecurl()
{
    curl -H "Content-Type: application/json" "$@"
}

submarine_move()
{
    echo "POST /api/submarine/move '$1'"
    ecurl -X POST -d "$1" http://$HOSTNAME/api/submarine/move
    sleep 1
}

add_fish()
{
    echo "POST /api/fish/add '$1'"
    ecurl -X POST -d "$1" http://$HOSTNAME/api/fish/add
    sleep 1
}

update_artifact()
{
    echo "POST /api/artifact/update '$1'"
    ecurl -X POST -d "$1" http://$HOSTNAME/api/artifact/update
    sleep 1
}

echo "GET /api/submarine"
ecurl http://$HOSTNAME/api/submarine
sleep 1

for i in $(seq 1 5); do
    submarine_move '{"x": 1, "y": 1}'
done

add_fish '{"x": 1, "y": 2}'
add_fish '{"x": 15, "y": 13}'
add_fish '{"x": 30, "y": 16}'
add_fish '{"x": 40, "y": 1}'
update_artifact '{"x": 10, "y": 13}'

for i in $(seq 1 5); do
    submarine_move '{"x": -1, "y": 0}'
done

update_artifact '{}'

update_artifact '{"x": 20, "y": 4}'
