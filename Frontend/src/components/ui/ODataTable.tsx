import { useEffect, useState, useCallback } from 'react';
import Table, { TableColumn } from './Table';
import Button from './Button';
import api from '../../utils/api';
import { buildODataQuery, ODataCollectionResponse } from '../../utils/odata';
import './ODataTable.css';

export interface ODataTableColumn<T> extends TableColumn<T> {
    /** Field name for OData sorting (if different from key) */
    sortField?: string;
    /** Whether this column is sortable */
    sortable?: boolean;
}

interface ODataTableProps<T> {
    /** OData entity set path (e.g., "/api/v2/News") */
    endpoint: string;
    /** Column definitions */
    columns: ODataTableColumn<T>[];
    /** Extract unique key from each item */
    keyExtractor: (item: T) => string;
    /** Optional static filter to apply (e.g., "ClubID eq 'abc'") */
    filter?: string;
    /** Optional expand clause */
    expand?: string | string[];
    /** Optional select clause */
    select?: string[];
    /** Page size for pagination */
    pageSize?: number;
    /** Additional CSS class */
    className?: string;
    /** Empty state message */
    emptyMessage?: string;
    /** Loading message */
    loadingMessage?: string;
    /** Error message */
    errorMessage?: string;
    /** Initial sort field */
    initialSortField?: string;
    /** Initial sort direction */
    initialSortDirection?: 'asc' | 'desc';
}

interface SortState {
    field: string | null;
    direction: 'asc' | 'desc';
}

/**
 * ODataTable - A table component with server-side pagination and sorting via OData
 * 
 * This component extends the base Table component to automatically handle:
 * - Server-side pagination using $skip and $top
 * - Server-side sorting using $orderby
 * - OData response format parsing (value, @odata.count)
 * 
 * @example
 * <ODataTable
 *   endpoint="/api/v2/News"
 *   filter="ClubID eq 'abc-123'"
 *   columns={columns}
 *   keyExtractor={(item) => item.ID}
 *   pageSize={10}
 * />
 */
function ODataTable<T>({
    endpoint,
    columns,
    keyExtractor,
    filter,
    expand,
    select,
    pageSize = 10,
    className = '',
    emptyMessage = 'No data available',
    loadingMessage = 'Loading...',
    errorMessage = 'Error loading data',
    initialSortField,
    initialSortDirection = 'desc',
}: ODataTableProps<T>) {
    const [data, setData] = useState<T[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [currentPage, setCurrentPage] = useState(0);
    const [totalCount, setTotalCount] = useState(0);
    const [sortState, setSortState] = useState<SortState>({
        field: initialSortField || null,
        direction: initialSortDirection,
    });

    const fetchData = useCallback(async () => {
        setLoading(true);
        setError(null);
        
        try {
            // Build OData query with pagination and sorting
            const queryOptions: Record<string, string | number | boolean | string[]> = {
                skip: currentPage * pageSize,
                top: pageSize,
                count: true, // Request total count
            };

            if (filter) {
                queryOptions.filter = filter;
            }

            if (expand) {
                queryOptions.expand = expand;
            }

            if (select) {
                queryOptions.select = select;
            }

            if (sortState.field) {
                queryOptions.orderby = `${sortState.field} ${sortState.direction}`;
            }

            const queryString = buildODataQuery(queryOptions);
            const response = await api.get<ODataCollectionResponse<T>>(`${endpoint}${queryString}`);
            
            setData(response.data.value || []);
            setTotalCount(response.data['@odata.count'] || 0);
        } catch (err) {
            console.error('ODataTable fetch error:', err);
            setError(err instanceof Error ? err.message : 'Failed to fetch data');
            setData([]);
            setTotalCount(0);
        } finally {
            setLoading(false);
        }
    }, [endpoint, filter, expand, select, currentPage, pageSize, sortState]);

    useEffect(() => {
        fetchData();
    }, [fetchData]);

    const handleSort = (field: string) => {
        setSortState((prev) => {
            // Toggle direction if same field, otherwise default to 'asc'
            if (prev.field === field) {
                return {
                    field,
                    direction: prev.direction === 'asc' ? 'desc' : 'asc',
                };
            }
            return {
                field,
                direction: 'asc',
            };
        });
        // Reset to first page when sorting changes
        setCurrentPage(0);
    };

    const totalPages = Math.ceil(totalCount / pageSize);
    const hasNextPage = currentPage < totalPages - 1;
    const hasPreviousPage = currentPage > 0;

    const handleNextPage = () => {
        if (hasNextPage) {
            setCurrentPage((prev) => prev + 1);
        }
    };

    const handlePreviousPage = () => {
        if (hasPreviousPage) {
            setCurrentPage((prev) => prev - 1);
        }
    };

    const handleFirstPage = () => {
        setCurrentPage(0);
    };

    const handleLastPage = () => {
        setCurrentPage(totalPages - 1);
    };

    // Enhance columns with sortable headers
    const enhancedColumns: TableColumn<T>[] = columns.map((col) => ({
        ...col,
        header: col.sortable ? (
            <button
                className="odata-table-sort-header"
                onClick={() => handleSort(col.sortField || col.key)}
                aria-label={`Sort by ${col.header}`}
            >
                {col.header}
                {sortState.field === (col.sortField || col.key) && (
                    <span className="sort-indicator">
                        {sortState.direction === 'asc' ? ' ▲' : ' ▼'}
                    </span>
                )}
            </button>
        ) : (
            col.header
        ),
    }));

    const paginationFooter = totalCount > 0 ? (
        <div className="odata-table-pagination">
            <div className="pagination-info">
                Showing {currentPage * pageSize + 1} - {Math.min((currentPage + 1) * pageSize, totalCount)} of {totalCount}
            </div>
            <div className="pagination-controls">
                <Button
                    variant="secondary"
                    size="sm"
                    onClick={handleFirstPage}
                    disabled={!hasPreviousPage}
                    aria-label="First page"
                >
                    ««
                </Button>
                <Button
                    variant="secondary"
                    size="sm"
                    onClick={handlePreviousPage}
                    disabled={!hasPreviousPage}
                    aria-label="Previous page"
                >
                    ‹
                </Button>
                <span className="page-indicator">
                    Page {currentPage + 1} of {totalPages}
                </span>
                <Button
                    variant="secondary"
                    size="sm"
                    onClick={handleNextPage}
                    disabled={!hasNextPage}
                    aria-label="Next page"
                >
                    ›
                </Button>
                <Button
                    variant="secondary"
                    size="sm"
                    onClick={handleLastPage}
                    disabled={!hasNextPage}
                    aria-label="Last page"
                >
                    »»
                </Button>
            </div>
        </div>
    ) : null;

    return (
        <div className={`odata-table-wrapper ${className}`}>
            <Table
                columns={enhancedColumns}
                data={data}
                keyExtractor={keyExtractor}
                loading={loading}
                error={error}
                emptyMessage={emptyMessage}
                loadingMessage={loadingMessage}
                errorMessage={errorMessage}
                footer={paginationFooter}
            />
        </div>
    );
}

export default ODataTable;
