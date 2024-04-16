import React, { useEffect, useState } from 'react';
import { Dialog, DialogTitle, DialogContent, DialogActions, TextField, Button, Select, MenuItem, InputLabel, FormControl } from '@mui/material';
import { AddRepositoryInput, RepositoryType, useAddRepositoryMutation } from '../../api/types';

interface AddRepositoryDialogProps {
    open: boolean;
    onClose: () => void;
    refetch: () => void;
}

const AddRepositoryDialog = ({ open, onClose, refetch }: AddRepositoryDialogProps) => {
    const [type, setType] = useState('');
    const [owner, setOwner] = useState('');
    const [name, setName] = useState('');
    const [addRepository, { data, loading, error }] = useAddRepositoryMutation();

    const handleSubmit = () => {
        // Prepare the input object
        const input: AddRepositoryInput = {
            type: type as RepositoryType,
            owner: owner,
            name: name,
        };
        addRepository({
            variables: {
                input: input,
            },
        }).then(() => {
            refetch(); // Refetch the data
        }).finally(() => {
            onClose(); // Close the dialog
        });
    };

    // Clear the form fields
    useEffect(() => {
        if (!open) {
            setType('');
            setOwner('');
            setName('');
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
                        onChange={(e) => setType(e.target.value)}
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
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Cancel</Button>
                <Button onClick={handleSubmit}>Add</Button>
            </DialogActions>
        </Dialog>
    );
};

export default AddRepositoryDialog;
