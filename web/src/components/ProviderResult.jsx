import { useRef, useState } from "react";
import { api } from "../api.js";
import { renderForceGraph } from "../graph.js";

export default function ProviderResult({ result, user, onError }) {
  const [showPath, setShowPath] = useState(false);
  const [loaded, setLoaded] = useState(false);
  const graphRef = useRef(null);

  async function togglePath() {
    const next = !showPath;
    setShowPath(next);
    if (!next || loaded) return;

    if (graphRef.current) graphRef.current.innerHTML = `<div class="spinner">Tracing your network…</div>`;
    try {
      const recs = await api(`/api/providers/${result.id}/recommendations`);
      const nodesMap = new Map();
      nodesMap.set(user.id, { id: user.id, name: user.name, isMe: true });
      nodesMap.set(result.id, { id: result.id, name: result.name, isProvider: true });
      const links = [];
      let found = false;

      for (const r of recs) {
        let steps = [];
        if (r.userId === user.id) {
          steps = [{ id: user.id, name: user.name }];
        } else {
          try {
            const res = await api(`/api/path?meId=${user.id}&recommenderId=${r.userId}`);
            steps = res.path || [];
          } catch { continue; }
        }
        if (!steps.length) continue;
        found = true;
        steps.forEach(s => { if (!nodesMap.has(s.id)) nodesMap.set(s.id, { id: s.id, name: s.name }); });
        for (let i = 0; i < steps.length - 1; i++) links.push({ source: steps[i].id, target: steps[i + 1].id, weight: 0.6 });
        links.push({ source: r.userId, target: result.id, weight: r.rating / 5, kind: "recommends" });
      }

      if (!found) {
        if (graphRef.current) graphRef.current.innerHTML = `<div class="empty">No one in your trust network recommended this — it showed up from outside your circle.</div>`;
      } else if (graphRef.current) {
        renderForceGraph(graphRef.current, Array.from(nodesMap.values()), links, {});
      }
      setLoaded(true);
    } catch (e) {
      if (graphRef.current) graphRef.current.innerHTML = `<div class="error-banner">Could not trace path: ${e.message}</div>`;
    }
  }

  async function recommend() {
    const rating = Number(prompt("Rate this provider 1-5:"));
    if (!rating || rating < 1 || rating > 5) return;
    try {
      await api(`/api/recommend`, {
        method: "POST", headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ userId: user.id, providerId: result.id, rating })
      });
      alert("Thanks — your recommendation was saved.");
    } catch (e) { onError("Could not save recommendation: " + e.message); }
  }

  return (
    <div className="result-card">
      <div className="result-top">
        <span className="name">{result.name}</span>
        <span className="score">{result.weightedScore.toFixed(1)} ★</span>
      </div>
      <div className="voices">{result.voices} {result.voices === 1 ? "voice" : "voices"} in your network</div>
      <div className="row" style={{ marginTop: 8 }}>
        <button className="secondary" onClick={togglePath}>How you got here</button>
        <button onClick={recommend}>Recommend this provider</button>
      </div>
      {showPath && <div className="graph-box small" ref={graphRef}></div>}
    </div>
  );
}