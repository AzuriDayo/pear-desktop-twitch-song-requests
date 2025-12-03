import { useEffect, useState, useRef } from "react";
import type { MusicPlayerState } from "../../app/api";
import { useAppSelector, useAppDispatch } from "../../app/hooks";
import {
	selectMusicState,
	updatePlayerState,
	setServiceHealth,
} from "./musicPlayerSlice";
import { Card, Alert } from "react-bootstrap";

export function MusicPlayer() {
	const playerState = useAppSelector(selectMusicState);
	const dispatch = useAppDispatch();
	const [error, setError] = useState<string | null>(null);
	const [connectionStatus, setConnectionStatus] = useState<'connecting' | 'connected' | 'backend_only' | 'disconnected' | 'error'>('connecting');
	const wsConnectedRef = useRef(false);

	// Connect to websocket for real-time updates
	useEffect(() => {
		// Prevent multiple connections due to React.StrictMode
		if (wsConnectedRef.current) {
			return;
		}

		wsConnectedRef.current = true;
		let ws: WebSocket | null = null;

		const connectWebSocket = () => {
			setConnectionStatus('connecting');
			ws = new WebSocket(`ws://${window.location.host}/api/v1/music/ws`);

			ws.onopen = () => {
				console.log("WebSocket connected for music updates");
			};

			ws.onmessage = (event) => {
				try {
					const data = JSON.parse(event.data);

					// Check if this is a connection status message
					if (data.frontend_connected !== undefined && data.pear_desktop_connected !== undefined) {
						// This is a connection status update
						const status = data as { frontend_connected: boolean; pear_desktop_connected: boolean };

						if (status.frontend_connected && status.pear_desktop_connected) {
							setConnectionStatus('connected');
							dispatch(setServiceHealth(true));
							setError(null);
						} else if (status.frontend_connected && !status.pear_desktop_connected) {
							setConnectionStatus('backend_only');
							dispatch(setServiceHealth(false));
							setError("Pear Desktop connection failed. Please check the API Server in Pear Desktop.");
						} else {
							setConnectionStatus('disconnected');
							dispatch(setServiceHealth(false));
							setError("Connection lost");
						}
					} else {
						// This is a music state update
						const musicData: MusicPlayerState = data;
						dispatch(updatePlayerState(musicData));
						dispatch(setServiceHealth(true));
						setError(null);
					}
				} catch (err) {
					console.error("Failed to parse WebSocket message:", err);
				}
			};

			ws.onerror = (error) => {
				console.error("WebSocket error:", error);
				dispatch(setServiceHealth(false));
				setError("Connection lost - attempting to reconnect...");
				setConnectionStatus('error');
			};

			ws.onclose = () => {
				console.log("WebSocket connection closed");
				dispatch(setServiceHealth(false));
				setError("Connection lost - attempting to reconnect...");
				setConnectionStatus('disconnected');
				setTimeout(connectWebSocket, 3000); // Reconnect after 3 seconds
			};
		};

		// Initial connection
		connectWebSocket();

		// Cleanup on unmount
		return () => {
			wsConnectedRef.current = false;
			if (ws) {
				ws.close();
			}
		};
	}, [dispatch]);


	const formatTime = (seconds: number): string => {
		const minutes = Math.floor(seconds / 60);
		const remainingSeconds = seconds % 60;
		return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
	};

	const openSongUrl = () => {
		if (playerState.url) {
			window.open(playerState.url, '_blank', 'noopener,noreferrer');
		}
	};

	const getConnectionStatusDisplay = () => {
		const getStatusConfig = () => {
			switch (connectionStatus) {
				case 'connecting':
					return {
						icon: '🔄',
						text: 'Connecting...',
						color: 'warning',
						pulse: true
					};
				case 'connected':
					return {
						icon: '🟢',
						text: 'Fully Connected',
						color: 'success',
						pulse: false
					};
				case 'backend_only':
					return {
						icon: '⚠️',
						text: 'Pear Desktop Connection Error',
						color: 'warning',
						pulse: false
					};
				case 'disconnected':
					return {
						icon: '🔴',
						text: 'Disconnected',
						color: 'danger',
						pulse: false
					};
				case 'error':
					return {
						icon: '🔴',
						text: 'Connection Error',
						color: 'danger',
						pulse: false
					};
				default:
					return {
						icon: '⚫',
						text: 'Unknown',
						color: 'secondary',
						pulse: false
					};
			}
		};

		const config = getStatusConfig();
		return (
			<div className={`connection-status ${config.pulse ? 'pulse' : ''}`}>
				<span className="status-icon">{config.icon}</span>
				<span className={`status-text text-${config.color}`}>{config.text}</span>
			</div>
		);
	};

	return (
		<Card className="mb-4 music-player-card">
			<Card.Header>
				<div className="d-flex align-items-center justify-content-between">
					<h5 className="mb-0">Pear Desktop Status</h5>
					{getConnectionStatusDisplay()}
				</div>
			</Card.Header>
			<Card.Body>
				{error && (
					<Alert variant="danger" className="mb-3 error-alert">
						<div className="alert-text">
							<strong>Connection Issue:</strong> {error}
						</div>
					</Alert>
				)}

				{!playerState.serviceHealthy && (
					<Alert variant="warning" className="mb-3">
						<strong>Music service is not available.</strong> Make sure
						the backend is running and the service is accessible.
					</Alert>
				)}

				<div className="music-info-display d-flex align-items-start">
							{playerState.imageSrc && (
								<img
									src={playerState.imageSrc}
									alt="Album art"
									className="album-art me-3 clickable"
									style={{ width: "60px", height: "60px", objectFit: "cover" }}
									onClick={openSongUrl}
								/>
							)}
							<div className="song-details flex-grow-1">
								<div className="playback-status mb-2">
									<strong>Status:</strong>{" "}
									<span className={`status-indicator ${playerState.isPlaying ? 'playing' : 'paused'}`}>
										{playerState.isPlaying ? "▶ Playing" : "⏸ Paused"}
									</span>
								</div>
								{playerState.currentSong && (
									<div className="song-title mb-1 clickable" onClick={openSongUrl}>
										🎵 {playerState.currentSong}
									</div>
								)}
								{playerState.artist && (
									<div className="song-artist mb-2">
										👤 {playerState.artist}
									</div>
								)}
								<div className="song-meta">
									{(playerState.elapsedSeconds !== undefined ||
										playerState.songDuration !== undefined) && (
										<div className="progress-info mb-1">
											⏱️ Progress: {formatTime(playerState.elapsedSeconds || 0)}
											{playerState.songDuration &&
												` / ${formatTime(playerState.songDuration)}`}
										</div>
									)}
									{playerState.volume !== undefined && (
										<div className="volume-info">
											🔊 Volume: {playerState.volume}%
										</div>
									)}
								</div>
							</div>
						</div>

			</Card.Body>
		</Card>
	);
}
