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
  const [sseEventSource, setSseEventSource] = useState(null);

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

  // ✅ Axios instance with dynamic token attachment
  const authAxios = axios.create({
    // baseURL: "http://localhost:8080",
    baseURL: "http://localhost:30081",
  });

  authAxios.interceptors.request.use(
    (config) => {
      const token = localStorage.getItem("token");
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      } else {
        console.warn("⚠️ No token found, request may fail.");
      }
      return config;
    },
    (error) => Promise.reject(error)
  );

  // ✅ Handle Unauthorized (401) Responses
  authAxios.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response && error.response.status === 401) {
        console.warn("🔴 Token expired or invalid, logging out...");
        localStorage.clear();
        navigate("/login");
      }
      return Promise.reject(error);
    }
  );

  const logoutUser = () => {
    localStorage.clear();
    navigate("/login");
  };

  const fetchUserProfile = async () => {
    try {
      // const response = await authAxios.get(`/user/${localStorage.getItem("user_id")}`);
      const response = await authAxios.get(`/user/${user_id}`);
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
         // 🔹 Initialize all stages as "Pending" first
        const initializedStages = response.data.map(stage => ({
          StageID: stage.StageID,
          Status: "Pending", // Assume "Pending" initially
        }));

        setSelectedPipelineStages(initializedStages);
        setSelectedPipelineId(pipelineId);
        setOpenStageModal(true); // ✅ Ensure modal opens
        setupSSE(pipelineId);
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
  
  // ✅ SSE with Token in Query Params (Recommended)
  const setupSSE = (pipelineId) => {
    if (sseEventSource) {
      sseEventSource.close();
    }

    const token = localStorage.getItem("token");
    // const eventSource = new EventSource(`http://localhost:8080/pipelines/${pipelineId}/stream?token=${token}`);
    const eventSource = new EventSource(`http://localhost:30081/pipelines/${pipelineId}/stream?token=${token}`);

    eventSource.onmessage = (event) => {
      const eventData = JSON.parse(event.data);
      console.log("🔄 SSE Event Received:", eventData);

      if (!eventData.status) {
        console.warn("⚠️ Invalid SSE Event:", eventData);
        return;
      }
      
      if (eventData.type === "pipeline") {
        setPipelines((prevPipelines) =>
          prevPipelines.map((pipeline) =>
            pipeline.PipelineID === eventData.pipeline_id ? { ...pipeline, Status: eventData.status } : pipeline
          )
        );
      }
      if (eventData.type === "stage") {
        setSelectedPipelineStages((prevStages) => {
            const existingStageIndex = prevStages.findIndex(stage => stage.StageID === eventData.stage_id);

            if (existingStageIndex !== -1) {
                // ✅ Update stage status (Pending → Running → Completed)
                return prevStages.map(stage =>
                    stage.StageID === eventData.stage_id ? { ...stage, Status: eventData.status } : stage
                );
            } else {
                // ✅ If this is the first stage and no stages exist, add it immediately
                if (prevStages.length === 0) {
                    console.log("🚀 First stage detected! Adding immediately as Running.");
                    return [{ StageID: eventData.stage_id, Status: eventData.status }];
                }

                // ✅ Add new stage dynamically
                return [...prevStages, { StageID: eventData.stage_id, Status: eventData.status }];
            }
        });
    }
    };

    eventSource.onerror = (error) => {
      console.error("❌ SSE Connection Error:", error);
      eventSource.close();
    };

    setSseEventSource(eventSource);
  };

  // ✅ Cleanup SSE on unmount
  useEffect(() => {
    return () => {
      if (sseEventSource) {
        console.log("❌ Closing SSE connection...");
        sseEventSource.close();
      }
    };
  }, []);

  return (
    <Container maxWidth="md">
      <AppBar position="static">
      <Toolbar sx={{ display: "flex", justifyContent: "space-between" }}>
        <Typography variant="h6">Dashboard</Typography>

        <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
          {/* Profile Menu */}
          <Button color="inherit" onClick={(e) => setAnchorEl(e.currentTarget)}>Profile</Button>
          <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={() => setAnchorEl(null)}>
            <MenuItem onClick={() => setProfileOpen(true)}>Edit Profile</MenuItem>
          </Menu>

          {/* 🚀 Logout Button */}
          <Button color="secondary" variant="contained" onClick={logoutUser}>
            Logout
          </Button>
        </Box>
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
                    {/* <Button variant="outlined" onClick={() => setOpenStageModal(true)}>Show Stages</Button> */}
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
              <TableCell>
                      <Typography sx={{ fontWeight: "bold", color: stage.Status === "Running" ? "blue" : stage.Status === "Completed" ? "green" : "gray" }}>
                        {stage.Status}
                      </Typography>
                    </TableCell>
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
