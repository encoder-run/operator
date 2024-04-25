import React, { useEffect, useState } from 'react';
import { Dialog, DialogTitle, DialogContent, DialogActions, TextField, Button, Select, MenuItem, InputLabel, FormControl } from '@mui/material';
import { AddRepositoryInput, RepositoryType, useAddRepositoryMutation } from '../../api/types';

interface AddRepositoryDialogProps {
    open: boolean;
    onClose: () => void;
    onSuccess(id: String): void;
    refetch: () => void;
}

const AddRepositoryDialog = ({ open, onClose, onSuccess, refetch }: AddRepositoryDialogProps) => {
    const [type, setType] = useState<RepositoryType>(RepositoryType.Github);
    const [owner, setOwner] = useState('');
    const [name, setName] = useState('');
    const [token, setToken] = useState('');
    const [branch, setBranch] = useState('main');
    const [addRepository, { data, loading, error }] = useAddRepositoryMutation();

    const handleSubmit = () => {
        // Prepare the input object
        const input: AddRepositoryInput = {
            type: type as RepositoryType,
            owner: owner,
            name: name,
            token: token,
            branch: branch,
        };
        addRepository({
            variables: {
                input: input,
            },
        }).then((resp) => {
            refetch(); // Refetch the data
            onClose(); // Close the dialog
            if (resp.data?.addRepository?.id) {
                onSuccess(resp.data.addRepository.id);
            }
        });
    };

    // Clear the form fields
    useEffect(() => {
        if (!open) {
            setType(RepositoryType.Github);
            setOwner('');
            setName('');
            setBranch('main');
            setToken('');
        }
    }, [open]);

    return (
        <Dialog open={open} onClose={onClose} sx={{width: "100hv"}}>
            <DialogTitle>Add New Repository</DialogTitle>
            <DialogContent sx={{minWidth: "600px"}}>
                <FormControl fullWidth margin="dense">
                    <InputLabel id="type-label">Type</InputLabel>
                    <Select
                        labelId="type-label"
                        value={type}
                        label="Type"
                        onChange={(e) => setType(e.target.value as RepositoryType)}
                    >
                        <MenuItem value="GITHUB">GitHub</MenuItem>
                        <MenuItem value="GITLAB">GitLab</MenuItem>
                        <MenuItem value="BITBUCKET">Bitbucket</MenuItem>
                    </Select>
                </FormControl>
                <TextField
                    margin="dense"
                    label="Owner"
                    type="text"
                    fullWidth
                    variant="outlined"
                    value={owner}
                    onChange={(e) => setOwner(e.target.value)}
                />
                <TextField
                    margin="dense"
                    label="Repository Name"
                    type="text"
                    fullWidth
                    variant="outlined"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                />
                <TextField
                    margin="dense"
                    label="Branch"
                    type="text"
                    fullWidth
                    variant="outlined"
                    value={branch}
                    onChange={(e) => setBranch(e.target.value)}
                />
                <TextField
                    margin="dense"
                    label="Token"
                    type="text"
                    fullWidth
                    variant="outlined"
                    value={token}
                    onChange={(e) => setToken(e.target.value)}
                />
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Cancel</Button>
                <Button onClick={handleSubmit}>Add</Button>
            </DialogActions>
        </Dialog>
    );
};

export default AddRepositoryDialog;
