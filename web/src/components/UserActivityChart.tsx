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
        '#FF6384',
        '#36A2EB',
        '#FFCE56',
        '#4BC0C0',
        '#9966FF',
      ],
      borderColor: [
        '#FF6384',
        '#36A2EB',
        '#FFCE56',
        '#4BC0C0',
        '#9966FF',
      ],
      borderWidth: 1,
    },
  ],
};

const options = {
  responsive: true,
  plugins: {
    legend: {
      position: 'right' as const,
    },
    title: {
      display: true,
      text: 'User Activity Distribution',
      font: {
        size: 16,
      },
    },
  },
};

const UserActivityChart: React.FC = () => {
  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-xl font-semibold mb-4">User Activity Overview</h2>
      <div className="h-[400px] flex items-center justify-center">
        <Pie data={userActivityData} options={options} />
      </div>
      <div className="mt-4 text-sm text-gray-600">
        <p>Total Activities: 100</p>
        <p>Last Updated: {new Date().toLocaleString()}</p>
      </div>
    </div>
  );
};

export default UserActivityChart; 