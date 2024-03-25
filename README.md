# encoder-run: Kubernetes Operator for Automatic Model Encoding

## Overview
`encoder-run` is a Kubernetes operator designed to automate the deployment and management of encoding models for generating and synchronizing source code embeddings. This project aims to ensure that vectors are up-to-date with the latest changes in source code repositories, facilitating improved search, similarity checks, RAG, and code analysis tasks.

### Features
- **Automatic Model Deployment**: Dynamically deploys encoding models in a Kubernetes environment as needed.
- **Continuous Synchronization**: Monitors source code repositories and triggers the encoding process for updated code, ensuring that embeddings remain current.
- **Scalable and Efficient**: Optimized for performance and scalability, handling multiple repositories and large codebases with ease.
- **Customizable Encoding Models**: Supports a variety of encoding models, allowing users to choose the best fit for their specific requirements.

## Getting Started

### Prerequisites