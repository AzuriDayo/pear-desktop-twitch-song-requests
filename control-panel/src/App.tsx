import { Routes, Route } from "react-router-dom";
import { Home } from "./pages/Home";
import { Connect } from "./pages/Connect";
import { OAuthProcessing } from "./pages/OAuthProcessing";
import { OAuthSuccess } from "./pages/OAuthSuccess";
import { OAuthError } from "./pages/OAuthError";
import { NavBar } from "./components/NavBar";
import "./App.css";

function App() {
	return (
		<>
			<NavBar />
		<Routes>
			<Route path="/" element={<Home />} />
			<Route path="/connect" element={<Connect />} />
			<Route path="/oauth/twitch-connect" element={<OAuthProcessing />} />
			<Route path="/oauth/success" element={<OAuthSuccess />} />
			<Route path="/oauth/error" element={<OAuthError />} />
		</Routes>
		</>
	);
}

export default App;
