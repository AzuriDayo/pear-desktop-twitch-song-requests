import { Link } from "react-router";
import { useAppSelector } from "../app/hooks";
import { useEffect, useState, useCallback } from "react";
import { defaultCmdPermissions } from "../features/twitchws/twitchSlice";

const urlPath = "/api/v1/settings";
const method = "PATCH";

const PERMISSION_OPTIONS = [
	{ value: 0, label: "Broadcaster" },
	{ value: 1, label: "Moderator" },
	{ value: 2, label: "VIP" },
	{ value: 3, label: "Subscriber" },
	{ value: 4, label: "Viewer (everyone)" },
];

const CMD_PERMISSION_KEYS = [
	{
		key: "cmd_permission_sr",
		label: "!sr (Song Request)",
		note: undefined,
	},
	{
		key: "cmd_permission_queue",
		label: "!queue",
		note: undefined,
	},
	{
		key: "cmd_permission_song",
		label: "!song (Now Playing)",
		note: undefined,
	},
	{
		key: "cmd_permission_delsong",
		label: "!delsong",
		note: "The original requester can always delete their own song, regardless of this setting.",
	},
] as const;

type CmdPermissionKey = (typeof CMD_PERMISSION_KEYS)[number]["key"];

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

	// Permission state — initialise from defaults, then overwrite with server values on load
	const [permissionValues, setPermissionValues] = useState<Record<CmdPermissionKey, number>>(
		defaultCmdPermissions as Record<CmdPermissionKey, number>,
	);
	const [permissionStatus, setPermissionStatus] = useState("");

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

	// Load current permission values from the server on mount
	useEffect(() => {
		void (async () => {
			try {
				const r = await fetch("/api/v1/settings");
				if (r.ok) {
					const data = (await r.json()) as {
						cmd_permissions?: Record<CmdPermissionKey, number>;
					};
					if (data.cmd_permissions) {
						setPermissionValues((prev) => ({ ...prev, ...data.cmd_permissions }));
					}
				}
			} catch {
				// silently ignore
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

	const savePermissions = useCallback(async () => {
		setPermissionStatus("");
		const body: Record<string, string> = {};
		for (const k of Object.keys(permissionValues) as CmdPermissionKey[]) {
			body[k] = String(permissionValues[k]);
		}
		try {
			const r = await fetch(urlPath, { method, body: JSON.stringify(body) });
			if (r.status >= 200 && r.status < 300) {
				setPermissionStatus("Permissions saved successfully!");
			} else {
				const msg = (await r.json()) as { error?: string };
				setPermissionStatus("Save failed: " + (msg.error ?? "unknown error"));
			}
		} catch {
			setPermissionStatus("Network error, please try again.");
		}
	}, [permissionValues]);

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

			{/* ── Command Permissions ─────────────────────────────────────────── */}
			<hr />
			<h2>Command Permissions</h2>
			<p>Set the minimum role required to use each chat command.</p>
			<form
				onSubmit={(e) => {
					e.preventDefault();
					void savePermissions();
				}}
			>
				<table style={{ borderSpacing: "12px 8px", borderCollapse: "separate" }}>
					<tbody>
						{CMD_PERMISSION_KEYS.map(({ key, label, note }) => (
							<tr key={key}>
								<td>
									<label htmlFor={`perm-${key}`}>
										<code>{label}</code>
									</label>
								</td>
								<td>
									<select
										id={`perm-${key}`}
										value={permissionValues[key]}
										onChange={(e) => {
											setPermissionValues((prev) => ({
												...prev,
												[key]: Number(e.target.value),
											}));
										}}
									>
										{PERMISSION_OPTIONS.map((opt) => (
											<option key={opt.value} value={opt.value}>
												{opt.label}
											</option>
										))}
									</select>
								</td>
								{note && (
									<td>
										<small style={{ color: "#aaa" }}>{note}</small>
									</td>
								)}
							</tr>
						))}
					</tbody>
				</table>
				<br />
				<button type="submit">Save Permissions</button>
			</form>
			{permissionStatus && <h3>{permissionStatus}</h3>}

			<br />
			<br />
			<br />
			<Link to="/">Back to home</Link>
		</>
	);
}
