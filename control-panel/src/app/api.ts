// API service for communicating with the backend
const API_BASE_URL = "http://" + window.location.host + "/api/v1";

export interface MusicPlayerState {
	isPlaying: boolean;
	currentSong?: string;
	artist?: string;
	url?: string;
	songDuration?: number;
	imageSrc?: string;
	elapsedSeconds?: number;
	volume?: number;
}

export interface ApiResponse<T> {
	data?: T;
	error?: string;
}

// Generic fetch wrapper with error handling
async function apiRequest<T>(
	endpoint: string,
	options?: RequestInit,
): Promise<ApiResponse<T>> {
	try {
		const response = await fetch(`${API_BASE_URL}${endpoint}`, {
			headers: {
				"Content-Type": "application/json",
				...options?.headers,
			},
			...options,
		});

		if (!response.ok) {
			return { error: `HTTP ${response.status}: ${response.statusText}` };
		}

		const data = await response.json();
		return { data };
	} catch (error) {
		return { error: error instanceof Error ? error.message : "Unknown error" };
	}
}

// Health check for music service
export async function checkServiceHealth(): Promise<boolean> {
	try {
		// The backend should have a health endpoint or we can use the music state endpoint as health check
		const response = await fetch(`${API_BASE_URL}/music/state`, {
			method: "GET",
			headers: {
				"Content-Type": "application/json",
			},
		});
		return response.ok;
	} catch {
		return false;
	}
}

// Get current music player state
export async function getMusicPlayerState(): Promise<
	ApiResponse<MusicPlayerState>
> {
	return apiRequest<MusicPlayerState>("/music/state");
}

// Set music player state
export async function setMusicPlayerState(
	state: MusicPlayerState,
): Promise<ApiResponse<any>> {
	return apiRequest("/music/state", {
		method: "POST",
		body: JSON.stringify(state),
	});
}

// WebSocket connection for real-time updates
export function connectMusicWebSocket(
	onMessage: (data: MusicPlayerState) => void,
	onError?: (error: Event) => void,
	onClose?: () => void,
): WebSocket | null {
	try {
		const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
		const wsUrl = `${protocol}//${window.location.host}/api/v1/music/ws`;

		const ws = new WebSocket(wsUrl);

		ws.onopen = () => {
			console.log("WebSocket connected for music updates");
		};

		ws.onmessage = (event) => {
			try {
				const data = JSON.parse(event.data);

				// Check if this is a connection status message
				if (data.frontend_connected !== undefined && data.pear_desktop_connected !== undefined) {
					// This is handled by the component directly - don't call onMessage
					console.log("Received connection status update:", data);
				} else {
					// This is a music state update
					const musicData: MusicPlayerState = data;
					onMessage(musicData);
				}
			} catch (err) {
				console.error("Failed to parse WebSocket message:", err);
			}
		};

		ws.onerror = (error) => {
			console.error("WebSocket error:", error);
			if (onError) onError(error);
		};

		ws.onclose = () => {
			console.log("WebSocket connection closed");
			if (onClose) onClose();
		};

		return ws;
	} catch (err) {
		console.error("Failed to create WebSocket connection:", err);
		return null;
	}
}
