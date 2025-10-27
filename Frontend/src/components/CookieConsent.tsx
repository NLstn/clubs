import React, { useState } from 'react';
import './CookieConsent.css';

const COOKIE_CONSENT_KEY = 'cookie-consent';

const CookieConsent: React.FC = () => {
  // Initialize state directly from localStorage
  const [isVisible, setIsVisible] = useState(() => {
    const consent = localStorage.getItem(COOKIE_CONSENT_KEY);
    return !consent;
  });

  const handleAccept = () => {
    localStorage.setItem(COOKIE_CONSENT_KEY, 'accepted');
    setIsVisible(false);
  };

  const handleLearnMore = () => {
    // For now, this could link to a privacy policy or simply provide more info
    // In a real implementation, this might open a modal or navigate to a privacy page
    alert('This website uses cookies to enhance your browsing experience and provide personalized content. By clicking "Accept", you consent to our use of cookies.');
  };

  if (!isVisible) {
    return null;
  }

  return (
    <div className="cookie-consent-banner" data-testid="cookie-consent-banner">
      <div className="cookie-consent-content">
        <p className="cookie-consent-text">
          We use cookies to improve your experience on our site. By using our site, you accept our use of cookies.
        </p>
        <div className="cookie-consent-actions">
          <button 
            className="cookie-consent-btn cookie-consent-btn-learn-more"
            onClick={handleLearnMore}
            type="button"
          >
            Learn More
          </button>
          <button 
            className="cookie-consent-btn cookie-consent-btn-accept"
            onClick={handleAccept}
            type="button"
          >
            Accept
          </button>
        </div>
      </div>
    </div>
  );
};

export default CookieConsent;