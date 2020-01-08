import { NextPage } from "next";
import { useRouter } from "next/router";
import { Search } from "../components/search";

const Index: NextPage = () => {
  const router = useRouter();
  const query = Array.isArray(router.query.q)
    ? router.query.q.join(",")
    : router.query.q;
  const variables = Array.isArray(router.query.variables)
    ? router.query.variables.join(",")
    : router.query.variables;

  const setQuery = (q: string) => {
    router.replace({
      pathname: "/",
      query: { q, variables }
    });
  };
  const setVariables = (v: string) => {
    router.replace({
      pathname: "/",
      query: { variables: v, q: query }
    });
  };

  return (
    <div className="container mx-auto max-w-md pt-12">
      <h1 className="mb-6 leading-none text-gray-700 font-light text-3xl">
        Welcome to the dgql demo.
      </h1>

      <Search
        query={query}
        setQuery={setQuery}
        variables={variables}
        setVariables={setVariables}
      />
    </div>
  );
};

export default Index;
