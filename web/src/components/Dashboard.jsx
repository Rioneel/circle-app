import { useState } from "react";
import TrustCircle from "./TrustCircle.jsx";
import Search from "./Search.jsx";

export default function Dashboard({ user, onLogout }) {
  const [error, setError] = useState(null);

  return (
    <div className="wrap">
      <div className="topbar">
        <div>
          <h1>Circle</h1>
          <div className="sub">Logged in as {user.name}</div>
        </div>
        <button className="secondary" onClick={onLogout}>Log out</button>
      </div>
      {error && <div className="error-banner">{error}</div>}
      <TrustCircle user={user} onError={setError} />
      <Search user={user} onError={setError} />
    </div>
  );
}