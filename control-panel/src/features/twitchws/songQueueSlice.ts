import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { RootState } from "../../app/store";

export interface SongQueueItem {
	requested_by: string;
	song: {
		title: string;
		artist: string;
		videoId: string;
		imageUrl: string;
	};
}

export interface ISongQueueState {
	isLoaded: boolean;
	queue: SongQueueItem[];
}

const initialState: ISongQueueState = {
	isLoaded: false,
	queue: [],
};

export const songQueueSlice = createSlice({
	name: "queuestate",
	initialState,
	reducers: {
		setQueueInfo: (
			state,
			action: PayloadAction<{ queue: SongQueueItem[] }>,
		) => {
			state.isLoaded = true;
			state.queue = action.payload.queue;
		},
	},
});

export const { setQueueInfo } = songQueueSlice.actions;

export const selectQueueState = (state: RootState) => state.songQueueState;

export default songQueueSlice.reducer;
