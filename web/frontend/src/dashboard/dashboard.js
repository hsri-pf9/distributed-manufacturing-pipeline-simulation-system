import React, { useState, useEffect } from "react";
import { AppBar, Toolbar, Typography, Button, Container, Box, Dialog, DialogTitle, DialogContent, TextField, DialogActions, Menu, MenuItem, ToggleButton, ToggleButtonGroup, IconButton } from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import RemoveIcon from "@mui/icons-material/Remove";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { Table, TableHead, TableBody, TableRow, TableCell } from "@mui/material";


const Dashboard = () => {
  const [user, setUser] = useState({ name: "", role: "", email: "" });
  const [pipelines, setPipelines] = useState([]);
  const [profileOpen, setProfileOpen] = useState(false);
  const [pipelineStages, setPipelineStages] = useState(1);
  const [isParallel, setIsParallel] = useState(false);
  const [anchorEl, setAnchorEl] = useState(null);
  const navigate = useNavigate();

  // Retrieve user_id from localStorage
  const user_id = localStorage.getItem("user_id");

  useEffect(() => {
    if (!user_id) {
      console.error("User ID not found! Redirecting to login.");
      navigate("/login");
      return;
    }
    fetchUserProfile();
    fetchUserPipelines();
  }, [user_id]);

  const fetchUserProfile = async () => {
    try {
      const response = await axios.get(`http://localhost:8080/user/${user_id}`);
      if (response.data) {
        setUser(response.data); // ✅ Always update UI state with fresh backend data
        localStorage.setItem("user_name", response.data.name); // ✅ Store updated name in localStorage
        localStorage.setItem("user_role", response.data.role);
      }
    } catch (error) {
      console.error("Failed to fetch user profile", error);
    }
  };
  

  const fetchUserPipelines = async () => {
    try {
      const response = await axios.get(`http://localhost:8080/pipelines?user_id=${user_id}`);
      
      console.log("Fetched Pipelines:", response.data); // ✅ Debugging log
      
      if (Array.isArray(response.data)) {
        setPipelines(response.data);
      } else {
        console.error("Unexpected response format:", response.data);
      }
    } catch (error) {
      console.error("Failed to fetch pipelines", error);
    }
  };
  

  const handleProfileSave = async () => {
    try {
      await axios.put(`http://localhost:8080/user/${user_id}`, {
        name: user.name,
        role: user.role,
      });
      setProfileOpen(false);
    } catch (error) {
      console.error("Failed to update profile", error);
    }
  };

  const handleCreatePipeline = async () => {
    try {
      await axios.post("http://localhost:8080/createpipelines", {
        stages: pipelineStages,
        is_parallel: isParallel,
        user_id: user_id,
      });
      fetchUserPipelines();
    } catch (error) {
      console.error("Failed to create pipeline", error);
    }
  };

  return (
    <Container maxWidth="md">
      <AppBar position="static">
        <Toolbar sx={{ display: "flex", justifyContent: "space-between" }}>
          <Typography variant="h6">Dashboard</Typography>
          <Button color="inherit" onClick={(e) => setAnchorEl(e.currentTarget)}>Profile</Button>
          <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={() => setAnchorEl(null)}>
            <MenuItem onClick={() => setProfileOpen(true)}>Edit Profile</MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>
      <Box sx={{ textAlign: "center", mt: 5, p: 4, boxShadow: 3, borderRadius: 2 }}>
        <Typography variant="h4">Welcome, {user.name || "User"}</Typography>
      </Box>



      <Box sx={{ mt: 5, p: 3, boxShadow: 3, borderRadius: 2 }}>
        <Typography variant="h5" sx={{ mb: 2 }}>
          Your Pipelines
        </Typography>

        {pipelines.length > 0 ? (
          <Table sx={{ minWidth: 500, border: "1px solid #ccc", borderRadius: "8px" }}>
            <TableHead>
              <TableRow sx={{ backgroundColor: "#f5f5f5" }}>
                <TableCell><strong>Pipeline ID</strong></TableCell>
                <TableCell><strong>Status</strong></TableCell>
                <TableCell><strong>Action</strong></TableCell>
              </TableRow>
            </TableHead>

            <TableBody>
              {pipelines.map((pipeline) => (
                <TableRow key={pipeline.PipelineID}>
                  <TableCell>{pipeline.PipelineID}</TableCell>
                  <TableCell>
                    <Typography 
                      sx={{
                        fontWeight: "bold",
                        color: pipeline.Status === "Running" ? "green" : "gray",
                      }}
                    >
                      {pipeline.Status}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Button
                      variant="contained"
                      color={pipeline.Status === "Running" ? "error" : "primary"}
                      onClick={() => handlePipelineAction(pipeline.PipelineID, pipeline.Status)}
                    >
                      {pipeline.Status === "Running" ? "Cancel Pipeline" : "Start Pipeline"}
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        ) : (
          <Typography>No pipelines created.</Typography>
        )}
      </Box>



      <Box sx={{ mt: 5, p: 3, boxShadow: 3, borderRadius: 2 }}>
        <Typography variant="h5">Create New Pipeline</Typography>
        <Box sx={{ display: "flex", alignItems: "center", gap: 2, mt: 2 }}>
          <IconButton onClick={() => setPipelineStages(Math.max(1, pipelineStages - 1))}>
            <RemoveIcon />
          </IconButton>
          <Typography variant="h6">{pipelineStages}</Typography>
          <IconButton onClick={() => setPipelineStages(pipelineStages + 1)}>
            <AddIcon />
          </IconButton>
        </Box>
        <Box sx={{ display: "flex", alignItems: "center", gap: 2, mt: 2 }}>
          <ToggleButtonGroup
            value={isParallel}
            exclusive
            onChange={() => setIsParallel(!isParallel)}
            sx={{ border: "1px solid #ccc", borderRadius: "8px", p: 1 }}
          >
            <ToggleButton value={false} selected={!isParallel}>Sequential</ToggleButton>
            <ToggleButton value={true} selected={isParallel}>Parallel</ToggleButton>
          </ToggleButtonGroup>
          <Button variant="contained" color="secondary" onClick={handleCreatePipeline}>
            Create Pipeline
          </Button>
        </Box>
      </Box>

      <Dialog open={profileOpen} onClose={() => setProfileOpen(false)}>
        <DialogTitle>Edit Profile</DialogTitle>
        <DialogContent>
          <TextField
            label="Name"
            fullWidth
            margin="normal"
            value={user.name}
            onChange={(e) => setUser({ ...user, name: e.target.value })}
          />
          <TextField
            label="Role"
            fullWidth
            margin="normal"
            value={user.role}
            onChange={(e) => setUser({ ...user, role: e.target.value })}
          />
          <TextField
            label="Email"
            fullWidth
            margin="normal"
            value={user.email}
            disabled
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setProfileOpen(false)}>Cancel</Button>
          <Button onClick={handleProfileSave} color="primary">Save</Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default Dashboard;
