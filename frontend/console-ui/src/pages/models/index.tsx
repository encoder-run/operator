import React, { useEffect, useState } from 'react';
import { Typography, Box, TextField, Button, IconButton } from '@mui/material';
import { DataGrid, GridColDef, GridRowSelectionModel } from '@mui/x-data-grid';
import LaunchIcon from '@mui/icons-material/Launch';
import { useNavigate } from 'react-router-dom'; // Import useNavigate from react-router-dom
import ConfirmDelete from '../../components/confirm-delete-dialog';
import AddModelDialog from './add-model-dialog';
import { Model, useDeleteModelMutation, useModelsQuery } from '../../api/types';

export default function ModelsPage() {
    const columns: GridColDef[] = [
        { field: 'id', headerName: 'ID', flex: 1, minWidth: 100 },
        { field: 'displayName', headerName: 'Name', flex: 1, minWidth: 150 },
        { field: 'type', headerName: 'Type', flex: 1, minWidth: 130 },
        { field: 'status', headerName: 'Status', flex: 1, minWidth: 150 },
        {
            field: 'actions',
            headerName: 'Actions',
            type: 'actions',
            width: 100,
            getActions: (params) => [
                <IconButton
                    color="primary"
                    aria-label="go to details"
                    size="small"
                    onClick={() => navigate(`/models/${params.id}`)} // Adjust the path as needed
                >
                    <LaunchIcon />
                </IconButton>,
            ],
        },
    ];

    const [searchQuery, setSearchQuery] = useState('');
    const { data, loading, error, refetch } = useModelsQuery(
        { fetchPolicy: 'network-only' }
    );
    const [rows, setRows] = useState(data?.models || []);
    const [deleteModel, { data: deleteData, loading: deleteLoading, error: deleteError }] = useDeleteModelMutation();
    const [selectedModels, setSelectedModels] = useState<GridRowSelectionModel>([]);
    const navigate = useNavigate();
    const [openConfirmDelete, setOpenConfirmDelete] = useState(false);
    const [openAddDialog, setOpenAddDialog] = useState(false);


    const handleDelete = () => {
        // Show the confirm dialog
        setOpenConfirmDelete(true);
    };

    const confirmDelete = () => {
        // Remove selected rows from the rows state
        selectedModels.forEach(async (id) => {
            deleteModel({ variables: { id: String(id) } }).then(() => {
                // TODO: Do this better
                refetch();
            });
        });
        // Close the dialog
        setOpenConfirmDelete(false);
        // Clear selection model
        setSelectedModels([]);
    };

    const handleSearch = (event: React.ChangeEvent<HTMLInputElement>) => {
        const query = event.target.value;
        setSearchQuery(query);
        const filteredRows = rows.filter(row =>
            row.displayName.toLowerCase().includes(query.toLowerCase()) ||
            row.type.toLowerCase().includes(query.toLowerCase())
        );
        setRows(filteredRows);
    };

    const onSuccessRedirect = (id: string) => {
        navigate(`/models/${id}`);
    }

    useEffect(() => {
        setRows(data?.models || []);
    }, [data]);

    return (
        <>
            <Box sx={{ p: 2, display: "flex", flexDirection: "column", height: "100%" }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                    <Typography variant="h4">Models Overview</Typography>
                    <Box>
                        <Button variant="contained" onClick={() => setOpenAddDialog(true)} sx={{ marginRight: 1 }}>
                            Add Model
                        </Button>
                        <Button
                            variant="outlined"
                            onClick={handleDelete}
                            disabled={selectedModels.length === 0}
                            color="error"
                        >
                            Delete Model
                        </Button>
                    </Box>
                </Box>
                <TextField
                    fullWidth
                    label="Search Models"
                    variant="outlined"
                    value={searchQuery}
                    onChange={handleSearch}
                    sx={{ mb: 2 }}
                />
                <Box sx={{ height: '100%', width: '100%' }}>
                    <DataGrid
                        rows={rows}
                        columns={columns}
                        initialState={{
                            pagination: {
                                paginationModel: {
                                    pageSize: 10,
                                },
                            },
                        }}
                        pageSizeOptions={[10]}
                        checkboxSelection
                        disableRowSelectionOnClick
                        disableColumnSelector
                        rowSelectionModel={selectedModels}
                        onRowSelectionModelChange={setSelectedModels}
                    />
                </Box>
                <ConfirmDelete
                    open={openConfirmDelete}
                    onClose={() => setOpenConfirmDelete(false)}
                    onConfirm={confirmDelete}
                    type={"model(s)"}
                />
            </Box>
            <Box>
                <AddModelDialog
                    open={openAddDialog}
                    onClose={() => setOpenAddDialog(false)}
                    onSuccess={onSuccessRedirect}
                    refetch={refetch}
                />
            </Box>
        </>
    );
}
