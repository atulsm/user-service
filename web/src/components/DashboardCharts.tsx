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
import UserActivityChart from './UserActivityChart';

// Uber's color palette
const uberColors = {
  primary: '#000000',
  secondary: '#6B6B6B',
  accent: '#000000',
  background: '#FFFFFF',
  success: '#00A870',
  warning: '#FFB800',
  error: '#FF3B30',
  info: '#007AFF',
  gray: {
    100: '#F7F7F7',
    200: '#E5E5E5',
    300: '#D4D4D4',
    400: '#A3A3A3',
    500: '#737373',
    600: '#525252',
    700: '#404040',
    800: '#262626',
    900: '#171717',
  }
};

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
        borderColor: uberColors.primary,
        backgroundColor: 'rgba(0, 0, 0, 0.1)',
        tension: 0.4,
        borderWidth: 2,
      },
      {
        label: 'New Users',
        data: userActivity.map(item => item.newUsers),
        borderColor: uberColors.success,
        backgroundColor: 'rgba(0, 168, 112, 0.1)',
        tension: 0.4,
        borderWidth: 2,
      },
    ],
  };

  const barChartData = {
    labels: userActivity.map(item => format(new Date(item.date), 'MMM d')),
    datasets: [
      {
        label: 'New Users',
        data: userActivity.map(item => item.newUsers),
        backgroundColor: uberColors.primary,
        borderRadius: 4,
      },
    ],
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'top' as const,
        labels: {
          font: {
            family: 'UberMove, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
            size: 12,
          },
          color: uberColors.gray[800],
          padding: 20,
        },
      },
      tooltip: {
        backgroundColor: uberColors.background,
        titleColor: uberColors.primary,
        bodyColor: uberColors.gray[800],
        borderColor: uberColors.gray[200],
        borderWidth: 1,
        padding: 12,
        boxPadding: 6,
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
        grid: {
          color: uberColors.gray[100],
        },
        ticks: {
          color: uberColors.gray[600],
          font: {
            family: 'UberMove, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
            size: 12,
          },
          callback: function(value: any) {
            return value.toLocaleString();
          }
        }
      },
      x: {
        grid: {
          display: false,
        },
        ticks: {
          color: uberColors.gray[600],
          font: {
            family: 'UberMove, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
            size: 12,
          },
        }
      },
    },
  };

  return (
    <Box sx={{ flexGrow: 1, p: 3, backgroundColor: uberColors.gray[100] }}>
      <Grid container spacing={3}>
        {/* Stats Cards */}
        <Grid item xs={12} md={4}>
          <Paper sx={{ 
            p: 3, 
            textAlign: 'center',
            borderRadius: '12px',
            boxShadow: 'none',
            border: `1px solid ${uberColors.gray[200]}`,
            backgroundColor: uberColors.background,
          }}>
            <Typography variant="h6" sx={{ color: uberColors.gray[600], mb: 1 }}>
              Total Users
            </Typography>
            <Typography variant="h4" sx={{ color: uberColors.primary, fontWeight: 500 }}>
              {userStats.totalUsers.toLocaleString()}
            </Typography>
          </Paper>
        </Grid>
        <Grid item xs={12} md={4}>
          <Paper sx={{ 
            p: 3, 
            textAlign: 'center',
            borderRadius: '12px',
            boxShadow: 'none',
            border: `1px solid ${uberColors.gray[200]}`,
            backgroundColor: uberColors.background,
          }}>
            <Typography variant="h6" sx={{ color: uberColors.gray[600], mb: 1 }}>
              Active Users
            </Typography>
            <Typography variant="h4" sx={{ color: uberColors.primary, fontWeight: 500 }}>
              {userStats.activeUsers.toLocaleString()}
            </Typography>
          </Paper>
        </Grid>
        <Grid item xs={12} md={4}>
          <Paper sx={{ 
            p: 3, 
            textAlign: 'center',
            borderRadius: '12px',
            boxShadow: 'none',
            border: `1px solid ${uberColors.gray[200]}`,
            backgroundColor: uberColors.background,
          }}>
            <Typography variant="h6" sx={{ color: uberColors.gray[600], mb: 1 }}>
              New Users
            </Typography>
            <Typography variant="h4" sx={{ color: uberColors.primary, fontWeight: 500 }}>
              {userStats.newUsers.toLocaleString()}
            </Typography>
          </Paper>
        </Grid>

        {/* User Activity Pie Chart */}
        <Grid item xs={12} md={6}>
          <UserActivityChart />
        </Grid>

        {/* Time Range Selector */}
        <Grid item xs={12}>
          <Paper sx={{ 
            p: 2, 
            borderRadius: '12px',
            boxShadow: 'none',
            border: `1px solid ${uberColors.gray[200]}`,
            backgroundColor: uberColors.background,
          }}>
            <Box sx={{ display: 'flex', gap: 2 }}>
              <FormControl sx={{ minWidth: 120 }}>
                <InputLabel sx={{ color: uberColors.gray[600] }}>Time Range</InputLabel>
                <Select
                  value={timeRange}
                  label="Time Range"
                  onChange={handleTimeRangeChange}
                  sx={{
                    '& .MuiOutlinedInput-notchedOutline': {
                      borderColor: uberColors.gray[200],
                    },
                    '&:hover .MuiOutlinedInput-notchedOutline': {
                      borderColor: uberColors.gray[300],
                    },
                    '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                      borderColor: uberColors.primary,
                    },
                  }}
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
                    InputLabelProps={{ 
                      shrink: true,
                      sx: { color: uberColors.gray[600] }
                    }}
                    sx={{
                      '& .MuiOutlinedInput-notchedOutline': {
                        borderColor: uberColors.gray[200],
                      },
                      '&:hover .MuiOutlinedInput-notchedOutline': {
                        borderColor: uberColors.gray[300],
                      },
                      '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                        borderColor: uberColors.primary,
                      },
                    }}
                  />
                  <TextField
                    label="End Date"
                    type="date"
                    value={format(endDate, 'yyyy-MM-dd')}
                    onChange={handleEndDateChange}
                    InputLabelProps={{ 
                      shrink: true,
                      sx: { color: uberColors.gray[600] }
                    }}
                    sx={{
                      '& .MuiOutlinedInput-notchedOutline': {
                        borderColor: uberColors.gray[200],
                      },
                      '&:hover .MuiOutlinedInput-notchedOutline': {
                        borderColor: uberColors.gray[300],
                      },
                      '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                        borderColor: uberColors.primary,
                      },
                    }}
                  />
                </>
              )}
            </Box>
          </Paper>
        </Grid>

        {/* Line Chart */}
        <Grid item xs={12} md={8}>
          <Paper sx={{ 
            p: 3, 
            height: 400,
            borderRadius: '12px',
            boxShadow: 'none',
            border: `1px solid ${uberColors.gray[200]}`,
            backgroundColor: uberColors.background,
          }}>
            <Typography variant="h6" sx={{ 
              color: uberColors.primary,
              mb: 3,
              fontWeight: 500,
            }}>
              User Activity Over Time
            </Typography>
            <Line data={lineChartData} options={chartOptions} />
          </Paper>
        </Grid>

        {/* Bar Chart */}
        <Grid item xs={12} md={4}>
          <Paper sx={{ 
            p: 3, 
            height: 400,
            borderRadius: '12px',
            boxShadow: 'none',
            border: `1px solid ${uberColors.gray[200]}`,
            backgroundColor: uberColors.background,
          }}>
            <Typography variant="h6" sx={{ 
              color: uberColors.primary,
              mb: 3,
              fontWeight: 500,
            }}>
              New User Registrations
            </Typography>
            <Bar data={barChartData} options={chartOptions} />
          </Paper>
        </Grid>
      </Grid>
    </Box>
  );
};

export default DashboardCharts; 