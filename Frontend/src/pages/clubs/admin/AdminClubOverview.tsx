import { useT } from '../../../hooks/useTranslation';
import styles from './AdminClubOverview.module.css';

interface Club {
    id: string;
    name: string;
    description: string;
    logo_url?: string;
    deleted?: boolean;
}

interface Metric {
    label: string;
    value: number | string;
}

interface AdminClubOverviewProps {
    club: Club;
    isOwner: boolean;
    logoUploading: boolean;
    logoError: string | null;
    onEdit: () => void;
    onDelete: () => void;
    onHardDelete: () => void;
    onLogoUpload: (event: React.ChangeEvent<HTMLInputElement>) => void;
    onLogoDelete: () => void;
    metrics?: Metric[];
}

const AdminClubOverview = ({
    club,
    isOwner,
    logoUploading,
    logoError,
    onEdit,
    onDelete,
    onHardDelete,
    onLogoUpload,
    onLogoDelete,
    metrics = []
}: AdminClubOverviewProps) => {
    const { t } = useT();

    return (
        <>
            <div className={styles.header}>
                <div className={styles.info}>
                    <div className={styles.logoSection}>
                        {club.logo_url ? (
                            <div className={styles.logoContainer}>
                                <img
                                    src={club.logo_url}
                                    alt={`${club.name} logo`}
                                    className={styles.logo}
                                />
                                {!club.deleted && (
                                    <div className={styles.logoActions}>
                                        <input
                                            type="file"
                                            id="logo-upload"
                                            accept="image/png,image/jpeg,image/jpg,image/webp"
                                            onChange={onLogoUpload}
                                            style={{ display: 'none' }}
                                        />
                                        <button
                                            onClick={() => document.getElementById('logo-upload')?.click()}
                                            className={`${styles.logoButton} ${styles.change}`}
                                            disabled={logoUploading}
                                        >
                                            {logoUploading ? t('common.uploading') || 'Uploading...' : t('common.change') || 'Change'}
                                        </button>
                                        <button
                                            onClick={onLogoDelete}
                                            className={`${styles.logoButton} ${styles.delete}`}
                                            disabled={logoUploading}
                                        >
                                            {t('common.delete')}
                                        </button>
                                    </div>
                                )}
                            </div>
                        ) : (
                            <div className={styles.placeholderContainer}>
                                <div
                                    className={styles.placeholder}
                                    onClick={!club.deleted ? () => document.getElementById('logo-upload')?.click() : undefined}
                                >
                                    {!club.deleted ? t('clubs.uploadLogoPrompt') || 'Click to upload logo' : t('clubs.noLogo') || 'No logo'}
                                </div>
                                {!club.deleted && (
                                    <input
                                        type="file"
                                        id="logo-upload"
                                        accept="image/png,image/jpeg,image/jpg,image/webp"
                                        onChange={onLogoUpload}
                                        style={{ display: 'none' }}
                                    />
                                )}
                            </div>
                        )}
                    </div>
                    <div className={styles.details}>
                        <h2>{club.name}</h2>
                        <p>{club.description}</p>
                        {logoError && <div className={styles.logoError}>{logoError}</div>}
                    </div>
                </div>
                <div className={styles.actions}>
                    {!club.deleted && (
                        <>
                            <button onClick={onEdit} className="button-accept">{t('clubs.editClub')}</button>
                            {isOwner && (
                                <button
                                    onClick={onDelete}
                                    className="button-cancel"
                                    style={{ marginLeft: '10px' }}
                                >
                                    {t('clubs.deleteClub')}
                                </button>
                            )}
                        </>
                    )}
                    {club.deleted && isOwner && (
                        <button
                            onClick={onHardDelete}
                            className="button-cancel"
                            style={{ backgroundColor: '#d32f2f', borderColor: '#d32f2f' }}
                        >
                            {t('clubs.hardDeleteClub')}
                        </button>
                    )}
                </div>
            </div>

            {metrics.length > 0 && (
                <div className={styles.metrics}>
                    {metrics.map((metric) => (
                        <div key={metric.label} className={styles.metricCard}>
                            <div className={styles.metricValue}>{metric.value}</div>
                            <div className={styles.metricLabel}>{metric.label}</div>
                        </div>
                    ))}
                </div>
            )}

            {club.deleted && (
                <div className={styles.deletedNotice}>
                    <strong>{t('clubs.clubDeleted')}</strong>
                </div>
            )}
        </>
    );
};

export default AdminClubOverview;
