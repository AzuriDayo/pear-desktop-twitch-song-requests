import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import type { RootState } from "../../app/store";

// Define a type for the slice state
interface IMusicPlayerState {
	isPlaying: boolean;
	currentSong?: string;
	artist?: string;
	url?: string;
	songDuration?: number;
	imageSrc?: string;
	elapsedSeconds?: number;
	volume?: number;
	videoId?: string;
	serviceHealthy: boolean;
	lastUpdated?: string;
}

const initialState: IMusicPlayerState = {
	isPlaying: false,
	serviceHealthy: false,
};

export const musicPlayerSlice = createSlice({
	name: "musicplayer",
	initialState,
	reducers: {
		setPlaying: (state, action: PayloadAction<boolean>) => {
			state.isPlaying = action.payload;
		},
		setCurrentSong: (state, action: PayloadAction<string>) => {
			state.currentSong = action.payload;
		},
		setVolume: (state, action: PayloadAction<number>) => {
			state.volume = action.payload;
		},
		setServiceHealth: (state, action: PayloadAction<boolean>) => {
			state.serviceHealthy = action.payload;
		},
		updatePlayerState: (
			state,
			action: PayloadAction<{
				isPlaying?: boolean;
				currentSong?: string;
				artist?: string;
				url?: string;
				songDuration?: number;
				imageSrc?: string;
				elapsedSeconds?: number;
				volume?: number;
				videoId?: string;
			}>,
		) => {
			if (action.payload.isPlaying !== undefined) {
				state.isPlaying = action.payload.isPlaying;
			}
			if (action.payload.currentSong !== undefined) {
				state.currentSong = action.payload.currentSong;
			}
			if (action.payload.artist !== undefined) {
				state.artist = action.payload.artist;
			}
			if (action.payload.url !== undefined) {
				state.url = action.payload.url;
			}
			if (action.payload.songDuration !== undefined) {
				state.songDuration = action.payload.songDuration;
			}
			if (action.payload.imageSrc !== undefined) {
				state.imageSrc = action.payload.imageSrc;
			}
			if (action.payload.elapsedSeconds !== undefined) {
				state.elapsedSeconds = action.payload.elapsedSeconds;
			}
			if (action.payload.volume !== undefined) {
				state.volume = action.payload.volume;
			}
			if (action.payload.videoId !== undefined) {
				state.videoId = action.payload.videoId;
			}
			state.lastUpdated = new Date().toISOString();
		},
	},
});

export const {
	setPlaying,
	setCurrentSong,
	setVolume,
	setServiceHealth,
	updatePlayerState,
} = musicPlayerSlice.actions;

export const selectMusicState = (state: RootState) => state.musicPlayer;
export const selectIsPlaying = (state: RootState) =>
	state.musicPlayer.isPlaying;
export const selectCurrentSong = (state: RootState) =>
	state.musicPlayer.currentSong;
export const selectArtist = (state: RootState) => state.musicPlayer.artist;
export const selectSongURL = (state: RootState) => state.musicPlayer.url;
export const selectSongDuration = (state: RootState) =>
	state.musicPlayer.songDuration;
export const selectImageSrc = (state: RootState) => state.musicPlayer.imageSrc;
export const selectElapsedSeconds = (state: RootState) =>
	state.musicPlayer.elapsedSeconds;
export const selectVolume = (state: RootState) => state.musicPlayer.volume;
export const selectVideoId = (state: RootState) => state.musicPlayer.videoId;
export const selectServiceHealth = (state: RootState) =>
	state.musicPlayer.serviceHealthy;

export default musicPlayerSlice.reducer;
