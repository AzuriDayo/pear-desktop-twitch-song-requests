import { configureStore } from "@reduxjs/toolkit";
import counterReducer from "../features/musicplayer/musicPlayerSlice";

const store = configureStore({
	reducer: {
		musicPlayer: counterReducer,
	},
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export default store;
