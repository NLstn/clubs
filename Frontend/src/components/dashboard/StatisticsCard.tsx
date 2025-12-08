import './StatisticsCard.css';

export interface StatisticsCardProps {
    title: string;
    value: number | string;
    icon?: string;
    subtitle?: string;
    trend?: {
        value: number;
        isPositive: boolean;
    };
    loading?: boolean;
}

const StatisticsCard = ({ title, value, icon, subtitle, trend, loading }: StatisticsCardProps) => {
    if (loading) {
        return (
            <div className="statistics-card loading">
                <div className="card-shimmer"></div>
            </div>
        );
    }

    return (
        <div className="statistics-card">
            <div className="card-header">
                <h3 className="card-title">{title}</h3>
                {icon && <span className="card-icon">{icon}</span>}
            </div>
            <div className="card-content">
                <div className="card-value">{value}</div>
                {subtitle && <div className="card-subtitle">{subtitle}</div>}
                {trend && (
                    <div className={`card-trend ${trend.isPositive ? 'positive' : 'negative'}`}>
                        <span className="trend-icon">{trend.isPositive ? '↗' : '↘'}</span>
                        <span className="trend-value">{Math.abs(trend.value)}</span>
                    </div>
                )}
            </div>
        </div>
    );
};

export default StatisticsCard;
