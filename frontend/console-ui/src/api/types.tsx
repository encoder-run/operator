import { gql } from '@apollo/client';
import * as React from 'react';
import * as Apollo from '@apollo/client';
import * as ApolloReactComponents from '@apollo/client/react/components';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
export type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>;
const defaultOptions = {} as const;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
};

export type AddModelInput = {
  huggingFace?: InputMaybe<HuggingFaceInput>;
  type: ModelType;
};

export type AddRepositoryInput = {
  name?: InputMaybe<Scalars['String']['input']>;
  owner?: InputMaybe<Scalars['String']['input']>;
  type?: InputMaybe<RepositoryType>;
  url?: InputMaybe<Scalars['String']['input']>;
};

export type AddStorageInput = {
  name: Scalars['String']['input'];
  type: StorageType;
};

export type HuggingFace = {
  __typename?: 'HuggingFace';
  name: Scalars['String']['output'];
  organization: Scalars['String']['output'];
};

export type HuggingFaceInput = {
  name: Scalars['String']['input'];
  organization: Scalars['String']['input'];
};

export type Model = {
  __typename?: 'Model';
  displayName: Scalars['String']['output'];
  huggingFace?: Maybe<HuggingFace>;
  id: Scalars['ID']['output'];
  status: ModelStatus;
  type: ModelType;
};

export enum ModelStatus {
  Deploying = 'DEPLOYING',
  Error = 'ERROR',
  NotDeployed = 'NOT_DEPLOYED',
  Ready = 'READY'
}

export enum ModelType {
  External = 'EXTERNAL',
  Huggingface = 'HUGGINGFACE',
  Openai = 'OPENAI'
}

export type Mutation = {
  __typename?: 'Mutation';
  addModel: Model;
  addRepository: Repository;
  addStorage: Storage;
  deleteModel: Model;
  deleteRepository: Repository;
  deleteStorage: Storage;
};


export type MutationAddModelArgs = {
  input: AddModelInput;
};


export type MutationAddRepositoryArgs = {
  input: AddRepositoryInput;
};


export type MutationAddStorageArgs = {
  input: AddStorageInput;
};


export type MutationDeleteModelArgs = {
  id: Scalars['ID']['input'];
};


export type MutationDeleteRepositoryArgs = {
  id: Scalars['ID']['input'];
};


export type MutationDeleteStorageArgs = {
  id: Scalars['ID']['input'];
};

export type Query = {
  __typename?: 'Query';
  models: Array<Model>;
  repositories: Array<Repository>;
  storages: Array<Storage>;
};

export type Repository = {
  __typename?: 'Repository';
  displayName: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  owner: Scalars['String']['output'];
  type: RepositoryType;
  url: Scalars['String']['output'];
};

export enum RepositoryType {
  Bitbucket = 'BITBUCKET',
  Github = 'GITHUB',
  Gitlab = 'GITLAB'
}

export type Storage = {
  __typename?: 'Storage';
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  status: StorageStatus;
  type: StorageType;
};

export enum StorageStatus {
  Deploying = 'DEPLOYING',
  Error = 'ERROR',
  NotDeployed = 'NOT_DEPLOYED',
  Ready = 'READY'
}

export enum StorageType {
  Elasticsearch = 'ELASTICSEARCH',
  Postgres = 'POSTGRES',
  Redis = 'REDIS'
}

export type ModelsQueryVariables = Exact<{ [key: string]: never; }>;


export type ModelsQuery = { __typename?: 'Query', models: Array<{ __typename?: 'Model', id: string, type: ModelType, status: ModelStatus, displayName: string, huggingFace?: { __typename?: 'HuggingFace', organization: string, name: string } | null }> };

export type AddModelMutationVariables = Exact<{
  input: AddModelInput;
}>;


