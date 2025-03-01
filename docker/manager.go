package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/shanurrahman/orchestrator/config"
	"github.com/shanurrahman/orchestrator/utils"
)

// DockerManager handles container operations
type DockerManager struct {
	cli            *client.Client
	cfg            *config.Config
	network        string
	containerStats containerStatusMap
}

// Update NewDockerManager
func NewDockerManager(cfg *config.Config) *DockerManager {
    cli, err := client.NewClientWithOpts(client.FromEnv)
    if err != nil {
        log.Printf("Error creating Docker client: %v", err)
        return nil
    }
    log.Println("Docker client initialized successfully")

    // Create Docker network if it doesn't exist
    networkName := "fabio_network"
    networks, err := cli.NetworkList(context.Background(), network.ListOptions{})
    if err != nil {
        log.Printf("Error listing networks: %v", err)
        return nil
    }

    networkExists := false
    for _, network := range networks {
        if network.Name == networkName {
            networkExists = true
            break
        }
    }

    if !networkExists {
        _, err := cli.NetworkCreate(context.Background(), networkName, network.CreateOptions{
            Driver: "bridge",
        })
        if err != nil {
            log.Printf("Error creating network: %v", err)
            return nil
        }
        log.Printf("Created Docker network: %s", networkName)
    }

    dm := &DockerManager{
        cli:     cli,
        cfg:     cfg,
        network: networkName,
        containerStats: containerStatusMap{
            statuses: make(map[string]*ContainerStatus),
        },
    }

    // Start the event listener with a background context
    dm.StartEventListener(context.Background())

    return dm
}

// Add these new methods
func (dm *DockerManager) GetContainerStatus(id string) *ContainerStatus {
    dm.containerStats.RLock()
    defer dm.containerStats.RUnlock()
    return dm.containerStats.statuses[id]
}

// Add this type definition near other types
type ContainerConfig struct {
    ImageID     string     `json:"imageId"`
    VNCConfig   config.VNCConfig  `json:"vncConfig,omitempty"`
}

// Update CreateContainerAsync to accept VNC configuration
func (dm *DockerManager) CreateContainerAsync(configObj ContainerConfig) (string, error) {
    // Find the requested image
    var selectedImage string
    for _, img := range availableImages {
        if img.ID == configObj.ImageID {
            selectedImage = img.Name
            break
        }
    }
    if selectedImage == "" {
        return "", fmt.Errorf("invalid image ID: %s", configObj.ImageID)
    }

    dm.containerStats.Lock()
    tempID := utils.GenerateID()[:12]
    dm.containerStats.statuses[tempID] = &ContainerStatus{
        Status:  "initializing",
        Message: "Starting container creation",
    }
    dm.containerStats.Unlock()

    go func() {
        // Fix: Use selectedImage and configObj.VNCConfig instead of undefined variables
        endpoints, err := dm.CreateContainer(selectedImage, configObj.VNCConfig)
        if err != nil {
            dm.containerStats.Lock()
            dm.containerStats.statuses[tempID].Status = "failed"
            dm.containerStats.statuses[tempID].Message = "Container creation failed"
            dm.containerStats.statuses[tempID].Error = err.Error()
            dm.containerStats.Unlock()
            return
        }

        // Instead of deleting tempID, update it with the container information
        dm.containerStats.Lock()
        dm.containerStats.statuses[tempID] = &ContainerStatus{
            Status:    "ready",
            Message:   "Container is ready",
            Endpoints: endpoints,
        }
        // Also store the status under the real container ID for future reference
        dm.containerStats.statuses[endpoints.ContainerID] = dm.containerStats.statuses[tempID]
        dm.containerStats.Unlock()
    }()

    return tempID, nil
}

