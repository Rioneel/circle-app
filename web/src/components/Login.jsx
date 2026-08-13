import { useEffect, useState } from "react";
import { api } from "../api.js";

export default function Login({ onLogin }) {
  const [users, setUsers] = useState(null);
  const [error, setError] = useState(null);

  useEffect(() => {
    api("/api/users").then(setUsers).catch(e => setError(e.message));
  }, []);

  return (
    <div className="wrap">
      <h1>Circle</h1>
      <div className="sub">Recommendations from people you actually trust. Pick a demo user to log in as.</div>
      <div className="card">
        <h2>Log in</h2>
        {error && <div className="error-banner">Could not reach the server: {error}</div>}
        {!error && !users && <div className="spinner">Loading users…</div>}
        {users && users.length === 0 && <div className="empty">No users found — did you run the seed script?</div>}
        {users && users.map(u => (
          <button key={u.id} className="user-btn" onClick={() => onLogin(u)}>
            {u.name} · {u.id}
          </button>
        ))}
      </div>
    </div>
  );
}