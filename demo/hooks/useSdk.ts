import { GraphQLClient } from "graphql-request";

import { getSdk } from "../client";

const API = "http://localhost:3000/api/graphql";

export const useSdk = () => getSdk(new GraphQLClient(API));
