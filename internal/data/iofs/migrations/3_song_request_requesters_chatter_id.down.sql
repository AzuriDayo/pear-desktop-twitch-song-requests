DROP INDEX IF EXISTS requester_user_id_idx;
ALTER TABLE song_request_requesters DROP is_ninja;
ALTER TABLE song_request_requesters DROP user_id;
DROP TABLE IF EXISTS data_transforms;