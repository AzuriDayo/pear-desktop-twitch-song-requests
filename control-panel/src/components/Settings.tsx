import { Link } from "react-router";
import { useAppSelector } from "../app/hooks";
import { useEffect, useState } from "react";

const urlPath = "/api/v1/settings";
const method = "PATCH";

export function Settings() {
	const twitchState = useAppSelector((state) => state.twitchState);
	const [twitchRewardId, setTwitchRewardId] = useState<{
		value: string | null;
		loaded: boolean;
	}>({ value: "", loaded: false });
	const [settings, setSettings] = useState<{ [key: string]: string }>({});
	const [status, setStatus] = useState("");
	const [availableRewards, setAvailableRewards] = useState<
		{ id: string; cost: number; name: string }[]
	>([]);

	useEffect(() => {
		(async () => {
			const r = await fetch("/api/v1/twitch/custom-rewards");
			const rews = await r.json();
			setAvailableRewards(rews);
		})();
	}, []);

	useEffect(() => {
		if (!twitchRewardId.loaded && twitchState.isLoaded) {
			const v = {
				value: twitchState.twitch_song_request_reward_id,
				loaded: true,
			};
			console.log(v);
			setTwitchRewardId(v);
		}
	}, [twitchState]);

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
			<h2>Configure Twitch Custom Reward</h2>
			<p>
				If you have not created a custom reward yet, create one on your{" "}
				<a
					target="_blank"
					href={`https://dashboard.twitch.tv/u/${twitchState.login}/viewer-rewards/channel-points/rewards`}
				>
					Creator Dashboard
				</a>{" "}
				and you must enable the toggle for <i>Require Viewer to Enter Text</i>.
			</p>
			<p>
				Once you have created your Custom Reward, you can refresh the dropdown
				below.
			</p>
			<br />
			<form
				onSubmit={(e) => {
					e.preventDefault();
					setSettings({
						twitch_song_request_reward_id: twitchRewardId.value ?? "",
					});
				}}
			>
				<label htmlFor="reward-id">Twitch Custom Reward: </label>
				<br />
				<div
					style={{
						display: "flex",
						flexDirection: "row",
						alignItems: "center",
						justifyContent: "center",
					}}
				>
					<select
						id="reward-id"
						disabled={!twitchRewardId.loaded}
						style={{ minWidth: "50vw", minHeight: "4vh" }}
						onChange={(e) => {
							const v = e.target.value;
							setTwitchRewardId({ value: v, loaded: true });
						}}
						value={twitchRewardId.value ?? ""}
					>
						<option key="0" value={""}>
							Select reward...
						</option>
						{availableRewards.map((r, i) => (
							<option key={i + 1} value={r.id}>
								{r.name} - {r.cost}
							</option>
						))}
					</select>
					<button
						disabled={!twitchRewardId.loaded}
						onClick={() => {
							(async () => {
								let v = {
									value: twitchRewardId.value,
									loaded: false,
								};
								setTwitchRewardId(v);
								setStatus("");
								try {
									const r = await fetch("/api/v1/twitch/custom-rewards");
									const rews = await r.json();
									setAvailableRewards(rews);
									setStatus("Done refresh");
								} catch (e) {
									setStatus("Failed refresh try again later");
								}
								v.loaded = true;
								setTwitchRewardId(v);
							})();
						}}
					>
						Refresh
					</button>
				</div>
				<br />
				<button type="submit">Save</button>
			</form>
			{status && <h3>{status}</h3>}
			<br />
			<br />
			<br />
			<Link to="/">Back to home</Link>
		</>
	);
}
