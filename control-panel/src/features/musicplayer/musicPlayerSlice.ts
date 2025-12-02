import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import type { RootState } from "../../app/store";

// Define a type for the slice state
interface IMusicPlayerState {
	isPlaying: boolean;
}

const initialState: IMusicPlayerState = {
	isPlaying: false,
};

export const musicPlayerSlice = createSlice({
	name: "musicplayer",
	initialState,
	reducers: {
		setPlaying: (state, action: PayloadAction<boolean>) => {
			state.isPlaying = action.payload;
		},
	},
});

export const { setPlaying } = musicPlayerSlice.actions;

export const selectMusicState = (state: RootState) =>
	state.musicPlayer.isPlaying;

export default musicPlayerSlice.reducer;
