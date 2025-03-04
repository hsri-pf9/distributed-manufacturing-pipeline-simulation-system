import React, { useState } from "react";
import { BrowserRouter as Router, Routes, Route, Link, useNavigate } from "react-router-dom";
import { TextField, Button, Select, MenuItem, Typography, Container, Box, AppBar, Toolbar } from "@mui/material";
import axios from "axios";
import Dashboard from "../dashboard/dashboard";


const AuthLayout = ({ children, title }) => {
  return (
    <Container maxWidth="sm">
      <Box sx={{ textAlign: "center", mt: 5, p: 4, boxShadow: 3, borderRadius: 2 }}>
        <Typography variant="h4" sx={{ mb: 2 }}>{title}</Typography>
        {children}
      </Box>
    </Container>
  );
};

const RegisterPage = ({ apiType }) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const navigate = useNavigate();

  const handleRegister = async () => {
    setMessage("");
    if (apiType === "rest") {
      try {
        await axios.post("http://localhost:8080/register", { email, password });
        setMessage("Registration successful! Please check your email to verify.");
        window.open("https://mail.google.com", "_blank");
      } catch {
        setMessage("Registration failed. Please try again.");
      }
    } else {
      console.log(`Run this gRPC command manually: grpcurl -plaintext -d '{"email": "${email}", "password": "${password}"}' localhost:50051 auth.AuthService/Register`);
      setMessage("Open your email and click the link to authenticate.");
    }
  };

  return (
    <AuthLayout title="Register">
      <TextField label="Email" fullWidth margin="normal" value={email} onChange={(e) => setEmail(e.target.value)} />
      <TextField label="Password" type="password" fullWidth margin="normal" value={password} onChange={(e) => setPassword(e.target.value)} />
      <Button variant="contained" color="primary" fullWidth onClick={handleRegister} sx={{ mt: 2 }}>Register</Button>
      {message && <Typography sx={{ mt: 2, color: "green" }}>{message}</Typography>}
      <Typography sx={{ mt: 2 }}>Already registered? <Link to="/login">Login</Link></Typography>
    </AuthLayout>
  );
};

const LoginPage = ({ apiType }) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const navigate = useNavigate();

  const handleLogin = async () => {
    setMessage("");
    if (apiType === "rest") {
      try {
        const response = await axios.post(
          "http://localhost:8080/login",
          { email, password },
          {
            headers: { "Content-Type": "application/json" },
            withCredentials: true, // ✅ Ensure cross-origin cookies are included
          }
        );
  
        // ✅ Extract user_id, email, and token correctly
        const { user_id, token } = response.data;
        console.log("Full Response from Backend:", response.data);
        if (!user_id || !email || !token) throw new Error("Invalid response from backend");
  
        console.log("User ID:", user_id);
  
        // ✅ Store user_id and token for session persistence
        localStorage.setItem("user_id", user_id);
        localStorage.setItem("email", email);
        localStorage.setItem("token", token);
  
        navigate("/dashboard");
      } catch (error) {
        console.error("Login failed:", error);
        setMessage("Login failed. Please check your credentials.");
      }
    } else {
      console.log(`Run this gRPC command manually: grpcurl -plaintext -d '{"email": "${email}", "password": "${password}"}' localhost:50051 auth.AuthService/Login`);
      setMessage("Check console for gRPC login command.");
    }
  };
  

  return (
    <AuthLayout title="Login">
      <TextField label="Email" fullWidth margin="normal" value={email} onChange={(e) => setEmail(e.target.value)} />
      <TextField label="Password" type="password" fullWidth margin="normal" value={password} onChange={(e) => setPassword(e.target.value)} />
      <Button variant="contained" color="primary" fullWidth onClick={handleLogin} sx={{ mt: 2 }}>Login</Button>
      {message && <Typography sx={{ mt: 2, color: "red" }}>{message}</Typography>}
      <Typography sx={{ mt: 2 }}>Do not have an account? <Link to="/register">Register</Link></Typography>
    </AuthLayout>
  );
};

const App = () => {
  const [apiType, setApiType] = useState("rest");

  return (
    <Router>
      <AppBar position="static">
        <Toolbar sx={{ display: "flex", justifyContent: "space-between" }}>
          <Typography variant="h6">Distributed Manufacturing System</Typography>
          <Select value={apiType} onChange={(e) => setApiType(e.target.value)} sx={{ color: "white", backgroundColor: "gray" }}>
            <MenuItem value="rest">REST API</MenuItem>
            <MenuItem value="grpc">gRPC</MenuItem>
          </Select>
        </Toolbar>
      </AppBar>
      <Routes>
        <Route path="/register" element={<RegisterPage apiType={apiType} />} />
        <Route path="/login" element={<LoginPage apiType={apiType} />} />
        <Route path="/dashboard/*" element={<Dashboard />} />
        <Route path="/" element={<RegisterPage apiType={apiType} />} />
      </Routes>
    </Router>
  );
};

export default App;