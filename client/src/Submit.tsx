import Button from "@mui/material/Button";

interface InputProps {
  inputDomain: string;
  setDomain: React.Dispatch<React.SetStateAction<string>>;
}

export default function Input({ inputDomain, setDomain }: InputProps) {
  const handleSubmit = () => {
    const trimmedInput = inputDomain.trim();
    if (trimmedInput) {
      setDomain(trimmedInput);
    }
  };
  return (
    <Button variant="outlined" onClick={handleSubmit}>
      Submit
    </Button>
  );
}
