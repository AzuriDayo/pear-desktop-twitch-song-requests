import { describe, it, expect } from "vite-plus/test";
import twitchReducer, { setTwitchInfo, defaultCmdPermissions } from "./twitchSlice";
import type { ITwitchState } from "./twitchSlice";

const initialState: ITwitchState = {
	isLoaded: false,
	expires_in: "",
	twitch_song_request_reward_id: "",
	login: "",
	login_bot: "",
	expires_in_bot: "",
	cmd_permissions: defaultCmdPermissions,
};

describe("twitchStateSlice", () => {
	it("returns the correct initial state", () => {
		const state = twitchReducer(undefined, { type: "@@INIT" });
		expect(state).toEqual(initialState);
		expect(state.isLoaded).toBe(false);
	});

	it("setTwitchInfo sets isLoaded to true", () => {
		const state = twitchReducer(undefined, setTwitchInfo({ login: "streamer" }));
		expect(state.isLoaded).toBe(true);
	});

	it("setTwitchInfo merges partial payload without wiping other fields", () => {
		const state = twitchReducer(
			undefined,
			setTwitchInfo({ login: "streamer", expires_in: "Mon, 01 Jan 2030 00:00:00 GMT" }),
		);
		expect(state.login).toBe("streamer");
		expect(state.expires_in).toBe("Mon, 01 Jan 2030 00:00:00 GMT");
		// Fields not in the payload should keep their initial values
		expect(state.login_bot).toBe("");
		expect(state.twitch_song_request_reward_id).toBe("");
	});

	it("setTwitchInfo updates cmd_permissions", () => {
		const newPerms = {
			cmd_permission_sr: 4,
			cmd_permission_queue: 4,
			cmd_permission_song: 4,
			cmd_permission_delsong: 0,
		};
		const state = twitchReducer(undefined, setTwitchInfo({ cmd_permissions: newPerms }));
		expect(state.cmd_permissions).toEqual(newPerms);
	});

	it("defaultCmdPermissions has correct default levels", () => {
		expect(defaultCmdPermissions.cmd_permission_sr).toBe(3); // subscriber
		expect(defaultCmdPermissions.cmd_permission_queue).toBe(4); // viewer
		expect(defaultCmdPermissions.cmd_permission_song).toBe(4); // viewer
		expect(defaultCmdPermissions.cmd_permission_delsong).toBe(1); // moderator
	});
});
