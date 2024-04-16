
import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  overwrite: true,
  schema: ["../../pkg/graph/*.graphqls"],
  documents: "src/**/*.graphql",
  generates: {
    'src/api/types.tsx': {
      plugins: [
        'typescript',
        "typescript-operations",
        "typescript-react-apollo"
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