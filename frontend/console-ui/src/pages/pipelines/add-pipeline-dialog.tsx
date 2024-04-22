import React, { useEffect, useState } from 'react';
import { Dialog, DialogTitle, DialogContent, DialogActions, FormControl, Button, Select, MenuItem, InputLabel, Typography, TextField } from '@mui/material';
import { useAddPipelineMutation, useRepositoriesQuery, useModelsQuery, useStoragesQuery, PipelineType } from '../../api/types';

interface AddPipelineDialogProps {
    open: boolean;
    onClose: () => void;
    onSuccess(id: string): void;
    refetch: () => void;
}

const AddPipelineDialog = ({ open, onClose, onSuccess, refetch }: AddPipelineDialogProps) => {
    const [name, setName] = useState('');
    const [type, setType] = useState<PipelineType>(PipelineType.RepositoryEmbeddings); // fixed type
    const [repositoryId, setRepositoryId] = useState('');
    const [modelId, setModelId] = useState('');
    const [storageId, setStorageId] = useState('');

    const { data: repositoriesData } = useRepositoriesQuery();
    const { data: modelsData } = useModelsQuery();
    const { data: storagesData } = useStoragesQuery();

    const [addPipeline, { loading, error }] = useAddPipelineMutation();

    const handleSubmit = () => {
        if (!loading) {
            const input = {
                name,
                type,
                repositoryEmbeddings: {
                    repositoryID: repositoryId,
                    modelID: modelId,
                    storageID: storageId,
                },
            };
            addPipeline({
                variables: { input },
            }).then((resp) => {
                if (resp.data?.addPipeline?.id) {
                    onSuccess(resp.data.addPipeline.id);
                }
                refetch();
                onClose();
            });
        }
    };

    useEffect(() => {
        if (!open) {
            setName('');
            setRepositoryId('');
            setModelId('');
            setStorageId('');
        }
    }, [open]);

    return (
        <Dialog open={open} onClose={onClose} fullWidth>
            <DialogTitle>Add New Pipeline</DialogTitle>
            <DialogContent>
                <TextField
                    fullWidth
                    margin="dense"
                    label="Name"
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                />
                <FormControl fullWidth margin="dense">
                    <InputLabel id="repository-label">Repository</InputLabel>
                    <Select
                        labelId="repository-label"
                        value={repositoryId}
                        label="Repository"
                        onChange={(e) => setRepositoryId(e.target.value)}
                    >
                        {repositoriesData?.repositories.length ? (
                            repositoriesData.repositories.map((repo) => (
                                <MenuItem key={repo.id} value={repo.id}>{repo.displayName}</MenuItem>
                            ))
                        ) : (
                            <MenuItem value="" disabled>Create Repository</MenuItem>
                        )}
                    </Select>
                </FormControl>
                <FormControl fullWidth margin="dense">
                    <InputLabel id="model-label">Model</InputLabel>
                    <Select
                        labelId="model-label"
                        value={modelId}
                        label="Model"
                        onChange={(e) => setModelId(e.target.value)}
                    >
                        {modelsData?.models.length ? (
                            modelsData.models.map((model) => (
                                <MenuItem key={model.id} value={model.id}>{model.displayName}</MenuItem>
                            ))
                        ) : (
                            <MenuItem value="" disabled>Create Model</MenuItem>
                        )}
                    </Select>
                </FormControl>
                <FormControl fullWidth margin="dense">
                    <InputLabel id="storage-label">Storage Reference</InputLabel>
                    <Select
                        labelId="storage-label"
                        value={storageId}
                        label="Storage Reference"
                        onChange={(e) => setStorageId(e.target.value)}
                    >
                        {storagesData?.storages.length ? (
                            storagesData.storages.map((storage) => (
                                <MenuItem key={storage.id} value={storage.id}>{storage.name}</MenuItem>
                            ))
                        ) : (
                            <MenuItem value="" disabled>Create Storage</MenuItem>
                        )}
                    </Select>
                </FormControl>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose} color="error">Cancel</Button>
                <Button onClick={handleSubmit} color="primary" disabled={loading}>Add</Button>
            </DialogActions>
        </Dialog>
    );
};

export default AddPipelineDialog;
