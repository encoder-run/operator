import React, { useState, useEffect } from 'react';
import { Dialog, DialogTitle, DialogContent, DialogActions, TextField, Button, Select, MenuItem, InputLabel, FormControl } from '@mui/material';
import { AddModelInput, ModelType, useAddModelMutation } from '../../api/types';
import { parse } from 'path';

interface AddModelDialogProps {
    open: boolean;
    onClose: () => void;
    onSuccess(id: String): void;
    refetch: () => void;
}

const AddModelDialog = ({ open, onClose, onSuccess, refetch }: AddModelDialogProps) => {
    const [modelType, setModelType] = useState<ModelType | ''>('');
    const [organization, setOrganization] = useState('');
    const [repoName, setRepoName] = useState('');
    const [maxSequenceLength, setMaxSequenceLength] = useState('');
    const [selectedModel, setSelectedModel] = useState('');
    const [addModel, { data, loading, error }] = useAddModelMutation();

    const handleSubmit = () => {
        if (modelType === '') {
            return;
        }
        const input: AddModelInput = {
            type: modelType,
        };

        if (modelType === ModelType.Huggingface) {
            input.huggingFace = { 
                organization: organization, 
                name: repoName,
                maxSequenceLength: parseInt(maxSequenceLength),
            };
        }

        addModel({
            variables: {
                input: input,
            },
        }).then((resp) => {
            refetch();
            onClose();
            if (resp.data?.addModel?.id) {
                onSuccess(resp.data.addModel.id);
            }
        });
    };

    useEffect(() => {
        if (!open) {
            setModelType('');
            setOrganization('');
            setRepoName('');
            setMaxSequenceLength('');
            setSelectedModel('');
        }
    }, [open]);

    const renderFormFields = () => {
        switch (modelType) {
            case 'HUGGINGFACE':
                return (
                    <>
                        <TextField
                            margin="dense"
                            label="Organization"
                            type="text"
                            fullWidth
                            variant="outlined"
                            value={organization}
                            onChange={(e) => setOrganization(e.target.value)}
                        />
                        <TextField
                            margin="dense"
                            label="Repository Name"
                            type="text"
                            fullWidth
                            variant="outlined"
                            value={repoName}
                            onChange={(e) => setRepoName(e.target.value)}
                        />
                        <TextField
                            margin="dense"
                            label="Max Sequence Length"
                            type="number"
                            fullWidth
                            variant="outlined"
                            value={maxSequenceLength}
                            onChange={(e) => setMaxSequenceLength(e.target.value)}
                        />
                    </>
                );
            case 'OPENAI':
                return (
                    <FormControl fullWidth margin="dense">
                        <InputLabel id="model-select-label">Model</InputLabel>
                        <Select
                            labelId="model-select-label"
                            value={selectedModel}
                            label="Model"
                            onChange={(e) => setSelectedModel(e.target.value)}
                        >
                            <MenuItem value="text-similarity-ada-001">text-similarity-ada-001</MenuItem>
                            <MenuItem value="code-davinci-002">code-davinci-002</MenuItem>
                        </Select>
                    </FormControl>
                );
            default:
                return null;
        }
    };

    return (
        <Dialog open={open} onClose={onClose} sx={{width: "100hv"}}>
            <DialogTitle>Add New Model</DialogTitle>
            <DialogContent sx={{minWidth: "600px"}}>
                <FormControl fullWidth margin="dense">
                    <InputLabel id="type-label">Model Type</InputLabel>
                    <Select
                        labelId="type-label"
                        value={modelType}
                        label="Model Type"
                        onChange={(e) => setModelType(e.target.value  as ModelType)}
                    >
                        <MenuItem value="HUGGINGFACE">Hugging Face</MenuItem>
                        <MenuItem value="OPENAI">OpenAI</MenuItem>
                    </Select>
                </FormControl>
                {renderFormFields()}
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Cancel</Button>
                <Button onClick={handleSubmit}>Add</Button>
            </DialogActions>
        </Dialog>
    );
};

export default AddModelDialog;
