INSERT INTO settings (key, value) VALUES ('cmd_permission_sr', '3') ON CONFLICT(key) DO NOTHING;
INSERT INTO settings (key, value) VALUES ('cmd_permission_queue', '4') ON CONFLICT(key) DO NOTHING;
INSERT INTO settings (key, value) VALUES ('cmd_permission_song', '4') ON CONFLICT(key) DO NOTHING;
INSERT INTO settings (key, value) VALUES ('cmd_permission_delsong', '1') ON CONFLICT(key) DO NOTHING;
