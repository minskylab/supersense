import type { AppProps } from "next/app";
import { ChakraProvider } from "@chakra-ui/core";

import { split, HttpLink } from "@apollo/client";
import { getMainDefinition } from "@apollo/client/utilities";
// import { WebSocketLink } from "apollo-link-ws";
import { WebSocketLink } from "@apollo/client/link/ws";
import { ApolloProvider, ApolloClient, InMemoryCache } from "@apollo/client";
import { useState, useEffect } from "react";

const hostname = process.browser ? window.location.host : "127.0.0.1:8080";

const httpLink = new HttpLink({
    uri: `http://${hostname}/graphql`,
    headers: {
        Origin: `http://${hostname}`,
    },
});

const wsLink = process.browser
    ? new WebSocketLink({
          uri: `ws://${hostname}/graphql`,
          options: {
              reconnect: true,
          },
      })
    : null;

const splitLink = process.browser
    ? split(
          ({ query }) => {
              const definition = getMainDefinition(query);
              return definition.kind === "OperationDefinition" && definition.operation === "subscription";
          },
          wsLink,
          httpLink
      )
    : httpLink;

const client = new ApolloClient({
    link: splitLink,
    cache: new InMemoryCache(),
});

const SupersenseApp = ({ Component, pageProps }: AppProps) => {
    return (
        <ChakraProvider resetCSS>
            <ApolloProvider client={client}>
                <Component {...pageProps} />
            </ApolloProvider>
        </ChakraProvider>
    );
};

// Only uncomment this method if you have blocking data requirements for
// every single page in your application. This disables the ability to
// perform automatic static optimization, causing every page in your app to
// be server-side rendered.
//
// SupersenseApp.getInitialProps = async (appContext: AppContext) => {
//   // calls page's `getInitialProps` and fills `appProps.pageProps`
//   const appProps = await App.getInitialProps(appContext);

//   return { ...appProps }
// }

export default SupersenseApp;
