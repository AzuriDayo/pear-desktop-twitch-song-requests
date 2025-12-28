ALTER TABLE song_request_requesters
ADD user_id TEXT NOT NULL DEFAULT "";
ALTER TABLE song_request_requesters
ADD is_ninja BOOLEAN NOT NULL DEFAULT false;
CREATE INDEX requester_user_id_idx ON song_request_requesters (user_id);
CREATE TABLE data_transforms (
    key TEXT PRIMARY KEY,
    value BOOLEAN NOT NULL
) WITHOUT ROWID;
INSERT INTO data_transforms
VALUES ("BACKPORT_REQUESTERS_USER_ID", false);