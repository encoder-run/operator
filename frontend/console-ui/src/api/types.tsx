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

export type AddModelDeploymentInput = {
  cpu: Scalars['String']['input'];
  id: Scalars['ID']['input'];
  memory: Scalars['String']['input'];
};

export type AddModelInput = {
  huggingFace?: InputMaybe<HuggingFaceInput>;
  type: ModelType;
};

export type AddPipelineDeploymentInput = {
  enabled: Scalars['Boolean']['input'];
  id: Scalars['ID']['input'];
};

export type AddPipelineInput = {
  name: Scalars['String']['input'];
  repositoryEmbeddings?: InputMaybe<AddRepositoryEmbeddingsInput>;
  type: PipelineType;
};

export type AddRepositoryEmbeddingsInput = {
  modelID: Scalars['ID']['input'];
  repositoryID: Scalars['ID']['input'];
  storageID: Scalars['ID']['input'];
};

export type AddRepositoryInput = {
  branch?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  owner?: InputMaybe<Scalars['String']['input']>;
  token?: InputMaybe<Scalars['String']['input']>;
  type?: InputMaybe<RepositoryType>;
  url?: InputMaybe<Scalars['String']['input']>;
};

export type AddStorageDeploymentInput = {
  cpu: Scalars['String']['input'];
  id: Scalars['ID']['input'];
  memory: Scalars['String']['input'];
};

export type AddStorageInput = {
  name: Scalars['String']['input'];
  type: StorageType;
};

export type HuggingFace = {
  __typename?: 'HuggingFace';
  maxSequenceLength: Scalars['Int']['output'];
  name: Scalars['String']['output'];
  organization: Scalars['String']['output'];
};

export type HuggingFaceInput = {
  maxSequenceLength: Scalars['Int']['input'];
  name: Scalars['String']['input'];
  organization: Scalars['String']['input'];
};

export type Model = {
  __typename?: 'Model';
  deployment?: Maybe<ModelDeployment>;
  displayName: Scalars['String']['output'];
  huggingFace?: Maybe<HuggingFace>;
  id: Scalars['ID']['output'];
  status: ModelStatus;
  type: ModelType;
};

export type ModelDeployment = {
  __typename?: 'ModelDeployment';
  cpu: Scalars['String']['output'];
  enabled: Scalars['Boolean']['output'];
  memory: Scalars['String']['output'];
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
  addModelDeployment: Model;
  addPipeline: Pipeline;
  addPipelineDeployment: Pipeline;
  addRepository: Repository;
  addStorage: Storage;
  addStorageDeployment: Storage;
  deleteModel: Model;
  deletePipeline: Pipeline;
  deleteRepository: Repository;
  deleteStorage: Storage;
};


export type MutationAddModelArgs = {
  input: AddModelInput;
};


export type MutationAddModelDeploymentArgs = {
  input: AddModelDeploymentInput;
};


export type MutationAddPipelineArgs = {
  input: AddPipelineInput;
};


export type MutationAddPipelineDeploymentArgs = {
  input: AddPipelineDeploymentInput;
};


export type MutationAddRepositoryArgs = {
  input: AddRepositoryInput;
};


export type MutationAddStorageArgs = {
  input: AddStorageInput;
};


export type MutationAddStorageDeploymentArgs = {
  input: AddStorageDeploymentInput;
};


export type MutationDeleteModelArgs = {
  id: Scalars['ID']['input'];
};


export type MutationDeletePipelineArgs = {
  id: Scalars['ID']['input'];
};


export type MutationDeleteRepositoryArgs = {
  id: Scalars['ID']['input'];
};


export type MutationDeleteStorageArgs = {
  id: Scalars['ID']['input'];
};

export type Pipeline = {
  __typename?: 'Pipeline';
  enabled: Scalars['Boolean']['output'];
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  repositoryEmbeddings?: Maybe<RepositoryEmbeddings>;
  status: PipelineStatus;
  type: PipelineType;
};

export type PipelineExecution = {
  __typename?: 'PipelineExecution';
  id: Scalars['ID']['output'];
  status: PipelineExecutionStatus;
};

export enum PipelineExecutionStatus {
  Active = 'ACTIVE',
  Failed = 'FAILED',
  Pending = 'PENDING',
  Succeeded = 'SUCCEEDED'
}

export enum PipelineStatus {
  Deploying = 'DEPLOYING',
  Error = 'ERROR',
  NotDeployed = 'NOT_DEPLOYED',
  Ready = 'READY'
}

export enum PipelineType {
  RepositoryEmbeddings = 'REPOSITORY_EMBEDDINGS'
}