func (dm *DockerManager) registerWithConsul(containerID string, containerIP string, consulAddr string) error {
    shortID := containerID[:12]
    
    // Register Chat Commands endpoint
    chatRegistration := ConsulServiceRegistration{
        Name:    fmt.Sprintf("chat-api-%s", shortID),
        ID:      fmt.Sprintf("chat-api-%s", shortID),
        Address: containerIP,
        Port:    3000,  // Changed from 8080 to 3000
        Tags:    []string{
            fmt.Sprintf("urlprefix-/%s/chat/ strip=/%s/chat/", shortID, shortID),
        },
    }
    chatRegistration.Check.HTTP = fmt.Sprintf("http://%s:3000/health", containerIP)  // Also update the health check URL
    chatRegistration.Check.Interval = "10s"

    // Register noVNC endpoint
    novncRegistration := ConsulServiceRegistration{
        Name:    fmt.Sprintf("novnc-%s", shortID),
        ID:      fmt.Sprintf("novnc-%s", shortID),
        Address: containerIP,
        Port:    6901,
        Tags:    []string{
            fmt.Sprintf("urlprefix-/%s/novnc/ strip=/%s/novnc/", shortID, shortID),
            fmt.Sprintf("urlprefix-/%s/novnc/websockify strip=/%s/novnc/websockify", shortID, shortID),
        },
    }
    novncRegistration.Check.TCP = fmt.Sprintf("%s:6901", containerIP)
    novncRegistration.Check.Interval = "10s"

    // Register VNC endpoint
    vncRegistration := ConsulServiceRegistration{
        Name:    fmt.Sprintf("vnc-%s", shortID),
        ID:      fmt.Sprintf("vnc-%s", shortID),
        Address: containerIP,
        Port:    5901,
        Tags:    []string{
            fmt.Sprintf("urlprefix-/%s/vnc/ strip=/%s/vnc/", shortID, shortID),
        },
    }
    vncRegistration.Check.TCP = fmt.Sprintf("%s:5901", containerIP)
    vncRegistration.Check.Interval = "10s"

    // Register all services
    for _, registration := range []ConsulServiceRegistration{chatRegistration, novncRegistration, vncRegistration} {
        jsonData, err := json.Marshal(registration)
        if err != nil {
            return fmt.Errorf("failed to marshal registration data: %v", err)
        }

        resp, err := http.DefaultClient.Do(&http.Request{
            Method: "PUT",
            URL:    &url.URL{Scheme: "http", Host: consulAddr, Path: "/v1/agent/service/register"},
            Body:   io.NopCloser(bytes.NewReader(jsonData)),
        })
        if err != nil {
            return fmt.Errorf("failed to register service: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            return fmt.Errorf("failed to register service, status: %d", resp.StatusCode)
        }
    }

    log.Printf("Successfully registered container %s services with Consul", containerID[:12])
    return nil
}

func (dm *DockerManager) ensureImageExists(imageName string) error {
    // Check if image exists locally
    _, _, err := dm.cli.ImageInspectWithRaw(context.Background(), imageName)
    if err == nil {
        // Image exists locally
        return nil
    }

    // Pull the image
    log.Printf("Pulling image: %s", imageName)
    reader, err := dm.cli.ImagePull(context.Background(), imageName, image.PullOptions{
        All: false,
    })
    if err != nil {
        return fmt.Errorf("failed to pull image: %v", err)
    }
    defer reader.Close()

    // Wait for the pull to complete
    _, err = io.Copy(io.Discard, reader)
    if err != nil {
        return fmt.Errorf("error while pulling image: %v", err)
    }

    return nil
}

// Update CreateContainer to accept VNC configuration
func (dm *DockerManager) CreateContainer(imageName string, vncConfig config.VNCConfig) (*ContainerEndpoints, error) {
    // Remove these lines
    // containerID := utils.GenerateID()
    // shortID := containerID[:12]
    log.Printf("Creating container using image: %s", imageName)

    if err := dm.ensureImageExists(imageName); err != nil {
        return nil, err
    }

    // Generate random password if not provided
    if vncConfig.Password == "" {
        vncConfig.Password = utils.GenerateID()[:12] // Use first 12 chars as password
    }

    // Create environment variables for VNC configuration
    env := []string{
        fmt.Sprintf("VNC_PW=%s", vncConfig.Password),
    }

    // Add optional VNC configurations only if they are set
    if vncConfig.Resolution != "" {
        env = append(env, fmt.Sprintf("VNC_RESOLUTION=%s", vncConfig.Resolution))
    } else {
        // Use the width and height from config
        env = append(env, fmt.Sprintf("VNC_RESOLUTION=%sx%s", 
            dm.cfg.ContainerEnvVars.Width, 
            dm.cfg.ContainerEnvVars.Height))
    }
    
    // Add the rest of the environment variables
    env = append(env, []string{
        fmt.Sprintf("PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH=%s", dm.cfg.ContainerEnvVars.PlaywrightChromiumPath),
        fmt.Sprintf("LOG_LEVEL=%s", dm.cfg.ContainerEnvVars.LogLevel),
        fmt.Sprintf("RABBITMQ_QUEUE=%s", dm.cfg.ContainerEnvVars.RabbitMQQueue),
        fmt.Sprintf("PORT=%s", dm.cfg.ContainerEnvVars.Port),
        fmt.Sprintf("RABBITMQ_USER=%s", dm.cfg.ContainerEnvVars.RabbitMQUser),
        fmt.Sprintf("RABBITMQ_PASSWORD=%s", dm.cfg.ContainerEnvVars.RabbitMQPassword),
        fmt.Sprintf("RABBITMQ_HOST=%s", dm.cfg.ContainerEnvVars.RabbitMQHost),
        fmt.Sprintf("RABBITMQ_PORT=%s", dm.cfg.ContainerEnvVars.RabbitMQPort),
    }...)
    
    // Add API keys only if they are set
    if dm.cfg.ContainerEnvVars.OpenAIAPIKey != "" {
        env = append(env, fmt.Sprintf("OPENAI_API_KEY=%s", dm.cfg.ContainerEnvVars.OpenAIAPIKey))
    }
    if dm.cfg.ContainerEnvVars.AnthropicAPIKey != "" {
        env = append(env, fmt.Sprintf("ANTHROPIC_API_KEY=%s", dm.cfg.ContainerEnvVars.AnthropicAPIKey))
    }
    
    if vncConfig.ColDepth != 0 {
        env = append(env, fmt.Sprintf("VNC_COL_DEPTH=%d", vncConfig.ColDepth))
    }
    if vncConfig.Display != "" {
        env = append(env, fmt.Sprintf("DISPLAY=%s", vncConfig.Display))
    }
    if vncConfig.ViewOnly {
        env = append(env, fmt.Sprintf("VNC_VIEW_ONLY=%v", vncConfig.ViewOnly))
    }

    hostConfig := &container.HostConfig{
        PortBindings: nat.PortMap{
            "3000/tcp": []nat.PortBinding{{HostPort: ""}}, // Chat API port (changed from 8080 to 3000)
            "6901/tcp": []nat.PortBinding{{HostPort: ""}}, // noVNC port
            "5901/tcp": []nat.PortBinding{{HostPort: ""}}, // VNC port
        },
        NetworkMode: container.NetworkMode(dm.network),
    }

    resp, err := dm.cli.ContainerCreate(
        context.Background(),
        &container.Config{
            Image: imageName,
            ExposedPorts: nat.PortSet{
                "3000/tcp": struct{}{}, // Changed from 8080 to 3000
                "6901/tcp": struct{}{},
                "5901/tcp": struct{}{},
            },
            Env: env,
        },
        hostConfig,
        nil,
        nil,
        "",
    )

    if err != nil {
        return nil, fmt.Errorf("failed to create container: %v", err)
    }

    // Use Docker's container ID
    containerID := resp.ID
    shortID := containerID[:12]

    // Start the container
    if err := dm.cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
        return nil, fmt.Errorf("failed to start container: %v", err)
    }

    // Get container IP address
    inspect, err := dm.cli.ContainerInspect(context.Background(), resp.ID)
    if err != nil {
        // Clean up the container if inspection fails
        dm.cli.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{Force: true})
        return nil, fmt.Errorf("failed to inspect container: %v", err)
    }

    containerIP := inspect.NetworkSettings.Networks[dm.network].IPAddress

    // Register services with Consul using environment variable for Consul address
    consulAddr := os.Getenv("CONSUL_HTTP_ADDR")
    if consulAddr == "" {
        consulAddr = "localhost:8500" // fallback to default
    }

    if err := dm.registerWithConsul(resp.ID, containerIP, consulAddr); err != nil {
        // Clean up the container if Consul registration fails
        dm.cli.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{Force: true})
        return nil, fmt.Errorf("container creation failed: unable to register with service discovery: %v", err)
    }

    // In CreateContainer method, update the endpoints
    endpoints := &ContainerEndpoints{
        ContainerID:  shortID,
        ChatAPIPath:  fmt.Sprintf("/%s/chat/", shortID),
        NoVNCPath:    fmt.Sprintf("/%s/novnc/vnc_lite.html?password=%s&path=%s/novnc/websockify", shortID, vncConfig.Password, shortID),
        VNCPath:      fmt.Sprintf("/%s/novnc/vnc.html?password=%s&path=%s/novnc/websockify", shortID, vncConfig.Password, shortID),
    }

    return endpoints, nil
}
// Add at the top after type definitions
type ImageInfo struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Category    string   `json:"category"`
    Tags        []string `json:"tags"`
}

var availableImages = []ImageInfo{
    // Generic Ubuntu images
    {
        ID:          "debian-chromium",
        Name:        "shanurcsenitap/chrome-desktop-playwright",
        Description: "Base Debian with VNC and Xfce",
        Category:    "Generic Debian",
        Tags:        []string{"debian", "base", "xfce", "xvfb"},
    },
}

// Add this method to DockerManager
func (dm *DockerManager) ListAvailableImages() []ImageInfo {
    return availableImages
}