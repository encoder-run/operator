import { ApolloServer } from '@apollo/server';
import { startStandaloneServer } from '@apollo/server/standalone';
import { loadFilesSync } from '@graphql-tools/load-files';
import { mergeTypeDefs } from '@graphql-tools/merge';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';
import { repositoryApi } from './resolvers/repository/index.js';
import { modelApi } from './resolvers/model/index.js';
import { storageApi } from './resolvers/storage/index.js';
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
        models: (parent, args, context, info) => {
            return modelApi.getModels();
        },
        storages: (parent, args, context, info) => {
            return storageApi.getStorages();
        }
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
        deleteModel: (parent, args, context, info) => {
            return modelApi.deleteModel(args.id);
        },
        addStorage: (parent, args, context, info) => {
            return storageApi.addStorage(args.input);
        },
        deleteStorage: (parent, args, context, info) => {
            return storageApi.deleteStorage(args.id);
        }
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