import { useEffect, useState } from "react";
import { Link } from "react-router";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { Login } from "../../wailsjs/go/main/App";
import "./ConnectWithTwitchEntry.css";

interface TwitchAuthResult {
	success: boolean;
	for_bot: boolean;
	error?: string;
}

function ConnectWithTwitchEntry(props: { forBot: boolean }) {
	const [status, setStatus] = useState("");
	const [error, setError] = useState("");
	const [isLoggingIn, setIsLoggingIn] = useState(false);

	useEffect(() => {
		const offSuccess = EventsOn("TWITCH_AUTH_SUCCESS", (result: TwitchAuthResult) => {
			if (result.for_bot !== props.forBot) return;
			setStatus("Connected successfully! You can return to the home page.");
			setError("");
			setIsLoggingIn(false);
		});
		const offError = EventsOn("TWITCH_AUTH_ERROR", (result: TwitchAuthResult) => {
			if (result.for_bot !== props.forBot) return;
			setError(result.error ?? "Twitch authorization failed");
			setStatus("");
			setIsLoggingIn(false);
		});
		return () => {
			offSuccess();
			offError();
		};
	}, [props.forBot]);

	const startAuth = async () => {
		setIsLoggingIn(true);
		setError("");
		setStatus("Opening your browser to sign in with Twitch…");
		try {
			await Login(props.forBot);
			setStatus("Connected successfully! You can return to the home page.");
		} catch (e) {
			setError(String(e));
			setStatus("");
		} finally {
			setIsLoggingIn(false);
		}
	};

	return (
		<div className="card">
			<h2>{props.forBot ? "Connect Twitch bot account" : "Connect Twitch main account"}</h2>
			<p>
				Click Login to open your browser and authorize this app with Twitch. After you approve
				access, you&apos;ll be redirected back to the app automatically.
			</p>

			<button
				type="button"
				className="twitch-connect-button"
				onClick={() => void startAuth()}
				disabled={isLoggingIn}
			>
				{props.forBot ? "Login bot account" : "Login"}
			</button>

			{isLoggingIn && (
				<p className="read-the-docs">
					Waiting for authorization… complete the login in your browser, then return here.
				</p>
			)}

			{status && <p className="auth-status ok">{status}</p>}
			{error && <p className="auth-status err">{error}</p>}

			<p>
				<Link to="/">Return home</Link>
			</p>
		</div>
	);
}

export default ConnectWithTwitchEntry;
