import React from "react";
import { Routes, Route, useNavigate } from "react-router-dom";
import { Container, Box, Typography, Button } from "@mui/material";

const DashboardHome = () => (
  <Typography variant="h4">Welcome to the Dashboard</Typography>
);

const Dashboard = () => {
  const navigate = useNavigate();

  const handleLogout = () => {
    navigate("/login");
  };

  return (
    <Container maxWidth="md">
      <Box sx={{ textAlign: "center", mt: 5, p: 4, boxShadow: 3, borderRadius: 2 }}>
        <Typography variant="h3">Dashboard</Typography>
        <Button variant="contained" color="secondary" sx={{ mt: 3 }} onClick={handleLogout}>
          Logout
        </Button>
        <Routes>
          <Route path="/" element={<DashboardHome />} />
        </Routes>
      </Box>
    </Container>
  );
};

export default Dashboard;
