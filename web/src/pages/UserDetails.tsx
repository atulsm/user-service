import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Container,
  Box,
  Typography,
  TextField,
  Button,
  Paper,
  Alert,
} from '@mui/material';
import { userService } from '../services/api';
import { User, UpdateUserRequest } from '../types/user';

const UserDetails = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [error, setError] = useState('');
  const [user, setUser] = useState<User | null>(null);
  const [formData, setFormData] = useState<UpdateUserRequest>({
    firstName: '',
    lastName: '',
    email: '',
  });

  useEffect(() => {
    const loadUser = async () => {
      try {
        if (id) {
          const data = await userService.getUser(id);
          setUser(data);
          setFormData({
            firstName: data.firstName,
            lastName: data.lastName,
            email: data.email,
          });
        }
      } catch (error) {
        setError('Failed to load user details');
      }
    };
    loadUser();
  }, [id]);

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
      if (id) {
        await userService.updateUser(id, formData);
        navigate('/users');
      }
    } catch (error) {
      setError('Failed to update user');
    }
  };

  if (!user) {
    return (
      <Container maxWidth="md" sx={{ mt: 4 }}>
        <Alert severity="error">User not found</Alert>
      </Container>
    );
  }

  return (
    <Container maxWidth="md" sx={{ mt: 4 }}>
      <Paper sx={{ p: 3 }}>
        <Typography component="h1" variant="h4" gutterBottom>
          Edit User
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
            type="email"
            value={formData.email}
            onChange={handleChange}
          />
          <Box sx={{ mt: 3, display: 'flex', gap: 2 }}>
            <Button
              type="submit"
              variant="contained"
              color="primary"
            >
              Save Changes
            </Button>
            <Button
              variant="outlined"
              onClick={() => navigate('/users')}
            >
              Cancel
            </Button>
          </Box>
        </Box>
      </Paper>
    </Container>
  );
};

export default UserDetails; 