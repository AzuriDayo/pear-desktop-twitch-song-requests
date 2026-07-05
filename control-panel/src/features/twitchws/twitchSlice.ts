import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { RootState } from "../../app/store";

export interface CmdPermissions {
	cmd_permission_sr: number;
	cmd_permission_queue: number;
	cmd_permission_song: number;
	cmd_permission_delsong: number;
}

// Define a type for the slice state
export interface ITwitchState {
	isLoaded: boolean;
	expires_in: string;
	twitch_song_request_reward_id: string;
	login: string;
	expires_in_bot: string;
	login_bot: string;
	cmd_permissions: CmdPermissions;
}

export const defaultCmdPermissions: CmdPermissions = {
	cmd_permission_sr: 3, // subscriber
	cmd_permission_queue: 4, // viewer
	cmd_permission_song: 4, // viewer
	cmd_permission_delsong: 1, // moderator
};

const initialState: ITwitchState = {
	isLoaded: false,
	expires_in: "",
	twitch_song_request_reward_id: "",
	login: "",
	login_bot: "",
	expires_in_bot: "",
	cmd_permissions: defaultCmdPermissions,
};

export const twitchStateSlice = createSlice({
	name: "twitchstate",
	initialState,
	reducers: {
		setTwitchInfo: (state, action: PayloadAction<Partial<ITwitchState>>) => {
			state.isLoaded = true;
			Object.assign(state, action.payload);
		},
	},
});

export const { setTwitchInfo } = twitchStateSlice.actions;

export const selectTwitchState = (state: RootState) => state.twitchState;

export default twitchStateSlice.reducer;
