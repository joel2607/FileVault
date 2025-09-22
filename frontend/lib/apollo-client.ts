import { ApolloClient, InMemoryCache, from, split } from "@apollo/client"
import { setContext } from "@apollo/client/link/context"
import { GraphQLWsLink } from "@apollo/client/link/subscriptions"
import { createClient } from "graphql-ws"
import { getMainDefinition } from "@apollo/client/utilities"
import { createUploadLink } from "apollo-upload-client";

console.log("Apollo Client URI being used:", process.env.NEXT_PUBLIC_GRAPHQL_ENDPOINT);

const uploadLink = createUploadLink({
  uri: process.env.NEXT_PUBLIC_GRAPHQL_ENDPOINT,
})

const authLink = setContext((_, { headers }) => {
  const token = typeof window !== "undefined" ? localStorage.getItem("token") : null
  return {
    headers: {
      ...headers,
      authorization: token ? `Bearer ${token}` : "",
    },
  }
})

const createWsLink = () => {
  if (typeof window === "undefined") {
    return null
  }
  return new GraphQLWsLink(
    createClient({
      url: process.env.NEXT_PUBLIC_GRAPHQL_WS_ENDPOINT!,
      connectionParams: () => {
        const token = localStorage.getItem("token")
        return {
          Authorization: token ? `Bearer ${token}` : "",
        }
      },
    }),
  )
}

const wsLink = createWsLink()

const splitLink =
  wsLink != null
    ? split(
        ({ query }) => {
          const definition = getMainDefinition(query)
          return definition.kind === "OperationDefinition" && definition.operation === "subscription"
        },
        wsLink,
        from([authLink, uploadLink]),
      )
    : from([authLink, uploadLink])

export const apolloClient = new ApolloClient({
  link: splitLink,
  cache: new InMemoryCache(),
  defaultOptions: {
    watchQuery: {
      errorPolicy: "all",
    },
    query: {
      errorPolicy: "all",
    },
  },
})
