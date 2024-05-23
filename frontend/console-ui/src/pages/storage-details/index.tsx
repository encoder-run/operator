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
import SettingsIcon from '@mui/icons-material/Settings';
import { useGetStorageQuery } from '../../api/types';
import AddStorageDeploymentDialog from './add-storage-deployment-dialog';
import toast from 'react-hot-toast';
import { useEffect, useState } from 'react';

export default function StorageDetailsPage() {
    const { storageId } = useParams<{ storageId: string }>(); // Use useParams to get the ID
    const { data, loading, error, refetch } = useGetStorageQuery({
        variables: { id: storageId ? storageId : "" },
        nextFetchPolicy: "network-only"
    });

    const [openAddDeploymentDialog, setOpenAddDeploymentDialog] = useState(false);
    const [isSubmittingOrDeploying, setIsSubmittingOrDeploying] = useState(data?.getStorage.status === "DEPLOYING" || false);

    useEffect(() => {
        let pollCount = 0;
        const maxPollCount = 10; // Maximum number of polling attempts
    
        const pollingInterval = setInterval(() => {
            if (isSubmittingOrDeploying) {
                refetch();
                pollCount += 1;
    
                // Check if the storage's status has changed to 'READY'
                if (data?.getStorage.status === "READY") {
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
        setIsSubmittingOrDeploying(data?.getStorage.status === "DEPLOYING" || false);
    }, [data?.getStorage.status]);

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
                        <Link component={RouterLink} to="/storage">
                            Storage
                        </Link>
                        <Typography color="text.primary">{data?.getStorage.name}</Typography>
                    </Breadcrumbs>
                </Box>
                <Paper elevation={3} sx={{ p: 2, mb: 2 }}>
                    <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                        <Typography variant="h5" gutterBottom>
                            Storage Details
                        </Typography>
                        <Button
                            startIcon={<EditIcon />}
                            onClick={() => console.log("Edit Storage Details")}
                        >
                            Edit
                        </Button>
                    </Box>
                    <Divider sx={{ my: 2 }} />
                    <Grid container spacing={2}>
                        <Grid item xs={12}>
                            <Typography variant="body1"><strong>Name:</strong> {data?.getStorage.name}</Typography>
                            <Typography variant="body1"><strong>Type:</strong> {data?.getStorage.type}</Typography>
                            <Typography variant="body1"><strong>Status:</strong> {data?.getStorage.status}</Typography>
                        </Grid>
                    </Grid>
                </Paper>
                
                <Paper elevation={3} sx={{ p: 2, display: "flex", flexDirection: "column" }}>
                    <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                        <Typography variant="h5" gutterBottom>
                            Deployment Details
                        </Typography>
                        {data?.getStorage?.status === "READY" ? (
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
                    {data?.getStorage.deployment ? (
                        <Grid container spacing={2}>
                            <Grid item xs={12}>
                                <Typography variant="body1"><strong>Enabled:</strong> {String(data?.getStorage?.deployment?.enabled)}</Typography>
                                <Typography variant="body1"><strong>Memory:</strong> {data?.getStorage?.deployment?.memory}</Typography>
                                <Typography variant="body1"><strong>CPU:</strong> {data?.getStorage?.deployment?.cpu}</Typography>
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
                <AddStorageDeploymentDialog
                    open={openAddDeploymentDialog}
                    onClose={() => setOpenAddDeploymentDialog(false)}
                    onSuccess={() => {
                        setIsSubmittingOrDeploying(true);
                        toast.success('Deployment created successfully',{
                            position: "top-center"
                          });
                    }}
                    refetch={refetch}
                    storageId={storageId ? storageId : ""}
                />
            </Box>
        </>
    );
}
