import { Storage, StorageStatus, StorageType } from "../../api/types.js";


let storages: Storage[] = [
    {
        id: '1',
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