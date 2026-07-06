import { useEffect } from "react";
import { useAppDispatch } from "../../app/hooks";
import { setTwitchInfo } from "./twitchSlice";
import { addSong, setQueueInfo, shiftQueue, type SongQueueItem } from "./songQueueSlice";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { GetQueue, GetTwitchInfo } from "../../wailsjs/go/main/App";

interface TwitchInfoEvent {
	expiry_date: string;
	reward_id: string;
	login: string;
	login_bot: string;
	expiry_date_bot: string;
	refresh_expiry_date: string;
	refresh_expiry_date_bot: string;
}

export function TwitchWS() {
	const dispatch = useAppDispatch();

	useEffect(() => {
		const applyTwitchInfo = (info: TwitchInfoEvent) => {
			dispatch(
				setTwitchInfo({
					expires_in: info.expiry_date,
					twitch_song_request_reward_id: info.reward_id,
					login: info.login,
					login_bot: info.login_bot,
					expires_in_bot: info.expiry_date_bot,
					refresh_expires_in: info.refresh_expiry_date,
					refresh_expires_in_bot: info.refresh_expiry_date_bot,
				}),
			);
		};

		// Seed initial state (replaces the push previously sent on websocket connect).
		void GetTwitchInfo().then(applyTwitchInfo);
		void GetQueue().then((q) => {
			dispatch(setQueueInfo({ song_queue: (q as unknown as SongQueueItem[]) ?? [] }));
		});

		// Live updates are delivered via the Wails runtime event system.
		const offTwitchInfo = EventsOn("TWITCH_INFO", (info: TwitchInfoEvent) => {
			applyTwitchInfo(info);
		});
		const offQueueInfo = EventsOn("QUEUE_INFO", (payload: { song_queue: SongQueueItem[] }) => {
			dispatch(setQueueInfo({ song_queue: payload.song_queue ?? [] }));
		});
		const offQueueAdd = EventsOn("QUEUE_ADD", (song: SongQueueItem) => {
			dispatch(addSong({ song }));
		});
		const offQueueShift = EventsOn("QUEUE_SHIFT", () => {
			dispatch(shiftQueue());
		});

		return () => {
			offTwitchInfo();
			offQueueInfo();
			offQueueAdd();
			offQueueShift();
		};
	}, [dispatch]);

	return <></>;
}
