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

export type UrlEntity = {
  __typename?: 'URLEntity';
  url: Scalars['String'];
  displayURL: Scalars['String'];
};


export type Query = {
  __typename?: 'Query';
  event: Event;
};


export type QueryEventArgs = {
  id: Scalars['String'];
};

export type MediaEntity = {
  __typename?: 'MediaEntity';
  url: Scalars['String'];
  type: Scalars['String'];
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

export type Subscription = {
  __typename?: 'Subscription';
  events?: Maybe<Event>;
};

export type Event = {
  __typename?: 'Event';
  id: Scalars['String'];
  createdAt: Scalars['Time'];
  emmitedAt: Scalars['Time'];
  title: Scalars['String'];
  message: Scalars['String'];
  entities: Entities;
  actor: Person;
  shareURL: Scalars['String'];
  sourceID: Scalars['String'];
  sourceName: Scalars['String'];
  eventKind: Scalars['String'];
};

export type Entities = {
  __typename?: 'Entities';
  tags: Array<Scalars['String']>;
  media: Array<MediaEntity>;
  urls: Array<UrlEntity>;
};

export type EventsStreamSubscriptionVariables = Exact<{ [key: string]: never; }>;


export type EventsStreamSubscription = (
  { __typename?: 'Subscription' }
  & { events?: Maybe<(
    { __typename?: 'Event' }
    & Pick<Event, 'id' | 'createdAt' | 'emmitedAt' | 'eventKind' | 'sourceName' | 'title' | 'message'>
    & { entities: (
      { __typename?: 'Entities' }
      & { media: Array<(
        { __typename?: 'MediaEntity' }
        & Pick<MediaEntity, 'url' | 'type'>
      )> }
    ), actor: (
      { __typename?: 'Person' }
      & Pick<Person, 'username' | 'photo' | 'name'>
    ) }
  )> }
);


export const EventsStreamDocument = gql`
    subscription EventsStream {
  events {
    id
    createdAt
    emmitedAt
    eventKind
    sourceName
    title
    entities {
      media {
        url
        type
      }
    }
    actor {
      username
      photo
      name
    }
    message
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