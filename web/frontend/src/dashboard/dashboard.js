import React, { useState, useEffect } from "react";
import { 
  AppBar, Toolbar, Typography, Button, Container, Box, Dialog, DialogTitle, DialogContent, DialogActions,
  Menu, MenuItem, ToggleButton, ToggleButtonGroup, IconButton, Table, TableHead, TableBody, TableRow, TableCell, TextField 
} from "@mui/material";
import AddIcon from "@mui/icons-material/Add";
import RemoveIcon from "@mui/icons-material/Remove";
import axios from "axios";
import { useNavigate } from "react-router-dom";


// ✅ Function to check token validity
const isTokenExpired = () => {
  const token = localStorage.getItem("token");
  if (!token) return true;
  try {
    const payload = JSON.parse(atob(token.split(".")[1]));
    return Date.now() >= payload.exp * 1000;
  } catch {
    return true;
  }
};

// ✅ Function to get user ID from the token
const getUserIdFromToken = () => {
  const token = localStorage.getItem("token");
  if (!token) return null;
  try {
    const payload = JSON.parse(atob(token.split(".")[1])); // Decode JWT
    return payload.sub; // Assuming 'sub' contains user_id
  } catch {
    return null;
  }
};

const Dashboard = () => {
  const [user, setUser] = useState({ name: "", role: "", email: "" });
  const [pipelines, setPipelines] = useState([]);
  const [profileOpen, setProfileOpen] = useState(false);
  const [pipelineStages, setPipelineStages] = useState(1);
  const [isParallel, setIsParallel] = useState(false);
  const [anchorEl, setAnchorEl] = useState(null);
  const navigate = useNavigate();
  const [stagesDialogOpen, setStagesDialogOpen] = useState(false);
  const [selectedPipelineStages, setSelectedPipelineStages] = useState([]);
  const [selectedPipelineId, setSelectedPipelineId] = useState(null);
  const [openStageModal, setOpenStageModal] = useState(false);

  // // Retrieve user_id from localStorage
  // const user_id = localStorage.getItem("user_id");

  // useEffect(() => {
  //   if (!user_id) {
  //     console.error("User ID not found! Redirecting to login.");
  //     navigate("/login");
  //     return;
  //   }
  //   fetchUserProfile();
  //   fetchUserPipelines();
  // }, [user_id]);

  // ✅ Get user ID from token
  const user_id = getUserIdFromToken();
  
  // ✅ Auto logout if token is expired
  useEffect(() => {
    if (isTokenExpired()) {
      console.warn("Token expired. Logging out...");
      localStorage.clear();
      navigate("/login");
      return;
    }
    fetchUserProfile();
    fetchUserPipelines();
  }, []);

  // ✅ Attach token to every API request
  const authAxios = axios.create({
    baseURL: "http://localhost:8080",
    headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
  });

  const logoutUser = () => {
    localStorage.clear();
    navigate("/login");
  };

  // const fetchUserProfile = async () => {
  //   try {
  //     const response = await axios.get(`http://localhost:8080/user/${user_id}`);
  //     if (response.data) {
  //       setUser(response.data);
  //       localStorage.setItem("user_name", response.data.name);
  //       localStorage.setItem("user_role", response.data.role);
  //     }
  //   } catch (error) {
  //     console.error("Failed to fetch user profile", error);
  //   }
  // };

  const fetchUserProfile = async () => {
    try {
      const response = await authAxios.get(`/user/${localStorage.getItem("user_id")}`);
      if (response.data) {
        setUser(response.data);
        localStorage.setItem("user_name", response.data.name);
        localStorage.setItem("user_role", response.data.role);
      }
    } catch (error) {
      console.error("Failed to fetch user profile", error);
      logoutUser();
    }
  };

  const fetchUserPipelines = async () => {
    try {
      const response = await authAxios.get(`/pipelines?user_id=${user_id}`);
      console.log("Fetched Pipelines:", response.data);
      if (Array.isArray(response.data)) {
        setPipelines(response.data);
      } else {
        console.error("Unexpected response format:", response.data);
      }
    } catch (error) {
      console.error("Failed to fetch pipelines", error);
    }
  };

  const handleCreatePipeline = async () => {
    try {
      await authAxios.post("/createpipelines", {
        stages: pipelineStages,
        is_parallel: isParallel,
        user_id: user_id,
      });
      fetchUserPipelines();
    } catch (error) {
      console.error("Failed to create pipeline", error);
    }
  };

  const handlePipelineAction = async (pipelineId, status) => {
    try {
      if (status === "Running") {
        await authAxios.post(`/pipelines/${pipelineId}/cancel`, {
          user_id: user_id,
          is_parallel: isParallel,
        });
      } else if (status === "Completed") {
        alert("Completed pipelines cannot be started again.");
        return;
      } else {
        await authAxios.post(`/pipelines/${pipelineId}/start`, {
          user_id: user_id,
          input: { raw_material: "Steel", quantity: 100 },
          is_parallel: isParallel,
        });
      }
      // fetchUserPipelines();
      // 🚀 **Update state instantly & reload status**
      setPipelines((prevPipelines) =>
        prevPipelines.map((pipeline) =>
          pipeline.PipelineID === pipelineId ? { ...pipeline, Status: "Running" } : pipeline
        )
      );
      setTimeout(fetchUserPipelines, 1000); // ✅ Refresh status after 1 sec
    } catch (error) {
      console.error("Failed to update pipeline status", error);
      // logoutUser();
    }
  };

  const fetchPipelineStages = async (pipelineId) => {
    try {
      console.log(`Fetching stages for pipeline: ${pipelineId}`); // ✅ Debugging log
  
      const response = await authAxios.get(`/pipelines/${pipelineId}/stages`);
      
      console.log("Stages Data:", response.data); // ✅ Check API response
  
      if (Array.isArray(response.data)) {
        setSelectedPipelineStages(response.data);
        setOpenStageModal(true); // ✅ Ensure modal opens
      } else {
        console.error("Unexpected response format:", response.data);
      }
    } catch (error) {
      console.error("Failed to fetch pipeline stages:", error);
      logoutUser();
    }
  };

  const handleProfileSave = async () => {
    try {
      await authAxios.put(`/user/${user_id}`, {
        name: user.name,
        role: user.role,
      });
      setProfileOpen(false);
    } catch (error) {
      console.error("Failed to update profile", error);
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

      {/* Pipelines Table */}
      <Box sx={{ mt: 5, p: 3, boxShadow: 3, borderRadius: 2 }}>
        <Typography variant="h5" sx={{ mb: 2 }}>Your Pipelines</Typography>
        {pipelines.length > 0 ? (
          <Table sx={{ minWidth: 500, border: "1px solid #ccc", borderRadius: "8px" }}>
            <TableHead>
              <TableRow sx={{ backgroundColor: "#f5f5f5" }}>
                <TableCell><strong>Pipeline ID</strong></TableCell>
                <TableCell><strong>Status</strong></TableCell>
                <TableCell><strong>Actions</strong></TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {pipelines.map((pipeline) => (
                <TableRow key={pipeline.PipelineID}>
                  <TableCell>{pipeline.PipelineID}</TableCell>
                  <TableCell>
                    <Typography sx={{ fontWeight: "bold", color: pipeline.Status === "Running" ? "green" : "gray" }}>
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
                    {pipeline.Status !== "Created" && (
                      <Button
                      variant="outlined"
                      sx={{ ml: 2 }}
                      onClick={() => {
                        console.log("Show Stages button clicked for pipeline:", pipeline.PipelineID); // ✅ Debug log
                        fetchPipelineStages(pipeline.PipelineID);
                      }}
                    >
                      Show Stages
                    </Button>
                    
                    )}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        ) : (
          <Typography>No pipelines created.</Typography>
        )}
      </Box>
      <Dialog open={openStageModal} onClose={() => setOpenStageModal(false)}>
  <DialogTitle>Pipeline Stages</DialogTitle>
  <DialogContent>
    {selectedPipelineStages.length > 0 ? (
      <Table>
        <TableHead>
          <TableRow>
            <TableCell><strong>Stage ID</strong></TableCell>
            <TableCell><strong>Status</strong></TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {selectedPipelineStages.map((stage) => (
            <TableRow key={stage.StageID}>
              <TableCell>{stage.StageID}</TableCell>
              <TableCell>{stage.Status}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    ) : (
      <Typography>No stages found for this pipeline.</Typography>
    )}
  </DialogContent>
  <DialogActions>
    <Button onClick={() => setOpenStageModal(false)}>Close</Button>
  </DialogActions>
</Dialog>

{/* Edit Profile Dialog */}
<Dialog open={profileOpen} onClose={() => setProfileOpen(false)}>
        <DialogTitle>Edit Profile</DialogTitle>
        <DialogContent>
          <TextField label="Name" fullWidth margin="normal" value={user.name} onChange={(e) => setUser({ ...user, name: e.target.value })} />
          <TextField label="Role" fullWidth margin="normal" value={user.role} onChange={(e) => setUser({ ...user, role: e.target.value })} />
          <TextField label="Email" fullWidth margin="normal" value={user.email} disabled />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setProfileOpen(false)}>Cancel</Button>
          <Button onClick={handleProfileSave} color="primary">Save</Button>
        </DialogActions>
      </Dialog>


      {/* Create New Pipeline */}
      {/* Create New Pipeline Section */}
    <Box sx={{ mt: 5, p: 3, boxShadow: 3, borderRadius: 2 }}>
      <Typography variant="h5" sx={{ mb: 2 }}>Create New Pipeline</Typography>

      {/* Number of Stages Title + Counter */}
      <Box sx={{ display: "flex", alignItems: "center", justifyContent: "center", gap: 3, mt: 2 }}>
        <Typography variant="h6" sx={{ minWidth: 150, textAlign: "right" }}>Number of Stages:</Typography>
        <IconButton onClick={() => setPipelineStages(Math.max(1, pipelineStages - 1))}>
          <RemoveIcon />
        </IconButton>
        <Typography variant="h6">{pipelineStages}</Typography>
        <IconButton onClick={() => setPipelineStages(pipelineStages + 1)}>
          <AddIcon />
        </IconButton>
      </Box>

      {/* Parallel/Sequential Selection & Create Button */}
      <Box sx={{ display: "flex", alignItems: "center", justifyContent: "center", gap: 3, mt: 3 }}>
        <ToggleButtonGroup value={isParallel} exclusive onChange={() => setIsParallel(!isParallel)}>
          <ToggleButton value={false}>Sequential</ToggleButton>
          <ToggleButton value={true}>Parallel</ToggleButton>
        </ToggleButtonGroup>
        
        <Button variant="contained" color="secondary" sx={{ px: 3, py: 1.2 }} onClick={handleCreatePipeline}>
          Create Pipeline
        </Button> 
      </Box>
    </Box>

    

    </Container>
  );
};

export default Dashboard;
