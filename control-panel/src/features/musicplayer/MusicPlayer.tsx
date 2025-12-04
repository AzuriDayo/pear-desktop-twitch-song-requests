import { useAppSelector } from "../../app/hooks";
import {
	selectMusicState,
} from "./musicPlayerSlice";
import {
	selectConnectionStatus,
} from "../websocket/websocketSlice";
import { useWebSocket } from "../websocket/websocketService";
import { Card, Alert } from "react-bootstrap";

export function MusicPlayer() {
	const playerState = useAppSelector(selectMusicState);
	const connectionStatus = useAppSelector(selectConnectionStatus);

	// Use the global websocket hook
	useWebSocket();

	const openSongUrl = () => {
		if (playerState.videoId) {
			// Construct YouTube Music URL using video ID
			const youtubeUrl = `https://music.youtube.com/watch?v=${playerState.videoId}`;
			window.open(youtubeUrl, '_blank', 'noopener,noreferrer');
		} else if (playerState.url) {
			// Fallback to existing URL if videoId is not available
			window.open(playerState.url, '_blank', 'noopener,noreferrer');
		}
	};

	const formatTime = (seconds: number): string => {
		const minutes = Math.floor(seconds / 60);
		const remainingSeconds = seconds % 60;
		return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
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
		<Card className="mb-4">
			<Card.Header>
				<div className="d-flex align-items-center justify-content-between">
					<div className="flex-grow-1 text-center">
						<h4 className="fancy-title mb-0">Pear Desktop</h4>
					</div>
					<div className="connection-status-wrapper">
						{getConnectionStatusDisplay()}
					</div>
				</div>
			</Card.Header>
			<Card.Body>
				{connectionStatus === 'backend_only' && (
					<Alert variant="warning" className="mb-3">
						<strong>Pear Desktop Connection Error:</strong> Backend connected but cannot reach Pear Desktop service. Please check the API Server in Pear Desktop.
					</Alert>
				)}

				<div className="music-player-layout">
					{/* Main Section: Song info on left, Album art on right */}
					<div className="main-content-section">
						<div className="song-info-left">
							{/* Status Line */}
							<div className="status-line">
								<span className="status-icon-small">
									{playerState.isPlaying ? "▶️" : "⏸️"}
								</span>
								<span className={`status-text ${playerState.isPlaying ? 'playing' : 'paused'}`}>
									{playerState.isPlaying ? "Now Playing" : "Paused"}
								</span>
							</div>

							{/* Song + Artist Line */}
							{(playerState.currentSong || playerState.artist) && (
								<div
									className="media-line clickable"
									onClick={openSongUrl}
								>
									{playerState.currentSong && (
										<span className="song-title-compact">{playerState.currentSong}</span>
									)}
									{playerState.currentSong && playerState.artist && (
										<span className="separator-compact"> - </span>
									)}
									{playerState.artist && (
										<span className="song-artist-compact">{playerState.artist}</span>
									)}
								</div>
							)}

							{/* Time Line */}
							{(playerState.elapsedSeconds !== undefined ||
								playerState.songDuration !== undefined) && (
								<div className="time-line">
									<span className="time-icon">⏱️</span>
									<span className="time-text">
										{formatTime(playerState.elapsedSeconds || 0)}
										{playerState.songDuration && (
											<span className="time-separator"> / </span>
										)}
										{playerState.songDuration && (
											<span className="time-total">
												{formatTime(playerState.songDuration)}
											</span>
										)}
									</span>
								</div>
							)}
						</div>

						<div className="album-right">
							{playerState.imageSrc && (
								<img
									src={playerState.imageSrc}
									alt="Album art"
									className="album-art clickable"
									onClick={openSongUrl}
								/>
							)}
						</div>
					</div>
				</div>

			</Card.Body>
		</Card>
	);
}
