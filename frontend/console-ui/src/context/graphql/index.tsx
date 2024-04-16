import { ApolloClient, InMemoryCache, createHttpLink } from "@apollo/client";

// Access the environment variable
const graphqlUri = process.env.REACT_APP_GRAPHQL_URI || "http://localhost:8080/graphql";

export const link = createHttpLink({
    uri: graphqlUri
});

export const client = new ApolloClient({
    cache: new InMemoryCache(),
    link,
});