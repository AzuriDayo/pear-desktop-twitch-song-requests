import { Link } from "react-router";
import { useAppSelector } from "./app/hooks";
import { useEffect, useState } from "react";
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

export function Home() {
	const twitchState = useAppSelector((state) => state.twitchState);
	const [userExpiryState, setUserExpiryState] = useState<EExpiryState>(
		expiryState.OK,
	);
	const [botExpiryState, setBotExpiryState] = useState<EExpiryState>(
		expiryState.OK,
	);
	useEffect(() => {
		if (twitchState.expires_in === "") {
			return;
		}

		const now = new Date();
		let expiry: Date;
		try {
			expiry = new Date(twitchState.expires_in);
		} catch (e) {
			return;
		}

		if (isAfter(now, expiry)) {
			setUserExpiryState(expiryState.EXPIRED);
		} else if (isAfter(addDays(now, 15), expiry)) {
			setUserExpiryState(expiryState.SOON);
		}

		if (twitchState.expires_in_bot === "") {
			return;
		}
		try {
			expiry = new Date(twitchState.expires_in_bot);
		} catch (e) {
			return;
		}

		if (isAfter(now, expiry)) {
			setBotExpiryState(expiryState.EXPIRED);
		} else if (isAfter(addDays(now, 15), expiry)) {
			setBotExpiryState(expiryState.SOON);
		}
	}, [twitchState]);
	return (
		<div>
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
	);
}
