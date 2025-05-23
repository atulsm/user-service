import React from 'react';
import { Box, Typography } from '@mui/material';
import DashboardCharts from '../components/DashboardCharts';

const Dashboard: React.FC = () => {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Dashboard
      </Typography>
      <DashboardCharts />
    </Box>
  );
};

export default Dashboard; 