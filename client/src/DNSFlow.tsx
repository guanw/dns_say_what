import { ReactFlow } from "@xyflow/react";

import "@xyflow/react/dist/style.css";
import { useMemo } from "react";

export default function DNSFlow({ trace }: { trace: string[] }) {
  const { nodes, edges } = useMemo(() => transformToFlow(trace), [trace]);

  return (
    <div
      style={{
        position: "absolute",
        top: 0,
        left: 0,
        width: "100vw",
        height: "100vh",
      }}
    >
      <ReactFlow nodes={nodes} edges={edges} fitView />
    </div>
  );
}

function transformToFlow(traceLines: string[]) {
  const nodes: any[] = [];
  const edges: any[] = [];

  let y = 0;
  let xSpacing = 200;
  let nodeCounter = 0;

  let rootNodesEachLevel = [];

  for (const line of traceLines) {
    if (!line.startsWith("Queried")) continue;

    const match = line.match(/Queried (.*?) → NS: (.*)/);
    if (!match) continue;

    const fromServer = match[1];
    const nsTargets = match[2]
      .split(",")
      .map((s) => s.trim())
      .filter(Boolean);

    const fromNodeId = `node-${nodeCounter++}`;
    // server node
    const serverNode = {
      id: fromNodeId,
      data: { label: fromServer },
      position: { x: 0, y },
    };
    nodes.push(serverNode);
    rootNodesEachLevel.push(serverNode);

    nsTargets.forEach((ns, index) => {
      const nsNodeId = `node-${ns}`;
      const nameRecordNode = {
        id: nsNodeId,
        data: { label: ns },
        position: { x: (index + 1) * xSpacing, y: y + 40 },
      };
      nodes.push(nameRecordNode);
      edges.push({
        id: `edge-${fromNodeId}-${nsNodeId}`,
        source: fromNodeId,
        target: nsNodeId,
      });
    });

    y += 150; // space vertically per query level
  }

  // Final A record line
  const finalLine = traceLines.find((line) =>
    line.startsWith("Final A record")
  );
  if (finalLine) {
    const finalMatch = finalLine.match(/Final A record from (.*): (.*)/);
    if (finalMatch) {
      const from = finalMatch[1];
      const to = finalMatch[2];

      const finalFromId = `node-${nodeCounter++}`;
      const finalToId = `node-${nodeCounter++}`;

      const finalServerNode = {
        id: finalFromId,
        data: { label: from },
        position: { x: 0, y },
      };
      nodes.push(finalServerNode);

      nodes.push({
        id: finalToId,
        data: { label: to },
        position: { x: xSpacing, y: y + 40 },
      });
      rootNodesEachLevel.push(finalServerNode);

      edges.push({
        id: `edge-${finalFromId}-${finalToId}`,
        source: finalFromId,
        target: finalToId,
      });
    }

    // concatenate all server nodes
    for (var i = 0; i < rootNodesEachLevel.length - 1; i++) {
      const parentNodeID = rootNodesEachLevel[i].id;
      const childNodeID = rootNodesEachLevel[i + 1].id;
      edges.push({
        id: `edge-${parentNodeID}-${childNodeID}`,
        source: parentNodeID,
        target: childNodeID,
        type: "step",
      });
    }
  }

  return { nodes, edges };
}
