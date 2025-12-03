import { useEffect, useState } from "react";
import { useAppSelector, useAppDispatch } from "../../app/hooks";
import {
	selectMusicState,
	updatePlayerState,
	setServiceHealth,
} from "./musicPlayerSlice";
import {
	getMusicPlayerState,
	setMusicPlayerState,
	checkServiceHealth,
	type MusicPlayerState,
} from "../../app/api";
import { Button, Card, Badge, Alert, Spinner } from "react-bootstrap";

export function MusicPlayer() {
	const playerState = useAppSelector(selectMusicState);
	const dispatch = useAppDispatch();
	const [isLoading, setIsLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);

	// Fetch music player state periodically
	useEffect(() => {
		const fetchState = async () => {
			const healthResult = await checkServiceHealth();
			dispatch(setServiceHealth(healthResult));

			if (healthResult) {
				const result = await getMusicPlayerState();
				if (result.data) {
					dispatch(updatePlayerState(result.data));
					setError(null);
				} else if (result.error) {
					setError(`Failed to fetch music state: ${result.error}`);
				}
			}
		};

		// Initial fetch
		fetchState();

		// Set up periodic fetching every 5 seconds
		const interval = setInterval(fetchState, 5000);

		return () => clearInterval(interval);
	}, [dispatch]);

	const handleTogglePlay = async () => {
		if (!playerState.serviceHealthy) return;

		setIsLoading(true);
		setError(null);

		const newState: MusicPlayerState = {
			isPlaying: !playerState.isPlaying,
			currentSong: playerState.currentSong,
			volume: playerState.volume,
		};

		const result = await setMusicPlayerState(newState);
		if (result.error) {
			setError(`Failed to update player state: ${result.error}`);
		} else {
			dispatch(updatePlayerState(newState));
		}

		setIsLoading(false);
	};

	const getHealthStatusBadge = () => {
		if (playerState.serviceHealthy) {
			return <Badge bg="success">Connected</Badge>;
		}
		return <Badge bg="danger">Disconnected</Badge>;
	};

	return (
		<Card className="mb-4">
			<Card.Header>
				<h5 className="mb-0 d-flex align-items-center justify-content-between">
					Music Player
					{getHealthStatusBadge()}
				</h5>
			</Card.Header>
			<Card.Body>
				{error && (
					<Alert variant="danger" className="mb-3">
						{error}
					</Alert>
				)}

				{!playerState.serviceHealthy && (
					<Alert variant="warning" className="mb-3">
						<strong>Pear Desktop service is not available.</strong> Make sure
						the backend is running and the service is accessible.
					</Alert>
				)}

				<div className="d-flex align-items-center mb-3">
					<Button
						variant={playerState.isPlaying ? "danger" : "success"}
						onClick={handleTogglePlay}
						disabled={!playerState.serviceHealthy || isLoading}
						className="me-3"
					>
						{isLoading ? (
							<>
								<Spinner
									as="span"
									animation="border"
									size="sm"
									role="status"
									aria-hidden="true"
									className="me-2"
								/>
								Updating...
							</>
						) : playerState.isPlaying ? (
							"⏸️ Pause"
						) : (
							"▶️ Play"
						)}
					</Button>

					<div className="flex-grow-1">
						<div className="d-flex align-items-start">
							{playerState.imageSrc && (
								<img
									src={playerState.imageSrc}
									alt="Album art"
									className="me-3"
									style={{ width: "60px", height: "60px", objectFit: "cover" }}
								/>
							)}
							<div className="flex-grow-1">
								<strong>Status:</strong>{" "}
								{playerState.isPlaying ? "Playing" : "Paused"}

								{playerState.currentSong && (
									<>
										<br />
										<strong>Song:</strong> {playerState.currentSong}
									</>
								)}

								{playerState.artist && (
									<>
										<br />
										<strong>Artist:</strong> {playerState.artist}
									</>
								)}

								{(playerState.elapsedSeconds !== undefined || playerState.songDuration !== undefined) && (
									<>
										<br />
										<strong>Progress:</strong>{" "}
										{playerState.elapsedSeconds || 0}s
										{playerState.songDuration && ` / ${playerState.songDuration}s`}
									</>
								)}

								{playerState.volume !== undefined && (
									<>
										<br />
										<strong>Volume:</strong> {playerState.volume}%
									</>
								)}

								{playerState.lastUpdated && (
									<>
										<br />
										<small className="text-muted">
											Last updated:{" "}
											{new Date(playerState.lastUpdated).toLocaleTimeString()}
										</small>
									</>
								)}
							</div>
						</div>
					</div>
				</div>

				<div className="text-muted small">
					Status updates automatically every 5 seconds when the service is
					available.
				</div>
			</Card.Body>
		</Card>
	);
}
