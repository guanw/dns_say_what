import DNSFlow from "./DNSFlow";

function App() {
  // TODO get trace from backend
  const dnsTraceSteps: string[] = [
    "Queried 198.41.0.4:53 → NS: l.gtld-servers.net., j.gtld-servers.net., h.gtld-servers.net., d.gtld-servers.net., b.gtld-servers.net., f.gtld-servers.net., k.gtld-servers.net., m.gtld-servers.net., i.gtld-servers.net., g.gtld-servers.net., a.gtld-servers.net., c.gtld-servers.net., e.gtld-servers.net.",
    "Queried 192.41.162.30:53 → NS: ns2.google.com., ns1.google.com., ns3.google.com., ns4.google.com.",
    "Queried 216.239.34.10:53 → NS:",
    "Final A record from 216.239.34.10:53: 142.250.65.174",
  ];
  return (
    <>
      <DNSFlow trace={dnsTraceSteps} />
    </>
  );
}

export default App;
