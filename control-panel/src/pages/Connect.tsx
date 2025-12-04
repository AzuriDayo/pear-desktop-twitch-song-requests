import { useEffect, useState } from "react";
import { Container, Row, Col, Card, Alert } from "react-bootstrap";
import { getTwitchStatus } from "../app/api";

type TwitchStatus = {
	authenticated: boolean;
	login: string;
	user_id: string;
	expires_at: string;
	expired: boolean;
};

export function Connect() {
	const [params, setParams] = useState<URLSearchParams>();
	const [twitchStatus, setTwitchStatus] = useState<TwitchStatus | null>(null);
	const [loading, setLoading] = useState(true);

	// on page load, set the oauth params and fetch twitch status
	useEffect(() => {
		const params = new URLSearchParams();
		params.append(
			"client_id",
			import.meta.env.VITE_TWITCH_CLIENT_ID || "7k7nl6w8e0owouonj7nb9g3k5s6gs5",
		);
		params.append(
			"redirect_uri",
			`${window.location.protocol}//${window.location.host}/oauth/twitch-connect`,
		);
		params.append("response_type", "token");
		params.append(
			"scope",
			[
				"chat:read",
				"chat:edit",
				"channel:moderate",
				"whispers:read",
				"whispers:edit",
				"moderator:manage:banned_users",
				"channel:read:redemptions",
				"user:read:chat",
				"user:write:chat",
				"user:bot",
			].join(" "),
		);
		setParams(params);

		// Fetch current Twitch authentication status
		getTwitchStatus()
			.then((response) => {
				if (response.data) {
					setTwitchStatus(response.data);
				}
			})
			.catch((error) => {
				console.error("Failed to fetch Twitch status:", error);
			})
			.finally(() => {
				setLoading(false);
			});
	}, []);

	return (
		<Container className="py-4">
			<Row>
				<Col>
					<h1 className="mb-4">Connect with Twitch</h1>

					<Card className="mb-4">
						<Card.Header>
							<h5 className="mb-0">Twitch Integration</h5>
						</Card.Header>
						<Card.Body>
							{loading ? (
								<p>Loading authentication status...</p>
							) : twitchStatus?.authenticated ? (
								<>
									<Alert variant="success" className="mb-3">
										<strong>✓ Connected to Twitch</strong>
									</Alert>
									<div className="mb-3">
										<p className="mb-1">
											<strong>Username:</strong> {twitchStatus.login}
										</p>
										{twitchStatus.expires_at && (
											<p className="mb-1">
												<strong>Token Expires:</strong>{" "}
												{new Date(twitchStatus.expires_at).toLocaleString()}
											</p>
										)}
									</div>
									{twitchStatus.expired && (
										<Alert variant="warning" className="mb-3">
											Your token has expired. Please refresh your authentication.
										</Alert>
									)}
									<a
										href={`https://id.twitch.tv/oauth2/authorize?${params}`}
										className="btn btn-warning"
									>
										Refresh Token
									</a>
								</>
							) : (
								<>
									<p>
										Connect your Twitch account to enable song requests from chat.
									</p>
									<a
										href={`https://id.twitch.tv/oauth2/authorize?${params}`}
										className="btn btn-primary"
									>
										Connect with Twitch
									</a>
								</>
							)}
						</Card.Body>
					</Card>
				</Col>
			</Row>
		</Container>
	);
}
