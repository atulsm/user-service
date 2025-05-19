import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Container,
  Paper,
  TextField,
  Button,
  Typography,
  Box,
  Alert,
} from '@mui/material';
import { userService, User, UpdateUserRequest, CreateUserRequest } from '../services/api';

const UserForm: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [error, setError] = useState('');
  const [formData, setFormData] = useState<CreateUserRequest | UpdateUserRequest>({
    firstName: '',
    lastName: '',
    email: '',
    phoneNumber: '',
    password: '',
  });

  useEffect(() => {
    if (id && id !== 'new') {
      loadUser();
    }
  }, [id]);

  const loadUser = async () => {
    try {
      const user = await userService.getUser(id!);
      setFormData({
        firstName: user.firstName,
        lastName: user.lastName,
        email: user.email,
        phoneNumber: user.phoneNumber || '',
      });
    } catch (err) {
      setError('Failed to load user');
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (id === 'new') {
        // For new user, ensure required fields are present
        const createData: CreateUserRequest = {
          firstName: formData.firstName || '',
          lastName: formData.lastName || '',
          email: formData.email || '',
          phoneNumber: formData.phoneNumber,
          password: (formData as CreateUserRequest).password || '',
        };
        await userService.createUser(createData);
      } else {
        await userService.updateUser(id!, formData as UpdateUserRequest);
      }
      navigate('/users');
    } catch (err) {
      setError('Failed to save user');
    }
  };

  return (
    <Container maxWidth="md">
      <Box sx={{ mt: 4 }}>
        <Paper sx={{ p: 3 }}>
          <Typography component="h1" variant="h4" gutterBottom>
            {id === 'new' ? 'Create User' : 'Edit User'}
          </Typography>
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}
          <Box component="form" onSubmit={handleSubmit}>
            <TextField
              margin="normal"
              required
              fullWidth
              id="firstName"
              label="First Name"
              name="firstName"
              value={formData.firstName}
              onChange={handleChange}
            />
            <TextField
              margin="normal"
              required
              fullWidth
              id="lastName"
              label="Last Name"
              name="lastName"
              value={formData.lastName}
              onChange={handleChange}
            />
            <TextField
              margin="normal"
              required
              fullWidth
              id="email"
              label="Email Address"
              name="email"
              autoComplete="email"
              value={formData.email}
              onChange={handleChange}
            />
            {id === 'new' && (
              <TextField
                margin="normal"
                required
                fullWidth
                id="password"
                label="Password"
                name="password"
                type="password"
                autoComplete="new-password"
                value={(formData as CreateUserRequest).password}
                onChange={handleChange}
                helperText="Password must be at least 8 characters long"
              />
            )}
            <TextField
              margin="normal"
              fullWidth
              id="phoneNumber"
              label="Phone Number"
              name="phoneNumber"
              autoComplete="tel"
              value={formData.phoneNumber}
              onChange={handleChange}
              placeholder="+1234567890"
            />
            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ mt: 3, mb: 2 }}
            >
              {id === 'new' ? 'Create' : 'Update'}
            </Button>
            <Button
              fullWidth
              variant="outlined"
              onClick={() => navigate('/users')}
              sx={{ mb: 2 }}
            >
              Cancel
            </Button>
          </Box>
        </Paper>
      </Box>
    </Container>
  );
};

export default UserForm; 