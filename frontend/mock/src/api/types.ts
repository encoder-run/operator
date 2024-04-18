export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
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

export type AddRepositoryInput = {
  name?: InputMaybe<Scalars['String']['input']>;
  owner?: InputMaybe<Scalars['String']['input']>;
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
  addRepository: Repository;
  addStorage: Storage;
  addStorageDeployment: Storage;
  deleteModel: Model;
  deleteRepository: Repository;
  deleteStorage: Storage;
};


export type MutationAddModelArgs = {
  input: AddModelInput;
};


export type MutationAddModelDeploymentArgs = {
  input: AddModelDeploymentInput;
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


export type MutationDeleteRepositoryArgs = {
  id: Scalars['ID']['input'];
};


export type MutationDeleteStorageArgs = {
  id: Scalars['ID']['input'];
};

export type Query = {
  __typename?: 'Query';
  getModel: Model;
  getRepository: Repository;
  getStorage: Storage;
  models: Array<Model>;
  repositories: Array<Repository>;
  storages: Array<Storage>;
};


export type QueryGetModelArgs = {
  id: Scalars['ID']['input'];
};


export type QueryGetRepositoryArgs = {
  id: Scalars['ID']['input'];
};


export type QueryGetStorageArgs = {
  id: Scalars['ID']['input'];
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
