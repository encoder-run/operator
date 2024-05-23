# encoder-run: Kubernetes Operator for Automatic Model Encoding
[![license badge](https://img.shields.io/badge/license-Apache--2.0-green.svg)](LICENSE.md)

## Overview
`encoder-run` is a Kubernetes operator designed to automate the deployment and management of encoding models for generating and synchronizing source code embeddings. This project aims to ensure that vectors are up-to-date with the latest changes in source code repositories, facilitating improved search, similarity checks, RAG, and code analysis tasks.

### Features
- **Automatic Model Deployment**: Dynamically deploys encoding models in a Kubernetes environment as needed.
- **Continuous Synchronization**: Monitors source code repositories and triggers the encoding process for updated code, ensuring that embeddings remain current.
- **Scalable and Efficient**: Optimized for performance and scalability, handling multiple repositories and large codebases with ease.
- **Customizable Encoding Models**: Supports a variety of encoding models, allowing users to choose the best fit for their specific requirements.

## Getting Started

### Prerequisites
Before setting up `encoder-run`, ensure the following prerequisites are met:
- **Golang**: Goland version [1.21.*+](https://go.dev/dl/) is required.
- **Docker**: Docker must be installed on your system as it is required to containerize the application components and dependencies. Visit [Docker's official site](https://docs.docker.com/get-docker/) for installation instructions.
- **Kubernetes Kind**: Kind is used to create a local Kubernetes cluster. If it is not already installed, you can run `make kind-install` which will check for its presence and install it if it's missing.
- **Node**: Node is used only if you are developing locally. Please install [v18.17.0](https://nodejs.org/en/download/package-manager).

### Deployment
After satisfying the prerequisites, use the following command to deploy the application:
```
make deploy
```
This command will set up the necessary Kubernetes configurations and launch all components. Once deployed, it will expose:
- **Console-UI** on `localhost:32081`
- **Gateway** on `localhost:32080`

### Running the Frontend with Mock Data
To run the frontend interface with mock data, follow these steps:
1. Navigate to the frontend console-UI directory:
   ```bash
   cd frontend/console-ui
   npm install && npm run generate
   cd ../mock
   npm install
   cd ../console-ui
   npm run local
   ```
This setup will spin up the UI on `localhost:3000` and the mock gateway on `localhost:4000`, simulating the full environment for development or testing.

## Community and Feedback
Encoder-run is an open-source project and we encourage and welcome contributions. If you wish to contribute, be sure to review our [contribution guidelines](CONTRIBUTING.md) and [code of conduct](CODE_OF_CONDUCT.md).

For problems with the installation and setup, discussions about how to best to use encoder-run please use our Github [issues](https://github.com/encoder-run/operator/issues) and make sure to include as much detail as possible.
