// collaboration.js
document.addEventListener('DOMContentLoaded', function () {
    const fileListSection = document.getElementById('file-list');
    const fileContent = document.getElementById('fileContent');
    const saveBtn = document.getElementById('saveBtn');
    const uploadForm = document.getElementById('upload-form');
    
    let currentFile = null;
  
    // Function to fetch the list of shared files
    async function fetchFiles() {
      try {
        const response = await fetch('/api/collaboration/files');
        if (!response.ok) {
          throw new Error('Failed to fetch files');
        }
  
        const data = await response.json();
        renderFileList(data.files);
      } catch (error) {
        console.error(error);
      }
    }
  
    // Function to render the file list
    function renderFileList(files) {
      fileListSection.innerHTML = ''; // Clear previous content
      files.forEach(file => {
        const fileItem = document.createElement('div');
        fileItem.classList.add('file-item');
        fileItem.innerHTML = `
          <strong>File:</strong> ${file.fileName} <br>
          <button onclick="selectFile('${file.fileName}')">Collaborate</button>
        `;
        fileListSection.appendChild(fileItem);
      });
    }
  
    // Function to upload a file
    uploadForm.addEventListener('submit', async function (event) {
      event.preventDefault();
  
      const fileInput = document.getElementById('fileInput');
      if (!fileInput.files.length) {
        return alert('Please select a file to upload.');
      }
  
      const formData = new FormData();
      formData.append('file', fileInput.files[0]);
  
      try {
        const response = await fetch('/api/collaboration/upload', {
          method: 'POST',
          body: formData,
        });
  
        if (!response.ok) {
          throw new Error('Failed to upload file');
        }
  
        alert('File uploaded successfully!');
        fetchFiles(); // Refresh the file list after upload
      } catch (error) {
        console.error(error);
      }
    });
  
    // Function to select a file and load its content for collaboration
    window.selectFile = async function (fileName) {
      try {
        const response = await fetch(`/api/collaboration/file/${fileName}`);
        if (!response.ok) {
          throw new Error('Failed to fetch file content');
        }
  
        const data = await response.json();
        currentFile = fileName;
        fileContent.value = data.content;
        saveBtn.disabled = false; // Enable the save button when a file is selected
      } catch (error) {
        console.error(error);
      }
    };
  
    // Function to save changes to the file
    saveBtn.addEventListener('click', async function () {
      if (!currentFile) return;
  
      const updatedContent = fileContent.value;
  
      try {
        const response = await fetch(`/api/collaboration/file/${currentFile}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ content: updatedContent }),
        });
  
        if (!response.ok) {
          throw new Error('Failed to save changes');
        }
  
        alert('File saved successfully!');
      } catch (error) {
        console.error(error);
      }
    });
  
    // Initial load: Fetch the list of shared files
    fetchFiles();
  });
  