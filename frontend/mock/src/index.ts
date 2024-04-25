import { ApolloServer } from '@apollo/server';
import { startStandaloneServer } from '@apollo/server/standalone';
import { loadFilesSync } from '@graphql-tools/load-files';
import { mergeTypeDefs } from '@graphql-tools/merge';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';
import { repositoryApi } from './resolvers/repository/index.js';
import { modelApi } from './resolvers/model/index.js';
import { storageApi } from './resolvers/storage/index.js';
import { pipelineApi } from './resolvers/pipeline/index.js';
import { searchApi } from './resolvers/search/index.js';
// Convert the URL to a directory name
const __dirname = dirname(fileURLToPath(import.meta.url));

const typesArray = loadFilesSync(join(__dirname, '../../../../pkg/graph'), {
    extensions: ['graphqls']
});
const typeDefs = mergeTypeDefs(typesArray);

// Resolvers define how to fetch the types defined in your schema.
// This resolver retrieves books from the "books" array above.
const resolvers = {
    Query: {
        repositories: (parent, args, context, info) => {
            return repositoryApi.getRepositories();
        },
        getRepository: (parent, args, context, info) => { 
            return repositoryApi.getRepository(args.id);
        },
        models: (parent, args, context, info) => {
            return modelApi.getModels();
        },
        getModel: (parent, args, context, info) => {
            return modelApi.getModel(args.id);
        },
        storages: (parent, args, context, info) => {
            return storageApi.getStorages();
        },
        getStorage: (parent, args, context, info) => {
            return storageApi.getStorage(args.id);
        },
        pipelines: (parent, args, context, info) => {
            return pipelineApi.getPipelines();
        },
        getPipeline: (parent, args, context, info) => {
            return pipelineApi.getPipeline(args.id);
        },
        getPipelineExecutions: (parent, args, context, info) => {
            return pipelineApi.getPipelineExecutions(args.id);
        },
        semanticSearch: (parent, args, context, info) => {
            return searchApi.search(args.query);
        },
    },
    Mutation: {
        addRepository: (parent, args, context, info) => {
            return repositoryApi.addRepository(args.input);
        },
        deleteRepository: (parent, args, context, info) => {
            return repositoryApi.deleteRepository(args.id);
        },
        addModel: (parent, args, context, info) => {
            return modelApi.addModel(args.input);
        },
        addModelDeployment: (parent, args, context, info) => {
            return modelApi.addDeployment(args.input);
        },
        deleteModel: (parent, args, context, info) => {
            return modelApi.deleteModel(args.id);
        },
        addStorage: (parent, args, context, info) => {
            return storageApi.addStorage(args.input);
        },
        addStorageDeployment: (parent, args, context, info) => {
            return storageApi.addDeployment(args.input);
        },
        deleteStorage: (parent, args, context, info) => {
            return storageApi.deleteStorage(args.id);
        },
        addPipeline: (parent, args, context, info) => {
            return pipelineApi.addPipeline(args.input);
        },
        addPipelineDeployment: (parent, args, context, info) => {
            return pipelineApi.addDeployment(args.input);
        },
    },
  };

  // The ApolloServer constructor requires two parameters: your schema
// definition and your set of resolvers.
const server = new ApolloServer({
    typeDefs,
    resolvers,
  });
  
  // Passing an ApolloServer instance to the `startStandaloneServer` function:
  //  1. creates an Express app
  //  2. installs your ApolloServer instance as middleware
  //  3. prepares your app to handle incoming requests
  const { url } = await startStandaloneServer(server, {
    listen: { port: 4000 },
  });
  
  console.log(`🚀  Server ready at: ${url}`);