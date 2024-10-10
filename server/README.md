# Nimbus Server

Nimbus is a secure and scalable platform designed for Industrial IoT environments, enabling device management, collaborative data sharing, secure file storage, and real-time analytics. The server integrates advanced technologies like Hyperswarm, Atsign authentication, NoPorts tunneling, Hypercore, and Autobase.

## Features

- **Peer Discovery and Communication**: Uses Hyperswarm for efficient peer discovery and decentralized communication.
- **Secure Authentication**: Implements Atsign authentication for secure and private interactions.
- **NoPorts Tunneling**: Enables tunnel connections without needing traditional port configurations.
- **Collaborative Data Sharing**: Uses Autobase for managing collaborative updates with multiple writers.
- **Immutable Data Storage**: Leverages Hypercore for secure and immutable data storage.

## Getting Started

### Prerequisites

- **Node.js**: Install Node.js (version 14 or above) from [Node.js Official Website](https://nodejs.org/).
- **NPM**: Node package manager (comes bundled with Node.js).
- **Git**: Ensure you have Git installed for cloning the repository.
- **Atsign Account**: You need an atSign for authentication purposes. You can register at [atsign.com](https://atsign.com).

### Installation

1. **Clone the Repository**:

    ```bash
    git clone https://github.com/yourusername/nimbus-server.git
    cd nimbus-server
    ```

2. **Install Dependencies**:

    Install the required Node.js packages:

    ```bash
    npm install
    ```

3. **Set Up Environment Variables**:

    Create a `.env` file in the root directory to store sensitive configuration values:

    ```plaintext
    PORT=3000
    ATSIGN="@youratsign"
    HYPERSWARM_TOPIC="your_hyperswarm_topic"
    NO_PORTS_API_KEY="your_noports_api_key"
    ```

    > **Note**: Replace `@youratsign`, `your_hyperswarm_topic`, and `your_noports_api_key` with your actual credentials.

4. **Start the Server**:

    You can start the Nimbus server using the following command:

    ```bash
    npm start
    ```

    This will start the server on the configured port (default is 3000).

### Usage

Once the server is running, you can access the following features:

- **Dashboard**: Access the web-based dashboard at `http://localhost:3000/ui/dashboard/index.html`.
- **API Endpoints**:
  - `/api/devices`: Manage device connections and interactions.
  - `/api/files`: Handle file uploads, downloads, and sharing.

### Testing

Nimbus server comes with a suite of unit and integration tests to ensure core functionalities are working as expected.

1. **Run Tests**:

    To run the test suite, use the following command:

    ```bash
    npm test
    ```

    This will execute all test cases in the `/test` directory.

### Core Technologies

- **Hyperswarm**: Decentralized peer-to-peer networking for efficient discovery and communication.
- **Atsign Authentication**: Secure, private communication leveraging the atProtocol.
- **NoPorts**: Tunneling solution without traditional port configurations, ensuring secure and hassle-free connections.
- **Hypercore**: Secure and immutable data storage technology, ensuring data integrity and persistence.
- **Autobase**: Multi-writer database for collaborative file sharing and updates.

### Folder Structure

```plaintext
├── /api                    # API logic for various core functionalities
│   ├── hyperswarm.js        # Peer discovery and networking via Hyperswarm
│   ├── atsignAuth.js        # Atsign authentication and encryption logic
│   ├── noPorts.js           # NoPorts tunneling logic
│   ├── hypercore.js         # Hypercore-based data storage
│   ├── autobase.js          # Autobase for collaborative data sharing
├── /routes                 # REST API routes
│   ├── deviceRoutes.js      # Routes for managing device connections
│   └── fileRoutes.js        # Routes for file uploads, downloads, and sharing
├── /public                 # Publicly accessible files (UI components)
│   ├── index.html           # Main landing page
│   ├── styles.css           # Global styles for UI
│   └── logo.png             # Nimbus branding logo
├── /test                   # Test cases for different functionalities
│   ├── testHyperswarm.js    # Tests for Hyperswarm peer discovery
│   ├── testAtsign.js        # Tests for Atsign authentication
│   ├── testNoPorts.js       # Tests for NoPorts tunneling
│   ├── testHypercore.js     # Tests for Hypercore data storage
│   └── testAutobase.js      # Tests for Autobase collaborative sharing
├── app.js                  # Main entry point for the server
├── package.json            # Project metadata and scripts
└── README.md               # Documentation for setup and usage
