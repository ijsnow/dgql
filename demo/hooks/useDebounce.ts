import { useState, useEffect } from "react";
import useTimeout from "@restart/hooks/useTimeout";

export function useDebounce<V extends any>(value: V, delay: number) {
  const [debouncedValue, setDebouncedValue] = useState(value);

  const timeout = useTimeout();

  useEffect(() => {
    timeout.set(() => {
      setDebouncedValue(value);
    }, delay);
  }, [value, delay]);

  return debouncedValue;
}
