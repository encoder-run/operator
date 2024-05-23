import { useParams, Link as RouterLink, useNavigate } from 'react-router-dom';
import NavigateNextIcon from '@mui/icons-material/NavigateNext';
import LaunchIcon from '@mui/icons-material/Launch';
import {
    Typography,
    Box,
    Button,
    CircularProgress,
    Paper,
    Grid,
    Divider,
    Breadcrumbs,
    Link,
    IconButton
} from '@mui/material';
import PowerSettingsNewIcon from '@mui/icons-material/PowerSettingsNew';
import SyncIcon from '@mui/icons-material/Sync';
import { useAddPipelineDeploymentMutation, useGetPipelineExecutionsQuery, useGetPipelineQuery, useTriggerPipelineMutation } from '../../api/types'; // Ensure this is set up in your GraphQL API file
import toast from 'react-hot-toast';
import { DataGrid, GridColDef } from '@mui/x-data-grid';
import { useEffect } from 'react';

export default function PipelineDetailsPage() {
    const { pipelineId } = useParams<{ pipelineId: string }>();
    const { data, loading, error, refetch } = useGetPipelineQuery({
        variables: { id: pipelineId ? pipelineId : "" },
        nextFetchPolicy: "network-only"
    });
    const navigate = useNavigate();
    const [addDeployment, { loading: addDeploymentLoading }] = useAddPipelineDeploymentMutation();
    const [triggerPipeline] = useTriggerPipelineMutation();
    const {
        data: executionsData,
        loading: executionsLoading,
        error: executionsError,
        refetch: executionsRefetch
    } = useGetPipelineExecutionsQuery({
        variables: { id: pipelineId ? pipelineId : "" }
    });

    const columns: GridColDef[] = [
        { field: 'id', headerName: 'ID', flex: 1, minWidth: 100 },
        {
            field: 'status',
            headerName: 'Status',
            flex: 1,
            minWidth: 130,
            renderCell: (params) => {
                return ["pending", "active"].includes(params.value) ?
                    <CircularProgress size={20} /> :
                    <Typography>{params.value}</Typography>;
            }
        }, 
        {
            field: 'actions',
            headerName: 'Actions',
            type: 'actions',
            width: 100,
            getActions: (params: any) => [
                <IconButton
                    disabled
                    color="primary"
                    aria-label="go to details"
                    size="small"
                    onClick={() => navigate(`/execution-details/${params.id}`)}
                >
                    <LaunchIcon />
                </IconButton>,
            ],
        },
    ];

    useEffect(() => {
        const interval = setInterval(() => {
            executionsRefetch();
        }, 10000);

        // Stop polling after 15 minutes
        const timeout = setTimeout(() => {
            clearInterval(interval);
        }, 900000);

        return () => {
            clearInterval(interval);
            clearTimeout(timeout);
        };
    }, [executionsRefetch]);

    const handleEnablePipeline = () => {
        if (!pipelineId) return;
        addDeployment({
            variables: { input: { id: pipelineId, enabled: true } }
        }).then(() => {
            refetch();
            // Logic to enable the pipeline should be implemented here
            // sleep for 1 second then refetch\
            toast.success('Pipeline enabled successfully',{
                position: "top-center"
              });
            setTimeout(() => {
                executionsRefetch();
            }, 500);

            console.log('Enable pipeline logic here');
        });

    };

    const handleTriggerPipeline = () => {
        if (!pipelineId) return;
        triggerPipeline({
            variables: { id: pipelineId }
        }).then(() => {
            setTimeout(() => {
                executionsRefetch();
            }, 500);
            toast.success('Pipeline triggered successfully',{
                position: "top-center"
              });
            console.log('Trigger pipeline logic here');
        });
    }

    if (loading) return <CircularProgress />;
    if (error) return <Typography>Error: {error.message}</Typography>;

    return (
        <>
            <Box sx={{ p: 2 }}>
                <Box sx={{ pb: 2 }}>
                    <Breadcrumbs aria-label="breadcrumb" separator={<NavigateNextIcon fontSize="small" />}>
                        <Link component={RouterLink} to="/pipelines">
                            Pipelines
                        </Link>
                        <Typography color="text.primary">{data?.getPipeline.name}</Typography>
                    </Breadcrumbs>
                </Box>
                <Paper elevation={3} sx={{ p: 2, mb: 2 }}>
                    <Typography variant="h5" gutterBottom>
                        Pipeline Details
                    </Typography>
                    <Divider sx={{ my: 2 }} />
                    <Grid container spacing={2}>
                        <Grid item xs={12}>
                            <Typography variant="body1"><strong>ID:</strong> {data?.getPipeline.id}</Typography>
                            <Typography variant="body1"><strong>Name:</strong> {data?.getPipeline.name}</Typography>
                            <Typography variant="body1"><strong>Status:</strong> {data?.getPipeline.status}</Typography>
                            {data?.getPipeline?.repositoryEmbeddings && (
                                <>
                                    <Typography variant="body1"><strong>Repository ID:</strong> {data.getPipeline.repositoryEmbeddings.repositoryID}</Typography>
                                    <Typography variant="body1"><strong>Model ID:</strong> {data.getPipeline.repositoryEmbeddings.modelID}</Typography>
                                    <Typography variant="body1"><strong>Storage ID:</strong> {data.getPipeline.repositoryEmbeddings.storageID}</Typography>
                                </>
                            )}
                        </Grid>
                    </Grid>
                </Paper>
                <Paper elevation={3} sx={{ p: 2, display: "flex", flexDirection: "column" }}>
                    <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                        <Typography variant="h5" gutterBottom>
                            Deployment Details
                        </Typography>
                        <Button
                            startIcon={<PowerSettingsNewIcon />}
                            onClick={handleEnablePipeline}
                            variant="contained"
                            color="primary"
                            disabled={data?.getPipeline.enabled}
                        >
                            Enable Pipeline
                        </Button>
                    </Box>
                    <Divider sx={{ my: 2 }} />
                    <Typography variant="body1"><strong>Enabled:</strong> {String(data?.getPipeline.enabled)}</Typography>
                </Paper>
                <Paper elevation={3} sx={{ p: 2, mt: 2 }}>
                    <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>

                        <Typography variant="h5" gutterBottom>
                            Pipeline Executions
                        </Typography>
                        <Button
                            startIcon={<SyncIcon />}
                            onClick={handleTriggerPipeline}
                            variant="contained"
                            color="primary"
                            disabled={!data?.getPipeline.enabled}
                        >
                            Manual Trigger
                        </Button>
                    </Box>
                    <Divider sx={{ my: 2 }} />
                    <Box sx={{ height: 400, width: '100%' }}>
                        <DataGrid
                            rows={executionsData?.getPipelineExecutions || []}
                            columns={columns}
                            initialState={{
                                pagination: {
                                    paginationModel: {
                                        pageSize: 5,
                                    },
                                },
                            }}
                            pageSizeOptions={[5]}
                            checkboxSelection
                            disableRowSelectionOnClick
                            disableColumnSelector
                        />
                    </Box>
                </Paper>
            </Box>
        </>
    );
}
