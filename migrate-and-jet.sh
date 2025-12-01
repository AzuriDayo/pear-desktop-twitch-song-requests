#!/bin/bash
~/go/bin/migrate -path internal/data/iofs/migrations -database sqlite3://pear-desktop-twitch-song-requests.db up
~/go/bin/jet -source=sqlite -dsn="file:pear-desktop-twitch-song-requests.db" -path=./gen -ignore-tables=schema_migrations
# we ignore schema_migrations because it's used by golang-migrate and uses an incompatible format
