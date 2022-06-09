#!/bin/sh

[ -f "./.env" ] && . ./.env
[ "$DISCORD_TOKEN" = "" ] && echo "Missing Discord Token" && exit 1
[ "$APPLICATION_ID" = "" ] && echo "Missing application ID" && exit 1

curl -i -X 'PUT' "https://discord.com/api/applications/$APPLICATION_ID/commands" \
    -H "Authorization: Bot $DISCORD_TOKEN" \
    -H 'Content-Type: application/json' \
    -H 'Accept: application/json' \
    -d '[]'