export type AddModelMutation = { __typename?: 'Mutation', addModel: { __typename?: 'Model', id: string, type: ModelType, status: ModelStatus, displayName: string, huggingFace?: { __typename?: 'HuggingFace', organization: string, name: string } | null } };

export type DeleteModelMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeleteModelMutation = { __typename?: 'Mutation', deleteModel: { __typename?: 'Model', id: string, type: ModelType, status: ModelStatus, displayName: string, huggingFace?: { __typename?: 'HuggingFace', organization: string, name: string } | null } };

export type RepositoriesQueryVariables = Exact<{ [key: string]: never; }>;


export type RepositoriesQuery = { __typename?: 'Query', repositories: Array<{ __typename?: 'Repository', id: string, type: RepositoryType, displayName: string, owner: string, name: string, url: string }> };

export type AddRepositoryMutationVariables = Exact<{
  input: AddRepositoryInput;
}>;


export type AddRepositoryMutation = { __typename?: 'Mutation', addRepository: { __typename?: 'Repository', id: string, type: RepositoryType, displayName: string, owner: string, name: string, url: string } };

export type DeleteRepositoryMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeleteRepositoryMutation = { __typename?: 'Mutation', deleteRepository: { __typename?: 'Repository', id: string, type: RepositoryType, displayName: string, owner: string, name: string, url: string } };

export type StoragesQueryVariables = Exact<{ [key: string]: never; }>;


export type StoragesQuery = { __typename?: 'Query', storages: Array<{ __typename?: 'Storage', id: string, type: StorageType, name: string, status: StorageStatus }> };

export type AddStorageMutationVariables = Exact<{
  input: AddStorageInput;
}>;


export type AddStorageMutation = { __typename?: 'Mutation', addStorage: { __typename?: 'Storage', id: string, type: StorageType, name: string, status: StorageStatus } };

export type DeleteStorageMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeleteStorageMutation = { __typename?: 'Mutation', deleteStorage: { __typename?: 'Storage', id: string, type: StorageType, name: string, status: StorageStatus } };


export const ModelsDocument = gql`
    query models {
  models {
    id
    type
    status
    displayName
    huggingFace {
      organization
      name
    }
  }
}
    `;
export type ModelsComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<ModelsQuery, ModelsQueryVariables>, 'query'>;

    export const ModelsComponent = (props: ModelsComponentProps) => (
      <ApolloReactComponents.Query<ModelsQuery, ModelsQueryVariables> query={ModelsDocument} {...props} />
    );
    

/**
 * __useModelsQuery__
 *
 * To run a query within a React component, call `useModelsQuery` and pass it any options that fit your needs.
 * When your component renders, `useModelsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useModelsQuery({
 *   variables: {
 *   },
 * });
 */
export function useModelsQuery(baseOptions?: Apollo.QueryHookOptions<ModelsQuery, ModelsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ModelsQuery, ModelsQueryVariables>(ModelsDocument, options);
      }
export function useModelsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ModelsQuery, ModelsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ModelsQuery, ModelsQueryVariables>(ModelsDocument, options);
        }
export function useModelsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<ModelsQuery, ModelsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<ModelsQuery, ModelsQueryVariables>(ModelsDocument, options);
        }
export type ModelsQueryHookResult = ReturnType<typeof useModelsQuery>;
export type ModelsLazyQueryHookResult = ReturnType<typeof useModelsLazyQuery>;
export type ModelsSuspenseQueryHookResult = ReturnType<typeof useModelsSuspenseQuery>;
export type ModelsQueryResult = Apollo.QueryResult<ModelsQuery, ModelsQueryVariables>;
export const AddModelDocument = gql`
    mutation addModel($input: AddModelInput!) {
  addModel(input: $input) {
    id
    type
    status
    displayName
    huggingFace {
      organization
      name
    }
  }
}
    `;
