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
	is_ninja: boolean;
}

export interface ISongQueueState {
	isLoaded: boolean;
	song_queue: SongQueueItem[];
}

const initialState: ISongQueueState = {
	isLoaded: false,
	song_queue: [],
};

export const songQueueSlice = createSlice({
	name: "queuestate",
	initialState,
	reducers: {
		setQueueInfo: (
			state,
			action: PayloadAction<{ song_queue: SongQueueItem[] }>,
		) => {
			state.isLoaded = true;
			if (!action.payload.song_queue) state.song_queue = [];
			else state.song_queue = action.payload.song_queue;
		},
		addSong: (state, action: PayloadAction<{ song: SongQueueItem }>) => {
			state.song_queue.push(action.payload.song);
		},
		shiftQueue: (state) => {
			state.song_queue.shift();
		},
	},
});

export const { setQueueInfo, addSong, shiftQueue } = songQueueSlice.actions;

export const selectQueueState = (state: RootState) => state.songQueueState;

export default songQueueSlice.reducer;
