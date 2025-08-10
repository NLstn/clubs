import { useT } from '../../../hooks/useTranslation';

interface Club {
    id: string;
    name: string;
    description: string;
    logo_url?: string;
    deleted?: boolean;
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
    openFinesCount: number;
    onCreateEvent: () => void;
    onCreateTeam?: () => void;
    teamsEnabled?: boolean;
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
    openFinesCount,
    onCreateEvent,
    onCreateTeam,
    teamsEnabled = false
}: AdminClubOverviewProps) => {
    const { t } = useT();

    return (
        <>
            <div className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
                <div className="flex gap-4">
                    <div>
                        {club.logo_url ? (
                            <div className="flex flex-col items-center">
                                <img
                                    src={club.logo_url}
                                    alt={`${club.name} logo`}
                                    className="h-24 w-24 rounded object-cover"
                                />
                                {!club.deleted && (
                                    <div className="mt-2 flex gap-2">
                                        <input
                                            type="file"
                                            id="logo-upload"
                                            accept="image/png,image/jpeg,image/jpg,image/webp"
                                            onChange={onLogoUpload}
                                            className="hidden"
                                        />
                                        <button
                                            onClick={() => document.getElementById('logo-upload')?.click()}
                                            className="button-accept"
                                            disabled={logoUploading}
                                        >
                                            {logoUploading
                                                ? t('common.uploading') || 'Uploading...'
                                                : t('common.change') || 'Change'}
                                        </button>
                                        <button
                                            onClick={onLogoDelete}
                                            className="button-cancel"
                                            disabled={logoUploading}
                                        >
                                            {t('common.delete')}
                                        </button>
                                    </div>
                                )}
                            </div>
                        ) : (
                            <div className="flex flex-col items-center">
                                <div
                                    className={`flex h-24 w-24 items-center justify-center border text-center text-sm text-gray-500 ${!club.deleted ? 'cursor-pointer' : ''}`}
                                    onClick={!club.deleted ? () => document.getElementById('logo-upload')?.click() : undefined}
                                >
                                    {!club.deleted
                                        ? t('clubs.uploadLogoPrompt') || 'Click to upload logo'
                                        : t('clubs.noLogo') || 'No logo'}
                                </div>
                                {!club.deleted && (
                                    <input
                                        type="file"
                                        id="logo-upload"
                                        accept="image/png,image/jpeg,image/jpg,image/webp"
                                        onChange={onLogoUpload}
                                        className="hidden"
                                    />
                                )}
                            </div>
                        )}
                    </div>
                    <div>
                        <h2 className="text-2xl font-bold">{club.name}</h2>
                        <p>{club.description}</p>
                        {logoError && <div className="text-red-500">{logoError}</div>}
                    </div>
                </div>
                <div className="flex flex-wrap gap-2">
                    {!club.deleted && (
                        <>
                            <button onClick={onEdit} className="button-accept">
                                {t('clubs.editClub')}
                            </button>
                            {isOwner && (
                                <button onClick={onDelete} className="button-cancel">
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

            <div className="mt-4 flex flex-wrap gap-2">
                <button onClick={onCreateEvent} className="button-accept">
                    {t('events.addEvent') || 'Add Event'}
                </button>
                {teamsEnabled && onCreateTeam && (
                    <button onClick={onCreateTeam} className="button-accept">
                        {t('teams.createTeam')}
                    </button>
                )}
            </div>

            <div className="mt-4 grid gap-4 sm:grid-cols-2 md:grid-cols-3">
                <div className="rounded border p-4 text-center">
                    <div className="text-2xl font-bold">{openFinesCount}</div>
                    <div className="text-sm text-gray-600">
                        {t('clubs.openFines') || 'Open fines'}
                    </div>
                </div>
            </div>

            {club.deleted && (
                <div className="mt-4 text-red-600">
                    <strong>{t('clubs.clubDeleted')}</strong>
                </div>
            )}
        </>
    );
};

export default AdminClubOverview;

