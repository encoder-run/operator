import React, { useEffect, useState } from 'react';
import { useParams, Link as RouterLink } from 'react-router-dom';
import NavigateNextIcon from '@mui/icons-material/NavigateNext';
import {
    Typography,
    Box,
    Button,
    CircularProgress,
    Paper,
    Grid,
    Divider,
    Breadcrumbs,
    Link
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import { useGetModelQuery } from '../../api/types';
import AddModelDeploymentDialog from './add-model-deployment-dialog';
import toast from 'react-hot-toast';

export default function ModelDetailsPage() {
    const { modelId } = useParams<{ modelId: string }>(); // Use useParams to get the ID
    const { data, loading, error, refetch } = useGetModelQuery({ 
        variables: { id: modelId ? modelId : "" }, 
        nextFetchPolicy: "network-only"
    });
    const [openAddDeploymentDialog, setOpenAddDeploymentDialog] = useState(false);
    const [isSubmittingOrDeploying, setIsSubmittingOrDeploying] = useState(data?.getModel.status === "DEPLOYING" || false);

    useEffect(() => {
        let pollCount = 0;
        const maxPollCount = 10; // Maximum number of polling attempts
    
        const pollingInterval = setInterval(() => {
            if (isSubmittingOrDeploying) {
                refetch();
                pollCount += 1;
    
                // Check if the model's status has changed to 'READY'
                if (data?.getModel.status === "READY") {
                    setIsSubmittingOrDeploying(false);
                    clearInterval(pollingInterval);
                }
    
                // If maximum poll count reached, show error and stop polling
                if (pollCount >= maxPollCount) {
                    toast.error('Deployment is taking longer than expected. Please check back later.',{
                        position: "top-center"
                      });
                    setIsSubmittingOrDeploying(false);
                    clearInterval(pollingInterval);
                }
            }
        }, 5000); // Poll every 5 seconds
    
        // Cleanup interval on component unmount
        return () => clearInterval(pollingInterval);
    }, [isSubmittingOrDeploying]);

    // useEffect to update if the status is currently DEPLOYING
    useEffect(() => {
        setIsSubmittingOrDeploying(data?.getModel.status === "DEPLOYING" || false);
    }, [data?.getModel.status]);


    const handleEditDeployment = () => {
        console.log('Edit deployment logic here');
    };

    if (loading) {
        return <CircularProgress />;
    }

    if (error) {
        return <Typography>Error: {error.message}</Typography>;
    }

    return (
        <>
            <Box sx={{ p: 2 }}>
                <Box sx={{ pb: 2 }}>
                    <Breadcrumbs aria-label="breadcrumb" separator={<NavigateNextIcon fontSize="small" />}>
                        <Link component={RouterLink} to="/models">
                            Models
                        </Link>
                        <Typography color="text.primary">{data?.getModel.displayName}</Typography>
                    </Breadcrumbs>
                </Box>
                <Paper elevation={3} sx={{ p: 2, mb: 2 }}>
                    <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                        <Typography variant="h5" gutterBottom>
                            Model Details
                        </Typography>
                        <Button
                            startIcon={<EditIcon />}
                            onClick={() => console.log("Edit Model Details")}
                        >
                            Edit
                        </Button>
                    </Box>
                    <Divider sx={{ my: 2 }} />
                    <Grid container spacing={2}>
                        <Grid item xs={12}>
                            <Typography variant="body1"><strong>Name:</strong> {data?.getModel.displayName}</Typography>
                            <Typography variant="body1"><strong>Type:</strong> {data?.getModel.type}</Typography>
                            <Typography variant="body1"><strong>Status:</strong> {data?.getModel.status}</Typography>
                        </Grid>
                    </Grid>
                </Paper>
                <Paper elevation={3} sx={{ p: 2, display: "flex", flexDirection: "column" }}>
                    <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                        <Typography variant="h5" gutterBottom>
                            Deployment Details
                        </Typography>
                        {data?.getModel?.status === "READY" ? (
                            <Button
                                startIcon={<EditIcon />}
                                onClick={handleEditDeployment}
                            >
                                Edit
                            </Button>
                        ) : isSubmittingOrDeploying ? (
                            <CircularProgress size={24} />
                        ) : null}
                    </Box>
                    <Divider sx={{ my: 2 }} />
                    {data?.getModel.deployment ? (
                        <Grid container spacing={2}>
                            <Grid item xs={12}>
                                <Typography variant="body1"><strong>Enabled:</strong> {String(data?.getModel?.deployment?.enabled)}</Typography>
                                <Typography variant="body1"><strong>Memory:</strong> {data?.getModel?.deployment?.memory}</Typography>
                                <Typography variant="body1"><strong>CPU:</strong> {data?.getModel?.deployment?.cpu}</Typography>
                            </Grid>
                        </Grid>
                    ) : (
                        <Box sx={{ display: "flex", justifyContent: "center" }}>
                            <Button onClick={() => setOpenAddDeploymentDialog(true)} sx={{ maxWidth: "50vh" }} variant="contained" startIcon={<CloudUploadIcon />}>
                                Create Deployment
                            </Button>
                        </Box>
                    )}
                </Paper>
            </Box>
            <Box>
                <AddModelDeploymentDialog
                    open={openAddDeploymentDialog}
                    onClose={() => setOpenAddDeploymentDialog(false)}
                    onSuccess={() => {
                        setIsSubmittingOrDeploying(true);
                        toast.success('Deployment created successfully',{
                            position: "top-center"
                          });
                    }}
                    refetch={refetch}
                    modelId={modelId ? modelId : ""}
                />
            </Box>
        </>
    );
}
