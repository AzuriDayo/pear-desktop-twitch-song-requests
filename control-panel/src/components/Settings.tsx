import { Link } from "react-router";
import { useAppSelector } from "../app/hooks";
import { useEffect, useState } from "react";

const urlPath = "/api/v1/settings";
const method = "PATCH";

export function Settings() {
	const twitchState = useAppSelector((state) => state.twitchState);
	const [twitchRewardId, setTwitchRewardId] = useState("");
	const [settings, setSettings] = useState<{ [key: string]: string }>({});
	const [status, setStatus] = useState("");

	useEffect(() => {
		if (
			twitchRewardId === "" &&
			twitchState.twitch_song_request_reward_id != ""
		) {
			setTwitchRewardId(twitchState.twitch_song_request_reward_id);
		}
	}, [twitchState.twitch_song_request_reward_id, twitchRewardId]);

	useEffect(() => {
		if (Object.keys(settings).length > 0) {
			fetch(urlPath, {
				method,
				body: JSON.stringify(settings),
			})
				.then((response) => {
					if (response.status >= 200 && response.status < 300) {
						setStatus("Settings saved successfully!");
						setSettings({});
						return Promise.resolve("");
					}
					return response.text();
				})
				.then((text) => {
					if (text == "") return;
					let msg = {
						["error"]: "",
					};
					try {
						if (text != "") {
							msg = JSON.parse(text);
							setStatus(
								"Settings save failed with error: " + (msg.error ?? ""),
							);
						}
					} catch (e) {
						console.log(e);
					}
				});
		}
	}, [settings]);

	return (
		<>
			<form
				onSubmit={(e) => {
					e.preventDefault();
					setSettings({
						twitch_song_request_reward_id: twitchRewardId,
					});
				}}
			>
				<label htmlFor="reward-id">Twitch Reward ID: </label>
				<input
					name="reward-id"
					type="text"
					onChange={(e) => {
						setTwitchRewardId(e.target.value);
					}}
					value={twitchRewardId}
					autoComplete="off"
				/>
				<br />
				<button type="submit">save</button>
			</form>
			{status && <h3>{status}</h3>}
			<br />
			<br />
			<br />
			<Link to="/">Back to home</Link>
		</>
	);
}
