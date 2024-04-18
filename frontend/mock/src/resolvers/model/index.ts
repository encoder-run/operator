import { AddModelDeploymentInput, Model, ModelStatus, ModelType } from "../../api/types.js";

let models: Model[] = [
    {
        id: '1',
        displayName: 'jinaai/jina-embeddings-v2-base-code',
        type: ModelType.Huggingface,
        status: ModelStatus.NotDeployed,
    },
    {
        id: '2',
        displayName: 'distilbert/distilbert-base-uncased',
        type: ModelType.Huggingface,
        status: ModelStatus.Ready,
        deployment: {
            enabled: true,
            cpu: '1',
            memory: '1G',
        },
    },
    {
        id: '3',
        displayName: 'openai/gpt-3',
        type: ModelType.External,
        status: ModelStatus.Deploying,
        deployment: {
            enabled: true,
            cpu: '1',
            memory: '1G',
        },
    },
    {
        id: '4',
        displayName: 'jinaai/jina-embeddings-v3-base-code',
        type: ModelType.Huggingface,
        status: ModelStatus.Deploying,
        deployment: {
            enabled: true,
            cpu: '1',
            memory: '1G',
        },
    }
];

class ModelApi {
    getModels() {
        return models;
    }

    getModel(id: string) {
        const model = models.find(model => model.id === id);
        if (!model) {
            throw new Error('Model not found');
        }
        return model;
    }

    addModel(input: any) {
        const newModel: Model = {
            id: String(models.length + 1),
            displayName: "Org+Repo",
            type: input.type,
            status: ModelStatus.NotDeployed,
        };
        models.push(newModel);
        return newModel;
    }

    addDeployment(input: AddModelDeploymentInput) {
        // Find the model
        const model = models.find(model => model.id === input.id);
        if (!model) {
            throw new Error('Model not found');
        }
        // Update the model deployment spec
        model.deployment = {
            enabled: true,
            cpu: input.cpu,
            memory: input.memory,
        };

        // We need to asynchronously update the model status so that it updates
        // after this call completes after 5 seconds.
        setTimeout(() => {
            model.status = ModelStatus.Ready;
        }, 5000);

        return model;
    }

    deleteModel(id: string) {
        const index = models.findIndex(model => model.id === id);
        if (index === -1) {
            throw new Error('Model not found');
        }
        const deletedModel = models.splice(index, 1)[0];
        return deletedModel;
    }
}

export const modelApi = new ModelApi();