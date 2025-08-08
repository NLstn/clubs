import React from 'react';
import './Table.css';

export interface TableColumn<T> {
    key: string;
    header: string;
    render: (item: T) => React.ReactNode;
    className?: string;
}

interface TableProps<T> {
    columns: TableColumn<T>[];
    data: T[];
    keyExtractor: (item: T) => string;
    className?: string;
    emptyMessage?: string;
    footer?: React.ReactNode;
    loading?: boolean;
    error?: string | null;
    loadingMessage?: string;
    errorMessage?: string;
}

function Table<T>({
    columns,
    data,
    keyExtractor,
    className = '',
    emptyMessage = 'No data available',
    footer,
    loading = false,
    error = null,
    loadingMessage = 'Loading...',
    errorMessage = 'Error loading data'
}: TableProps<T>) {
    // Handle loading state
    if (loading) {
        return <div className="table-loading-text">{loadingMessage}</div>;
    }

    // Handle error state
    if (error) {
        return <div className="table-error-text">{errorMessage}</div>;
    }

    // Normalize data to avoid runtime crashes when null/undefined is passed
    const safeData: T[] = Array.isArray(data) ? data : [];

    return (
        <div className={`table-container ${className}`}>
            <table className="reusable-table">
                <thead>
                    <tr>
                        {columns.map((column) => (
                            <th key={column.key} className={column.className}>
                                {column.header}
                            </th>
                        ))}
                    </tr>
                </thead>
                <tbody>
                    {safeData.length === 0 ? (
                        <tr>
                            <td colSpan={columns.length} className="table-empty-text" style={{ textAlign: 'center', fontStyle: 'italic' }}>
                                {emptyMessage}
                            </td>
                        </tr>
                    ) : (
                        safeData.map((item) => (
                            <tr key={keyExtractor(item)}>
                                {columns.map((column) => (
                                    <td key={column.key} className={column.className}>
                                        {column.render(item)}
                                    </td>
                                ))}
                            </tr>
                        ))
                    )}
                </tbody>
            </table>
            <div className="table-footer">
                {footer}
            </div>
        </div>
    );
}

export default Table;
