import first from "lodash/first";

import { useNodes } from "./useNodes";

export const useNode = (id: string) => {
  const { nodes, typeInfo, error } = useNodes(`id:${id}`);

  return { node: nodes ? first(nodes) : undefined, typeInfo, error };
};
