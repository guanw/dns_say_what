import TextField from "@mui/material/TextField";

interface InputProps {
  domain: string;
  setDomain: React.Dispatch<React.SetStateAction<string>>;
}

export default function Input({ domain, setDomain }: InputProps) {
  return (
    <TextField
      label="domain name"
      variant="standard"
      sx={{ width: "100%" }}
      value={domain}
      onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
        setDomain(event.target.value);
      }}
    />
  );
}
