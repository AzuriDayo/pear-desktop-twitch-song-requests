import { useEffect, useRef } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  setConnecting,
  setConnected,
  setDisconnected,
  setConnectionError,
  updateConnectionStatus,
  handleWebSocketMessage,
} from './websocketSlice';
import { updatePlayerState } from '../musicplayer/musicPlayerSlice';
import type { AppDispatch, RootState } from '../../app/store';
import type { WebSocketMessage, ConnectionStatusMessage, MusicPlayerStateMessage } from './websocketSlice';

// WebSocket service hook for managing connections across the app
export const useWebSocket = () => {
  const dispatch = useDispatch<AppDispatch>();
  const wsConnectedRef = useRef(false);
  const wsRef = useRef<WebSocket | null>(null);

  const { isConnected, connectionStatus } = useSelector((state: RootState) => ({
    isConnected: state.websocket.isConnected,
    connectionStatus: state.websocket.connectionStatus,
  }));

  const connectWebSocket = () => {
    // Prevent multiple connections due to React.StrictMode
    if (wsConnectedRef.current) {
      return;
    }

    dispatch(setConnecting());
    wsRef.current = new WebSocket(`ws://${window.location.host}/api/v1/music/ws`);

    wsRef.current.onopen = () => {
      console.log("WebSocket connected for real-time updates");
      dispatch(setConnected(wsRef.current!));
      wsConnectedRef.current = true;
    };

    wsRef.current.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);

        // Check if this is a connection status message
        if (data.frontend_connected !== undefined && data.pear_desktop_connected !== undefined) {
          // This is a connection status update
          dispatch(updateConnectionStatus(data as ConnectionStatusMessage));
        } else {
          // This is a music state update
          const musicData: MusicPlayerStateMessage = data;
          dispatch(updatePlayerState(musicData));
          dispatch(handleWebSocketMessage(data as WebSocketMessage));
        }
      } catch (err) {
        console.error("Failed to parse WebSocket message:", err);
      }
    };

    wsRef.current.onerror = (error) => {
      console.error("WebSocket error:", error);
      dispatch(setConnectionError("Connection lost - attempting to reconnect..."));
    };

    wsRef.current.onclose = () => {
      console.log("WebSocket connection closed");
      dispatch(setDisconnected());
      wsConnectedRef.current = false;

      // Attempt to reconnect after 3 seconds
      setTimeout(() => {
        if (!wsConnectedRef.current) {
          connectWebSocket();
        }
      }, 3000);
    };
  };

  const disconnectWebSocket = () => {
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
    wsConnectedRef.current = false;
    dispatch(setDisconnected());
  };

  // Auto-connect on mount, disconnect on unmount
  useEffect(() => {
    connectWebSocket();

    return () => {
      disconnectWebSocket();
    };
  }, []);

  return {
    isConnected,
    connectionStatus,
    connectWebSocket,
    disconnectWebSocket,
  };
};

// Hook for sending websocket messages (future use)
export const useWebSocketSend = () => {
  const { wsConnection } = useSelector((state: RootState) => ({
    wsConnection: state.websocket.wsConnection,
  }));

  const sendMessage = (message: WebSocketMessage) => {
    if (wsConnection && wsConnection.readyState === WebSocket.OPEN) {
      wsConnection.send(JSON.stringify(message));
    } else {
      console.warn('WebSocket not connected, cannot send message:', message);
    }
  };

  return { sendMessage };
};