export type AddModelMutationFn = Apollo.MutationFunction<AddModelMutation, AddModelMutationVariables>;
export type AddModelComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<AddModelMutation, AddModelMutationVariables>, 'mutation'>;

    export const AddModelComponent = (props: AddModelComponentProps) => (
      <ApolloReactComponents.Mutation<AddModelMutation, AddModelMutationVariables> mutation={AddModelDocument} {...props} />
    );
    

/**
 * __useAddModelMutation__
 *
 * To run a mutation, you first call `useAddModelMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddModelMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addModelMutation, { data, loading, error }] = useAddModelMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useAddModelMutation(baseOptions?: Apollo.MutationHookOptions<AddModelMutation, AddModelMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddModelMutation, AddModelMutationVariables>(AddModelDocument, options);
      }
export type AddModelMutationHookResult = ReturnType<typeof useAddModelMutation>;
export type AddModelMutationResult = Apollo.MutationResult<AddModelMutation>;
export type AddModelMutationOptions = Apollo.BaseMutationOptions<AddModelMutation, AddModelMutationVariables>;
export const DeleteModelDocument = gql`
    mutation deleteModel($id: ID!) {
  deleteModel(id: $id) {
    id
    type
    status
    displayName
    huggingFace {
      organization
      name
    }
  }
}
    `;
export type DeleteModelMutationFn = Apollo.MutationFunction<DeleteModelMutation, DeleteModelMutationVariables>;
export type DeleteModelComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<DeleteModelMutation, DeleteModelMutationVariables>, 'mutation'>;

    export const DeleteModelComponent = (props: DeleteModelComponentProps) => (
      <ApolloReactComponents.Mutation<DeleteModelMutation, DeleteModelMutationVariables> mutation={DeleteModelDocument} {...props} />
    );
    

/**
 * __useDeleteModelMutation__
 *
 * To run a mutation, you first call `useDeleteModelMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteModelMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteModelMutation, { data, loading, error }] = useDeleteModelMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useDeleteModelMutation(baseOptions?: Apollo.MutationHookOptions<DeleteModelMutation, DeleteModelMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteModelMutation, DeleteModelMutationVariables>(DeleteModelDocument, options);
      }
export type DeleteModelMutationHookResult = ReturnType<typeof useDeleteModelMutation>;
export type DeleteModelMutationResult = Apollo.MutationResult<DeleteModelMutation>;
export type DeleteModelMutationOptions = Apollo.BaseMutationOptions<DeleteModelMutation, DeleteModelMutationVariables>;
export const RepositoriesDocument = gql`
    query repositories {
  repositories {
    id
    type
    displayName
    owner
    name
    url
  }
}
    `;
export type RepositoriesComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<RepositoriesQuery, RepositoriesQueryVariables>, 'query'>;

    export const RepositoriesComponent = (props: RepositoriesComponentProps) => (
      <ApolloReactComponents.Query<RepositoriesQuery, RepositoriesQueryVariables> query={RepositoriesDocument} {...props} />
    );
    

/**
 * __useRepositoriesQuery__
 *
 * To run a query within a React component, call `useRepositoriesQuery` and pass it any options that fit your needs.
 * When your component renders, `useRepositoriesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useRepositoriesQuery({
 *   variables: {
 *   },
 * });
 */
export function useRepositoriesQuery(baseOptions?: Apollo.QueryHookOptions<RepositoriesQuery, RepositoriesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<RepositoriesQuery, RepositoriesQueryVariables>(RepositoriesDocument, options);
      }
export function useRepositoriesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<RepositoriesQuery, RepositoriesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<RepositoriesQuery, RepositoriesQueryVariables>(RepositoriesDocument, options);
        }
export function useRepositoriesSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<RepositoriesQuery, RepositoriesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<RepositoriesQuery, RepositoriesQueryVariables>(RepositoriesDocument, options);
        }
