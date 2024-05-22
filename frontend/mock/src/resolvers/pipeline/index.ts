import { Pipeline, PipelineExecution, PipelineExecutionStatus, PipelineStatus, PipelineType } from "../../api/types.js";


let pipelines: Pipeline[] = [
    {
        id: '1',
        name: 'pipeline-1',
        enabled: true,
        type: PipelineType.RepositoryEmbeddings,
        repositoryEmbeddings: {
            repositoryID: '1',
            modelID: '1',
            storageID: '1',
        },
        status: PipelineStatus.Ready,
    },
    {
        id: '2',
        name: 'pipeline-2',
        enabled: false,
        type: PipelineType.RepositoryEmbeddings,
        repositoryEmbeddings: {
            repositoryID: '2',
            modelID: '2',
            storageID: '2',
        },
        status: PipelineStatus.Ready,
    },
    {
        id: '3',
        name: 'pipeline-3',
        enabled: true,
        type: PipelineType.RepositoryEmbeddings,
        repositoryEmbeddings: {
            repositoryID: '3',
            modelID: '3',
            storageID: '3',
        },
        status: PipelineStatus.Ready,
    },
    {
        id: '4',
        name: 'pipeline-4',
        enabled: false,
        type: PipelineType.RepositoryEmbeddings,
        repositoryEmbeddings: {
            repositoryID: '4',
            modelID: '4',
            storageID: '4',
        },
        status: PipelineStatus.Error,
    }
];

let pipelineExecutions: PipelineExecution[] = [
    {
        id: '1',
        status: PipelineExecutionStatus.Succeeded,
    }
];

class PipelineApi {
    getPipelines() {
        return pipelines;
    }

    getPipeline(id: string) {
        const pipeline = pipelines.find(pipeline => pipeline.id === id);
        if (!pipeline) {
            throw new Error('Pipeline not found');
        }
        return pipeline;
    }

    getPipelineExecutions(id: string) {
        return pipelineExecutions;
    }

    addPipeline(input: any) {
        const newPipeline: Pipeline = {
            id: String(pipelines.length + 1),
            name: input.name,
            enabled: false,
            type: input.type,
            status: PipelineStatus.Ready,
            repositoryEmbeddings: {
                repositoryID: input.repositoryEmbeddings.repositoryID,
                modelID: input.repositoryEmbeddings.modelID,
                storageID: input.repositoryEmbeddings.storageID,
            },
        };
        pipelines.push(newPipeline);
        return newPipeline;
    }

    addDeployment(input: any) {
        // Find the pipeline
        const pipeline = pipelines.find(pipeline => pipeline.id === input.id);
        if (!pipeline) {
            throw new Error('Pipeline not found');
        }
        // Update the pipeline deployment spec
        pipeline.enabled = true;
        pipeline.status = PipelineStatus.Ready;


        return pipeline;
    }

    trigger(id: any) {
        // Find the pipeline
        const pipeline = pipelines.find(pipeline => pipeline.id === id);
        if (!pipeline) {
            throw new Error('Pipeline not found');
        }
        // Add a new pipeline execution
        const newExecution: PipelineExecution = {
            id: String(pipelineExecutions.length + 1),
            status: PipelineExecutionStatus.Pending,
        };

        pipelineExecutions.push(newExecution);

        // Return the execution
        return newExecution;
    }
}

export const pipelineApi = new PipelineApi();