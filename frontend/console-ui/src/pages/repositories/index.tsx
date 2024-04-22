import React, { useEffect, useState } from 'react';
import { Typography, Box, TextField, Button, IconButton } from '@mui/material';
import { DataGrid, GridColDef, GridRowSelectionModel } from '@mui/x-data-grid';
import LaunchIcon from '@mui/icons-material/Launch';
import { useNavigate } from 'react-router-dom';
import ConfirmDelete from '../../components/confirm-delete-dialog';
import AddRepositoryDialog from './add-repository-dialog';
import { useDeleteRepositoryMutation, useRepositoriesQuery } from '../../api/types';

export default function RepositoriesPage() {
    const columns: GridColDef[] = [
        { field: 'id', headerName: 'ID', flex: 1, minWidth: 100 },
        { field: 'displayName', headerName: 'Display Name', flex: 1, minWidth: 130 },
        { field: 'type', headerName: 'Type', flex: 1, minWidth: 130 },
        { field: 'owner', headerName: 'Owner', flex: 1, minWidth: 150 },
        { field: 'name', headerName: 'Repository', flex: 1, minWidth: 150 },
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
                    onClick={() => navigate(`/repositories/${params.id}`)} // Adjust the path as needed
                >
                    <LaunchIcon />
                </IconButton>,
            ],
        },
    ];

    const [searchQuery, setSearchQuery] = useState('');
    const { data, loading, error, refetch } = useRepositoriesQuery(
        { fetchPolicy: 'network-only' }
    );
    const [repositories, setRepositories] = useState(data?.repositories || []);
    const [deleteRepository, { data: deleteData, loading: deleteLoading, error: deleteError }] = useDeleteRepositoryMutation();
    const [selectedRepositories, setSelectedRepositories] = useState<GridRowSelectionModel>([]);
    const navigate = useNavigate();
    const [openConfirmDelete, setOpenConfirmDelete] = useState(false);
    // Inside RepositoriesPage component
    const [openAddDialog, setOpenAddDialog] = useState(false);

    const handleShowDeleteDialog = () => {
        // Show the confirm dialog
        setOpenConfirmDelete(true);
    };

    const confirmDelete = () => {
        // Actual deletion logic here, after confirmation
        console.log('Delete confirmed for selected rows:', selectedRepositories);
        // Delete each of the selected repositories
        selectedRepositories.forEach(async (id) => {
            await deleteRepository({ variables: { id: String(id) } });
        });

        refetch();
        // Close the dialog
        setOpenConfirmDelete(false);
        // Clear selection model
        setSelectedRepositories([]);
    };

    const handleSearch = (event: React.ChangeEvent<HTMLInputElement>) => {
        const query = event.target.value;
        setSearchQuery(query);
        const filteredRepositories = repositories.filter(repo =>
            repo.name.toLowerCase().includes(query.toLowerCase()) ||
            repo.owner.toLowerCase().includes(query.toLowerCase()) ||
            repo.type.toLowerCase().includes(query.toLowerCase())
        );
        setRepositories(filteredRepositories);
    };

    const onSuccessRedirect = (id: string) => {
        navigate(`/repositories/${id}`);
    }

    // Watch for changes in the data
    useEffect(() => {
        if (data) {
            setRepositories(data.repositories);
        }
    }, [data]);

    return (
        <>
            <Box sx={{ p: 2, display: "flex", flexDirection: "column", height: "100%" }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                    <Typography variant="h4">Repositories Overview</Typography>
                    <Box>
                        <Button variant="contained" onClick={() => setOpenAddDialog(true)} sx={{ marginRight: 1 }}>
                            Add Repository
                        </Button>
                        <Button
                            variant="outlined"
                            onClick={handleShowDeleteDialog}
                            disabled={selectedRepositories.length === 0}
                            color="error"
                        >
                            Delete Repository
                        </Button>
                    </Box>
                </Box>
                <TextField
                    fullWidth
                    label="Search Repositories"
                    variant="outlined"
                    value={searchQuery}
                    onChange={handleSearch}
                    sx={{ mb: 2 }}
                />
                <Box sx={{ height: '100%', width: '100%' }}>
                    <DataGrid
                        rows={repositories}
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
                        rowSelectionModel={selectedRepositories}
                        onRowSelectionModelChange={setSelectedRepositories}
                    />
                </Box>
                <ConfirmDelete
                    open={openConfirmDelete}
                    onClose={() => setOpenConfirmDelete(false)}
                    onConfirm={confirmDelete}
                    type={"repository(s)"}
                />
            </Box>
            <AddRepositoryDialog
                open={openAddDialog}
                onClose={() => setOpenAddDialog(false)}
                onSuccess={onSuccessRedirect}
                refetch={refetch}
            />
        </>
    );
}
