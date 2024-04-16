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

export default function StorageDetailsPage() {
    const { storageId } = useParams<{ storageId: string }>(); // Use useParams to get the ID
    const { data, loading, error } = useGetStorageQuery({ variables: { id: storageId ? storageId : ""} });

    const handleDeploy = () => {
        console.log('Deploy storage logic here');
        // Actual deployment logic to be implemented
    };

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
                    {data?.getStorage?.status === "READY" && (
                        <Button
                            startIcon={<SettingsIcon />}
                            onClick={handleEditDeployment}
                        >
                            Edit
                        </Button>
                    )}
                </Box>
                <Divider sx={{ my: 2 }} />
                {data?.getStorage?.status === "READY" || data?.getStorage?.status === "DEPLOYING" ? (
                    <Grid container spacing={2}>
                        <Grid item xs={12}>
                            {/* Deployment specifics like Replicas, Memory, CPU to be implemented */}
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
