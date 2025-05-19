import React, { useEffect, useState } from 'react';
import {
  Box,
  Grid,
  Paper,
  Typography,
  Card,
  CardContent,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Divider,
  Button,
  CircularProgress,
  useTheme,
} from '@mui/material';
import {
  People as PeopleIcon,
  PersonAdd as PersonAddIcon,
  AccessTime as AccessTimeIcon,
  CheckCircle as CheckCircleIcon,
  Warning as WarningIcon,
  Error as ErrorIcon,
} from '@mui/icons-material';
import { userService } from '../services/api';
import { useNavigate } from 'react-router-dom';

interface UserStats {
  totalUsers: number;
  activeUsers: number;
  inactiveUsers: number;
}

interface RecentActivity {
  id: string;
  type: 'user_created' | 'user_updated' | 'user_deleted';
  description: string;
  timestamp: string;
  user: {
    id: string;
    email: string;
  };
}

const Dashboard: React.FC = () => {
  const [stats, setStats] = useState<UserStats | null>(null);
  const [recentActivity, setRecentActivity] = useState<RecentActivity[]>([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  const theme = useTheme();

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        // In a real application, these would be actual API calls
        // For now, we'll use mock data
        setStats({
          totalUsers: 150,
          activeUsers: 120,
          inactiveUsers: 30,
        });

        setRecentActivity([
          {
            id: '1',
            type: 'user_created',
            description: 'New user created',
            timestamp: new Date().toISOString(),
            user: { id: '1', email: 'john.doe@example.com' },
          },
          {
            id: '2',
            type: 'user_updated',
            description: 'User profile updated',
            timestamp: new Date(Date.now() - 3600000).toISOString(),
            user: { id: '2', email: 'jane.smith@example.com' },
          },
          {
            id: '3',
            type: 'user_deleted',
            description: 'User account deleted',
            timestamp: new Date(Date.now() - 7200000).toISOString(),
            user: { id: '3', email: 'bob.johnson@example.com' },
          },
        ]);

        setLoading(false);
      } catch (error) {
        console.error('Error fetching dashboard data:', error);
        setLoading(false);
      }
    };

    fetchDashboardData();
  }, []);

  const getActivityIcon = (type: string) => {
    switch (type) {
      case 'user_created':
        return <PersonAddIcon color="success" />;
      case 'user_updated':
        return <CheckCircleIcon color="info" />;
      case 'user_deleted':
        return <ErrorIcon color="error" />;
      default:
        return <AccessTimeIcon />;
    }
  };

  if (loading) {
    return (
      <Box
        display="flex"
        justifyContent="center"
        alignItems="center"
        minHeight="80vh"
      >
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Dashboard
      </Typography>

      {/* Quick Actions */}
      <Box sx={{ mb: 4 }}>
        <Button
          variant="contained"
          startIcon={<PersonAddIcon />}
          onClick={() => navigate('/users/new')}
          sx={{ mr: 2 }}
        >
          Create New User
        </Button>
        <Button
          variant="outlined"
          startIcon={<PeopleIcon />}
          onClick={() => navigate('/users')}
        >
          View All Users
        </Button>
      </Box>

      {/* Statistics Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={2}>
                <PeopleIcon color="primary" sx={{ mr: 1 }} />
                <Typography variant="h6">Total Users</Typography>
              </Box>
              <Typography variant="h3">{stats?.totalUsers}</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={2}>
                <CheckCircleIcon color="success" sx={{ mr: 1 }} />
                <Typography variant="h6">Active Users</Typography>
              </Box>
              <Typography variant="h3">{stats?.activeUsers}</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={2}>
                <WarningIcon color="warning" sx={{ mr: 1 }} />
                <Typography variant="h6">Inactive Users</Typography>
              </Box>
              <Typography variant="h3">{stats?.inactiveUsers}</Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Recent Activity */}
      <Paper sx={{ p: 2 }}>
        <Typography variant="h6" gutterBottom>
          Recent Activity
        </Typography>
        <List>
          {recentActivity.map((activity, index) => (
            <React.Fragment key={activity.id}>
              <ListItem>
                <ListItemIcon>{getActivityIcon(activity.type)}</ListItemIcon>
                <ListItemText
                  primary={activity.description}
                  secondary={`${activity.user.email} - ${new Date(
                    activity.timestamp
                  ).toLocaleString()}`}
                />
              </ListItem>
              {index < recentActivity.length - 1 && <Divider />}
            </React.Fragment>
          ))}
        </List>
      </Paper>
    </Box>
  );
};

export default Dashboard; 