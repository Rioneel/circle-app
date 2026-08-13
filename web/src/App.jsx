import { useState } from "react";
import Login from "./components/Login.jsx";
import Dashboard from "./components/Dashboard.jsx";

export default function App() {
  const [user, setUser] = useState(() => {
    const saved = localStorage.getItem("circle_user");
    return saved ? JSON.parse(saved) : null;
  });

  function handleLogin(u) {
    localStorage.setItem("circle_user", JSON.stringify(u));
    setUser(u);
  }
  function handleLogout() {
    localStorage.removeItem("circle_user");
    setUser(null);
  }

  return user ? <Dashboard user={user} onLogout={handleLogout} /> : <Login onLogin={handleLogin} />;
}