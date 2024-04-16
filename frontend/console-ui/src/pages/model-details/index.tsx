import React, { useState, useEffect } from 'react';
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

// Mock function to simulate fetching model details
const fetchModelDetails = (modelId: string) => {
    return new Promise((resolve) => {
        setTimeout(() => {
            resolve({
                id: modelId,
                name: 'GPT-3',
                type: 'Transformer',
                modelSize: 175,
                deployed: false, // Assume not deployed for initial state
                deploymentDetails: {
                    replicas: 2,
                    memory: '8GB',
                    cpu: '4 cores',
                },
            });
        }, 1000);
    });
};

export default function ModelDetailsPage() {
    const { modelId } = useParams();
    const [modelDetails, setModelDetails] = useState<any>(null);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        setLoading(true);
        if (!modelId) {
            return;
        }
        fetchModelDetails(modelId).then((details) => {
            setModelDetails(details);
            setLoading(false);
        });
    }, [modelId]);

    const handleDeploy = () => {
        // Implement the deployment logic here
        console.log('Deploy model logic here');
        // Simulate deployment success
        setModelDetails((prevDetails: any) => ({
            ...prevDetails,
            deployed: true,
        }));
    };

    const handleEditDeployment = () => {
        // Placeholder for edit deployment logic
        console.log('Edit deployment logic here');
    };

    if (loading) {
        return <CircularProgress />;
    }

    return (
        <Box sx={{ p: 2 }}>
            <Box sx={{pb: 2}}>
                <Breadcrumbs aria-label="breadcrumb" separator={<NavigateNextIcon fontSize="small" />}>
                    <Link component={RouterLink} to="/models">
                        Models
                    </Link>
                    <Typography color="text.primary">{modelDetails?.name}</Typography>
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
                        <Typography variant="body1"><strong>Name:</strong> {modelDetails?.name}</Typography>
                        <Typography variant="body1"><strong>Type:</strong> {modelDetails?.type}</Typography>
                        <Typography variant="body1"><strong>Model Size:</strong> {modelDetails?.modelSize} GB</Typography>
                    </Grid>
                </Grid>
            </Paper>

            <Paper elevation={3} sx={{ p: 2, display: "flex", flexDirection: "column" }}>
                <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                    <Typography variant="h5" gutterBottom>
                        Deployment Details
                    </Typography>
                    {modelDetails?.deployed && (
                        <Button
                            startIcon={<SettingsIcon />}
                            onClick={handleEditDeployment}
                        >
                            Edit
                        </Button>
                    )}
                </Box>
                <Divider sx={{ my: 2 }} />
                {modelDetails?.deployed ? (
                    <Grid container spacing={2}>
                        <Grid item xs={12}>
                            <Typography variant="body1"><strong>Replicas:</strong> {modelDetails?.deploymentDetails.replicas}</Typography>
                            <Typography variant="body1"><strong>Memory:</strong> {modelDetails?.deploymentDetails.memory}</Typography>
                            <Typography variant="body1"><strong>CPU:</strong> {modelDetails?.deploymentDetails.cpu}</Typography>
                        </Grid>
                    </Grid>
                ) : (
                    <Box sx={{display: "flex", justifyContent: "center"}}>
                    <Button sx={{maxWidth: "50vh"}} variant="contained" startIcon={<CloudUploadIcon />} onClick={handleDeploy}>
                        Create Deployment
                    </Button>
                    </Box>
                )}
            </Paper>
        </Box>
    );
}
