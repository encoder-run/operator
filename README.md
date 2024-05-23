# encoder-run: Kubernetes Operator for Automatic Model Encoding

## Overview
`encoder-run` is a Kubernetes operator designed to automate the deployment and management of encoding models for generating and synchronizing source code embeddings. This project aims to ensure that vectors are up-to-date with the latest changes in source code repositories, facilitating improved search, similarity checks, RAG, and code analysis tasks.

### Features
- **Automatic Model Deployment**: Dynamically deploys encoding models in a Kubernetes environment as needed.
- **Continuous Synchronization**: Monitors source code repositories and triggers the encoding process for updated code, ensuring that embeddings remain current.
- **Scalable and Efficient**: Optimized for performance and scalability, handling multiple repositories and large codebases with ease.
- **Customizable Encoding Models**: Supports a variety of encoding models, allowing users to choose the best fit for their specific requirements.

## Getting Started

### Notes
- Ensure you have plenty of available memory to run the kubernetes setup. Check Docker Desktop resources for memory and make sure its 24GB+. Ensure disk is set reasonably as well due to the nature of large model weights and large docker images.
- Dimensions for embeddings is currently hard-coded to 768 so ensure the transformer you're deploying matches.
- Local embedding time is much longer than using a GPU based setup. Sometimes the models first request/response batch is slow but subsequent requests will be faster and more consistent. Feel free to run this setup on a machine with GPU access to speed up inference time.

### Prerequisites
Before setting up `encoder-run`, ensure the following prerequisites are met:
- **Golang**: Golang version 1.21.*+ is required.
- **Docker**: Docker must be installed on your system as it is required to containerize the application components and dependencies. Visit [Docker's official site](https://docs.docker.com/get-docker/) for installation instructions.
- **Kubernetes Kind**: `kind` is used to create a local Kubernetes cluster. If it is not already installed, you can run `make kind-install` which will check for its presence and install it if it's missing.
- **Node**: `node` is used only if you are developing locally. Install v18.17.0.

### Deployment
Use the following command to deploy the application:
```bash
make deploy
```
This command will set up the necessary Kubernetes configurations and launch all components. Once deployed, it will expose:
- **Console-UI** on `localhost:32081`
- **Gateway** on `localhost:32080`

## Roadmap
Here's a short list of planned features and enhancements for `encoder-run`:

- [X] Redis deployments
- [X] Postgres external db
- [ ] Postgres deployments
- [ ] ElasticSearch deployments/external
- [ ] Custom chunking strategies based on AST, etc

Feel free to contribute or suggest new features by opening an issue or submitting a pull request on our GitHub repository.

## Contributing

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
This setup will spin up the UI on localhost:3000 and the mock gateway on localhost:4000, simulating the full environment for development or testing.

