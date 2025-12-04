import { useEffect, useState } from "react";
import { Container, Row, Col, Card, Spinner, Alert } from "react-bootstrap";
import { useNavigate } from "react-router-dom";
import { storeOAuthToken } from "../app/api";

export function OAuthProcessing() {
	const [error, setError] = useState<string | null>(null);
	const navigate = useNavigate();

	useEffect(() => {
		const processOAuth = async () => {
			try {
				// Check if we have OAuth token in URL fragment (after implicit grant redirect)
				const hashFragment = window.location.hash.substring(1);
				const hashParams = new URLSearchParams(hashFragment);
				const accessToken = hashParams.get("access_token");
				const tokenType = hashParams.get("token_type");
				let expiresIn = hashParams.get("expires_in");

				if (!accessToken || tokenType !== "bearer") {
					throw new Error("Invalid OAuth response: missing or invalid access token");
				}

				// Parse expires_in from the OAuth response
				let expiresInSeconds = 3600; // Default to 1 hour

				if (expiresIn !== null && expiresIn !== undefined && expiresIn !== "") {
					const parsed = parseInt(expiresIn, 10);
					if (!isNaN(parsed) && parsed > 0) {
						expiresInSeconds = parsed;
					}
				}

				// Send the OAuth token to the backend for storage
				await storeOAuthToken({
					access_token: accessToken,
					token_type: tokenType,
					expires_in: expiresInSeconds,
				});

				// Success! Redirect to success page
				navigate("/oauth/success");
			} catch (err) {
				console.error("OAuth processing failed:", err);
				setError(err instanceof Error ? err.message : "Unknown error occurred");
				// Redirect to error page after a short delay
				setTimeout(() => {
					navigate("/oauth/error");
				}, 3000);
			}
		};

		processOAuth();
	}, [navigate]);

	return (
		<Container className="py-4">
			<Row>
				<Col>
					<Card className="mb-4">
						<Card.Header>
							<h5 className="mb-0">Processing Twitch Authentication</h5>
						</Card.Header>
						<Card.Body className="text-center">
							{error ? (
								<>
									<Alert variant="danger" className="mb-3">
										<strong>Authentication Failed:</strong> {error}
									</Alert>
									<p>Redirecting to error page...</p>
								</>
							) : (
								<>
									<Spinner animation="border" role="status" className="mb-3">
										<span className="visually-hidden">Processing...</span>
									</Spinner>
									<h6>Connecting your Twitch account...</h6>
									<p className="text-muted">Please wait while we complete the authentication process.</p>
								</>
							)}
						</Card.Body>
					</Card>
				</Col>
			</Row>
		</Container>
	);
}
