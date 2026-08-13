import { useState } from "react";
import { api } from "../api.js";
import ProviderResult from "./ProviderResult.jsx";

const CATEGORIES = ["Plumber","Electrician","Dentist","AC Repair","Carpenter","Tutor","Pest Control","House Cleaning"];
const AREAS = ["Hanamkonda","Warangal","Hyderabad"];

export default function Search({ user, onError }) {
  const [category, setCategory] = useState(CATEGORIES[0]);
  const [area, setArea] = useState("");
  const [results, setResults] = useState(null);
  const [loading, setLoading] = useState(false);

  async function runSearch() {
    setLoading(true);
    setResults(null);
    try {
      const params = new URLSearchParams({ userId: user.id, category });
      if (area) params.set("area", area);
      setResults(await api(`/api/search?` + params.toString()));
    } catch (e) {
      onError("Search failed: " + e.message);
    }
    setLoading(false);
  }

  return (
    <div className="card">
      <h2>Find a service provider</h2>
      <div className="row">
        <select value={category} onChange={e => setCategory(e.target.value)}>
          {CATEGORIES.map(c => <option key={c}>{c}</option>)}
        </select>
        <select value={area} onChange={e => setArea(e.target.value)}>
          <option value="">Any area</option>
          {AREAS.map(a => <option key={a}>{a}</option>)}
        </select>
        <button onClick={runSearch}>Search</button>
      </div>
      <div style={{ marginTop: 14 }}>
        {loading && <div className="spinner">Searching your network…</div>}
        {results && results.length === 0 && (
          <div className="empty">Nobody in your trust network (within 3 hops) has recommended a {category.toLowerCase()} yet.</div>
        )}
        {results && results.map(r => <ProviderResult key={r.id} result={r} user={user} onError={onError} />)}
      </div>
    </div>
  );
}