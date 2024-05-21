import { AddStorageDeploymentInput, Storage, StorageStatus, StorageType } from "../../api/types.js";


let storages: Storage[] = [
    {
        id: '10',
        type: StorageType.Redis,
        name: 'Redis-dev',
        status: StorageStatus.Ready,
    },
    {
        id: '2',
        type: StorageType.Postgres,
        name: 'Postgres-prod',
        status: StorageStatus.Deploying,
    },
    {
        id: '3',
        type: StorageType.Elasticsearch,
        name: 'Elasticsearch-test',
        status: StorageStatus.NotDeployed,
    },
];


class StorageApi {
    getStorages() {
        return storages;
    }

    getStorage(id: string) {
        const storage = storages.find(storage => storage.id === id);
        if (!storage) {
            throw new Error('Storage not found');
        }
        return storage;
    }

    addStorage(input: any) {
        const newStorage: Storage = {
            id: String(storages.length + 1),
            type: input.type,
            name: input.name,
            status: StorageStatus.NotDeployed,
        };
        storages.push(newStorage);
        return newStorage;
    }

    addDeployment(input: AddStorageDeploymentInput) {
        // Find the storage
        const storage = storages.find(storage => storage.id === input.id);
        if (!storage) {
            throw new Error('Storage not found');
        }
        // Update the storage deployment spec
        storage.deployment = {
            enabled: true,
            cpu: input.cpu,
            memory: input.memory,
        };

        // We need to asynchronously update the storage status so that it updates
        // after this call completes after 5 seconds.
        setTimeout(() => {
            storage.status = StorageStatus.Ready;
        }, 5000);

        return storage;
    }

    deleteStorage(id: string) {
        const index = storages.findIndex(storage => storage.id === id);
        if (index === -1) {
            throw new Error('Storage not found');
        }
        const deletedStorage = storages.splice(index, 1)[0];
        return deletedStorage;
    }
}

export const storageApi = new StorageApi();