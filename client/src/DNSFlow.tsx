import { ReactFlow } from "@xyflow/react";

import "@xyflow/react/dist/style.css";
import { useMemo } from "react";

export default function DNSFlow({ trace }: { trace: string[] }) {
  // TODO use trace from backend to generate nodes and edges
  //   const { nodes, edges } = useMemo(() => transformToFlow(trace), [trace]);
  const nodes = [
    {
      id: "q1",
      type: "default",
      position: { x: 0, y: 0 },
      data: { label: "198.41.0.4:53" },
    },
    {
      id: "n1",
      position: { x: 200, y: 40 },
      data: { label: "l.gtld-servers.net." },
    },
    {
      id: "n2",
      position: { x: 400, y: 40 },
      data: { label: "j.gtld-servers.net." },
    },
    {
      id: "q2",
      position: { x: 0, y: 150 },
      data: { label: "192.41.162.30:53" },
    },
    {
      id: "n3",
      position: { x: 200, y: 190 },
      data: { label: "ns2.google.com." },
    },
    {
      id: "n4",
      position: { x: 400, y: 190 },
      data: { label: "ns1.google.com." },
    },
    {
      id: "q3",
      position: { x: 0, y: 300 },
      data: { label: "216.239.34.10:53" },
    },
    {
      id: "a1",
      position: { x: 0, y: 450 },
      data: { label: "142.250.65.174" },
    },
  ];

  const edges = [
    { id: "q1-n1", source: "q1", target: "n1" },
    { id: "q1-n2", source: "q1", target: "n2" },
    { id: "q2-n3", source: "q2", target: "n3" },
    { id: "q2-n4", source: "q2", target: "n4" },
    { id: "q1-q2", source: "q1", target: "q2" },
    { id: "q2-q3", source: "q2", target: "q3" },
    { id: "q3-a1", source: "q3", target: "a1" },
  ];
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
