import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import App from "./App.tsx";
import store from "./app/store";
import { Provider } from "react-redux";
import { BrowserRouter, Routes, Route } from "react-router-dom";

const root = ReactDOM.createRoot(
	document.getElementById("root") as HTMLElement,
);

root.render(
	<React.StrictMode>
		<Provider store={store}>
			<BrowserRouter>
				<Routes>
					<Route path="/" element={<App />}/>
					<Route path="/settings" element={<></>}/>
					<Route path="/queue" element={<></>}/>
					<Route path="/oauth/twitch" element={<></>}/>
				</Routes>
			</BrowserRouter>
		</Provider>
	</React.StrictMode>,
);
