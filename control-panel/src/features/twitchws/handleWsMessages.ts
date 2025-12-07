import { Dispatch, UnknownAction } from "@reduxjs/toolkit";
import { setTwitchInfo } from "./twitchSlice";

export const handleWsMessages = (
	data: string,
	dispatch: Dispatch<UnknownAction>,
) => {
	// Change this later its ugly af

	const d: MsgTwitchInfo = JSON.parse(data);
	console.log(d);
	dispatch(
		setTwitchInfo({
			expires_in: d.expiry_date,
			twitch_song_request_reward_id: d.reward_id,
			login: d.login,
		}),
	);
};

interface MsgTwitchInfo {
	type: "TWITCH_INFO";
	login: string;
	expiry_date: string;
	stream_online: string;
	reward_id: string;
}
