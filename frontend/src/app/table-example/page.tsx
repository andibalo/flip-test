'use client';

import { Table } from '@/components/Table';
import { Column, Row } from '@/components/Table/types';
import styles from './example.module.css';

export default function TableExample() {
    // Sample data
    const sampleData: Row[] = [
        { id: 1, name: 'John Doe', email: 'john@example.com', status: 'SUCCESS', amount: 1500000 },
        { id: 2, name: 'Jane Smith', email: 'jane@example.com', status: 'FAILED', amount: 750000 },
        { id: 3, name: 'Bob Johnson', email: 'bob@example.com', status: 'PENDING', amount: 2000000 },
        { id: 4, name: 'Alice Brown', email: 'alice@example.com', status: 'SUCCESS', amount: 500000 },
    ];

    // Column definitions
    const columns: Column[] = [
        {
            key: 'id',
            label: 'ID',
            sortable: true,
            width: '80px',
            align: 'center',
        },
        {
            key: 'name',
            label: 'Name',
            sortable: true,
        },
        {
            key: 'email',
            label: 'Email',
            sortable: true,
        },
        {
            key: 'status',
            label: 'Status',
            sortable: true,
            render: (value) => (
                <span
                    className={`${styles.badge} ${value === 'SUCCESS'
                            ? styles.success
                            : value === 'FAILED'
                                ? styles.error
                                : styles.pending
                        }`}
                >
                    {value}
                </span>
            ),
        },
        {
            key: 'amount',
            label: 'Amount',
            sortable: true,
            align: 'right',
            render: (value) => `Rp ${value.toLocaleString('id-ID')}`,
        },
    ];

    // Row decoration based on status
    const getRowDecoration = (row: Row) => {
        if (row.status === 'FAILED') {
            return {
                backgroundColor: '#fef2f2',
                borderColor: '#fecaca',
            };
        }
        if (row.status === 'PENDING') {
            return {
                backgroundColor: '#fffbeb',
                borderColor: '#fde68a',
            };
        }
        return undefined;
    };

    // Action button
    const renderActionButton = (rowData: Row) => (
        <button
            className={styles.actionButton}
            onClick={() => alert(`Viewing details for ${rowData.name}`)}
        >
            View Details
        </button>
    );

    return (
        <div className={styles.container}>
            <div className={styles.header}>
                <h1>Table Component Example</h1>
                <p>A beautiful, sortable table with row decorations and actions</p>
            </div>

            <Table
                data={sampleData}
                columns={columns}
                rowDecoration={getRowDecoration}
                actionButton={renderActionButton}
                initialSort={{ key: 'id', direction: 'asc' }}
            />

            <div className={styles.emptyExample}>
                <h2>Empty State Example</h2>
                <Table
                    data={[]}
                    columns={columns}
                />
            </div>
        </div>
    );
}
