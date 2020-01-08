import useSWR from "swr";
import { request } from "graphql-request";

import { useDebounce } from "./useDebounce";

export const useNodes = (query: string, variables: any) => {
  const debouncedQuery = useDebounce(query, 250);

  return useSWR(debouncedQuery ? `nodes/${debouncedQuery}` : null, () =>
    request("/api/graph", query, variables)
  );
};
