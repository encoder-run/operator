import React from 'react';
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
import { useGetRepositoryQuery } from '../../api/types';  // Adjust this import based on your file structure

export default function RepositoryDetailsPage() {
    const { repositoryId } = useParams<{ repositoryId: string }>(); // Use useParams to get the ID
    const { data, loading, error } = useGetRepositoryQuery({ variables: { id: repositoryId ? repositoryId : "" } });

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
                    <Link component={RouterLink} to="/repositories">
                        Repositories
                    </Link>
                    <Typography color="text.primary">{data?.getRepository.displayName}</Typography>
                </Breadcrumbs>
            </Box>
            <Paper elevation={3} sx={{ p: 2, mb: 2 }}>
                <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                    <Typography variant="h5" gutterBottom>
                        Repository Details
                    </Typography>
                    <Button
                        startIcon={<EditIcon />}
                        onClick={() => console.log("Edit Repository Details")}
                    >
                        Edit
                    </Button>
                </Box>
                <Divider sx={{ my: 2 }} />
                <Grid container spacing={2}>
                    <Grid item xs={12}>
                        <Typography variant="body1"><strong>Display Name:</strong> {data?.getRepository.displayName}</Typography>
                        <Typography variant="body1"><strong>Type:</strong> {data?.getRepository.type}</Typography>
                        <Typography variant="body1"><strong>Owner:</strong> {data?.getRepository.owner}</Typography>
                        <Typography variant="body1"><strong>Repository:</strong> {data?.getRepository.name}</Typography>
                        <Typography variant="body1"><strong>URL:</strong> <Link href={data?.getRepository.url} target="_blank">{data?.getRepository.url}</Link></Typography>
                    </Grid>
                </Grid>
            </Paper>
        </Box>
    );
}
