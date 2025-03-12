import React, { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, Link, useNavigate } from "react-router-dom";
import { TextField, Button, Select, MenuItem, Typography, Container, Box, AppBar, Toolbar } from "@mui/material";
import axios from "axios";
import Dashboard from "../dashboard/dashboard";

// ✅ Function to decode JWT token and extract user_id
const decodeToken = (token) => {
  try {
    const payload = JSON.parse(atob(token.split(".")[1]));
    return payload;
  } catch (error) {
    console.error("Failed to decode token:", error);
    return null;
  }
};

// ✅ Function to check if token is expired
const isTokenExpired = (token) => {
  const payload = decodeToken(token);
  if (!payload || !payload.exp) return true;
  return Date.now() >= payload.exp * 1000;
};

// ✅ Check and redirect if token is expired
const checkSession = (navigate) => {
  const token = localStorage.getItem("token");
  if (!token || isTokenExpired(token)) {
    console.warn("Token expired. Logging out...");
    localStorage.clear();
    navigate("/login");
  }
};


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
    // if (apiType === "rest") {
      try {
        await axios.post("http://localhost:8080/register", { email, password });
        setMessage("Registration successful! Please check your email to verify.");
        window.open("https://mail.google.com", "_blank");
      } catch {
        setMessage("Registration failed. Please try again.");
      }
    // } else {
    //   console.log(`Run this gRPC command manually: grpcurl -plaintext -d '{"email": "${email}", "password": "${password}"}' localhost:50051 auth.AuthService/Register`);
    //   setMessage("Open your email and click the link to authenticate.");
    // }
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

  useEffect(() => {
    checkSession(navigate);  // ✅ Auto logout if token expired
  }, []);

  const handleLogin = async () => {
    setMessage("");
    // if (apiType === "rest") {
      try {
        // const response = await axios.post(
        //   "http://localhost:8080/login",
        //   { email, password },
        //   {
        //     headers: { "Content-Type": "application/json" },
        //     withCredentials: true, // ✅ Ensure cross-origin cookies are included
        //   }
        // );
        const response = await axios.post("http://localhost:8080/login", { email, password });
  
        // // ✅ Extract user_id, email, and token correctly
        // const { user_id, token } = response.data;
        // console.log("Full Response from Backend:", response.data);
        // if (!user_id || !email || !token) throw new Error("Invalid response from backend");
  
        // console.log("User ID:", user_id);
  
        // // ✅ Store user_id and token for session persistence
        // localStorage.setItem("user_id", user_id);
        // localStorage.setItem("email", email);
        // localStorage.setItem("token", token);

        const { token } = response.data;
        if (!token) throw new Error("Token not received");

        // ✅ Store token in localStorage
        localStorage.setItem("token", token);

        // ✅ Decode token to extract user ID (JWT tokens are usually Base64 encoded)
        const payload = JSON.parse(atob(token.split(".")[1]));
        localStorage.setItem("user_id", payload.sub);
  
        navigate("/dashboard");
      } catch (error) {
        console.error("Login failed:", error);
        setMessage("Login failed. Please check your credentials.");
      }
    // } else {
    //   console.log(`Run this gRPC command manually: grpcurl -plaintext -d '{"email": "${email}", "password": "${password}"}' localhost:50051 auth.AuthService/Login`);
    //   setMessage("Check console for gRPC login command.");
    // }
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
  // const [apiType, setApiType] = useState("rest");

  return (
    <Router>
      <AppBar position="static">
        <Toolbar sx={{ display: "flex", justifyContent: "space-between" }}>
          <Typography variant="h6">Distributed Manufacturing System</Typography>
          {/* <Select value={apiType} onChange={(e) => setApiType(e.target.value)} sx={{ color: "white", backgroundColor: "gray" }}>
            <MenuItem value="rest">REST API</MenuItem>
            <MenuItem value="grpc">gRPC</MenuItem>
          </Select> */}
        </Toolbar>
      </AppBar>
      {/* <Routes>
        <Route path="/register" element={<RegisterPage apiType={apiType} />} />
        <Route path="/login" element={<LoginPage apiType={apiType} />} />
        <Route path="/dashboard/*" element={<Dashboard />} />
        <Route path="/" element={<RegisterPage apiType={apiType} />} />
      </Routes> */}
      <Routes>
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/dashboard/*" element={<Dashboard />} />
        <Route path="/" element={<RegisterPage />} />
      </Routes>
    </Router>
  );
};

export default App;