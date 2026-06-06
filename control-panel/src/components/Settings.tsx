import { Link } from "react-router";
import { useAppSelector } from "../app/hooks";
import { useEffect, useState, useCallback } from "react";

const urlPath = "/api/v1/settings";
const method = "PATCH";

export function Settings() {
	const twitchState = useAppSelector((state) => state.twitchState);

	// Track the user's in-progress selection separately from the persisted value.
	// When null, the current persisted value from Redux is used (no useEffect needed).
	const [selectedRewardId, setSelectedRewardId] = useState<string | null>(null);
	const rewardId =
		selectedRewardId ?? (twitchState.isLoaded ? twitchState.twitch_song_request_reward_id : "");

	const [isRefreshing, setIsRefreshing] = useState(false);
	const [settings, setSettings] = useState<{ [key: string]: string }>({});
	const [status, setStatus] = useState("");
	const [availableRewards, setAvailableRewards] = useState<
		{ id: string; cost: number; name: string }[]
	>([]);

	const fetchRewards = useCallback(async () => {
		setIsRefreshing(true);
		setStatus("");
		try {
			const r = await fetch("/api/v1/twitch/custom-rewards");
			setAvailableRewards(await r.json());
		} catch {
			setStatus("Failed to refresh, try again later");
		} finally {
			setIsRefreshing(false);
		}
	}, []);

	useEffect(() => {
		void (async () => {
			setIsRefreshing(true);
			try {
				const r = await fetch("/api/v1/twitch/custom-rewards");
				setAvailableRewards(await r.json());
			} catch {
				// silently ignore initial load failures
			} finally {
				setIsRefreshing(false);
			}
		})();
	}, []);

	useEffect(() => {
		if (Object.keys(settings).length > 0) {
			void fetch(urlPath, {
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
					try {
						if (text != "") {
							const msg: { error?: string } = JSON.parse(text);
							setStatus("Settings save failed with error: " + (msg.error ?? ""));
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
			<p>Once you have created your Custom Reward, you can refresh the dropdown below.</p>
			<br />
			<form
				onSubmit={(e) => {
					e.preventDefault();
					setSettings({
						twitch_song_request_reward_id: rewardId ?? "",
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
						disabled={!twitchState.isLoaded || isRefreshing}
						style={{ minWidth: "50vw", minHeight: "4vh" }}
						onChange={(e) => {
							setSelectedRewardId(e.target.value);
						}}
						value={rewardId ?? ""}
					>
						<option key="0" value={""}>
							Select reward...
						</option>
						{availableRewards.map((r) => (
							<option key={r.id} value={r.id}>
								{r.name} - {r.cost}
							</option>
						))}
					</select>
					<button
						disabled={!twitchState.isLoaded || isRefreshing}
						onClick={() => {
							setStatus("");
							void fetchRewards().then(() => setStatus("Done refresh"));
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
