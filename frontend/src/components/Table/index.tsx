'use client';

import React, { useState, useCallback } from 'react';
import { FiChevronUp, FiChevronDown, FiTable } from 'react-icons/fi';
import styles from './Table.module.css';
import { TableProps, SortConfig } from './types';
import { TableHead } from './TableHead';
import { TableHeader } from './TableHeader';
import { TableBody } from './TableBody';
import { TableRow } from './TableRow';
import { TableCell } from './TableCell';
import { PaginationArrow } from './PaginationArrow';
import classNames from 'classnames';

export const Table: React.FC<TableProps> = ({
    data = [],
    columns = [],
    containerStyle,
    actionButton,
    ListEmptyComponent,
    rowDecoration,
    initialSort,
    onSort,
    showPaginationArrow = false,
    page = 1,
    setPage,
    totalData = 0,
    totalPage = 0,
    pageSize = 10,
    paginationLayout = 'default',
    isLoading = false,
    useClientSideSort = false,
}) => {
    const [sortConfig, setSortConfig] = useState<SortConfig | undefined>(initialSort);

    const handleSort = useCallback(
        async (columnKey: string) => {
            const column = columns.find((col) => col.key === columnKey);
            if (!column?.sortable) return;

            const newDirection: 'asc' | 'desc' =
                sortConfig?.key === columnKey && sortConfig.direction === 'asc' ? 'desc' : 'asc';

            const newSortConfig: SortConfig = {
                key: columnKey,
                direction: newDirection,
            };

            setSortConfig(newSortConfig);

            if (onSort) {
                await onSort(newSortConfig);
            }
        },
        [columns, sortConfig, onSort]
    );

    const sortedData = React.useMemo(() => {
        if (!sortConfig || onSort) return data;

        const sorted = [...data].sort((a, b) => {
            const aValue = a[sortConfig.key];
            const bValue = b[sortConfig.key];

            if (aValue === bValue) return 0;

            const comparison = aValue < bValue ? -1 : 1;
            return sortConfig.direction === 'asc' ? comparison : -comparison;
        });

        return sorted;
    }, [data, sortConfig, onSort]);

    const getSortIcon = (columnKey: string) => {
        if (sortConfig?.key !== columnKey) {
            return (
                <span className={styles.sortIcon}>
                    <FiChevronUp size={14} />
                </span>
            );
        }

        return (
            <span className={styles.sortIconActive}>
                {sortConfig.direction === 'asc' ? (
                    <FiChevronUp size={14} />
                ) : (
                    <FiChevronDown size={14} />
                )}
            </span>
        );
    };

    if (!isLoading && data.length === 0) {
        return (
            <div className={styles.emptyContainer} style={containerStyle}>
                {ListEmptyComponent || (
                    <div className={styles.emptyState}>
                        <FiTable size={64} strokeWidth={1.5} />
                        <h3>No Data Available</h3>
                        <p>There are no records to display at the moment.</p>
                    </div>
                )}
            </div>
        );
    }

    return (
        <div className={styles.tableContainer} style={containerStyle}>
            <div className={styles.tableWrapper}>
                <table className={styles.table}>
                    <TableHeader>
                        <TableRow>
                            {columns.map((column) => (
                                <TableHead key={column.key} className={column.sortable ? styles.sortable : ''}>
                                    <div
                                        className={classNames(styles.headerContent, sortConfig?.key === column.key && styles.headerContentActive)}
                                        onClick={() => column.sortable && handleSort(column.key)}
                                        style={{
                                            cursor: column.sortable ? 'pointer' : 'default',
                                            width: column.width,
                                        }}
                                    >
                                        <span>{column.label}</span>
                                        {column.sortable && getSortIcon(column.key)}
                                    </div>
                                </TableHead>
                            ))}
                            {actionButton && <TableHead>Actions</TableHead>}
                        </TableRow>
                    </TableHeader>
                    {isLoading && (
                        <div className={styles.loadingOverlay}>
                            <div className={styles.spinner} />
                        </div>
                    )}
                    {
                        !isLoading && (
                            <TableBody>

                                {
                                    !isLoading && useClientSideSort ? (
                                        sortedData.map((row, index) => {
                                            const decoration = rowDecoration?.(row, index);

                                            return (
                                                <TableRow key={index} decoration={decoration}>
                                                    {columns.map((column) => (
                                                        <TableCell key={column.key} align={column.align}>
                                                            {column.render
                                                                ? column.render(row[column.key], row)
                                                                : row[column.key]}
                                                        </TableCell>
                                                    ))}
                                                    {actionButton && (
                                                        <TableCell align="center" className={styles.actionCell}>
                                                            {actionButton(row)}
                                                        </TableCell>
                                                    )}
                                                </TableRow>
                                            );
                                        })
                                    ) : (
                                        data.map((row, index) => {
                                            const decoration = rowDecoration?.(row, index);

                                            return (
                                                <TableRow key={index} decoration={decoration}>
                                                    {columns.map((column) => (
                                                        <TableCell key={column.key} align={column.align}>
                                                            {column.render
                                                                ? column.render(row[column.key], row)
                                                                : row[column.key]}
                                                        </TableCell>
                                                    ))}
                                                    {actionButton && (
                                                        <TableCell align="center" className={styles.actionCell}>
                                                            {actionButton(row)}
                                                        </TableCell>
                                                    )}
                                                </TableRow>
                                            );
                                        })
                                    )
                                }
                            </TableBody>
                        )
                    }
                </table>
            </div>
            {showPaginationArrow && setPage && (
                <div className={styles.paginationWrapper}>
                    <PaginationArrow
                        page={page}
                        setPage={setPage}
                        totalData={totalData}
                        totalCurrentData={data.length}
                        pageSize={pageSize}
                        totalPage={totalPage}
                        layout={paginationLayout}
                    />
                </div>
            )}
        </div>
    );
};

export { TableHead } from './TableHead';
export { TableHeader } from './TableHeader';
export { TableBody } from './TableBody';
export { TableRow } from './TableRow';
export { TableCell } from './TableCell';

export default Table;
