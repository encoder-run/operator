import React, { useState } from 'react';
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, Select, MenuItem, InputLabel, FormControl, Grid } from '@mui/material';
import { useAddStorageDeploymentMutation } from '../../api/types';

interface AddStorageDeploymentDialogProps {
    open: boolean;
    onClose: () => void;
    onSuccess(): void;
    refetch: () => void;
    storageId: string;
}

const AddStorageDeploymentDialog = ({ open, onClose, onSuccess, refetch, storageId }: AddStorageDeploymentDialogProps) => {
    const [cpu, setCpu] = useState('');
    const [memory, setMemory] = useState('');
    const [addDeployment, { loading, error }] = useAddStorageDeploymentMutation();

    const handleSubmit = async () => {
        try {
            await addDeployment({
                variables: {
                    input: {
                        id: storageId,
                        cpu,
                        memory
                    }
                }
            });
            refetch(); // Refetch the data
            onClose(); // Close the dialog on successful submission
            onSuccess(); // Call the onSuccess callback
        } catch (e) {
            console.error("Error submitting form: ", e);
        }
    };

    if (loading) return <p>Loading...</p>;
    if (error) return <p>Error: {error.message}</p>;

    return (
        <Dialog open={open} onClose={onClose} sx={{ width: '100hv' }}>
            <DialogTitle>Add Storage Deployment</DialogTitle>
            <DialogContent sx={{ minWidth: "600px" }}>
                <Grid container spacing={2} sx={{p: 1}}>
                    <Grid item xs={12}>
                        <FormControl fullWidth>
                            <InputLabel id="cpu-label">CPU</InputLabel>
                            <Select
                                labelId="cpu-label"
                                value={cpu}
                                onChange={e => setCpu(e.target.value)}
                                label="CPU"
                            >
                                <MenuItem value="0.5">0.5 vCPUs</MenuItem>
                                <MenuItem value="1">1 vCPU</MenuItem>
                                <MenuItem value="2">2 vCPUs</MenuItem>
                            </Select>
                        </FormControl>
                    </Grid>
                    <Grid item xs={12}>
                        <FormControl fullWidth>
                            <InputLabel id="memory-label">Memory</InputLabel>
                            <Select
                                labelId="memory-label"
                                value={memory}
                                onChange={e => setMemory(e.target.value)}
                                label="Memory"
                            >
                                <MenuItem value="2G">2 GB</MenuItem>
                                <MenuItem value="4G">4 GB</MenuItem>
                                <MenuItem value="8G">8 GB</MenuItem>
                                <MenuItem value="16G">16 GB</MenuItem>
                                <MenuItem value="32G">32 GB</MenuItem>
                            </Select>
                        </FormControl>
                    </Grid>
                </Grid>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Cancel</Button>
                <Button onClick={handleSubmit} color="primary">
                    Add
                </Button>
            </DialogActions>
        </Dialog>
    );
};

export default AddStorageDeploymentDialog;