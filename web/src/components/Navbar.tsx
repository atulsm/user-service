import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Box,
  IconButton,
  Menu,
  MenuItem,
  Avatar,
  Tooltip,
  Divider,
  useTheme,
  useMediaQuery,
  Drawer,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
} from '@mui/material';
import {
  AccountCircle,
  People as PeopleIcon,
  Logout as LogoutIcon,
  Menu as MenuIcon,
  Dashboard as DashboardIcon,
  Settings as SettingsIcon,
  Add as AddIcon,
} from '@mui/icons-material';
import { authService } from '../services/api';

const Navbar: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const userEmail = localStorage.getItem('userEmail');

  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = async () => {
    handleClose();
    await authService.logout();
    navigate('/login');
  };

  const handleProfile = () => {
    handleClose();
    navigate('/users/profile');
  };

  const handleSettings = () => {
    handleClose();
    navigate('/settings');
  };

  const handleCreateUser = () => {
    navigate('/users/new');
  };

  const isActive = (path: string) => location.pathname === path;

  const navItems = [
    { text: 'Dashboard', icon: <DashboardIcon />, path: '/dashboard' },
    { text: 'Users', icon: <PeopleIcon />, path: '/users' },
  ];

  const renderNavItems = () => (
    <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
      {navItems.map((item) => (
        <Button
          key={item.path}
          color="inherit"
          onClick={() => navigate(item.path)}
          startIcon={item.icon}
          sx={{
            backgroundColor: isActive(item.path) ? 'rgba(255, 255, 255, 0.1)' : 'transparent',
            '&:hover': {
              backgroundColor: 'rgba(255, 255, 255, 0.2)',
            },
          }}
        >
          {item.text}
        </Button>
      ))}
      <Button
        color="inherit"
        onClick={handleCreateUser}
        startIcon={<AddIcon />}
        sx={{
          backgroundColor: 'rgba(255, 255, 255, 0.1)',
          '&:hover': {
            backgroundColor: 'rgba(255, 255, 255, 0.2)',
          },
        }}
      >
        Create User
      </Button>
    </Box>
  );

  const renderMobileMenu = () => (
    <Drawer
      anchor="left"
      open={mobileMenuOpen}
      onClose={() => setMobileMenuOpen(false)}
    >
      <Box sx={{ width: 250 }}>
        <List>
          {navItems.map((item) => (
            <ListItem
              button
              key={item.path}
              onClick={() => {
                navigate(item.path);
                setMobileMenuOpen(false);
              }}
              selected={isActive(item.path)}
            >
              <ListItemIcon>{item.icon}</ListItemIcon>
              <ListItemText primary={item.text} />
            </ListItem>
          ))}
          <ListItem button onClick={handleCreateUser}>
            <ListItemIcon><AddIcon /></ListItemIcon>
            <ListItemText primary="Create User" />
          </ListItem>
        </List>
        <Divider />
        <List>
          <ListItem button onClick={handleProfile}>
            <ListItemIcon><AccountCircle /></ListItemIcon>
            <ListItemText primary="Profile" />
          </ListItem>
          <ListItem button onClick={handleSettings}>
            <ListItemIcon><SettingsIcon /></ListItemIcon>
            <ListItemText primary="Settings" />
          </ListItem>
          <ListItem button onClick={handleLogout}>
            <ListItemIcon><LogoutIcon /></ListItemIcon>
            <ListItemText primary="Logout" />
          </ListItem>
        </List>
      </Box>
    </Drawer>
  );

  return (
    <AppBar 
      position="static" 
      elevation={0} 
      sx={{ 
        borderBottom: '1px solid rgba(0, 0, 0, 0.12)',
        backgroundColor: 'white',
        color: 'text.primary',
      }}
    >
      <Toolbar>
        {isMobile && (
          <IconButton
            edge="start"
            color="inherit"
            aria-label="menu"
            onClick={() => setMobileMenuOpen(true)}
            sx={{ mr: 2 }}
          >
            <MenuIcon />
          </IconButton>
        )}
        <Typography
          variant="h6"
          component="div"
          sx={{ 
            flexGrow: 1, 
            display: 'flex', 
            alignItems: 'center',
            color: 'primary.main',
            fontWeight: 'bold',
          }}
        >
          <PeopleIcon sx={{ mr: 1 }} />
          User Service
        </Typography>
        {!isMobile && renderNavItems()}
        <Box sx={{ ml: 2 }}>
          <Tooltip title="Account settings">
            <IconButton
              size="large"
              aria-label="account of current user"
              aria-controls="menu-appbar"
              aria-haspopup="true"
              onClick={handleMenu}
              color="inherit"
            >
              {userEmail ? (
                <Avatar sx={{ width: 32, height: 32, bgcolor: 'primary.main' }}>
                  {userEmail[0].toUpperCase()}
                </Avatar>
              ) : (
                <AccountCircle />
              )}
            </IconButton>
          </Tooltip>
          <Menu
            id="menu-appbar"
            anchorEl={anchorEl}
            anchorOrigin={{
              vertical: 'bottom',
              horizontal: 'right',
            }}
            keepMounted
            transformOrigin={{
              vertical: 'top',
              horizontal: 'right',
            }}
            open={Boolean(anchorEl)}
            onClose={handleClose}
          >
            <MenuItem onClick={handleProfile}>
              <ListItemIcon>
                <AccountCircle fontSize="small" />
              </ListItemIcon>
              Profile
            </MenuItem>
            <MenuItem onClick={handleSettings}>
              <ListItemIcon>
                <SettingsIcon fontSize="small" />
              </ListItemIcon>
              Settings
            </MenuItem>
            <Divider />
            <MenuItem onClick={handleLogout}>
              <ListItemIcon>
                <LogoutIcon fontSize="small" />
              </ListItemIcon>
              Logout
            </MenuItem>
          </Menu>
        </Box>
      </Toolbar>
      {renderMobileMenu()}
    </AppBar>
  );
};

export default Navbar; 