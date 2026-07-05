import { describe, it, expect } from "vite-plus/test";
import songQueueReducer, {
	setQueueInfo,
	addSong,
	shiftQueue,
	removeSongAtIndex,
} from "./songQueueSlice";
import type { SongQueueItem } from "./songQueueSlice";

const makeSong = (title: string): SongQueueItem => ({
	requested_by: "user1",
	song: { title, artist: "Artist", videoId: "abc123", imageUrl: "" },
	is_ninja: false,
});

describe("songQueueSlice", () => {
	it("returns correct initial state", () => {
		const state = songQueueReducer(undefined, { type: "@@INIT" });
		expect(state.isLoaded).toBe(false);
		expect(state.song_queue).toEqual([]);
	});

	it("setQueueInfo sets the queue and marks isLoaded", () => {
		const songs = [makeSong("Song A"), makeSong("Song B")];
		const state = songQueueReducer(undefined, setQueueInfo({ song_queue: songs }));
		expect(state.isLoaded).toBe(true);
		expect(state.song_queue).toHaveLength(2);
		expect(state.song_queue[0].song.title).toBe("Song A");
	});

	it("setQueueInfo with null/missing payload results in empty array", () => {
		// eslint-disable-next-line @typescript-eslint/no-explicit-any
		const state = songQueueReducer(undefined, setQueueInfo({ song_queue: null as any }));
		expect(state.song_queue).toEqual([]);
	});

	it("addSong appends a song to the queue", () => {
		const initial = songQueueReducer(undefined, setQueueInfo({ song_queue: [makeSong("Song A")] }));
		const state = songQueueReducer(initial, addSong({ song: makeSong("Song B") }));
		expect(state.song_queue).toHaveLength(2);
		expect(state.song_queue[1].song.title).toBe("Song B");
	});

	it("shiftQueue removes the first item", () => {
		const initial = songQueueReducer(
			undefined,
			setQueueInfo({ song_queue: [makeSong("Song A"), makeSong("Song B")] }),
		);
		const state = songQueueReducer(initial, shiftQueue());
		expect(state.song_queue).toHaveLength(1);
		expect(state.song_queue[0].song.title).toBe("Song B");
	});

	it("shiftQueue on empty queue does not throw", () => {
		const state = songQueueReducer(undefined, shiftQueue());
		expect(state.song_queue).toEqual([]);
	});

	it("removeSongAtIndex removes item at the given index", () => {
		const initial = songQueueReducer(
			undefined,
			setQueueInfo({
				song_queue: [makeSong("Song A"), makeSong("Song B"), makeSong("Song C")],
			}),
		);
		const state = songQueueReducer(initial, removeSongAtIndex({ index: 1 }));
		expect(state.song_queue).toHaveLength(2);
		expect(state.song_queue[0].song.title).toBe("Song A");
		expect(state.song_queue[1].song.title).toBe("Song C");
	});

	it("removeSongAtIndex at index 0 removes the first item", () => {
		const initial = songQueueReducer(
			undefined,
			setQueueInfo({ song_queue: [makeSong("Song A"), makeSong("Song B")] }),
		);
		const state = songQueueReducer(initial, removeSongAtIndex({ index: 0 }));
		expect(state.song_queue).toHaveLength(1);
		expect(state.song_queue[0].song.title).toBe("Song B");
	});
});
