
import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  overwrite: true,
  schema: ["../../pkg/graph/*.graphqls"],
  documents: "src/**/*.graphql",
  // generates: {
  //   "./src/__generated__/": {
  //     preset: "client",
  //     plugins: []
  //   }
  // },
  generates: {
    'src/api/types.ts': {
      plugins: [
        'typescript',
        // "typescript-operations",
        // "typescript-react-apollo"
      ],
      config: {
        withHooks: true,
        withComponent: true,
      }
    }
  },
  ignoreNoDocuments: true,
};

export default config;