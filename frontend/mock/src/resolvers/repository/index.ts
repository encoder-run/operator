import { AddRepositoryInput, Repository, RepositoryType } from "../../api/types.js";

let repositories: Repository[] = [
    {
        id: '1',
        name: 'frontend',
        owner: 'facebook',
        type: RepositoryType.Github,
        displayName: 'facebook/frontend',
        url: 'github.com/facebook/frontend'
    },
    {
        id: '2',
        name: 'console-ui',
        owner: 'facebook',
        type: RepositoryType.Github,
        displayName: 'facebook/console-ui',
        url: 'github.com/facebook/console-ui'
    },
    {
        id: '3',
        name: 'mock',
        owner: 'facebook',
        type: RepositoryType.Github,
        displayName: 'facebook/mock',
        url: 'github.com/facebook/mock'
    },
    {
        id: '4',
        name: 'console-api',
        owner: 'facebook',
        type: RepositoryType.Github,
        displayName: 'facebook/console-api',
        url: 'github.com/facebook/console-api'
    },
    {
        id: '5',
        name: 'console-server',
        owner: 'facebook',
        type: RepositoryType.Gitlab,
        displayName: 'facebook/console-server',
        url: 'gitlab.com/facebook/console-server'
    },
    {
        id: '6',
        name: 'mock',
        owner: 'facebook',
        type: RepositoryType.Bitbucket,
        displayName: 'facebook/mock',
        url: 'bitbucket.com/facebook/mock'
    },
];


class RepositoryApi {
    getRepositories() {
        return repositories;
    }

    addRepository(input: AddRepositoryInput) {
        const newRepo: Repository = {
            id: String(repositories.length + 1),
            name: input.name,
            owner: input.owner,
            type: input.type,
            displayName: `${input.owner}/${input.name}`,
            url: `${input.type}.com/${input.owner}/${input.name}`
        };
        repositories.push(newRepo);
        return newRepo;
    }

    deleteRepository(id: string) {
        const index = repositories.findIndex(repo => repo.id === id);
        if (index !== -1) {
            const deleted = repositories.splice(index, 1);
            return deleted[0];
        }
    }
}

export const repositoryApi = new RepositoryApi();
