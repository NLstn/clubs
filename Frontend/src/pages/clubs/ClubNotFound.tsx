import { useNavigate } from 'react-router-dom';
import Layout from '../../components/layout/Layout';
import { useT } from '../../hooks/useTranslation';
import './ClubNotFound.css';

interface ClubNotFoundProps {
  clubId?: string;
  title?: string;
  message?: string;
}

const ClubNotFound: React.FC<ClubNotFoundProps> = ({ 
  clubId, 
  title = 'Club Not Found',
  message = 'The club you are looking for does not exist or has been deleted.'
}) => {
  const { t } = useT();
  const navigate = useNavigate();

  return (
    <Layout title={title}>
      <div className="club-not-found">
        <div className="club-not-found-content">
          <div className="club-not-found-icon">
            <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <circle cx="12" cy="12" r="10"/>
              <path d="m9 9 6 6"/>
              <path d="m15 9-6 6"/>
            </svg>
          </div>
          
          <h1 className="club-not-found-title">{title}</h1>
          
          <p className="club-not-found-message">{message}</p>
          
          {clubId && (
            <p className="club-not-found-id">
              {t('clubs.clubId')}: <code>{clubId}</code>
            </p>
          )}
          
          <div className="club-not-found-actions">
            <button 
              onClick={() => navigate('/clubs')}
              className="button-primary"
            >
              {t('clubs.viewAllClubs')}
            </button>
            
            <button 
              onClick={() => navigate('/')}
              className="button-secondary"
            >
              {t('common.goToDashboard')}
            </button>
            
            <button 
              onClick={() => navigate('/createClub')}
              className="button-accent"
            >
              {t('clubs.createClub')}
            </button>
          </div>
          
          <div className="club-not-found-help">
            <h3>{t('clubs.possibleReasons')}</h3>
            <ul>
              <li>{t('clubs.clubDeleted')}</li>
              <li>{t('clubs.noLongerMember')}</li>
              <li>{t('clubs.invalidLink')}</li>
              <li>{t('clubs.temporaryIssue')}</li>
            </ul>
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default ClubNotFound;
