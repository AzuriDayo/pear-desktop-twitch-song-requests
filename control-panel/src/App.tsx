import { useEffect, useState } from "react";
import "./App.css";

function App() {
  const [params, setParams] = useState<URLSearchParams>();
  useEffect(() => {
    const params = new URLSearchParams();
    params.append("client_id", "9ewyamhdh6hhn9w4zjg0dodw3u1gq1");
    params.append("redirect_uri", "http://localhost:5173/");
    params.append("response_type", "code");
    params.append(
      "scope",
      ["chat:read", "chat:edit", "channel:moderate", "whispers:read", "whispers:edit", "moderator:manage:banned_users", "channel:read:redemptions", "user:read:chat", "user:write:chat", "user:bot"].join(' ')
    );
    setParams(params);
  }, []);
  /*
    Once you have the code, POST this:
    curl --location 'https://id.twitch.tv/oauth2/token' \
        --header 'Content-Type: application/x-www-form-urlencoded' \
        --data-urlencode 'client_id=9ewyamhdh6hhn9w4zjg0dodw3u1gq1' \
        --data-urlencode 'client_secret=TO_BE_FILLED' \
        --data-urlencode 'code=TO_BE_FILLED_WITH_CODE' \
        --data-urlencode 'grant_type=authorization_code' \
        --data-urlencode 'redirect_uri=http://localhost:5173/'

  */

  return (
    <>
      <a href={`https://id.twitch.tv/oauth2/authorize?${params}`}>
        Connect with Twitch
      </a>
    </>
  );
}

export default App;
