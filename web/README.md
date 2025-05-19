# User Service Web UI

A React + TypeScript web application for managing users, built with Material-UI.

## Features

- User authentication (login/register)
- List all users
- Create new users
- Edit existing users
- Delete users
- Responsive design

## Prerequisites

- Node.js (v14 or higher)
- npm (v6 or higher)

## Installation

1. Install dependencies:
   ```bash
   npm install
   ```

2. Create a `.env` file in the root directory with the following content:
   ```
   REACT_APP_API_URL=http://localhost:8080/api/v1
   ```

## Development

To start the development server:

```bash
npm start
```

The application will be available at `http://localhost:3000`.

## Building for Production

To create a production build:

```bash
npm run build
```

The build artifacts will be stored in the `build/` directory.

## Project Structure

```
src/
  ├── components/     # Reusable components
  ├── pages/         # Page components
  ├── services/      # API services
  ├── types/         # TypeScript type definitions
  └── utils/         # Utility functions
```

## API Integration

The web UI communicates with the User Service API endpoints:

- Authentication:
  - POST /api/v1/auth/login
  - POST /api/v1/auth/register

- User Management:
  - GET /api/v1/users
  - GET /api/v1/users/:id
  - POST /api/v1/users
  - PUT /api/v1/users/:id
  - DELETE /api/v1/users/:id

## Technologies Used

- React
- TypeScript
- Material-UI
- React Router
- Axios
