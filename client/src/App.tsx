import DNSFlow from "./DNSFlow";
import Input from "./Input";
import Stack from "@mui/material/Stack";
import Submit from "./Submit";
import { useState } from "react";
import { useDNSQuery } from "./Hooks/useDNSQuery";
import CircularProgress from "@mui/material/CircularProgress";
import Box from "@mui/material/Box";

function App() {
  const [inputDomain, setInputDomain] = useState("");
  const [domain, setDomain] = useState("");

  const { data, error, isLoading } = useDNSQuery(domain);
  return (
    <Stack spacing={2}>
      <Stack direction="row" spacing={2} sx={{ width: "100%" }}>
        <Input domain={inputDomain} setDomain={setInputDomain} />
        <Submit inputDomain={inputDomain} setDomain={setDomain} />
      </Stack>
      {isLoading && <CircularProgress />}
      {error instanceof Error && <Box color="red">Error: {error.message}</Box>}
      {data && <DNSFlow trace={data} />}
    </Stack>
  );
}

export default App;
