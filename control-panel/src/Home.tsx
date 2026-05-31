import { Link } from "react-router";
import { useAppSelector } from "./app/hooks";
import { useMemo } from "react";
import { isAfter, addDays } from "date-fns";

const expiryState = {
	OK: "OK",
	SOON: "SOON",
	EXPIRED: "EXPIRED",
} as const;

const getExpiryStateEmoji = (state: EExpiryState): string => {
	switch (state) {
		case expiryState.OK:
			return "🟢";
		case expiryState.EXPIRED:
			return "❌";
		case expiryState.SOON:
			return "⚠️";
	}
};

type EExpiryState = (typeof expiryState)[keyof typeof expiryState];
// type TExpiryState = keyof typeof expiryState;

function computeExpiryState(expiresIn: string): EExpiryState {
	if (expiresIn === "") return expiryState.OK;
	let expiry: Date;
	try {
		expiry = new Date(expiresIn);
	} catch {
		return expiryState.OK;
	}
	const now = new Date();
	if (isAfter(now, expiry)) return expiryState.EXPIRED;
	if (isAfter(addDays(now, 15), expiry)) return expiryState.SOON;
	return expiryState.OK;
}

export function Home() {
	const twitchState = useAppSelector((state) => state.twitchState);

	const userExpiryState = useMemo(
		() => computeExpiryState(twitchState.expires_in),
		[twitchState.expires_in],
	);

	const botExpiryState = useMemo(
		() => computeExpiryState(twitchState.expires_in_bot),
		[twitchState.expires_in_bot],
	);

	return twitchState.isLoaded ? (
		<div>
			<Link to="/queue">View song queue and history</Link>
			<br />
			<br />
			<br />
			<Link to="/oauth/twitch-connect">
				{twitchState.login !== ""
					? "Refresh Twitch token"
					: "Connect with twitch"}
			</Link>
			<h3>
				{twitchState.expires_in == ""
					? "No Twitch token configured"
					: "Twitch token for " +
						twitchState.login +
						" expires on " +
						twitchState.expires_in +
						" " +
						getExpiryStateEmoji(userExpiryState)}
			</h3>
			<br />
			<Link to="/oauth/twitch-connect-bot">
				{twitchState.login_bot !== ""
					? "Refresh Twitch bot token"
					: "Connect twitch bot account"}
			</Link>
			<h3>
				{twitchState.expires_in_bot == ""
					? "No bot Twitch token configured"
					: "Twitch token for " +
						twitchState.login_bot +
						" expires on " +
						twitchState.expires_in_bot +
						" " +
						getExpiryStateEmoji(botExpiryState)}
			</h3>
			<br />
			<br />
			<br />
			<Link to="/settings">Configure settings</Link>
		</div>
	) : (
		<h3>Loading...</h3>
	);
}
