import React, { useState } from "react";
import { TextField, Button, Select, MenuItem, Typography, Container, Box } from "@mui/material";
import axios from "axios";

const AuthPage = () => {
  const [apiType, setApiType] = useState("rest");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");

  const handleRegister = async () => {
    setMessage("");
    if (apiType === "rest") {
      try {
        const response = await axios.post("http://localhost:8080/register", {
          email,
          password,
        });
        setMessage(response.data.message);
      } catch (error) {
        setMessage("Registration failed. Please try again.");
      }
    } else {
      const grpcurl = `grpcurl -plaintext -d '{"email": "${email}", "password": "${password}"}' localhost:50051 auth.AuthService/Register`;
      console.log("Run this gRPC command manually:", grpcurl);
      setMessage("Open your email and click the link to authenticate.");
    }
  };

  const handleLogin = async () => {
    setMessage("");
    if (apiType === "rest") {
      try {
        const response = await axios.post("http://localhost:8080/login", {
          email,
          password,
        });
        setMessage("Login successful!");
      } catch (error) {
        setMessage("Login failed. Please check your credentials.");
      }
    } else {
      const grpcurl = `grpcurl -plaintext -d '{"email": "${email}", "password": "${password}"}' localhost:50051 auth.AuthService/Login`;
      console.log("Run this gRPC command manually:", grpcurl);
      setMessage("Check console for gRPC login command.");
    }
  };

  return (
    <Container maxWidth="sm">
      <Box sx={{ textAlign: "center", mt: 5 }}>
        <Typography variant="h4">Select API Type</Typography>
        <Select value={apiType} onChange={(e) => setApiType(e.target.value)} fullWidth>
          <MenuItem value="rest">REST API</MenuItem>
          <MenuItem value="grpc">gRPC</MenuItem>
        </Select>
        <Typography variant="h5" sx={{ mt: 3 }}>Authentication</Typography>
        <TextField label="Email" fullWidth margin="normal" value={email} onChange={(e) => setEmail(e.target.value)} />
        <TextField label="Password" type="password" fullWidth margin="normal" value={password} onChange={(e) => setPassword(e.target.value)} />
        <Button variant="contained" color="primary" fullWidth onClick={handleRegister} sx={{ mt: 2 }}>Register</Button>
        <Button variant="contained" color="secondary" fullWidth onClick={handleLogin} sx={{ mt: 2 }}>Login</Button>
        {message && <Typography sx={{ mt: 2, color: "green" }}>{message}</Typography>}
      </Box>
    </Container>
  );
};

export default AuthPage;
