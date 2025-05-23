import React from 'react';
import { Pie } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  ArcElement,
  Tooltip,
  Legend,
  Title,
} from 'chart.js';

// Register ChartJS components
ChartJS.register(ArcElement, Tooltip, Legend, Title);

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

// Dummy data for user activities
const userActivityData = {
  labels: [
    'Profile Updates',
    'Password Changes',
    'Login Attempts',
    'API Calls',
    'File Uploads',
  ],
  datasets: [
    {
      data: [30, 15, 25, 20, 10],
      backgroundColor: [
        uberColors.primary,
        uberColors.success,
        uberColors.warning,
        uberColors.info,
        uberColors.error,
      ],
      borderColor: uberColors.background,
      borderWidth: 2,
      hoverOffset: 4,
    },
  ],
};

const options = {
  responsive: true,
  plugins: {
    legend: {
      position: 'right' as const,
      labels: {
        font: {
          family: 'UberMove, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
          size: 12,
        },
        color: uberColors.gray[800],
        padding: 20,
      },
    },
    title: {
      display: true,
      text: 'User Activity Distribution',
      font: {
        family: 'UberMove, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
        size: 16,
        weight: 'normal' as const,
      },
      color: uberColors.primary,
      padding: {
        top: 10,
        bottom: 20,
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
      usePointStyle: true,
      callbacks: {
        label: function(context: any) {
          return `${context.label}: ${context.raw}%`;
        }
      }
    },
  },
};

const UserActivityChart: React.FC = () => {
  return (
    <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
      <h2 className="text-xl font-medium mb-6 text-gray-900">User Activity Overview</h2>
      <div className="h-[400px] flex items-center justify-center">
        <Pie data={userActivityData} options={options} />
      </div>
      <div className="mt-6 text-sm text-gray-600 border-t border-gray-100 pt-4">
        <div className="flex justify-between items-center">
          <p>Total Activities: 100</p>
          <p>Last Updated: {new Date().toLocaleString()}</p>
        </div>
      </div>
    </div>
  );
};

export default UserActivityChart; 