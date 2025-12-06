import { useEffect, useState } from "react";
import { useNavigate, Link } from "react-router";

export function ProcessTwitchOAuth() {
	const navigate = useNavigate();

	const [errorObj, setErrorObj] = useState<{
		error: string;
	} | null>(null);

	useEffect(() => {
		// check initial state
		/*
			example fragment: #access_token=73d0f8mkabpbmjp921asv2jaidwxn&scope=channel%3Amanage%3Apolls+channel%3Aread%3Apolls&state=c3ab8aa609ea11e793ae92361f002671&token_type=bearer
			example error: ?error=redirect_mismatch&error_description=Parameter+redirect_uri+does+not+match+registered+URI
		*/
		if (window.location.search) {
			try {
				const queryParams = new URLSearchParams(window.location.search);
				if (queryParams.has("error") && queryParams.has("error_description")) {
					setErrorObj({
						error:
							(queryParams.get("error") ?? "") +
							": " +
							(queryParams.get("error_description") ?? ""),
					});
					return;
				}
			} catch (e) {
				// do nothing
				return;
			}
		} else if (window.location.hash) {
			try {
				const hashParams = new URLSearchParams(
					window.location.hash.substring(1),
				);
				if (hashParams.has("access_token")) {
					const obj: {
						access_token: string;
						scope: string;
						state?: string;
						token_type: string;
					} = {
						access_token: hashParams.get("access_token") ?? "",
						scope: hashParams.get("access_token") ?? "",
						token_type: hashParams.get("token_type") ?? "",
					};
					if (hashParams.has("state")) {
						obj.state = hashParams.get("state") ?? "";
					}
					fetch("/api/v1/twitch-oauth", {
						body: JSON.stringify(obj),
						method: "POST",
					})
						.then((response) => {
							if (response.status === 200) {
								navigate("/oauth/twitch-success");
							} else {
								response.text().then((v) => {
									if (v.length > 0) {
										try {
											const j = JSON.parse(v);
											setErrorObj(j);
										} catch (e) {
											console.error("failed parsing json from server response");
										}
									} else {
										console.log("server response status only", response.status);
										setErrorObj({
											error: "Token invalid",
										});
									}
								});
							}
						})
						.catch((e) => {
							console.log(e);
							setErrorObj({
								error: "big failure " + e,
							});
						});
				}
			} catch (e) {
				console.error(e);
				return;
			}
		} else {
			navigate("/oauth/twitch-connect");
		}
	}, []);

	return (
		<div>
			{errorObj ? (
				<>
					<h3>{errorObj.error}</h3>
					<br />
					<Link to={"/oauth/twitch-connect"}>Retry</Link>
				</>
			) : (
				<h2>Working...</h2>
			)}
		</div>
	);
}
