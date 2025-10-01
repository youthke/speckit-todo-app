import React from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import GoogleSignupButton from '../components/GoogleSignupButton';
import './SignupPage.css';

const SignupPage: React.FC = () => {
  const [searchParams] = useSearchParams();
  const error = searchParams.get('error');

  return (
    <div className="signup-page">
      <div className="signup-container">
        <h1>Sign Up</h1>
        <p className="signup-description">
          Create an account to start managing your tasks
        </p>

        <div className="signup-form">
          <GoogleSignupButton />

          {error === 'rate_limit_exceeded' && (
            <div className="error-message">
              Too many signup attempts. Please try again later.
            </div>
          )}

          {error === 'authentication_failed' && (
            <div className="error-message">
              Authentication failed. Please try again.
            </div>
          )}
        </div>

        <div className="signup-footer">
          <p>
            Already have an account?{' '}
            <Link to="/login" className="login-link">
              Log in
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
};

export default SignupPage;
