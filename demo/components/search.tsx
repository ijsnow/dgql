import { FunctionComponent } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSearch } from "@fortawesome/free-solid-svg-icons";

import { useNodes } from "../hooks/useNodes";

interface Props {
  query: string;
  variables: string;
  setQuery: (query: string) => void;
  setVariables: (variables: string) => void;
}

export const Search: FunctionComponent<Props> = ({
  query,
  variables,
  setQuery,
  setVariables
}) => {
  const res = useNodes(query, variables);

  return (
    <div>
      <div className="relative">
        <span className="pointer-events-none absolute inset-y-0 left-0 pl-4 flex items-center">
          <FontAwesomeIcon icon={faSearch} className="text-gray-600" />
        </span>
        <input
          onChange={({ target: { value } }) => setQuery(value)}
          value={query}
          placeholder="Enter your search query"
          className="transition bg-white shadow-md focus:outline-0 border border-transparent placeholder-gray-600 rounded-lg py-2 pr-4 pl-10 block w-full appearance-none leading-normal ds-input text-gray-700"
        />
      </div>

      {/* {nodes && nodes.map(node => <Node node={node} typeInfo={typeInfo} />)} */}
    </div>
  );
};
