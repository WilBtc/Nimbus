Here's the updated README content for the NimBus Edge Server, incorporating the latest project details and providing clear, concise instructions for setting up and using the server:

---

# NimBus Edge Server

NimBus is a secure edge server designed for Industrial IoT, built on top of the atProtocol using the at_server core. It extends the capabilities of at_server with advanced security features, real-time analytics, and industry-specific routing and access control, making it ideal for modern industrial applications.

## Features
- **Secure Communication**: Leverages the atProtocol for secure, end-to-end encrypted communication, ensuring data integrity and privacy.
- **Edge Analytics**: Processes data in real-time at the edge, providing instant insights and enabling automated decision-making based on critical thresholds and anomaly detection.
- **Advanced Routing**: Implements industrial-grade routing mechanisms that prioritize critical data, optimize traffic flow, and enhance system performance.
- **Access Control**: Manages permissions with role-based access control, ensuring that only authorized devices and users can interact with the server.
- **Intrusion Detection System (IDS)**: Includes anomaly detection and traffic inspection capabilities to identify and mitigate potential security threats.

## Getting Started

### Prerequisites
- **Go**: Version 1.18 or higher is required. You can download and install Go from [golang.org](https://golang.org/dl/).
- **at_server**: NimBus integrates seamlessly with the latest at_server components from the DESS framework. Clone and set up at_server from the official [at_server repository](https://github.com/atsign-foundation/dess).
- **Docker** (Optional but recommended for production deployments): Follow the instructions at [docker.com](https://docs.docker.com/get-docker/) to install Docker on your system.

### Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/yourusername/nimbus_server.git
   cd nimbus_server/server
   ```

2. **Set Up Dependencies**:
   Ensure all dependencies are installed by running:
   ```bash
   go mod tidy
   ```

3. **Configure Environment Variables**:
   Configure the required environment variables in your `.env` file or export them directly in your shell. Essential variables include:
   - `ATSIGN`: Your atSign, e.g., `@nimbus`.
   - `ROOT_DOMAIN`: The domain for your secondary server, e.g., `root.atsign.org`.
   - `SECRET`: The secret key for cram authentication.
   - `EMAIL`: The email address for managing SSL certificates.

4. **Build the Server**:
   Compile the NimBus Edge Server:
   ```bash
   go build -o nimbus_server src/main.go
   ```

5. **Run the Server**:
   Start the NimBus server with:
   ```bash
   ./nimbus_server
   ```

### Configuration

NimBus can be configured through environment variables or by modifying the `config.go` file directly. Key configuration options include:
- **Security Settings**: Adjust the security level and enable/disable IDS features.
- **Routing Preferences**: Set routing priorities for different data types.
- **Analytics Thresholds**: Define thresholds for anomaly detection and automated decision-making.

### Usage

- **Monitoring**: Use built-in logging to monitor server operations. Logs are output to the console and can be directed to a file for persistence.
- **Access Control Management**: Manage device permissions using the integrated access control module. Role-based access ensures secure interactions.
- **Real-Time Analytics**: Leverage the analytics engine to process and visualize data insights directly at the edge.

### Troubleshooting

- **Server Not Starting**: Check if all required environment variables are set correctly and the ports are available.
- **Configuration Errors**: Review the logs for any configuration errors and adjust settings in the `config.go` file or environment variables.
- **Security Issues**: Ensure your encryption settings and IDS configurations are properly initialized and validated during startup.

### Contributing

We welcome contributions to NimBus! Please fork the repository and submit pull requests for any improvements or feature additions.

### License

NimBus is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

### Contact

For more information, questions, or issues, please contact our support team at [w.aroca@insaing.com](w.aroca@insaing.com).

---