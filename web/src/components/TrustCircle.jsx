import { useEffect, useRef, useState, useCallback } from "react";
import { api } from "../api.js";
import { renderForceGraph } from "../graph.js";

export default function TrustCircle({ user, onError }) {
  const graphRef = useRef(null);
  const [trusted, setTrusted] = useState(null);
  const [allUsers, setAllUsers] = useState([]);
  const [selected, setSelected] = useState(null);
  const [addTarget, setAddTarget] = useState("");

  const load = useCallback(async () => {
    try {
      const [t, u] = await Promise.all([api(`/api/trust/${user.id}`), api(`/api/users`)]);
      setTrusted(t);
      setAllUsers(u);
      setSelected(null);
    } catch (e) {
      onError("Could not load your circle: " + e.message);
    }
  }, [user.id, onError]);

  useEffect(() => { load(); }, [load]);

  useEffect(() => {
    if (!trusted || !graphRef.current) return;
    if (trusted.length === 0) {
      graphRef.current.innerHTML = `<div class="empty">You haven't trusted anyone yet — add someone below.</div>`;
      return;
    }
    const nodes = [{ id: user.id, name: user.name, isMe: true },
      ...trusted.map(t => ({ id: t.ID, name: t.Name, weight: t.Weight }))];
    const links = trusted.map(t => ({ source: user.id, target: t.ID, weight: t.Weight }));
    renderForceGraph(graphRef.current, nodes, links, { onClick: (d) => { if (!d.isMe) setSelected(d.id); } });
  }, [trusted, user]);

  async function adjust(action) {
    const t = trusted.find(x => x.ID === selected);
    try {
      if (action === "remove") {
        await api(`/api/trust/${user.id}/${selected}`, { method: "DELETE" });
      } else {
        let weight = t.Weight + (action === "inc" ? 0.1 : -0.1);
        weight = Math.max(0, Math.min(1, weight));
        await api(`/api/trust`, {
          method: "PATCH", headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ fromId: user.id, toId: selected, weight })
        });
      }
      load();
    } catch (e) { onError("Could not update trust: " + e.message); }
  }

  async function addTrust() {
    if (!addTarget) return;
    try {
      await api(`/api/trust`, {
        method: "POST", headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ fromId: user.id, toId: addTarget, weight: 0.8 })
      });
      setAddTarget("");
      load();
    } catch (e) { onError("Could not add trust: " + e.message); }
  }

  const selectedUser = selected && trusted ? trusted.find(t => t.ID === selected) : null;
  const trustedIds = new Set((trusted || []).map(t => t.ID));
  const candidates = allUsers.filter(u => u.id !== user.id && !trustedIds.has(u.id));

  return (
    <div className="card">
      <h2>Your trust circle</h2>
      <div className="graph-box" ref={graphRef}>{!trusted && <div className="spinner">Loading…</div>}</div>
      <div className="legend">
        <span><span className="dot" style={{ background: "#5b8def" }}></span>You</span>
        <span><span className="dot" style={{ background: "#2a3a5c" }}></span>Trusted</span>
        <span>Drag nodes · click a node to adjust or remove trust</span>
      </div>
      {selectedUser && (
        <div className="selected-panel active">
          <div className="row" style={{ justifyContent: "space-between" }}>
            <strong>{selectedUser.Name}</strong>
            <span>{Math.round(selectedUser.Weight * 100)}% trust</span>
          </div>
          <div className="row" style={{ marginTop: 8 }}>
            <button className="secondary" onClick={() => adjust("dec")}>- 10%</button>
            <button className="secondary" onClick={() => adjust("inc")}>+ 10%</button>
            <button className="danger" onClick={() => adjust("remove")}>Remove trust</button>
          </div>
        </div>
      )}
      <div className="row" style={{ marginTop: 12 }}>
        <select value={addTarget} onChange={e => setAddTarget(e.target.value)}>
          <option value="">Select a person…</option>
          {candidates.map(u => <option key={u.id} value={u.id}>{u.name}</option>)}
        </select>
        <button onClick={addTrust}>Trust this person</button>
      </div>
    </div>
  );
}