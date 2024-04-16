import React, { useState } from 'react';
import { Typography, Box, TextField, List, ListItem, ListItemText, Divider } from '@mui/material';

// Dummy data for initial demonstration
const initialCodeChunks = [
    { chunkId: 1, content: 'Example code snippet 1', startCol: 1, endCol: 5, org: 'OpenAI', repo: 'GPT' },
    { chunkId: 2, content: 'Example code snippet 2', startCol: 10, endCol: 15, org: 'Google', repo: 'TensorFlow' },
    { chunkId: 3, content: 'Example code snippet 3', startCol: 5, endCol: 12, org: 'Facebook', repo: 'React' },
    // Add more code chunks here
];

export default function CodeSearchPage() {
    const [searchQuery, setSearchQuery] = useState('');
    const [codeChunks, setCodeChunks] = useState(initialCodeChunks);

    const handleSearch = (event: React.ChangeEvent<HTMLInputElement>) => {
        const query = event.target.value.toLowerCase();
        setSearchQuery(query);
        const filteredCodeChunks = initialCodeChunks.filter(chunk =>
            chunk.content.toLowerCase().includes(query) ||
            chunk.org.toLowerCase().includes(query) ||
            chunk.repo.toLowerCase().includes(query)
        );
        setCodeChunks(filteredCodeChunks);
    };

    return (
        <Box sx={{ p: 2, display: "flex", flexDirection: "column", height: "100%" }}>
            <Typography variant="h4" sx={{ mb: 2 }}>Source Code Search</Typography>
            <TextField
                fullWidth
                label="Search"
                variant="outlined"
                value={searchQuery}
                onChange={handleSearch}
                sx={{ mb: 2 }}
            />
            <List sx={{ width: '100%', bgcolor: 'background.paper' }}>
                {codeChunks.map((chunk, index) => (
                    <React.Fragment key={chunk.chunkId}>
                        {index > 0 && <Divider component="li" />}
                        <ListItem alignItems="flex-start">
                            <ListItemText
                                primary={`Chunk ID: ${chunk.chunkId}, Org: ${chunk.org}, Repo: ${chunk.repo}`}
                                secondary={
                                    <>
                                        <Typography
                                            sx={{ display: 'inline' }}
                                            component="span"
                                            variant="body2"
                                            color="text.primary"
                                        >
                                            Content: 
                                        </Typography>
                                        {` ${chunk.content}`}
                                        <br />
                                        {`Start/End Columns: ${chunk.startCol} to ${chunk.endCol}`}
                                    </>
                                }
                            />
                        </ListItem>
                    </React.Fragment>
                ))}
            </List>
        </Box>
    );
}
