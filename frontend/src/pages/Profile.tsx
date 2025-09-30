import React from 'react';
import { useAuth } from '../hooks/useAuth';

const Profile: React.FC = () => {
  const { user } = useAuth();

  if (!user) {
    return (
      <div className="profile-page">
        <h1>Profile</h1>
        <p>No user information available.</p>
      </div>
    );
  }

  return (
    <div className="profile-page">
      <header className="profile-header">
        <h1>My Profile</h1>
      </header>

      <div className="profile-content">
        <div className="profile-card">
          <div className="profile-field">
            <label>Name:</label>
            <span>{user.name}</span>
          </div>

          <div className="profile-field">
            <label>Email:</label>
            <span>{user.email}</span>
          </div>

          <div className="profile-field">
            <label>User ID:</label>
            <span>{user.id}</span>
          </div>
        </div>

        <div className="profile-actions">
          <button className="action-button secondary" disabled>
            Edit Profile (Coming Soon)
          </button>
        </div>
      </div>
    </div>
  );
};

export default Profile;