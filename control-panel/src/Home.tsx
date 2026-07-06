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

function renderTwitchAccountStatus(
	login: string,
	accessExpiry: string,
	refreshExpiry: string,
	accessExpiryState: EExpiryState,
	refreshExpiryState: EExpiryState,
	notConfiguredLabel: string,
) {
	if (login === "" && accessExpiry === "" && refreshExpiry === "") {
		return <h3>{notConfiguredLabel}</h3>;
	}

	// Device-code flow: refresh token present — show re-auth deadline only.
	if (refreshExpiry !== "") {
		return (
			<h3>
				Twitch connected as {login} — re-authenticate by {refreshExpiry} (30-day inactivity limit){" "}
				{getExpiryStateEmoji(refreshExpiryState)}
			</h3>
		);
	}

	// Legacy implicit grant: no refresh token — show access token expiry.
	if (accessExpiry !== "") {
		return (
			<h3>
				Twitch token for {login} expires on {accessExpiry} {getExpiryStateEmoji(accessExpiryState)}
			</h3>
		);
	}

	return <h3>Twitch connected as {login}</h3>;
}

export function Home() {
	const twitchState = useAppSelector((state) => state.twitchState);

	const userAccessExpiryState = useMemo(
		() => computeExpiryState(twitchState.expires_in),
		[twitchState.expires_in],
	);

	const botAccessExpiryState = useMemo(
		() => computeExpiryState(twitchState.expires_in_bot),
		[twitchState.expires_in_bot],
	);

	const userRefreshExpiryState = useMemo(
		() => computeExpiryState(twitchState.refresh_expires_in),
		[twitchState.refresh_expires_in],
	);

	const botRefreshExpiryState = useMemo(
		() => computeExpiryState(twitchState.refresh_expires_in_bot),
		[twitchState.refresh_expires_in_bot],
	);

	return twitchState.isLoaded ? (
		<div>
			<Link to="/queue">View song queue and history</Link>
			<br />
			<br />
			<br />
			<Link to="/oauth/twitch-connect">
				{twitchState.login !== "" ? "Re-login with Twitch" : "Login with Twitch"}
			</Link>
			{renderTwitchAccountStatus(
				twitchState.login,
				twitchState.expires_in,
				twitchState.refresh_expires_in,
				userAccessExpiryState,
				userRefreshExpiryState,
				"No Twitch token configured",
			)}
			<br />
			<Link to="/oauth/twitch-connect-bot">
				{twitchState.login_bot !== "" ? "Re-login Twitch bot" : "Login Twitch bot account"}
			</Link>
			{renderTwitchAccountStatus(
				twitchState.login_bot,
				twitchState.expires_in_bot,
				twitchState.refresh_expires_in_bot,
				botAccessExpiryState,
				botRefreshExpiryState,
				"No bot Twitch token configured",
			)}
			<br />
			<br />
			<br />
			<Link to="/settings">Configure settings</Link>
		</div>
	) : (
		<h3>Loading...</h3>
	);
}
