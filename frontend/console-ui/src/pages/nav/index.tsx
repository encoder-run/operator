import { Box, Typography, Divider, List, ListItem, ListItemButton, ListItemText } from '@mui/material';
import { Outlet, Link } from 'react-router-dom';
import { MagnifyingGlassIcon } from "@heroicons/react/24/outline";
import { WrenchIcon } from "@heroicons/react/24/outline";
import { ServerStackIcon } from "@heroicons/react/24/outline";
import { CodeBracketSquareIcon } from "@heroicons/react/24/outline";
import { VariableIcon } from "@heroicons/react/24/outline";
import { CubeTransparentIcon } from "@heroicons/react/24/outline";
import logoImage from '../../assets/logo.png';

export default function Sidebar() {
    return (
        <>
            <Box sx={{
                width: '18rem',
                backgroundColor: '#f7f7f7',
                borderRight: '1px solid #e3e3e3',
                display: 'flex',
                flexDirection: 'column',
                p: 2,
            }}>
                {/* Top Nav Items */}
                <Box sx={{ display: 'flex', flexDirection: 'column' }}>
                    <Typography variant="h6" sx={{ pb: 1 }}>
                        <Box sx={{ display: "flex", justifyContent: "center", alignItems: "center" }}>
                            <img src={logoImage} alt="Logo" style={{ width: '50px' }} />
                            <Typography variant="h6" sx={{ display: 'inline', ml: 1 }}>encoder.run</Typography>
                        </Box>
                    </Typography>
                    <Divider sx={{ my: 1 }} />
                    <List>
                        <ListItem disablePadding>
                            <ListItemButton component={Link} to="/search">
                                <Box
                                    width={24}
                                    height={24}>
                                    <MagnifyingGlassIcon />
                                </Box>
                                <ListItemText sx={{ p: 1 }} primary="Search" />
                            </ListItemButton>
                        </ListItem>
                    </List>
                </Box>

                {/* Bottom Nav Items */}
                <Box sx={{ marginTop: 'auto', pb: 2 }}>
                    {/* Divider */}
                    <Divider sx={{ my: 1 }} />
                    <List>
                        <ListItem disablePadding>
                            <ListItemButton component={Link} to="/api">
                                <Box
                                    width={24}
                                    height={24}>
                                    <WrenchIcon />
                                </Box>
                                <ListItemText sx={{ p: 1 }} primary="API" />
                            </ListItemButton>
                        </ListItem>
                        <ListItem disablePadding>
                            <ListItemButton component={Link} to="/pipelines">
                                <Box
                                    width={24}
                                    height={24}>
                                    <CubeTransparentIcon />
                                </Box>
                                <ListItemText sx={{ p: 1 }} primary="Pipelines" />
                            </ListItemButton>
                        </ListItem>
                        <ListItem disablePadding>
                            <ListItemButton component={Link} to="/models">
                                <Box
                                    width={24}
                                    height={24}>
                                    <VariableIcon />
                                </Box>
                                <ListItemText sx={{ p: 1 }} primary="Models" />
                            </ListItemButton>
                        </ListItem>
                        <ListItem disablePadding>
                            <ListItemButton component={Link} to="/repositories">
                                <Box
                                    width={24}
                                    height={24}>
                                    <CodeBracketSquareIcon />
                                </Box>
                                <ListItemText sx={{ p: 1 }} primary="Repositories" />
                            </ListItemButton>
                        </ListItem>
                        <ListItem disablePadding>
                            <ListItemButton component={Link} to="/storage">
                                <Box
                                    width={24}
                                    height={24}>
                                    <ServerStackIcon />
                                </Box>
                                <ListItemText sx={{ p: 1 }} primary="Storage" />
                            </ListItemButton>
                        </ListItem>
                    </List>
                </Box>
            </Box>
            <Box sx={{ flexGrow: 1, width: "100%" }}>
                <Outlet/>
            </Box>
        </>
    );
}
