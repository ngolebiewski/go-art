import React, { useState, type ChangeEvent, type FormEvent } from 'react';
import axios from 'axios';

// --- CONSTANTS ---
const UPLOAD_URL = '/api/artworks';
// const TEST_ARTIST_ID = 1; // Assuming Artist ID 1 exists from manual insertion

const ArtworkUploader: React.FC = () => {
  const [title, setTitle] = useState<string>('Untitled')
  const [grade, setGrade] = useState('3');
  const [school, setSchool] = useState('');
  const [description, setDescription] = useState<string>('')
  const [artistId, setArtistId] = useState('1');
  const [file, setFile] = useState<File | null>(null);
  const [progress, setProgress] = useState<number>(0);
  const [message, setMessage] = useState<string>('');
  const [status, setStatus] = useState<'idle' | 'uploading' | 'success' | 'error'>('idle');
  const [imageId, setImageId] = useState<number | null>(null);

  // Handle form field changes
  const handleFieldChange = (e: ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    if (name === 'title') setTitle(value);
    if (name === 'grade') setGrade(value);
    if (name === 'school') setSchool(value);
    if (name === 'description') setDescription(value);
    if (name === 'artistId') setArtistId(value);
  };

  // Handle file selection
  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      setFile(e.target.files[0]);
      setMessage(`File selected: ${e.target.files[0].name}`);
      setStatus('idle');
      setProgress(0);
    } else {
      setFile(null);
    }
  };

  // Handle combined submission
  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();

    if (!title || !file) {
      setMessage('Please provide a Title and select an Image.');
      return;
    }

    setStatus('uploading');
    setMessage('Uploading, creating artwork, and processing image...');

    const formData = new FormData();

    // 1. Append Text Data (must match r.FormValue in Go handler)
    formData.append('title', title);
    formData.append('artist_id', artistId);
    formData.append('grade', grade);
    formData.append('school', school);
    formData.append('description', description)

    // 2. Append File Data (must match r.FormFile in Go handler)
    formData.append('image', file);

    try {
      const response = await axios.post(UPLOAD_URL, formData, {

        headers: {
          'Content-Type': 'multipart/form-data', // Essential for FormData
        },
        onUploadProgress: (p) => {
          // Note: Progress bar only tracks file upload, not server processing time.
          if (p.total) {
            const percent = Math.round((p.loaded * 100) / p.total);
            setProgress(percent);
          }
        },
      });
      setImageId(response.data.data.image_id); // <--- Add this line
      setStatus('success');
      // Assuming the success message structure is correct
      setMessage(`SUCCESS! Artwork ID: ${response.data.data.artwork_id}, Image ID: ${response.data.data.image_id}`);
      setTitle(''); // Clear form fields on success
      setSchool('');
      setGrade('');
      setDescription('');
      setFile(null);
      setProgress(0);

    } catch (error) {
      setStatus('error');
      if (axios.isAxiosError(error) && error.response) {
        setMessage(`HTTP Error (${error.response.status}): ${error.response.data.Error || error.message}`);
      } else {
        setMessage('Network Error. Is the Go server running?');
      }
    }
  };

  const getStatusColor = () => {
    switch (status) {
      case 'uploading': return 'blue';
      case 'success': return 'green';
      case 'error': return 'red';
      default: return 'gray';
    }
  };

  return (
    <div style={{ padding: '20px', maxWidth: '500px', margin: '20px auto', border: `1px solid ${getStatusColor()}`, borderRadius: '8px' }}>
      <h2>New Artwork & Image Upload</h2>
      <form onSubmit={handleSubmit}>
        {/* Artwork Fields */}
        <label style={{ display: 'block', marginBottom: '10px' }}>
          Title (Required):
          <input
            type="text"
            name="title"
            value={title}
            onChange={handleFieldChange}
            style={{ width: '100%', padding: '8px', boxSizing: 'border-box', marginTop: '4px' }}
          />
        </label>
        <label style={{ display: 'block', marginBottom: '10px' }}>
          Artist ID (Required -- will be a drop down of available artists):
          <input
            type="text"
            name="artistId"
            value={artistId}
            onChange={handleFieldChange}
            style={{ width: '100%', padding: '8px', boxSizing: 'border-box', marginTop: '4px' }}
          />
        </label>
        <label style={{ display: 'block', marginBottom: '10px' }}>
          School:
          <input
            type="text"
            name="school"
            value={school}
            onChange={handleFieldChange}
            style={{ width: '100%', padding: '8px', boxSizing: 'border-box', marginTop: '4px' }}
          />
        </label>

        <label style={{ display: 'block', marginBottom: '10px' }}>
          Grade:
          <input
            type="text"
            name="grade"
            value={grade}
            onChange={handleFieldChange}
            style={{ width: '100%', padding: '8px', boxSizing: 'border-box', marginTop: '4px' }}
          />
        </label>
        <label style={{ display: 'block', marginBottom: '10px' }}>
          Description:
          <input
            type="text"
            name="description"
            value={description}
            onChange={handleFieldChange}
            style={{ width: '100%', padding: '8px', boxSizing: 'border-box', marginTop: '4px' }}
          />
        </label>
        {/* Displaying static artist ID for clarity */}
        <p style={{ fontSize: '0.9em', color: 'gray' }}>Using static Artist ID: **{artistId}**</p>

        <hr style={{ margin: '20px 0' }} />

        {/* File Input */}
        <label style={{ display: 'block', marginBottom: '15px' }}>
          Image File (Required):
          <input
            type="file"
            accept="image/jpeg,image/png,image/gif"
            onChange={handleFileChange}
            style={{ display: 'block', marginTop: '4px' }}
            disabled={status === 'uploading'}
          />
        </label>

        {/* Upload Button */}
        <button
          type="submit"
          disabled={!title || !file || status === 'uploading'}
          style={{
            padding: '10px 15px',
            backgroundColor: status === 'uploading' ? 'lightgray' : '#007bff',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: status === 'uploading' ? 'not-allowed' : 'pointer'
          }}
        >
          {status === 'uploading' ? 'Processing...' : 'Create Artwork & Upload Image'}
        </button>
      </form>

      <hr style={{ margin: '15px 0' }} />

      {/* Progress/Status */}
      {status === 'uploading' && (
        <div style={{ width: '100%', backgroundColor: '#f3f3f3', borderRadius: '4px', marginTop: '10px' }}>
          <div
            style={{
              height: '20px',
              width: `${progress}%`,
              backgroundColor: 'skyblue',
              textAlign: 'center',
              color: 'white',
              borderRadius: '4px'
            }}
          >
            {progress}%
          </div>
        </div>
      )}

      <p style={{ marginTop: '10px', color: getStatusColor() }}>
        Status: {message}
      </p>

      {/* Display image ONLY if status is 'success' AND we have an imageId */}
      {status === 'success' && imageId && (
        <div style={{ marginTop: '15px' }}>
          {/* The src must point to your retrieval route */}
          <img
            src={`/api/artworks/images/${imageId}/thumb`}
            alt={`Uploaded Image ID ${imageId}`}
            style={{ maxWidth: '200px', height: 'auto', border: '1px solid #ccc' }}
          />
          <p style={{ fontSize: '0.8em', color: 'gray' }}>
            Preview of Image ID: {imageId}
          </p>
        </div>
      )}



    </div>
  );
};

export default ArtworkUploader;