import React, { useState, useEffect } from 'react';
import {
  Box,
  Paper,
  Typography,
  Grid,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  SelectChangeEvent,
  TextField,
} from '@mui/material';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import { Line, Bar } from 'react-chartjs-2';
import { format, subDays, startOfDay, endOfDay, eachDayOfInterval } from 'date-fns';
import { userService } from '../services/api';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend
);

interface UserActivity {
  date: string;
  newUsers: number;
  activeUsers: number;
}

interface UserStats {
  totalUsers: number;
  activeUsers: number;
  newUsers: number;
}

// Generate dummy data for the selected date range
const generateDummyData = (start: Date, end: Date): UserActivity[] => {
  const days = eachDayOfInterval({ start, end });
  const baseActiveUsers = 150;
  const baseNewUsers = 10;

  return days.map((date, index) => {
    // Add some randomness to make the data look more realistic
    const randomFactor = 0.8 + Math.random() * 0.4; // Random factor between 0.8 and 1.2
    const trendFactor = 1 + (index / days.length) * 0.3; // Gradual upward trend
    const weekdayFactor = date.getDay() === 0 || date.getDay() === 6 ? 0.7 : 1; // Weekend effect

    return {
      date: date.toISOString(),
      activeUsers: Math.round(baseActiveUsers * randomFactor * trendFactor * weekdayFactor),
      newUsers: Math.round(baseNewUsers * randomFactor * trendFactor * weekdayFactor),
    };
  });
};

// Generate dummy stats
const generateDummyStats = (activityData: UserActivity[]): UserStats => {
  const latestData = activityData[activityData.length - 1];
  const totalUsers = 250 + Math.floor(Math.random() * 50); // Random total between 250-300

  return {
    totalUsers,
    activeUsers: latestData.activeUsers,
    newUsers: latestData.newUsers,
  };
};

const DashboardCharts: React.FC = () => {
  const [timeRange, setTimeRange] = useState('7d');
  const [startDate, setStartDate] = useState<Date>(subDays(new Date(), 7));
  const [endDate, setEndDate] = useState<Date>(new Date());
  const [userActivity, setUserActivity] = useState<UserActivity[]>([]);
  const [userStats, setUserStats] = useState<UserStats>({
    totalUsers: 0,
    activeUsers: 0,
    newUsers: 0,
  });

  useEffect(() => {
    // Generate dummy data instead of making API calls
    const dummyActivityData = generateDummyData(startDate, endDate);
    setUserActivity(dummyActivityData);
    setUserStats(generateDummyStats(dummyActivityData));
  }, [timeRange, startDate, endDate]);

  const handleTimeRangeChange = (event: SelectChangeEvent) => {
    const range = event.target.value;
    setTimeRange(range);
    const now = new Date();
    
    switch (range) {
      case '7d':
        setStartDate(subDays(now, 7));
        break;
      case '30d':
        setStartDate(subDays(now, 30));
        break;
      case '90d':
        setStartDate(subDays(now, 90));
        break;
      default:
        setStartDate(subDays(now, 7));
    }
    setEndDate(now);
  };

  const handleStartDateChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const date = new Date(event.target.value);
    if (!isNaN(date.getTime())) {
      setStartDate(date);
    }
  };

  const handleEndDateChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const date = new Date(event.target.value);
    if (!isNaN(date.getTime())) {
      setEndDate(date);
    }
  };

  const lineChartData = {
    labels: userActivity.map(item => format(new Date(item.date), 'MMM d')),
    datasets: [
      {
        label: 'Active Users',
        data: userActivity.map(item => item.activeUsers),
        borderColor: '#8884d8',
        backgroundColor: 'rgba(136, 132, 216, 0.5)',
        tension: 0.4,
      },
      {
        label: 'New Users',
        data: userActivity.map(item => item.newUsers),
        borderColor: '#82ca9d',
        backgroundColor: 'rgba(130, 202, 157, 0.5)',
        tension: 0.4,
      },
    ],
  };

  const barChartData = {
    labels: userActivity.map(item => format(new Date(item.date), 'MMM d')),
    datasets: [
      {
        label: 'New Users',
        data: userActivity.map(item => item.newUsers),
        backgroundColor: 'rgba(130, 202, 157, 0.8)',
      },
    ],
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'top' as const,
      },
      tooltip: {
        callbacks: {
          label: function(context: any) {
            let label = context.dataset.label || '';
            if (label) {
              label += ': ';
            }
            if (context.parsed.y !== null) {
              label += context.parsed.y.toLocaleString();
            }
            return label;
          }
        }
      }
    },
    scales: {
      y: {
        beginAtZero: true,
        ticks: {
          callback: function(value: any) {
            return value.toLocaleString();
          }
        }
      },
    },
  };

  return (
    <Box sx={{ flexGrow: 1, p: 3 }}>
      <Grid container spacing={3}>
        {/* Stats Cards */}
        <Grid item xs={12} md={4}>
          <Paper sx={{ p: 2, textAlign: 'center' }}>
            <Typography variant="h6" color="text.secondary">
              Total Users
            </Typography>
            <Typography variant="h4">{userStats.totalUsers.toLocaleString()}</Typography>
          </Paper>
        </Grid>
        <Grid item xs={12} md={4}>
          <Paper sx={{ p: 2, textAlign: 'center' }}>
            <Typography variant="h6" color="text.secondary">
              Active Users
            </Typography>
            <Typography variant="h4">{userStats.activeUsers.toLocaleString()}</Typography>
          </Paper>
        </Grid>
        <Grid item xs={12} md={4}>
          <Paper sx={{ p: 2, textAlign: 'center' }}>
            <Typography variant="h6" color="text.secondary">
              New Users
            </Typography>
            <Typography variant="h4">{userStats.newUsers.toLocaleString()}</Typography>
          </Paper>
        </Grid>

        {/* Time Range Selector */}
        <Grid item xs={12}>
          <Box sx={{ display: 'flex', gap: 2, mb: 2 }}>
            <FormControl sx={{ minWidth: 120 }}>
              <InputLabel>Time Range</InputLabel>
              <Select
                value={timeRange}
                label="Time Range"
                onChange={handleTimeRangeChange}
              >
                <MenuItem value="7d">Last 7 Days</MenuItem>
                <MenuItem value="30d">Last 30 Days</MenuItem>
                <MenuItem value="90d">Last 90 Days</MenuItem>
                <MenuItem value="custom">Custom Range</MenuItem>
              </Select>
            </FormControl>
            {timeRange === 'custom' && (
              <>
                <TextField
                  label="Start Date"
                  type="date"
                  value={format(startDate, 'yyyy-MM-dd')}
                  onChange={handleStartDateChange}
                  InputLabelProps={{ shrink: true }}
                />
                <TextField
                  label="End Date"
                  type="date"
                  value={format(endDate, 'yyyy-MM-dd')}
                  onChange={handleEndDateChange}
                  InputLabelProps={{ shrink: true }}
                />
              </>
            )}
          </Box>
        </Grid>

        {/* User Activity Line Chart */}
        <Grid item xs={12}>
          <Paper sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
              User Activity Over Time
            </Typography>
            <Box sx={{ height: 400 }}>
              <Line data={lineChartData} options={chartOptions} />
            </Box>
          </Paper>
        </Grid>

        {/* New Users Bar Chart */}
        <Grid item xs={12}>
          <Paper sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
              New Users Distribution
            </Typography>
            <Box sx={{ height: 400 }}>
              <Bar data={barChartData} options={chartOptions} />
            </Box>
          </Paper>
        </Grid>
      </Grid>
    </Box>
  );
};

export default DashboardCharts; 