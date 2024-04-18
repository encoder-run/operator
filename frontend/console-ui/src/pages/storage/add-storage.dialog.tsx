import React, { useEffect, useState } from 'react';
import { Dialog, DialogTitle, DialogContent, DialogActions, TextField, Button, Select, MenuItem, InputLabel, FormControl } from '@mui/material';
import { StorageType, useAddStorageMutation } from '../../api/types'; // Adjust the import according to your API file structure

interface AddStorageDialogProps {
    open: boolean;
    onClose: () => void;
    refetch: () => void;
    onSuccess(id: String): void;
}

const AddStorageDialog = ({ open, onClose, onSuccess, refetch }: AddStorageDialogProps) => {
    const [type, setType] = useState<StorageType | null>();
    const [name, setName] = useState('');
    const [addStorage, { data, loading, error }] = useAddStorageMutation();

    const handleSubmit = () => {
        if (!type || !name) {
            return;
        }
        const input = {
            type,
            name,
        };
        addStorage({
            variables: {
                input: input,
            },
        }).then((resp) => {
            refetch(); // Refetch the data
            onClose(); // Close the dialog
            if (resp.data?.addStorage?.id) {
                onSuccess(resp.data.addStorage.id);
            }
        });
    };

    useEffect(() => {
        if (!open) {
            setType(null);
            setName('');
        }
    }, [open]);

    return (
        <Dialog open={open} onClose={onClose} sx={{ width: "100hv" }}>
            <DialogTitle>Add New Storage</DialogTitle>
            <DialogContent sx={{ minWidth: "600px" }}>
                <FormControl fullWidth margin="dense">
                    <InputLabel id="type-label">Type</InputLabel>
                    <Select
                        labelId="type-label"
                        value={type}
                        label="Type"
                        onChange={(e) => setType(e.target.value as StorageType)}
                    >
                        <MenuItem value="REDIS">Redis</MenuItem>
                        <MenuItem value="POSTGRES">Postgres</MenuItem>
                        <MenuItem value="ELASTICSEARCH">ElasticSearch</MenuItem>
                    </Select>
                </FormControl>
                <TextField
                    margin="dense"
                    label="Name"
                    type="text"
                    fullWidth
                    variant="outlined"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                />
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Cancel</Button>
                <Button onClick={handleSubmit} disabled={loading}>Add</Button>
            </DialogActions>
        </Dialog>
    );
};

export default AddStorageDialog;
