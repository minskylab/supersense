import { gql } from '@apollo/client';
import * as React from 'react';
import * as Apollo from '@apollo/client';
import * as ApolloReactComponents from '@apollo/client/react/components';
export type Maybe<T> = T | null;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  Time: any;
};

export type Subscription = {
  __typename?: 'Subscription';
  eventStream: Event;
};


export type SubscriptionEventStreamArgs = {
  filter?: Maybe<EventStreamFilter>;
};

export type SuperHeader = {
  __typename?: 'SuperHeader';
  buffer: Scalars['Int'];
  title: Scalars['String'];
  hashtag: Scalars['String'];
  brand: Scalars['String'];
};

export type EventStreamFilter = {
  sources: Array<Scalars['String']>;
};

export type PersonDraft = {
  name: Scalars['String'];
  photo?: Maybe<Scalars['String']>;
  username?: Maybe<Scalars['String']>;
};


export type Event = {
  __typename?: 'Event';
  id: Scalars['String'];
  createdAt: Scalars['Time'];
  emittedAt: Scalars['Time'];
  title: Scalars['String'];
  message: Scalars['String'];
  entities: Entities;
  actor: Person;
  shareURL: Scalars['String'];
  sourceID: Scalars['String'];
  sourceName: Scalars['String'];
  eventKind: Scalars['String'];
};

export type Mutation = {
  __typename?: 'Mutation';
  emit: Scalars['String'];
};


export type MutationEmitArgs = {
  token: Scalars['String'];
  draft: EventDraft;
};

export type AuthResponse = {
  __typename?: 'AuthResponse';
  jwtToken: Scalars['String'];
  expirateAt: Scalars['Time'];
};

export type Entities = {
  __typename?: 'Entities';
  tags: Array<Scalars['String']>;
  media: Array<MediaEntity>;
  urls: Array<UrlEntity>;
};

export type Person = {
  __typename?: 'Person';
  name: Scalars['String'];
  photo: Scalars['String'];
  owner: Scalars['String'];
  email?: Maybe<Scalars['String']>;
  profileURL?: Maybe<Scalars['String']>;
  username?: Maybe<Scalars['String']>;
};

export type EventDraft = {
  title?: Maybe<Scalars['String']>;
  message: Scalars['String'];
  actor?: Maybe<PersonDraft>;
  kind?: Maybe<Scalars['String']>;
  shareURL?: Maybe<Scalars['String']>;
  entities?: Maybe<EntitiesDraft>;
};

export type MediaEntity = {
  __typename?: 'MediaEntity';
  url: Scalars['String'];
  type: Scalars['String'];
};

export type Query = {
  __typename?: 'Query';
  sharedBoard: Array<Event>;
  header: SuperHeader;
};


export type QuerySharedBoardArgs = {
  buffer: Scalars['Int'];
};

export type UrlEntity = {
  __typename?: 'URLEntity';
  url: Scalars['String'];
  displayURL: Scalars['String'];
};

export type MediaEntityDraft = {
  url: Scalars['String'];
  type: Scalars['String'];
};

export type UrlEntityDraft = {
  url: Scalars['String'];
  displayURL: Scalars['String'];
};

export type EntitiesDraft = {
  tags: Array<Scalars['String']>;
  media: Array<MediaEntityDraft>;
  urls: Array<UrlEntityDraft>;
};

export type HeaderQueryVariables = Exact<{ [key: string]: never; }>;


export type HeaderQuery = (
  { __typename?: 'Query' }
  & { header: (
    { __typename?: 'SuperHeader' }
    & Pick<SuperHeader, 'buffer' | 'title' | 'hashtag' | 'brand'>
  ) }
);

export type SharedBoardQueryVariables = Exact<{
  size: Scalars['Int'];
}>;


export type SharedBoardQuery = (
  { __typename?: 'Query' }
  & { sharedBoard: Array<(
    { __typename?: 'Event' }
    & Pick<Event, 'id' | 'createdAt' | 'emittedAt' | 'title' | 'message' | 'shareURL' | 'sourceName' | 'sourceID' | 'eventKind'>
    & { entities: (
      { __typename?: 'Entities' }
      & Pick<Entities, 'tags'>
      & { urls: Array<(
        { __typename?: 'URLEntity' }
        & Pick<UrlEntity, 'displayURL' | 'url'>
      )>, media: Array<(
        { __typename?: 'MediaEntity' }
        & Pick<MediaEntity, 'type' | 'url'>
      )> }
    ), actor: (
      { __typename?: 'Person' }
      & Pick<Person, 'name' | 'username' | 'photo' | 'profileURL'>
    ) }
  )> }
);

export type EventsStreamSubscriptionVariables = Exact<{ [key: string]: never; }>;


export type EventsStreamSubscription = (
  { __typename?: 'Subscription' }
  & { eventStream: (
    { __typename?: 'Event' }
    & Pick<Event, 'id' | 'createdAt' | 'emittedAt' | 'title' | 'message' | 'shareURL' | 'sourceName' | 'sourceID' | 'eventKind'>
    & { entities: (
      { __typename?: 'Entities' }
      & Pick<Entities, 'tags'>
      & { urls: Array<(
        { __typename?: 'URLEntity' }
        & Pick<UrlEntity, 'displayURL' | 'url'>
      )>, media: Array<(
        { __typename?: 'MediaEntity' }
        & Pick<MediaEntity, 'type' | 'url'>
      )> }
    ), actor: (
      { __typename?: 'Person' }
      & Pick<Person, 'name' | 'username' | 'photo' | 'profileURL'>
    ) }
  ) }
);


