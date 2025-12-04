import { configureStore } from "@reduxjs/toolkit";
import musicPlayerReducer from "../features/musicplayer/musicPlayerSlice";
import websocketReducer from "../features/websocket/websocketSlice";

const store = configureStore({
	reducer: {
		musicPlayer: musicPlayerReducer,
		websocket: websocketReducer,
	},
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export default store;
