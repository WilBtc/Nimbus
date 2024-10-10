// dashboard.js
document.addEventListener('DOMContentLoaded', function () {
    const deviceList = document.getElementById('device-list');
    const networkInfo = document.getElementById('network-info');
  
    // Function to fetch device status
    async function fetchDeviceStatus() {
      try {
        const response = await fetch('/api/devices/list');
        if (!response.ok) {
          throw new Error('Failed to fetch device status');
        }
  
        const data = await response.json();
        renderDeviceStatus(data.devices);
      } catch (error) {
        console.error(error);
      }
    }
  
    // Function to render device status
    function renderDeviceStatus(devices) {
      deviceList.innerHTML = ''; // Clear previous content
      devices.forEach(device => {
        const deviceItem = document.createElement('div');
        deviceItem.classList.add('device-item');
        deviceItem.innerHTML = `
          <strong>Device ID:</strong> ${device.deviceId} <br>
          <strong>Status:</strong> ${device.status} <br>
          <strong>Last Active:</strong> ${device.lastActive}
        `;
        deviceList.appendChild(deviceItem);
      });
    }
  
    // Function to fetch network status
    async function fetchNetworkStatus() {
      try {
        const response = await fetch('/api/network/status');
        if (!response.ok) {
          throw new Error('Failed to fetch network status');
        }
  
        const data = await response.json();
        renderNetworkStatus(data.network);
      } catch (error) {
        console.error(error);
      }
    }
  
    // Function to render network status
    function renderNetworkStatus(network) {
      networkInfo.innerHTML = ''; // Clear previous content
      const networkItem = document.createElement('div');
      networkItem.classList.add('network-item');
      networkItem.innerHTML = `
        <strong>IP Address:</strong> ${network.ipAddress} <br>
        <strong>Uptime:</strong> ${network.uptime} <br>
        <strong>Connected Devices:</strong> ${network.connectedDevices}
      `;
      networkInfo.appendChild(networkItem);
    }
  
    // Fetch device status and network status every 10 seconds
    setInterval(() => {
      fetchDeviceStatus();
      fetchNetworkStatus();
    }, 10000);
  
    // Initial load
    fetchDeviceStatus();
    fetchNetworkStatus();
  });
  