export type RepositoriesQueryHookResult = ReturnType<typeof useRepositoriesQuery>;
export type RepositoriesLazyQueryHookResult = ReturnType<typeof useRepositoriesLazyQuery>;
export type RepositoriesSuspenseQueryHookResult = ReturnType<typeof useRepositoriesSuspenseQuery>;
export type RepositoriesQueryResult = Apollo.QueryResult<RepositoriesQuery, RepositoriesQueryVariables>;
export const AddRepositoryDocument = gql`
    mutation addRepository($input: AddRepositoryInput!) {
  addRepository(input: $input) {
    id
    type
    displayName
    owner
    name
    url
  }
}
    `;
export type AddRepositoryMutationFn = Apollo.MutationFunction<AddRepositoryMutation, AddRepositoryMutationVariables>;
export type AddRepositoryComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<AddRepositoryMutation, AddRepositoryMutationVariables>, 'mutation'>;

    export const AddRepositoryComponent = (props: AddRepositoryComponentProps) => (
      <ApolloReactComponents.Mutation<AddRepositoryMutation, AddRepositoryMutationVariables> mutation={AddRepositoryDocument} {...props} />
    );
    

/**
 * __useAddRepositoryMutation__
 *
 * To run a mutation, you first call `useAddRepositoryMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddRepositoryMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addRepositoryMutation, { data, loading, error }] = useAddRepositoryMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useAddRepositoryMutation(baseOptions?: Apollo.MutationHookOptions<AddRepositoryMutation, AddRepositoryMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddRepositoryMutation, AddRepositoryMutationVariables>(AddRepositoryDocument, options);
      }
export type AddRepositoryMutationHookResult = ReturnType<typeof useAddRepositoryMutation>;
export type AddRepositoryMutationResult = Apollo.MutationResult<AddRepositoryMutation>;
export type AddRepositoryMutationOptions = Apollo.BaseMutationOptions<AddRepositoryMutation, AddRepositoryMutationVariables>;
export const DeleteRepositoryDocument = gql`
    mutation deleteRepository($id: ID!) {
  deleteRepository(id: $id) {
    id
    type
    displayName
    owner
    name
    url
  }
}
    `;
export type DeleteRepositoryMutationFn = Apollo.MutationFunction<DeleteRepositoryMutation, DeleteRepositoryMutationVariables>;
export type DeleteRepositoryComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<DeleteRepositoryMutation, DeleteRepositoryMutationVariables>, 'mutation'>;

    export const DeleteRepositoryComponent = (props: DeleteRepositoryComponentProps) => (
      <ApolloReactComponents.Mutation<DeleteRepositoryMutation, DeleteRepositoryMutationVariables> mutation={DeleteRepositoryDocument} {...props} />
    );
    

/**
 * __useDeleteRepositoryMutation__
 *
 * To run a mutation, you first call `useDeleteRepositoryMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteRepositoryMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteRepositoryMutation, { data, loading, error }] = useDeleteRepositoryMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useDeleteRepositoryMutation(baseOptions?: Apollo.MutationHookOptions<DeleteRepositoryMutation, DeleteRepositoryMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteRepositoryMutation, DeleteRepositoryMutationVariables>(DeleteRepositoryDocument, options);
      }
export type DeleteRepositoryMutationHookResult = ReturnType<typeof useDeleteRepositoryMutation>;
export type DeleteRepositoryMutationResult = Apollo.MutationResult<DeleteRepositoryMutation>;
export type DeleteRepositoryMutationOptions = Apollo.BaseMutationOptions<DeleteRepositoryMutation, DeleteRepositoryMutationVariables>;
export const StoragesDocument = gql`
    query storages {
  storages {
    id
    type
    name
    status
  }
}
    `;
export type StoragesComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<StoragesQuery, StoragesQueryVariables>, 'query'>;

    export const StoragesComponent = (props: StoragesComponentProps) => (
      <ApolloReactComponents.Query<StoragesQuery, StoragesQueryVariables> query={StoragesDocument} {...props} />
    );
    

/**
 * __useStoragesQuery__
 *
 * To run a query within a React component, call `useStoragesQuery` and pass it any options that fit your needs.
 * When your component renders, `useStoragesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useStoragesQuery({
 *   variables: {
 *   },
 * });
 */
