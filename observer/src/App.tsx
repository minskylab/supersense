import React from "react";
import { BrowserRouter as Router, Route } from "react-router-dom";
import { ObserverPage } from "@app/components/pages";

import { ThemeProvider } from "theme-ui";
import theme from "./theme";

import { split, HttpLink } from "@apollo/client";
import { getMainDefinition } from "@apollo/client/utilities";
import { WebSocketLink } from "@apollo/client/link/ws";

import { ApolloProvider, ApolloClient, InMemoryCache } from "@apollo/client";

const httpLink = new HttpLink({
    uri: "http://localhost:4000/graphql",
    headers: {
        Origin: "http://localhost:4000",
    },
});

const wsLink = new WebSocketLink({
    uri: `ws://localhost:4000/graphql`,
    options: {
        reconnect: true,
    },
});

const splitLink = split(
    ({ query }) => {
        const definition = getMainDefinition(query);
        return (
            definition.kind === "OperationDefinition" &&
            definition.operation === "subscription"
        );
    },
    wsLink,
    httpLink,
);

const client = new ApolloClient({
    link: splitLink,
    cache: new InMemoryCache(),
});

interface AppProps {}

function App({}: AppProps) {
    return (
        <ThemeProvider theme={theme}>
            <ApolloProvider client={client}>
                <Router>
                    <Route path={"/"}>
                        <ObserverPage />
                    </Route>
                </Router>
            </ApolloProvider>
        </ThemeProvider>
    );
}

export default App;
