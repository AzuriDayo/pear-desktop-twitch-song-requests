import { Link } from "react-router";

export function ProcessTwitchOAuth() {
	return (
		<div>
			<h2>Twitch login</h2>
			<p>
				Use the connect pages in this app to log in with Twitch. Your browser opens for
				authorization and the app updates automatically when you finish.
			</p>
			<br />
			<Link to="/oauth/twitch-connect">Connect main account</Link>
			<br />
			<Link to="/oauth/twitch-connect-bot">Connect bot account</Link>
			<br />
			<br />
			<Link to="/">Return home</Link>
		</div>
	);
}
