import React, { useState } from 'react';
import { Box, TextField, List, Typography, Button, Paper } from '@mui/material';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { materialLight } from 'react-syntax-highlighter/dist/esm/styles/prism';
import { SearchResult, useSemanticSearchLazyQuery } from '../../api/types';

// Helper to get language from file path
const getLanguage = (path: string) => {
    const extension = path.split('.').pop();
    switch (extension) {
        case 'js':
            return 'javascript';
        case 'py':
            return 'python';
        case 'java':
            return 'java';
        default:
            return 'plaintext';
    }
}

// Helper to process content and line number
const processContent = (content: string, startLine: number) => {
    if (content.startsWith('\n')) {
        return { content: content.slice(1), startLine: startLine + 1 };
    }
    return { content, startLine };
}

export default function CodeSearchPage() {
    const [searchQuery, setSearchQuery] = useState('');
    const [codeChunks, setCodeChunks] = useState<SearchResult[]>([]);
    const [expandedId, setExpandedId] = useState<string | null>(null);
    const [search, { data, loading, error }] = useSemanticSearchLazyQuery();

    const toggleExpand = (id: string) => {
        setExpandedId(expandedId === id ? null : id);
    };

    const handleSearch = (event?: React.FormEvent<HTMLFormElement> | React.KeyboardEvent<HTMLDivElement>) => {
        if (event) event.preventDefault(); // Prevent default behavior if called from form submit or key press
        setCodeChunks([]); // Clear existing chunks before fetching new ones
        search({
            variables: {
                query: {
                    query: searchQuery,
                    page: 1,
                    limit: 5,
                }
            },
            fetchPolicy: 'network-only',
        }).then((res) => {
            if (res.data) {
                setCodeChunks(res.data.semanticSearch);
            }
        });
    }

    const handleKeyPress = (event: React.KeyboardEvent<HTMLDivElement>) => {
        if (event.key === 'Enter') {
            handleSearch(event);
        }
    };

    return (
        <Box sx={{ p: 2, display: "flex", flexDirection: "column", height: "100%" }}>
            <Typography variant="h4" sx={{ mb: 2 }}>Search</Typography>
            <Box
                component="form"
                onSubmit={handleSearch}
                sx={{ display: "flex" }}
            >
                <TextField
                    fullWidth
                    label="Query"
                    variant="outlined"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    onKeyPress={handleKeyPress}
                    sx={{ mb: 2 }}
                />
                <Button
                    type="submit"
                    variant="contained"
                    color="primary"
                    sx={{ mb: 2, ml: 2 }}
                >
                    Search
                </Button>
            </Box>
            <List sx={{ width: '100%', bgcolor: 'background.paper' }}>
                {codeChunks.map((result) => {
                    const { content, startLine } = processContent(result.content, result.startLine);
                    const uniqueKey = `${result.owner}-${result.repo}-${result.path}-${result.hash}-${result.chunkID}`;
                    return (
                        <React.Fragment key={uniqueKey}>
                            <Paper variant="outlined" sx={{ p: 2, mb: 2 }}>
                                <Typography variant="subtitle1">
                                    {`${result.owner}/${result.repo} - ${result.path}`}
                                </Typography>
                                <SyntaxHighlighter
                                    language={getLanguage(result.path)}
                                    style={materialLight}
                                    showLineNumbers
                                    wrapLines
                                    startingLineNumber={startLine}
                                    lineProps={{ style: { wordBreak: 'break-all', whiteSpace: 'pre-wrap' } }}
                                >
                                    {expandedId === result.id ? content : content.split('\n').slice(0, 5).join('\n')}
                                </SyntaxHighlighter>
                                {content.split('\n').length > 5 && (
                                    <Button onClick={() => toggleExpand(result.id)} fullWidth>
                                        {expandedId === result.id ? 'Collapse' : 'Expand'}
                                    </Button>
                                )}
                            </Paper>
                        </React.Fragment>
                    );
                })}
            </List>
        </Box>
    );
}