export function useStoragesQuery(baseOptions?: Apollo.QueryHookOptions<StoragesQuery, StoragesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<StoragesQuery, StoragesQueryVariables>(StoragesDocument, options);
      }
export function useStoragesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<StoragesQuery, StoragesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<StoragesQuery, StoragesQueryVariables>(StoragesDocument, options);
        }
export function useStoragesSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<StoragesQuery, StoragesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<StoragesQuery, StoragesQueryVariables>(StoragesDocument, options);
        }
export type StoragesQueryHookResult = ReturnType<typeof useStoragesQuery>;
export type StoragesLazyQueryHookResult = ReturnType<typeof useStoragesLazyQuery>;
export type StoragesSuspenseQueryHookResult = ReturnType<typeof useStoragesSuspenseQuery>;
export type StoragesQueryResult = Apollo.QueryResult<StoragesQuery, StoragesQueryVariables>;
export const AddStorageDocument = gql`
    mutation addStorage($input: AddStorageInput!) {
  addStorage(input: $input) {
    id
    type
    name
    status
  }
}
    `;
export type AddStorageMutationFn = Apollo.MutationFunction<AddStorageMutation, AddStorageMutationVariables>;
export type AddStorageComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<AddStorageMutation, AddStorageMutationVariables>, 'mutation'>;

    export const AddStorageComponent = (props: AddStorageComponentProps) => (
      <ApolloReactComponents.Mutation<AddStorageMutation, AddStorageMutationVariables> mutation={AddStorageDocument} {...props} />
    );
    

/**
 * __useAddStorageMutation__
 *
 * To run a mutation, you first call `useAddStorageMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddStorageMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addStorageMutation, { data, loading, error }] = useAddStorageMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useAddStorageMutation(baseOptions?: Apollo.MutationHookOptions<AddStorageMutation, AddStorageMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddStorageMutation, AddStorageMutationVariables>(AddStorageDocument, options);
      }
export type AddStorageMutationHookResult = ReturnType<typeof useAddStorageMutation>;
export type AddStorageMutationResult = Apollo.MutationResult<AddStorageMutation>;
export type AddStorageMutationOptions = Apollo.BaseMutationOptions<AddStorageMutation, AddStorageMutationVariables>;
export const DeleteStorageDocument = gql`
    mutation deleteStorage($id: ID!) {
  deleteStorage(id: $id) {
    id
    type
    name
    status
  }
}
    `;
export type DeleteStorageMutationFn = Apollo.MutationFunction<DeleteStorageMutation, DeleteStorageMutationVariables>;
export type DeleteStorageComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<DeleteStorageMutation, DeleteStorageMutationVariables>, 'mutation'>;

    export const DeleteStorageComponent = (props: DeleteStorageComponentProps) => (
      <ApolloReactComponents.Mutation<DeleteStorageMutation, DeleteStorageMutationVariables> mutation={DeleteStorageDocument} {...props} />
    );
    

/**
 * __useDeleteStorageMutation__
 *
 * To run a mutation, you first call `useDeleteStorageMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteStorageMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteStorageMutation, { data, loading, error }] = useDeleteStorageMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useDeleteStorageMutation(baseOptions?: Apollo.MutationHookOptions<DeleteStorageMutation, DeleteStorageMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteStorageMutation, DeleteStorageMutationVariables>(DeleteStorageDocument, options);
      }
export type DeleteStorageMutationHookResult = ReturnType<typeof useDeleteStorageMutation>;
export type DeleteStorageMutationResult = Apollo.MutationResult<DeleteStorageMutation>;
export type DeleteStorageMutationOptions = Apollo.BaseMutationOptions<DeleteStorageMutation, DeleteStorageMutationVariables>;