// logViewer.js
document.addEventListener('DOMContentLoaded', function () {
    const logList = document.getElementById('log-list');
    const filterForm = document.getElementById('filter-form');
  
    // Function to fetch logs with optional filters
    async function fetchLogs(filters = {}) {
      try {
        const queryParams = new URLSearchParams(filters).toString();
        const response = await fetch(`/api/logs?${queryParams}`);
        if (!response.ok) {
          throw new Error('Failed to fetch logs');
        }
  
        const data = await response.json();
        renderLogs(data.logs);
      } catch (error) {
        console.error(error);
      }
    }
  
    // Function to render logs
    function renderLogs(logs) {
      logList.innerHTML = ''; // Clear previous content
      logs.forEach(log => {
        const logItem = document.createElement('div');
        logItem.classList.add('log-item');
        logItem.innerHTML = `
          <strong>Timestamp:</strong> ${log.timestamp} <br>
          <strong>Level:</strong> ${log.level} <br>
          <strong>Message:</strong> ${log.message}
        `;
        logList.appendChild(logItem);
      });
    }
  
    // Handle form submission to filter logs
    filterForm.addEventListener('submit', function (event) {
      event.preventDefault();
      
      const logLevel = document.getElementById('logLevel').value;
      const startDate = document.getElementById('startDate').value;
      const endDate = document.getElementById('endDate').value;
  
      const filters = {};
      if (logLevel) filters.logLevel = logLevel;
      if (startDate) filters.startDate = startDate;
      if (endDate) filters.endDate = endDate;
  
      fetchLogs(filters);
    });
  
    // Initial load: Fetch all logs
    fetchLogs();
  });
  