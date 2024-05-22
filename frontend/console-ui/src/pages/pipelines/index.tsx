import React, { useEffect, useState } from 'react';
import { Typography, Box, TextField, Button, IconButton } from '@mui/material';
import { DataGrid, GridColDef, GridRowSelectionModel } from '@mui/x-data-grid';
import LaunchIcon from '@mui/icons-material/Launch';
import { useNavigate } from 'react-router-dom';
import ConfirmDelete from '../../components/confirm-delete-dialog';
import AddPipelineDialog from './add-pipeline-dialog';
import { useDeletePipelineMutation, usePipelinesQuery } from '../../api/types';

export default function PipelinesPage() {
    const columns: GridColDef[] = [
        { field: 'id', headerName: 'ID', flex: 1, minWidth: 100 },
        { field: 'name', headerName: 'Name', flex: 1, minWidth: 100 },
        { field: 'status', headerName: 'Status', flex: 1, minWidth: 130 },
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
                    onClick={() => navigate(`/pipelines/${params.id}`)}
                >
                    <LaunchIcon />
                </IconButton>,
            ],
        },
    ];

    const [searchQuery, setSearchQuery] = useState('');
    const { data, loading, error, refetch } = usePipelinesQuery(
        { fetchPolicy: 'network-only' }
    );
    const [pipelines, setPipelines] = useState(data?.pipelines || []);
    const [deletePipeline, { data: deleteData, loading: deleteLoading, error: deleteError }] = useDeletePipelineMutation();
    const [selectedPipelines, setSelectedPipelines] = useState<GridRowSelectionModel>([]);
    const navigate = useNavigate();
    const [openConfirmDelete, setOpenConfirmDelete] = useState(false);
    const [openAddDialog, setOpenAddDialog] = useState(false);
    const [fetchCount, setFetchCount] = useState(0);

    useEffect(() => {
        const interval = setInterval(() => {
            if (fetchCount < 10) {
                refetch();
                setFetchCount(fetchCount + 1);
            } else {
                clearInterval(interval);
            }
        }, 5000);

        return () => clearInterval(interval);
    }, [fetchCount]);

    const handleShowDeleteDialog = () => {
        setOpenConfirmDelete(true);
    };

    const confirmDelete = () => {
        selectedPipelines.forEach(async (id) => {
            deletePipeline({ variables: { id: String(id) } }).then(() => {
                // TODO: Do this better
                refetch();
            });
        });
        setOpenConfirmDelete(false);
        setSelectedPipelines([]);
    };

    const handleSearch = (event: React.ChangeEvent<HTMLInputElement>) => {
        const query = event.target.value;
        setSearchQuery(query);
        const filteredPipelines = pipelines.filter(pipeline =>
            pipeline.name.toLowerCase().includes(query.toLowerCase())
        );
        setPipelines(filteredPipelines);
    };

    useEffect(() => {
        if (data) {
            setPipelines(data.pipelines);
        }
    }, [data]);
    

    return (
        <>
            <Box sx={{ p: 2, display: "flex", flexDirection: "column", height: "100%" }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                    <Typography variant="h4">Pipelines Overview</Typography>
                    <Box>
                        <Button variant="contained" onClick={() => setOpenAddDialog(true)} sx={{ marginRight: 1 }}>
                            Add Pipeline
                        </Button>
                        <Button
                            variant="outlined"
                            onClick={handleShowDeleteDialog}
                            disabled={selectedPipelines.length === 0}
                            color="error"
                        >
                            Delete Pipeline
                        </Button>
                    </Box>
                </Box>
                <TextField
                    fullWidth
                    label="Search Pipelines"
                    variant="outlined"
                    value={searchQuery}
                    onChange={handleSearch}
                    sx={{ mb: 2 }}
                />
                <Box sx={{ height: '100%', width: '100%' }}>
                    <DataGrid
                        rows={pipelines}
                        columns={columns}
                        checkboxSelection
                        disableColumnSelector
                        rowSelectionModel={selectedPipelines}
                        onRowSelectionModelChange={setSelectedPipelines}
                    />
                </Box>
                <ConfirmDelete
                    open={openConfirmDelete}
                    onClose={() => setOpenConfirmDelete(false)}
                    onConfirm={confirmDelete}
                    type={"pipeline(s)"}
                />
            </Box>
            <AddPipelineDialog
                open={openAddDialog}
                onClose={() => setOpenAddDialog(false)}
                onSuccess={(id: string) => navigate(`/pipelines/${id}`)}
                refetch={refetch}
            />
        </>
    );
}
