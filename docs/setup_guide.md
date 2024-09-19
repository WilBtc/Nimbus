Here's a detailed `setup_guide.md` document to guide Insa Automations engineers through the setup process of the Nimbus Edge Server, including the integration of DESS components. This guide will cover the prerequisites, installation, configuration, and verification steps necessary to get the system running smoothly.

### **`setup_guide.md` - Nimbus Edge Server Setup Guide**

# Nimbus Edge Server Setup Guide

This guide provides step-by-step instructions to set up the Nimbus Edge Server using the DESS (Decentralized Edge Secondary Server) components. Follow the steps carefully to ensure a successful installation and configuration.

## Prerequisites

Before you begin, ensure that the following prerequisites are met:

1. **Operating System**: A Linux-based server (Ubuntu, Debian, CentOS, or similar) with root access.
2. **Go**: Version 1.18 or higher. Install Go from [golang.org](https://golang.org/dl/).
3. **Docker**: Docker and Docker Compose must be installed to run DESS components. Install Docker from [docker.com](https://docs.docker.com/get-docker/).
4. **Git**: Git for cloning repositories. Install from [git-scm.com](https://git-scm.com/downloads).
5. **atSign**: A registered atSign for the server setup. Obtain an atSign from [atsign.com](https://atsign.com/).
6. **Environment Variables**: Securely store sensitive information such as API keys and secrets in environment variables or configuration files as instructed.

## Step 1: Clone the Nimbus Project Repository

Start by cloning the Nimbus project repository from GitHub:

```bash
git clone https://github.com/yourusername/nimbus_server.git
cd nimbus_server/server
```

## Step 2: Install Dependencies

Install the required Go dependencies using the Go module system:

```bash
# Ensure you are inside the /server directory
go mod tidy
```

This command will download and install all necessary dependencies specified in the `go.mod` file, including DESS components.

## Step 3: Configure the Environment Variables

Set up the environment variables needed for the Nimbus server. You can use a `.env` file in the `/config` directory or directly export the variables in your terminal session.

Create a `.env` file inside the `/config` directory:

```bash
# Create the .env file
touch config/keys.env
```

Edit the `keys.env` file and add the necessary environment variables:

```plaintext
ATSIGN="@nimbus1_01"
ROOT_DOMAIN="nimbus.example.com"
SERVER_PORT=6464
CRAM_SECRET="your_cram_secret_here"
SECRET_KEY="your_secret_key_here"
ATSIGN_API_KEY="dc34ac3e-00e6-4d49-a557-0e3945372ed9"
EMAIL="admin@nimbus.example.com"
```

> **Note**: Ensure that the `.env` file is included in `.gitignore` to prevent it from being committed to version control.

## Step 4: Set Up DESS Components

### 4.1 Install DESS (Decentralized Edge Secondary Server)

Follow the instructions to install DESS on your server:

1. **Download and Install DESS**:
   ```bash
   curl -fsSL https://getdess.atsign.com | sudo bash
   ```

2. **Create a New DESS Instance**:
   ```bash
   sudo dess-create
   ```

   During this process, you will be prompted to provide the atSign, domain, port, and CRAM key.

### 4.2 Configure DESS for Nimbus

Modify the DESS configuration file to match your Nimbus setup requirements. Typical configurations are found in the `/etc/dess` directory. Ensure the correct paths are set for SSL certificates and storage locations.

### 4.3 Verify DESS Installation

Check that DESS is running correctly by listing the Docker services:

```bash
sudo docker ps
```

You should see the DESS service listed and running without errors.

## Step 5: Configure Nimbus Server

Edit the `config.yaml` file in the `/config` directory to match your server's configuration. This file will include paths, security settings, logging configurations, and other essential parameters.

### Sample `config.yaml`

```yaml
server:
  at_sign: "@nimbus1_01"
  root_domain: "nimbus.example.com"
  host: "0.0.0.0"
  port: 6464

paths:
  storage_path: "./storage"
  commit_log_path: "./commitLog"

security:
  auth_required: true
  encryption_config: "default_encryption"
  security_level: 2
  cram_secret: "your_cram_secret_here"
  email: "admin@nimbus.example.com"

logging:
  log_level: "INFO"
  log_file_path: "./logs/nimbus.log"
  rotate_logs: true
  rotation_interval: "24h"
```

## Step 6: Start the Nimbus Edge Server

Start the Nimbus server using the following command:

```bash
go run src/main.go
```

Check the server logs for any startup errors. Ensure that all configurations are correctly applied and that the server is listening on the specified port.

## Step 7: Verify Server Functionality

1. **Health Check**: Verify that the server is healthy by checking the logs or using health endpoints if configured.
2. **Security Checks**: Ensure that the server’s security settings, such as authentication and encryption, are functioning correctly.

## Step 8: Enable Monitoring and Alerts

Configure monitoring and alerting settings in the `config.yaml` and ensure that email notifications are set up correctly for critical alerts.

## Troubleshooting

- **Common Errors**: Check logs in `/logs/nimbus.log` for any common startup or runtime errors.
- **Permissions**: Ensure that all file paths have the correct permissions and are accessible by the server process.
- **Network Issues**: Verify that the specified port is open and accessible through your firewall settings.

## Conclusion

The Nimbus Edge Server is now set up and running with integrated DESS components for secure, real-time, and industrial-specific operations. For further assistance or to contribute to the project, refer to the `contributing.md` in the `/docs` directory.

If you encounter issues during setup, consult the `security_considerations.md` for additional security tips and best practices.

---

For questions or further help, please contact the Insa Automations engineering support team or visit our internal wiki for more resources.

---

**End of Setup Guide**

If you need further details or additional guidance on specific steps, feel free to reach out!