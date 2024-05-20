import React, { useEffect, useState } from 'react';
import { Dialog, DialogTitle, DialogContent, DialogActions, TextField, Button, Select, MenuItem, InputLabel, FormControl, FormGroup, FormControlLabel, Switch } from '@mui/material';
import { PostgresInput, StorageType, useAddStorageMutation } from '../../api/types'; // Adjust the import according to your API file structure

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
    const [postgresConfig, setPostgresConfig] = useState<PostgresInput>({
        external: true,
        host: '',
        port: 5432,
        username: '',
        password: '',
        database: '',
        SSLMode: 'disable',
        timezone: 'America/Los_Angeles'
    });

    const handlePostgresConfigChange = (prop: any) => (event: any) => {
        setPostgresConfig({
            ...postgresConfig,
            [prop]: event.target.value
        });
    };

    const handleSubmit = () => {
        if (!type || !name) {
            return;
        }
        const input = {
            type,
            name,
            ...(type === 'POSTGRES' && { postgres: { ...postgresConfig } })
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
            setPostgresConfig({
                external: true,
                host: '',
                username: '',
                password: '',
                database: '',
                port: 5432,
                SSLMode: 'disable',
                timezone: 'America/Los_Angeles'
            });
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
                        value={type || ''}
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
                {type === 'POSTGRES' && (
                    <>
                        <FormGroup>
                            <FormControlLabel control={<Switch disabled defaultChecked />} label="External" />
                        </FormGroup>
                        <TextField
                            margin="dense"
                            label="Host"
                            type="text"
                            fullWidth
                            variant="outlined"
                            value={postgresConfig.host}
                            onChange={handlePostgresConfigChange('host')}
                        />
                        <TextField
                            margin="dense"
                            label="Port"
                            type="number"
                            fullWidth
                            variant="outlined"
                            value={postgresConfig.port}
                            onChange={handlePostgresConfigChange('port')}
                        />
                        <TextField
                            margin="dense"
                            label="Username"
                            type="text"
                            fullWidth
                            variant="outlined"
                            value={postgresConfig.username}
                            onChange={handlePostgresConfigChange('username')}
                        />
                        <TextField
                            margin="dense"
                            label="Password"
                            type="password"
                            fullWidth
                            variant="outlined"
                            value={postgresConfig.password}
                            onChange={handlePostgresConfigChange('password')}
                        />
                        <TextField
                            margin="dense"
                            label="Database"
                            type="text"
                            fullWidth
                            variant="outlined"
                            value={postgresConfig.database}
                            onChange={handlePostgresConfigChange('database')}
                        />
                        <TextField
                            margin="dense"
                            label="SSL Mode"
                            type="text"
                            fullWidth
                            variant="outlined"
                            value={postgresConfig.SSLMode}
                            onChange={handlePostgresConfigChange('SSLMode')}
                        />
                    </>
                )}
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Cancel</Button>
                <Button onClick={handleSubmit} disabled={loading}>Add</Button>
            </DialogActions>
        </Dialog>
    );
};

export default AddStorageDialog;