export type Query = {
  __typename?: 'Query';
  getModel: Model;
  getPipeline: Pipeline;
  getPipelineExecutions: Array<PipelineExecution>;
  getRepository: Repository;
  getStorage: Storage;
  models: Array<Model>;
  pipelines: Array<Pipeline>;
  repositories: Array<Repository>;
  semanticSearch: Array<SearchResult>;
  storages: Array<Storage>;
};


export type QueryGetModelArgs = {
  id: Scalars['ID']['input'];
};


export type QueryGetPipelineArgs = {
  id: Scalars['ID']['input'];
};


export type QueryGetPipelineExecutionsArgs = {
  id: Scalars['ID']['input'];
};


export type QueryGetRepositoryArgs = {
  id: Scalars['ID']['input'];
};


export type QueryGetStorageArgs = {
  id: Scalars['ID']['input'];
};


export type QuerySemanticSearchArgs = {
  query: QueryInput;
};

export type QueryInput = {
  limit?: InputMaybe<Scalars['Int']['input']>;
  page?: InputMaybe<Scalars['Int']['input']>;
  query: Scalars['String']['input'];
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

export type RepositoryEmbeddings = {
  __typename?: 'RepositoryEmbeddings';
  modelID: Scalars['ID']['output'];
  repositoryID: Scalars['ID']['output'];
  storageID: Scalars['ID']['output'];
};

export enum RepositoryType {
  Bitbucket = 'BITBUCKET',
  Github = 'GITHUB',
  Gitlab = 'GITLAB'
}

export type SearchResult = {
  __typename?: 'SearchResult';
  chunkID: Scalars['Int']['output'];
  content: Scalars['String']['output'];
  endIndex: Scalars['Int']['output'];
  hash: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  owner: Scalars['String']['output'];
  path: Scalars['String']['output'];
  repo: Scalars['String']['output'];
  score: Scalars['Float']['output'];
  startIndex: Scalars['Int']['output'];
  startLine: Scalars['Int']['output'];
};

export type Storage = {
  __typename?: 'Storage';
  deployment?: Maybe<StorageDeployment>;
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  status: StorageStatus;
  type: StorageType;
};

export type StorageDeployment = {
  __typename?: 'StorageDeployment';
  cpu: Scalars['String']['output'];
  enabled: Scalars['Boolean']['output'];
  memory: Scalars['String']['output'];
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


export type ModelsQuery = { __typename?: 'Query', models: Array<{ __typename?: 'Model', id: string, type: ModelType, status: ModelStatus, displayName: string, huggingFace?: { __typename?: 'HuggingFace', organization: string, name: string, maxSequenceLength: number } | null, deployment?: { __typename?: 'ModelDeployment', enabled: boolean, cpu: string, memory: string } | null }> };

export type GetModelQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type GetModelQuery = { __typename?: 'Query', getModel: { __typename?: 'Model', id: string, type: ModelType, status: ModelStatus, displayName: string, huggingFace?: { __typename?: 'HuggingFace', organization: string, name: string, maxSequenceLength: number } | null, deployment?: { __typename?: 'ModelDeployment', enabled: boolean, cpu: string, memory: string } | null } };

export type AddModelMutationVariables = Exact<{
  input: AddModelInput;
}>;


export type AddModelMutation = { __typename?: 'Mutation', addModel: { __typename?: 'Model', id: string, type: ModelType, status: ModelStatus, displayName: string, huggingFace?: { __typename?: 'HuggingFace', organization: string, name: string, maxSequenceLength: number } | null, deployment?: { __typename?: 'ModelDeployment', enabled: boolean, cpu: string, memory: string } | null } };

export type AddModelDeploymentMutationVariables = Exact<{
  input: AddModelDeploymentInput;
}>;


export type AddModelDeploymentMutation = { __typename?: 'Mutation', addModelDeployment: { __typename?: 'Model', id: string, type: ModelType, status: ModelStatus, displayName: string, huggingFace?: { __typename?: 'HuggingFace', organization: string, name: string, maxSequenceLength: number } | null, deployment?: { __typename?: 'ModelDeployment', enabled: boolean, cpu: string, memory: string } | null } };

export type DeleteModelMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeleteModelMutation = { __typename?: 'Mutation', deleteModel: { __typename?: 'Model', id: string, type: ModelType, status: ModelStatus, displayName: string, huggingFace?: { __typename?: 'HuggingFace', organization: string, name: string, maxSequenceLength: number } | null, deployment?: { __typename?: 'ModelDeployment', enabled: boolean, cpu: string, memory: string } | null } };

export type PipelinesQueryVariables = Exact<{ [key: string]: never; }>;


export type PipelinesQuery = { __typename?: 'Query', pipelines: Array<{ __typename?: 'Pipeline', id: string, name: string, type: PipelineType, enabled: boolean, status: PipelineStatus, repositoryEmbeddings?: { __typename?: 'RepositoryEmbeddings', repositoryID: string, modelID: string, storageID: string } | null }> };

export type GetPipelineQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type GetPipelineQuery = { __typename?: 'Query', getPipeline: { __typename?: 'Pipeline', id: string, name: string, type: PipelineType, enabled: boolean, status: PipelineStatus, repositoryEmbeddings?: { __typename?: 'RepositoryEmbeddings', repositoryID: string, modelID: string, storageID: string } | null } };

export type GetPipelineExecutionsQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type GetPipelineExecutionsQuery = { __typename?: 'Query', getPipelineExecutions: Array<{ __typename?: 'PipelineExecution', id: string, status: PipelineExecutionStatus }> };

export type AddPipelineMutationVariables = Exact<{
  input: AddPipelineInput;
}>;


export type AddPipelineMutation = { __typename?: 'Mutation', addPipeline: { __typename?: 'Pipeline', id: string, name: string, type: PipelineType, enabled: boolean, status: PipelineStatus, repositoryEmbeddings?: { __typename?: 'RepositoryEmbeddings', repositoryID: string, modelID: string, storageID: string } | null } };

export type AddPipelineDeploymentMutationVariables = Exact<{
  input: AddPipelineDeploymentInput;
}>;


export type AddPipelineDeploymentMutation = { __typename?: 'Mutation', addPipelineDeployment: { __typename?: 'Pipeline', id: string, name: string, type: PipelineType, enabled: boolean, status: PipelineStatus, repositoryEmbeddings?: { __typename?: 'RepositoryEmbeddings', repositoryID: string, modelID: string, storageID: string } | null } };

export type DeletePipelineMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeletePipelineMutation = { __typename?: 'Mutation', deletePipeline: { __typename?: 'Pipeline', id: string, name: string, type: PipelineType, enabled: boolean, status: PipelineStatus, repositoryEmbeddings?: { __typename?: 'RepositoryEmbeddings', repositoryID: string, modelID: string, storageID: string } | null } };

export type RepositoriesQueryVariables = Exact<{ [key: string]: never; }>;


export type RepositoriesQuery = { __typename?: 'Query', repositories: Array<{ __typename?: 'Repository', id: string, type: RepositoryType, displayName: string, owner: string, name: string, url: string }> };

export type GetRepositoryQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type GetRepositoryQuery = { __typename?: 'Query', getRepository: { __typename?: 'Repository', id: string, type: RepositoryType, displayName: string, owner: string, name: string, url: string } };

export type AddRepositoryMutationVariables = Exact<{
  input: AddRepositoryInput;
}>;


export type AddRepositoryMutation = { __typename?: 'Mutation', addRepository: { __typename?: 'Repository', id: string, type: RepositoryType, displayName: string, owner: string, name: string, url: string } };

export type DeleteRepositoryMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeleteRepositoryMutation = { __typename?: 'Mutation', deleteRepository: { __typename?: 'Repository', id: string, type: RepositoryType, displayName: string, owner: string, name: string, url: string } };

export type SemanticSearchQueryVariables = Exact<{
  query: QueryInput;
}>;


export type SemanticSearchQuery = { __typename?: 'Query', semanticSearch: Array<{ __typename?: 'SearchResult', id: string, chunkID: number, content: string, hash: string, path: string, owner: string, repo: string, startIndex: number, startLine: number, endIndex: number, score: number }> };

export type StoragesQueryVariables = Exact<{ [key: string]: never; }>;


export type StoragesQuery = { __typename?: 'Query', storages: Array<{ __typename?: 'Storage', id: string, type: StorageType, name: string, status: StorageStatus, deployment?: { __typename?: 'StorageDeployment', enabled: boolean, cpu: string, memory: string } | null }> };

export type GetStorageQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type GetStorageQuery = { __typename?: 'Query', getStorage: { __typename?: 'Storage', id: string, type: StorageType, name: string, status: StorageStatus, deployment?: { __typename?: 'StorageDeployment', enabled: boolean, cpu: string, memory: string } | null } };

export type AddStorageMutationVariables = Exact<{
  input: AddStorageInput;
}>;


export type AddStorageMutation = { __typename?: 'Mutation', addStorage: { __typename?: 'Storage', id: string, type: StorageType, name: string, status: StorageStatus, deployment?: { __typename?: 'StorageDeployment', enabled: boolean, cpu: string, memory: string } | null } };

export type AddStorageDeploymentMutationVariables = Exact<{
  input: AddStorageDeploymentInput;
}>;


export type AddStorageDeploymentMutation = { __typename?: 'Mutation', addStorageDeployment: { __typename?: 'Storage', id: string, type: StorageType, name: string, status: StorageStatus, deployment?: { __typename?: 'StorageDeployment', enabled: boolean, cpu: string, memory: string } | null } };

export type DeleteStorageMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeleteStorageMutation = { __typename?: 'Mutation', deleteStorage: { __typename?: 'Storage', id: string, type: StorageType, name: string, status: StorageStatus, deployment?: { __typename?: 'StorageDeployment', enabled: boolean, cpu: string, memory: string } | null } };


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
      maxSequenceLength
    }
    deployment {
      enabled
      cpu
      memory
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
export const GetModelDocument = gql`
    query getModel($id: ID!) {
  getModel(id: $id) {
    id
    type
    status
    displayName
    huggingFace {
      organization
      name
      maxSequenceLength
    }
    deployment {
      enabled
      cpu
      memory
    }
  }
}
    `;
export type GetModelComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<GetModelQuery, GetModelQueryVariables>, 'query'> & ({ variables: GetModelQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const GetModelComponent = (props: GetModelComponentProps) => (
      <ApolloReactComponents.Query<GetModelQuery, GetModelQueryVariables> query={GetModelDocument} {...props} />
    );
    

/**
 * __useGetModelQuery__
 *
 * To run a query within a React component, call `useGetModelQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetModelQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetModelQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetModelQuery(baseOptions: Apollo.QueryHookOptions<GetModelQuery, GetModelQueryVariables> & ({ variables: GetModelQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetModelQuery, GetModelQueryVariables>(GetModelDocument, options);
      }
export function useGetModelLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetModelQuery, GetModelQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetModelQuery, GetModelQueryVariables>(GetModelDocument, options);
        }
export function useGetModelSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetModelQuery, GetModelQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetModelQuery, GetModelQueryVariables>(GetModelDocument, options);
        }
export type GetModelQueryHookResult = ReturnType<typeof useGetModelQuery>;
export type GetModelLazyQueryHookResult = ReturnType<typeof useGetModelLazyQuery>;
export type GetModelSuspenseQueryHookResult = ReturnType<typeof useGetModelSuspenseQuery>;
export type GetModelQueryResult = Apollo.QueryResult<GetModelQuery, GetModelQueryVariables>;
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
      maxSequenceLength
    }
    deployment {
      enabled
      cpu
      memory
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
export const AddModelDeploymentDocument = gql`
    mutation addModelDeployment($input: AddModelDeploymentInput!) {
  addModelDeployment(input: $input) {
    id
    type
    status
    displayName
    huggingFace {
      organization
      name
      maxSequenceLength
    }
    deployment {
      enabled
      cpu
      memory
    }
  }
}
    `;
export type AddModelDeploymentMutationFn = Apollo.MutationFunction<AddModelDeploymentMutation, AddModelDeploymentMutationVariables>;
export type AddModelDeploymentComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<AddModelDeploymentMutation, AddModelDeploymentMutationVariables>, 'mutation'>;

    export const AddModelDeploymentComponent = (props: AddModelDeploymentComponentProps) => (
      <ApolloReactComponents.Mutation<AddModelDeploymentMutation, AddModelDeploymentMutationVariables> mutation={AddModelDeploymentDocument} {...props} />
    );
    

/**
 * __useAddModelDeploymentMutation__
 *
 * To run a mutation, you first call `useAddModelDeploymentMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddModelDeploymentMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addModelDeploymentMutation, { data, loading, error }] = useAddModelDeploymentMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useAddModelDeploymentMutation(baseOptions?: Apollo.MutationHookOptions<AddModelDeploymentMutation, AddModelDeploymentMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddModelDeploymentMutation, AddModelDeploymentMutationVariables>(AddModelDeploymentDocument, options);
      }
export type AddModelDeploymentMutationHookResult = ReturnType<typeof useAddModelDeploymentMutation>;
export type AddModelDeploymentMutationResult = Apollo.MutationResult<AddModelDeploymentMutation>;
export type AddModelDeploymentMutationOptions = Apollo.BaseMutationOptions<AddModelDeploymentMutation, AddModelDeploymentMutationVariables>;
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
      maxSequenceLength
    }
    deployment {
      enabled
      cpu
      memory
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
export const PipelinesDocument = gql`
    query pipelines {
  pipelines {
    id
    name
    type
    enabled
    status
    repositoryEmbeddings {
      repositoryID
      modelID
      storageID
    }
  }
}
    `;
export type PipelinesComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<PipelinesQuery, PipelinesQueryVariables>, 'query'>;

    export const PipelinesComponent = (props: PipelinesComponentProps) => (
      <ApolloReactComponents.Query<PipelinesQuery, PipelinesQueryVariables> query={PipelinesDocument} {...props} />
    );
    

/**
 * __usePipelinesQuery__
 *
 * To run a query within a React component, call `usePipelinesQuery` and pass it any options that fit your needs.
 * When your component renders, `usePipelinesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = usePipelinesQuery({
 *   variables: {
 *   },
 * });
 */
export function usePipelinesQuery(baseOptions?: Apollo.QueryHookOptions<PipelinesQuery, PipelinesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<PipelinesQuery, PipelinesQueryVariables>(PipelinesDocument, options);
      }
export function usePipelinesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<PipelinesQuery, PipelinesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<PipelinesQuery, PipelinesQueryVariables>(PipelinesDocument, options);
        }
export function usePipelinesSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<PipelinesQuery, PipelinesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<PipelinesQuery, PipelinesQueryVariables>(PipelinesDocument, options);
        }
export type PipelinesQueryHookResult = ReturnType<typeof usePipelinesQuery>;
export type PipelinesLazyQueryHookResult = ReturnType<typeof usePipelinesLazyQuery>;
export type PipelinesSuspenseQueryHookResult = ReturnType<typeof usePipelinesSuspenseQuery>;
export type PipelinesQueryResult = Apollo.QueryResult<PipelinesQuery, PipelinesQueryVariables>;
export const GetPipelineDocument = gql`
    query getPipeline($id: ID!) {
  getPipeline(id: $id) {
    id
    name
    type
    enabled
    status
    repositoryEmbeddings {
      repositoryID
      modelID
      storageID
    }
  }
}
    `;
export type GetPipelineComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<GetPipelineQuery, GetPipelineQueryVariables>, 'query'> & ({ variables: GetPipelineQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const GetPipelineComponent = (props: GetPipelineComponentProps) => (
      <ApolloReactComponents.Query<GetPipelineQuery, GetPipelineQueryVariables> query={GetPipelineDocument} {...props} />
    );
    

/**
 * __useGetPipelineQuery__
 *
 * To run a query within a React component, call `useGetPipelineQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetPipelineQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetPipelineQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetPipelineQuery(baseOptions: Apollo.QueryHookOptions<GetPipelineQuery, GetPipelineQueryVariables> & ({ variables: GetPipelineQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPipelineQuery, GetPipelineQueryVariables>(GetPipelineDocument, options);
      }
export function useGetPipelineLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPipelineQuery, GetPipelineQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPipelineQuery, GetPipelineQueryVariables>(GetPipelineDocument, options);
        }
export function useGetPipelineSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetPipelineQuery, GetPipelineQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetPipelineQuery, GetPipelineQueryVariables>(GetPipelineDocument, options);
        }
export type GetPipelineQueryHookResult = ReturnType<typeof useGetPipelineQuery>;
export type GetPipelineLazyQueryHookResult = ReturnType<typeof useGetPipelineLazyQuery>;
export type GetPipelineSuspenseQueryHookResult = ReturnType<typeof useGetPipelineSuspenseQuery>;
export type GetPipelineQueryResult = Apollo.QueryResult<GetPipelineQuery, GetPipelineQueryVariables>;
export const GetPipelineExecutionsDocument = gql`
    query getPipelineExecutions($id: ID!) {
  getPipelineExecutions(id: $id) {
    id
    status
  }
}
    `;
export type GetPipelineExecutionsComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<GetPipelineExecutionsQuery, GetPipelineExecutionsQueryVariables>, 'query'> & ({ variables: GetPipelineExecutionsQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const GetPipelineExecutionsComponent = (props: GetPipelineExecutionsComponentProps) => (
      <ApolloReactComponents.Query<GetPipelineExecutionsQuery, GetPipelineExecutionsQueryVariables> query={GetPipelineExecutionsDocument} {...props} />
    );
    

/**
 * __useGetPipelineExecutionsQuery__
 *
 * To run a query within a React component, call `useGetPipelineExecutionsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetPipelineExecutionsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetPipelineExecutionsQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetPipelineExecutionsQuery(baseOptions: Apollo.QueryHookOptions<GetPipelineExecutionsQuery, GetPipelineExecutionsQueryVariables> & ({ variables: GetPipelineExecutionsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPipelineExecutionsQuery, GetPipelineExecutionsQueryVariables>(GetPipelineExecutionsDocument, options);
      }
export function useGetPipelineExecutionsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPipelineExecutionsQuery, GetPipelineExecutionsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPipelineExecutionsQuery, GetPipelineExecutionsQueryVariables>(GetPipelineExecutionsDocument, options);
        }
export function useGetPipelineExecutionsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetPipelineExecutionsQuery, GetPipelineExecutionsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetPipelineExecutionsQuery, GetPipelineExecutionsQueryVariables>(GetPipelineExecutionsDocument, options);
        }
export type GetPipelineExecutionsQueryHookResult = ReturnType<typeof useGetPipelineExecutionsQuery>;
export type GetPipelineExecutionsLazyQueryHookResult = ReturnType<typeof useGetPipelineExecutionsLazyQuery>;
export type GetPipelineExecutionsSuspenseQueryHookResult = ReturnType<typeof useGetPipelineExecutionsSuspenseQuery>;
export type GetPipelineExecutionsQueryResult = Apollo.QueryResult<GetPipelineExecutionsQuery, GetPipelineExecutionsQueryVariables>;
export const AddPipelineDocument = gql`
    mutation addPipeline($input: AddPipelineInput!) {
  addPipeline(input: $input) {
    id
    name
    type
    enabled
    status
    repositoryEmbeddings {
      repositoryID
      modelID
      storageID
    }
  }
}
    `;
export type AddPipelineMutationFn = Apollo.MutationFunction<AddPipelineMutation, AddPipelineMutationVariables>;
export type AddPipelineComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<AddPipelineMutation, AddPipelineMutationVariables>, 'mutation'>;

    export const AddPipelineComponent = (props: AddPipelineComponentProps) => (
      <ApolloReactComponents.Mutation<AddPipelineMutation, AddPipelineMutationVariables> mutation={AddPipelineDocument} {...props} />
    );
    

/**
 * __useAddPipelineMutation__
 *
 * To run a mutation, you first call `useAddPipelineMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddPipelineMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addPipelineMutation, { data, loading, error }] = useAddPipelineMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useAddPipelineMutation(baseOptions?: Apollo.MutationHookOptions<AddPipelineMutation, AddPipelineMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddPipelineMutation, AddPipelineMutationVariables>(AddPipelineDocument, options);
      }
export type AddPipelineMutationHookResult = ReturnType<typeof useAddPipelineMutation>;
export type AddPipelineMutationResult = Apollo.MutationResult<AddPipelineMutation>;
export type AddPipelineMutationOptions = Apollo.BaseMutationOptions<AddPipelineMutation, AddPipelineMutationVariables>;
export const AddPipelineDeploymentDocument = gql`
    mutation addPipelineDeployment($input: AddPipelineDeploymentInput!) {
  addPipelineDeployment(input: $input) {
    id
    name
    type
    enabled
    status
    repositoryEmbeddings {
      repositoryID
      modelID
      storageID
    }
  }
}
    `;
export type AddPipelineDeploymentMutationFn = Apollo.MutationFunction<AddPipelineDeploymentMutation, AddPipelineDeploymentMutationVariables>;
export type AddPipelineDeploymentComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<AddPipelineDeploymentMutation, AddPipelineDeploymentMutationVariables>, 'mutation'>;

    export const AddPipelineDeploymentComponent = (props: AddPipelineDeploymentComponentProps) => (
      <ApolloReactComponents.Mutation<AddPipelineDeploymentMutation, AddPipelineDeploymentMutationVariables> mutation={AddPipelineDeploymentDocument} {...props} />
    );
    

/**
 * __useAddPipelineDeploymentMutation__
 *
 * To run a mutation, you first call `useAddPipelineDeploymentMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddPipelineDeploymentMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addPipelineDeploymentMutation, { data, loading, error }] = useAddPipelineDeploymentMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useAddPipelineDeploymentMutation(baseOptions?: Apollo.MutationHookOptions<AddPipelineDeploymentMutation, AddPipelineDeploymentMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddPipelineDeploymentMutation, AddPipelineDeploymentMutationVariables>(AddPipelineDeploymentDocument, options);
      }
export type AddPipelineDeploymentMutationHookResult = ReturnType<typeof useAddPipelineDeploymentMutation>;
export type AddPipelineDeploymentMutationResult = Apollo.MutationResult<AddPipelineDeploymentMutation>;
export type AddPipelineDeploymentMutationOptions = Apollo.BaseMutationOptions<AddPipelineDeploymentMutation, AddPipelineDeploymentMutationVariables>;
export const DeletePipelineDocument = gql`
    mutation deletePipeline($id: ID!) {
  deletePipeline(id: $id) {
    id
    name
    type
    enabled
    status
    repositoryEmbeddings {
      repositoryID
      modelID
      storageID
    }
  }
}
    `;
export type DeletePipelineMutationFn = Apollo.MutationFunction<DeletePipelineMutation, DeletePipelineMutationVariables>;
export type DeletePipelineComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<DeletePipelineMutation, DeletePipelineMutationVariables>, 'mutation'>;

    export const DeletePipelineComponent = (props: DeletePipelineComponentProps) => (
      <ApolloReactComponents.Mutation<DeletePipelineMutation, DeletePipelineMutationVariables> mutation={DeletePipelineDocument} {...props} />
    );
    

/**
 * __useDeletePipelineMutation__
 *
 * To run a mutation, you first call `useDeletePipelineMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeletePipelineMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deletePipelineMutation, { data, loading, error }] = useDeletePipelineMutation({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useDeletePipelineMutation(baseOptions?: Apollo.MutationHookOptions<DeletePipelineMutation, DeletePipelineMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeletePipelineMutation, DeletePipelineMutationVariables>(DeletePipelineDocument, options);
      }
export type DeletePipelineMutationHookResult = ReturnType<typeof useDeletePipelineMutation>;
export type DeletePipelineMutationResult = Apollo.MutationResult<DeletePipelineMutation>;
export type DeletePipelineMutationOptions = Apollo.BaseMutationOptions<DeletePipelineMutation, DeletePipelineMutationVariables>;
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
export const GetRepositoryDocument = gql`
    query getRepository($id: ID!) {
  getRepository(id: $id) {
    id
    type
    displayName
    owner
    name
    url
  }
}
    `;
export type GetRepositoryComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<GetRepositoryQuery, GetRepositoryQueryVariables>, 'query'> & ({ variables: GetRepositoryQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const GetRepositoryComponent = (props: GetRepositoryComponentProps) => (
      <ApolloReactComponents.Query<GetRepositoryQuery, GetRepositoryQueryVariables> query={GetRepositoryDocument} {...props} />
    );
    

/**
 * __useGetRepositoryQuery__
 *
 * To run a query within a React component, call `useGetRepositoryQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetRepositoryQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetRepositoryQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetRepositoryQuery(baseOptions: Apollo.QueryHookOptions<GetRepositoryQuery, GetRepositoryQueryVariables> & ({ variables: GetRepositoryQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetRepositoryQuery, GetRepositoryQueryVariables>(GetRepositoryDocument, options);
      }
export function useGetRepositoryLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetRepositoryQuery, GetRepositoryQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetRepositoryQuery, GetRepositoryQueryVariables>(GetRepositoryDocument, options);
        }
export function useGetRepositorySuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetRepositoryQuery, GetRepositoryQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetRepositoryQuery, GetRepositoryQueryVariables>(GetRepositoryDocument, options);
        }
export type GetRepositoryQueryHookResult = ReturnType<typeof useGetRepositoryQuery>;
export type GetRepositoryLazyQueryHookResult = ReturnType<typeof useGetRepositoryLazyQuery>;
export type GetRepositorySuspenseQueryHookResult = ReturnType<typeof useGetRepositorySuspenseQuery>;
export type GetRepositoryQueryResult = Apollo.QueryResult<GetRepositoryQuery, GetRepositoryQueryVariables>;
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
export const SemanticSearchDocument = gql`
    query semanticSearch($query: QueryInput!) {
  semanticSearch(query: $query) {
    id
    chunkID
    content
    hash
    path
    owner
    repo
    startIndex
    startLine
    endIndex
    score
  }
}
    `;
export type SemanticSearchComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<SemanticSearchQuery, SemanticSearchQueryVariables>, 'query'> & ({ variables: SemanticSearchQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const SemanticSearchComponent = (props: SemanticSearchComponentProps) => (
      <ApolloReactComponents.Query<SemanticSearchQuery, SemanticSearchQueryVariables> query={SemanticSearchDocument} {...props} />
    );
    

/**
 * __useSemanticSearchQuery__
 *
 * To run a query within a React component, call `useSemanticSearchQuery` and pass it any options that fit your needs.
 * When your component renders, `useSemanticSearchQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useSemanticSearchQuery({
 *   variables: {
 *      query: // value for 'query'
 *   },
 * });
 */
export function useSemanticSearchQuery(baseOptions: Apollo.QueryHookOptions<SemanticSearchQuery, SemanticSearchQueryVariables> & ({ variables: SemanticSearchQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<SemanticSearchQuery, SemanticSearchQueryVariables>(SemanticSearchDocument, options);
      }
export function useSemanticSearchLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<SemanticSearchQuery, SemanticSearchQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<SemanticSearchQuery, SemanticSearchQueryVariables>(SemanticSearchDocument, options);
        }
export function useSemanticSearchSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<SemanticSearchQuery, SemanticSearchQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<SemanticSearchQuery, SemanticSearchQueryVariables>(SemanticSearchDocument, options);
        }
export type SemanticSearchQueryHookResult = ReturnType<typeof useSemanticSearchQuery>;
export type SemanticSearchLazyQueryHookResult = ReturnType<typeof useSemanticSearchLazyQuery>;
export type SemanticSearchSuspenseQueryHookResult = ReturnType<typeof useSemanticSearchSuspenseQuery>;
export type SemanticSearchQueryResult = Apollo.QueryResult<SemanticSearchQuery, SemanticSearchQueryVariables>;
export const StoragesDocument = gql`
    query storages {
  storages {
    id
    type
    name
    status
    deployment {
      enabled
      cpu
      memory
    }
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
export const GetStorageDocument = gql`
    query getStorage($id: ID!) {
  getStorage(id: $id) {
    id
    type
    name
    status
    deployment {
      enabled
      cpu
      memory
    }
  }
}
    `;
export type GetStorageComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<GetStorageQuery, GetStorageQueryVariables>, 'query'> & ({ variables: GetStorageQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const GetStorageComponent = (props: GetStorageComponentProps) => (
      <ApolloReactComponents.Query<GetStorageQuery, GetStorageQueryVariables> query={GetStorageDocument} {...props} />
    );
    

/**
 * __useGetStorageQuery__
 *
 * To run a query within a React component, call `useGetStorageQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetStorageQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetStorageQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetStorageQuery(baseOptions: Apollo.QueryHookOptions<GetStorageQuery, GetStorageQueryVariables> & ({ variables: GetStorageQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetStorageQuery, GetStorageQueryVariables>(GetStorageDocument, options);
      }
export function useGetStorageLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetStorageQuery, GetStorageQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetStorageQuery, GetStorageQueryVariables>(GetStorageDocument, options);
        }
export function useGetStorageSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetStorageQuery, GetStorageQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetStorageQuery, GetStorageQueryVariables>(GetStorageDocument, options);
        }
export type GetStorageQueryHookResult = ReturnType<typeof useGetStorageQuery>;
export type GetStorageLazyQueryHookResult = ReturnType<typeof useGetStorageLazyQuery>;
export type GetStorageSuspenseQueryHookResult = ReturnType<typeof useGetStorageSuspenseQuery>;
export type GetStorageQueryResult = Apollo.QueryResult<GetStorageQuery, GetStorageQueryVariables>;
export const AddStorageDocument = gql`
    mutation addStorage($input: AddStorageInput!) {
  addStorage(input: $input) {
    id
    type
    name
    status
    deployment {
      enabled
      cpu
      memory
    }
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
export const AddStorageDeploymentDocument = gql`
    mutation addStorageDeployment($input: AddStorageDeploymentInput!) {
  addStorageDeployment(input: $input) {
    id
    type
    name
    status
    deployment {
      enabled
      cpu
      memory
    }
  }
}
    `;
export type AddStorageDeploymentMutationFn = Apollo.MutationFunction<AddStorageDeploymentMutation, AddStorageDeploymentMutationVariables>;
export type AddStorageDeploymentComponentProps = Omit<ApolloReactComponents.MutationComponentOptions<AddStorageDeploymentMutation, AddStorageDeploymentMutationVariables>, 'mutation'>;

    export const AddStorageDeploymentComponent = (props: AddStorageDeploymentComponentProps) => (
      <ApolloReactComponents.Mutation<AddStorageDeploymentMutation, AddStorageDeploymentMutationVariables> mutation={AddStorageDeploymentDocument} {...props} />
    );
    

/**
 * __useAddStorageDeploymentMutation__
 *
 * To run a mutation, you first call `useAddStorageDeploymentMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddStorageDeploymentMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addStorageDeploymentMutation, { data, loading, error }] = useAddStorageDeploymentMutation({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useAddStorageDeploymentMutation(baseOptions?: Apollo.MutationHookOptions<AddStorageDeploymentMutation, AddStorageDeploymentMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddStorageDeploymentMutation, AddStorageDeploymentMutationVariables>(AddStorageDeploymentDocument, options);
      }
export type AddStorageDeploymentMutationHookResult = ReturnType<typeof useAddStorageDeploymentMutation>;
export type AddStorageDeploymentMutationResult = Apollo.MutationResult<AddStorageDeploymentMutation>;
export type AddStorageDeploymentMutationOptions = Apollo.BaseMutationOptions<AddStorageDeploymentMutation, AddStorageDeploymentMutationVariables>;
export const DeleteStorageDocument = gql`
    mutation deleteStorage($id: ID!) {
  deleteStorage(id: $id) {
    id
    type
    name
    status
    deployment {
      enabled
      cpu
      memory
    }
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