import { describe, it, expect, vi } from "vite-plus/test";
import { handleWsMessages } from "./handleWsMessages";
import { setTwitchInfo } from "./twitchSlice";
import { addSong, setQueueInfo, shiftQueue } from "./songQueueSlice";

describe("handleWsMessages", () => {
	it("dispatches setTwitchInfo for TWITCH_INFO messages", () => {
		const dispatch = vi.fn();
		const msg = JSON.stringify({
			type: "TWITCH_INFO",
			login: "streamer",
			login_bot: "mybot",
			expiry_date: "Mon, 01 Jan 2030 00:00:00 GMT",
			expiry_date_bot: "Mon, 01 Jan 2030 00:00:00 GMT",
			stream_online: false,
			reward_id: "reward-123",
		});

		handleWsMessages(msg, dispatch);

		expect(dispatch).toHaveBeenCalledOnce();
		expect(dispatch).toHaveBeenCalledWith(
			setTwitchInfo({
				expires_in: "Mon, 01 Jan 2030 00:00:00 GMT",
				twitch_song_request_reward_id: "reward-123",
				login: "streamer",
				login_bot: "mybot",
				expires_in_bot: "Mon, 01 Jan 2030 00:00:00 GMT",
			}),
		);
	});

	it("dispatches setQueueInfo for QUEUE_INFO messages", () => {
		const dispatch = vi.fn();
		const queue = [
			{
				requested_by: "user1",
				song: { title: "Song A", artist: "Artist", videoId: "abc", imageUrl: "" },
				is_ninja: false,
			},
		];
		const msg = JSON.stringify({ type: "QUEUE_INFO", song_queue: queue });

		handleWsMessages(msg, dispatch);

		expect(dispatch).toHaveBeenCalledOnce();
		expect(dispatch).toHaveBeenCalledWith(setQueueInfo({ song_queue: queue }));
	});

	it("dispatches shiftQueue for QUEUE_SHIFT messages", () => {
		const dispatch = vi.fn();
		const msg = JSON.stringify({ type: "QUEUE_SHIFT" });

		handleWsMessages(msg, dispatch);

		expect(dispatch).toHaveBeenCalledOnce();
		expect(dispatch).toHaveBeenCalledWith(shiftQueue());
	});

	it("dispatches addSong for QUEUE_ADD messages", () => {
		const dispatch = vi.fn();
		const song = {
			requested_by: "user2",
			song: { title: "Song B", artist: "Artist", videoId: "xyz", imageUrl: "" },
			is_ninja: false,
		};
		const msg = JSON.stringify({ type: "QUEUE_ADD", song });

		handleWsMessages(msg, dispatch);

		expect(dispatch).toHaveBeenCalledOnce();
		expect(dispatch).toHaveBeenCalledWith(addSong({ song }));
	});

	it("does not dispatch for unknown message types", () => {
		const dispatch = vi.fn();
		const msg = JSON.stringify({ type: "UNKNOWN_TYPE" });

		handleWsMessages(msg, dispatch);

		expect(dispatch).not.toHaveBeenCalled();
	});
});
