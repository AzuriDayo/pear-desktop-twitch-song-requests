import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import type { RootState } from "../../app/store";

// Connection status types
export type ConnectionStatus = 'connecting' | 'connected' | 'backend_only' | 'disconnected' | 'error';

// Message types that can be received
export interface WebSocketMessage {
  type: string;
  [key: string]: any;
}

// Connection status from backend
export interface ConnectionStatusMessage {
  frontend_connected: boolean;
  pear_desktop_connected: boolean;
}

// Music player state (subset for websocket messages)
export interface MusicPlayerStateMessage {
  isPlaying?: boolean;
  currentSong?: string;
  artist?: string;
  url?: string;
  songDuration?: number;
  imageSrc?: string;
  elapsedSeconds?: number;
  volume?: number;
}

// Define a type for the slice state
interface IWebSocketState {
  // Connection state
  connectionStatus: ConnectionStatus;
  isConnected: boolean;
  wsConnection: WebSocket | null;

  // Connection health
  frontendConnected: boolean;
  pearDesktopConnected: boolean;

  // Message history (for debugging/monitoring)
  lastMessage?: WebSocketMessage;
  messageCount: number;

  // Error handling
  error?: string;
  reconnectAttempts: number;
}

const initialState: IWebSocketState = {
  connectionStatus: 'disconnected',
  isConnected: false,
  wsConnection: null,
  frontendConnected: false,
  pearDesktopConnected: false,
  messageCount: 0,
  reconnectAttempts: 0,
};

export const websocketSlice = createSlice({
  name: "websocket",
  initialState,
  reducers: {
    // Connection management
    setConnecting: (state) => {
      state.connectionStatus = 'connecting';
      state.isConnected = false;
      state.error = undefined;
    },
    setConnected: (state, action: PayloadAction<WebSocket>) => {
      state.connectionStatus = 'connected';
      state.isConnected = true;
      state.wsConnection = action.payload;
      state.error = undefined;
      state.reconnectAttempts = 0;
    },
    setDisconnected: (state) => {
      state.connectionStatus = 'disconnected';
      state.isConnected = false;
      state.wsConnection = null;
      state.frontendConnected = false;
      state.pearDesktopConnected = false;
    },
    setConnectionError: (state, action: PayloadAction<string>) => {
      state.connectionStatus = 'error';
      state.isConnected = false;
      state.wsConnection = null;
      state.error = action.payload;
      state.reconnectAttempts += 1;
    },

    // Connection status updates
    updateConnectionStatus: (state, action: PayloadAction<ConnectionStatusMessage>) => {
      const { frontend_connected, pear_desktop_connected } = action.payload;
      state.frontendConnected = frontend_connected;
      state.pearDesktopConnected = pear_desktop_connected;

      // Update overall connection status based on both connections
      if (frontend_connected && pear_desktop_connected) {
        state.connectionStatus = 'connected';
      } else if (frontend_connected && !pear_desktop_connected) {
        state.connectionStatus = 'backend_only';
      } else {
        state.connectionStatus = 'disconnected';
      }
    },

    // Message handling
    handleWebSocketMessage: (state, action: PayloadAction<WebSocketMessage>) => {
      const message = action.payload;
      state.lastMessage = message;
      state.messageCount += 1;

      // Handle different message types
      switch (message.type) {
        case 'connection_status':
          // Connection status updates are handled by updateConnectionStatus
          break;
        case 'music_state':
          // Music state updates are handled by the musicPlayer slice
          break;
        default:
          console.log('Unhandled websocket message type:', message.type);
      }
    },

    // Utility actions
    incrementReconnectAttempts: (state) => {
      state.reconnectAttempts += 1;
    },
    clearError: (state) => {
      state.error = undefined;
    },
    resetMessageCount: (state) => {
      state.messageCount = 0;
    },
  },
});

export const {
  setConnecting,
  setConnected,
  setDisconnected,
  setConnectionError,
  updateConnectionStatus,
  handleWebSocketMessage,
  incrementReconnectAttempts,
  clearError,
  resetMessageCount,
} = websocketSlice.actions;

// Selectors
export const selectWebSocketState = (state: RootState) => state.websocket;
export const selectConnectionStatus = (state: RootState) => state.websocket.connectionStatus;
export const selectIsWebSocketConnected = (state: RootState) => state.websocket.isConnected;
export const selectFrontendConnected = (state: RootState) => state.websocket.frontendConnected;
export const selectPearDesktopConnected = (state: RootState) => state.websocket.pearDesktopConnected;
export const selectWebSocketError = (state: RootState) => state.websocket.error;
export const selectReconnectAttempts = (state: RootState) => state.websocket.reconnectAttempts;
export const selectLastMessage = (state: RootState) => state.websocket.lastMessage;
export const selectMessageCount = (state: RootState) => state.websocket.messageCount;

export default websocketSlice.reducer;
