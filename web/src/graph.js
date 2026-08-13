import * as d3 from "d3";

export function renderForceGraph(container, nodes, links, opts = {}) {
  container.innerHTML = "";
  const width = container.clientWidth || 600;
  const height = container.clientHeight || 300;
  const svg = d3.select(container).append("svg").attr("viewBox", [0, 0, width, height]);

  const simulation = d3.forceSimulation(nodes)
    .force("link", d3.forceLink(links).id(d => d.id).distance(d => 70 + (1 - (d.weight ?? 0.5)) * 60))
    .force("charge", d3.forceManyBody().strength(-260))
    .force("center", d3.forceCenter(width / 2, height / 2))
    .force("collide", d3.forceCollide(30));

  const link = svg.append("g").selectAll("line").data(links).join("line")
    .attr("stroke", d => d.kind === "recommends" ? "#4caf7d" : "#5b8def")
    .attr("stroke-dasharray", d => d.kind === "recommends" ? "4,3" : null)
    .attr("stroke-opacity", d => 0.3 + (d.weight ?? 0.5) * 0.6)
    .attr("stroke-width", d => 1 + (d.weight ?? 0.5) * 4);

  const node = svg.append("g").selectAll("g").data(nodes).join("g").style("cursor", opts.onClick ? "pointer" : "default");

  node.append("circle")
    .attr("r", d => d.isMe ? 20 : (d.isProvider ? 17 : 13))
    .attr("fill", d => d.isMe ? "#5b8def" : (d.isProvider ? "#4caf7d" : "#2a3a5c"))
    .attr("stroke", "#e8eaed").attr("stroke-width", 1.5);

  node.append("text")
    .text(d => d.name)
    .attr("x", 0)
    .attr("y", d => (d.isMe ? 20 : (d.isProvider ? 17 : 13)) + 14)
    .attr("text-anchor", "middle").attr("fill", "#e8eaed").attr("font-size", "11px");

  if (opts.onClick) node.on("click", (event, d) => opts.onClick(d));

  node.call(d3.drag()
    .on("start", (event, d) => { if (!event.active) simulation.alphaTarget(0.3).restart(); d.fx = d.x; d.fy = d.y; })
    .on("drag", (event, d) => { d.fx = event.x; d.fy = event.y; })
    .on("end", (event, d) => { if (!event.active) simulation.alphaTarget(0); d.fx = null; d.fy = null; }));

  simulation.on("tick", () => {
    link.attr("x1", d => d.source.x).attr("y1", d => d.source.y)
        .attr("x2", d => d.target.x).attr("y2", d => d.target.y);
    node.attr("transform", d => `translate(${d.x},${d.y})`);
  });
  return simulation;
}