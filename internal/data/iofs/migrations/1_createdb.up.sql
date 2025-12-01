CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
) WITHOUT ROWID;
CREATE TABLE song_requests (
    id INTEGER PRIMARY KEY,
    user_id TEXT NOT NULL,
    song_title TEXT NOT NULL,
    artist_name TEXT,
    video_id TEXT,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    removed_at TEXT DEFAULT NULL
);