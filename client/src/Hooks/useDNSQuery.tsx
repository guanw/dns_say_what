import { useQuery } from "@tanstack/react-query";

export function useDNSQuery(domain: string) {
  return useQuery({
    queryKey: ["dns-trace", domain],
    queryFn: async () => {
      const res = await fetch(`/trace?domain=${encodeURIComponent(domain)}`);
      if (!res.ok) throw new Error("Failed to fetch DNS trace");
      const text = await res.text();
      return text.split("\n");
    },
    enabled: !!domain, // only run if domain is not empty
  });
}
