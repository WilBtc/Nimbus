**Nimbus Edge Server Setup Guide**, including DESS integration for engineers to follow during the installation and configuration process:

---

### **`setup_guide.md` - Nimbus Edge Server Setup Guide**

# **Nimbus Edge Server Setup Guide**

This guide provides detailed instructions to set up the Nimbus Edge Server, incorporating **DESS (Decentralized Edge Secondary Server)** components. Follow these steps to ensure a successful installation and configuration of the system.

## **Prerequisites**

Before starting the setup process, ensure you have the following:

1. **Operating System**: A Linux-based server (e.g., Ubuntu, Debian, CentOS) with root access.
2. **Go**: Version 1.18 or higher. Install Go from [golang.org](https://golang.org/dl/).
3. **Docker**: Docker and Docker Compose must be installed to run DESS. Install Docker from [docker.com](https://docs.docker.com/get-docker/).
4. **Git**: Git for cloning repositories. Install from [git-scm.com](https://git-scm.com/downloads).
5. **atSign**: A registered atSign for the server setup. Obtain one from [atsign.com](https://atsign.com/).
6. **Environment Variables**: Store sensitive information like API keys and secrets securely in environment variables or configuration files.

## **Step 1: Clone the Nimbus Project Repository**

Start by cloning the Nimbus project repository from GitHub:

```bash
git clone https://github.com/yourusername/nimbus_server.git
cd nimbus_server/server
```

## **Step 2: Install Dependencies**

Use the Go module system to install the required dependencies:

```bash
# Ensure you are inside the /server directory
go mod tidy
```

This will download and install all necessary dependencies as specified in the `go.mod` file, including DESS components.

## **Step 3: Configure the Environment Variables**

Set up the environment variables for the Nimbus server. You can either export them directly in your terminal or create a `.env` file.

Create a `.env` file inside the `/config` directory:

```bash
# Create the .env file
touch config/keys.env
```

Edit the `keys.env` file to add the required environment variables:

```plaintext
ATSIGN="@nimbus1_01"
ROOT_DOMAIN="nimbus.example.com"
SERVER_PORT=6464
CRAM_SECRET="your_cram_secret_here"
SECRET_KEY="your_secret_key_here"
ATSIGN_API_KEY="dc34ac3e-00e6-4d49-a557-0e3945372ed9"
EMAIL="admin@nimbus.example.com"
```

> **Note**: Make sure `.env` is listed in `.gitignore` to prevent sensitive data from being committed to version control.

## **Step 4: Set Up DESS Components**

### **4.1 Install DESS**

To install DESS, follow the steps below:

1. **Download and Install DESS**:
   ```bash
   curl -fsSL https://getdess.atsign.com | sudo bash
   ```

2. **Create a New DESS Instance**:
   ```bash
   sudo dess-create
   ```

   During the setup, provide the atSign, root domain, port, and CRAM key as prompted.

### **4.2 Configure DESS for Nimbus**

After installing DESS, modify its configuration file (typically located in `/etc/dess`) to match the Nimbus server’s setup. Ensure the correct paths for SSL certificates, storage, and logs are specified.

### **4.3 Verify DESS Installation**

To verify that DESS is running correctly, list the Docker services:

```bash
sudo docker ps
```

You should see the DESS service listed and running without any errors.

## **Step 5: Configure the Nimbus Server**

Edit the `config.yaml` file located in the `/config` directory to customize the server’s paths, security, and logging settings.

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

## **Step 6: Start the Nimbus Edge Server**

Start the Nimbus server using the following command:

```bash
go run src/main.go
```

Check the logs for any startup errors, and confirm that the server is listening on the specified port.

## **Step 7: Verify Server Functionality**

1. **Health Check**: Use server logs or configured health endpoints to ensure the server is running properly.
2. **Security Checks**: Confirm that all security measures, such as authentication and encryption, are active.

## **Step 8: Enable Monitoring and Alerts**

To ensure proper monitoring and alerting, configure the relevant parameters in `config.yaml`. Email notifications should be set up to send alerts for critical events such as anomalies or security incidents.

## **Troubleshooting**

- **Common Errors**: Check the `/logs/nimbus.log` file for detailed error messages during startup or operation.
- **Permissions**: Verify that all necessary file paths (storage, logs, etc.) have the correct permissions and are accessible.
- **Network Issues**: Ensure the specified port is open and available through your server's firewall settings.

## **Conclusion**

Your Nimbus Edge Server is now successfully set up and integrated with DESS components. This setup provides a secure, real-time environment for Industrial IoT operations. For further assistance, refer to the `contributing.md` file or reach out to the project maintainers.

If you experience issues during the setup process, consult the `security_considerations.md` document for additional best practices and security configurations.

---

**End of Setup Guide**

For any additional details or clarification on specific steps, feel free to reach out to the engineering support team or visit the internal project wiki.

---

Let me know if you need more details or further assistance!