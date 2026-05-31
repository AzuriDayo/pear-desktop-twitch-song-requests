import { useEffect, useState } from "react";
import { useAppDispatch, useAppSelector } from "../../app/hooks";
import { handleWsMessages } from "./handleWsMessages";

export function TwitchWS() {
	const [reconnect, setReconnect] = useState(1);
	const twitchState = useAppSelector((state) => state.twitchState);
	const dispatch = useAppDispatch();

	useEffect(() => {
		const wsUrl = `ws://${twitchState.hostname}/api/v1/ws`;
		console.log("Starting Twitch WebSocket...");

		let ws: WebSocket;
		try {
			ws = new WebSocket(wsUrl);
		} catch (err) {
			console.error("Failed to create WebSocket connection:", err);
			console.log("Attempting to re-connect to Twitch in 3s..");
			const timer = setTimeout(() => setReconnect((c) => c + 1), 3000);
			return () => clearTimeout(timer);
		}

		ws.onopen = () => {
			console.log("Twitch WebSocket connected for app updates");
		};

		ws.onmessage = (event) => {
			if (event.type === "message") {
				handleWsMessages(event.data as string, dispatch);
			} else {
				console.log("TWITCH_WS bin_data", event);
			}
		};

		ws.onerror = (error) => {
			console.error("WebSocket error:", error);
		};

		let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
		ws.onclose = () => {
			console.log("Connection Closed, will reconnect in 3s...");
			reconnectTimer = setTimeout(() => setReconnect((c) => c + 1), 3000);
		};

		return () => {
			ws.onclose = null;
			if (
				ws.readyState === WebSocket.OPEN ||
				ws.readyState === WebSocket.CONNECTING
			) {
				ws.close();
			}
			if (reconnectTimer !== null) clearTimeout(reconnectTimer);
		};
	}, [reconnect, dispatch, twitchState.hostname]);

	return <></>;
}
