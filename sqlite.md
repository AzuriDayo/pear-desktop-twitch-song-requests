# SQLite readme
Browse data with:
```sh
sqlite-utils tables pear-desktop-twitch-song-requests.db
```

running migrations and jet
```sh
# ensure ~/go/bin is in path

# migrate
go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@v4
migrate -path internal/data/iofs/migrations -database sqlite3://pear-desktop-twitch-song-requests.db up

# jet
go install github.com/go-jet/jet/v2/cmd/jet@latest
jet -source=sqlite -dsn="file:pear-desktop-twitch-song-requests.db" -path=./gen -ignore-tables=schema_migrations
```

See also `migrage-and-jet.sh` script.