export const HeaderDocument = gql`
    query Header {
  header {
    buffer
    title
    hashtag
    brand
  }
}
    `;
export type HeaderComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<HeaderQuery, HeaderQueryVariables>, 'query'>;

    export const HeaderComponent = (props: HeaderComponentProps) => (
      <ApolloReactComponents.Query<HeaderQuery, HeaderQueryVariables> query={HeaderDocument} {...props} />
    );
    

/**
 * __useHeaderQuery__
 *
 * To run a query within a React component, call `useHeaderQuery` and pass it any options that fit your needs.
 * When your component renders, `useHeaderQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useHeaderQuery({
 *   variables: {
 *   },
 * });
 */
export function useHeaderQuery(baseOptions?: Apollo.QueryHookOptions<HeaderQuery, HeaderQueryVariables>) {
        return Apollo.useQuery<HeaderQuery, HeaderQueryVariables>(HeaderDocument, baseOptions);
      }
export function useHeaderLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<HeaderQuery, HeaderQueryVariables>) {
          return Apollo.useLazyQuery<HeaderQuery, HeaderQueryVariables>(HeaderDocument, baseOptions);
        }
export type HeaderQueryHookResult = ReturnType<typeof useHeaderQuery>;
export type HeaderLazyQueryHookResult = ReturnType<typeof useHeaderLazyQuery>;
export type HeaderQueryResult = Apollo.QueryResult<HeaderQuery, HeaderQueryVariables>;
export const SharedBoardDocument = gql`
    query SharedBoard($size: Int!) {
  sharedBoard(buffer: $size) {
    id
    createdAt
    emittedAt
    title
    message
    shareURL
    sourceName
    sourceID
    eventKind
    entities {
      tags
      urls {
        displayURL
        url
      }
      media {
        type
        url
      }
    }
    actor {
      name
      username
      photo
      profileURL
    }
  }
}
    `;
export type SharedBoardComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<SharedBoardQuery, SharedBoardQueryVariables>, 'query'> & ({ variables: SharedBoardQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const SharedBoardComponent = (props: SharedBoardComponentProps) => (
      <ApolloReactComponents.Query<SharedBoardQuery, SharedBoardQueryVariables> query={SharedBoardDocument} {...props} />
    );
    

/**
 * __useSharedBoardQuery__
 *
 * To run a query within a React component, call `useSharedBoardQuery` and pass it any options that fit your needs.
 * When your component renders, `useSharedBoardQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useSharedBoardQuery({
 *   variables: {
 *      size: // value for 'size'
 *   },
 * });
 */
export function useSharedBoardQuery(baseOptions?: Apollo.QueryHookOptions<SharedBoardQuery, SharedBoardQueryVariables>) {
        return Apollo.useQuery<SharedBoardQuery, SharedBoardQueryVariables>(SharedBoardDocument, baseOptions);
      }
export function useSharedBoardLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<SharedBoardQuery, SharedBoardQueryVariables>) {
          return Apollo.useLazyQuery<SharedBoardQuery, SharedBoardQueryVariables>(SharedBoardDocument, baseOptions);
        }
export type SharedBoardQueryHookResult = ReturnType<typeof useSharedBoardQuery>;
export type SharedBoardLazyQueryHookResult = ReturnType<typeof useSharedBoardLazyQuery>;
export type SharedBoardQueryResult = Apollo.QueryResult<SharedBoardQuery, SharedBoardQueryVariables>;
export const EventsStreamDocument = gql`
    subscription EventsStream {
  eventStream {
    id
    createdAt
    emittedAt
    title
    message
    shareURL
    sourceName
    sourceID
    eventKind
    entities {
      tags
      urls {
        displayURL
        url
      }
      media {
        type
        url
      }
    }
    actor {
      name
      username
      photo
      profileURL
    }
  }
}
    `;
export type EventsStreamComponentProps = Omit<ApolloReactComponents.SubscriptionComponentOptions<EventsStreamSubscription, EventsStreamSubscriptionVariables>, 'subscription'>;

    export const EventsStreamComponent = (props: EventsStreamComponentProps) => (
      <ApolloReactComponents.Subscription<EventsStreamSubscription, EventsStreamSubscriptionVariables> subscription={EventsStreamDocument} {...props} />
    );
    

/**
 * __useEventsStreamSubscription__
 *
 * To run a query within a React component, call `useEventsStreamSubscription` and pass it any options that fit your needs.
 * When your component renders, `useEventsStreamSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useEventsStreamSubscription({
 *   variables: {
 *   },
 * });
 */
export function useEventsStreamSubscription(baseOptions?: Apollo.SubscriptionHookOptions<EventsStreamSubscription, EventsStreamSubscriptionVariables>) {
        return Apollo.useSubscription<EventsStreamSubscription, EventsStreamSubscriptionVariables>(EventsStreamDocument, baseOptions);
      }
export type EventsStreamSubscriptionHookResult = ReturnType<typeof useEventsStreamSubscription>;
export type EventsStreamSubscriptionResult = Apollo.SubscriptionResult<EventsStreamSubscription>;