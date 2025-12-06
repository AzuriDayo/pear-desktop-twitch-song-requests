import { Link } from "react-router";

export function ProcessTwitchSuccess() {
	return (
		<div>
			twitch connect successful
			<br />
			<Link to="/">Return home</Link>
		</div>
	);
}
