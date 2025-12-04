import { Routes, Route } from "react-router-dom";
import { Home } from "./pages/Home";
import { Connect } from "./pages/Connect";
import { NavBar } from "./components/NavBar";
import "./App.css";

function App() {
	return (
		<>
			<NavBar />
		<Routes>
			<Route path="/" element={<Home />} />
			<Route path="/connect" element={<Connect />} />
		</Routes>
		</>
	);
}

export default App;
