'use client';

import { CsvUploader } from '@/components/CsvUploader';
import { Container } from '@/components/layout/Container';
import { Table } from '@/components/Table';
import { Column } from '@/components/Table/types';
import { Transaction } from '@/services/transaction/types';
import { Badge } from '@/components/ui/Badge';
import { DEFAULT_PAGE_SIZE, useUploadBankStatement } from './hooks';
import 'react-toastify/dist/ReactToastify.css';
import styles from './UploadBankStatement.module.css';
import { COL_WIDTH_LARGE, COL_WIDTH_SMALL, COL_WIDTH_XSMALL } from '@/components/Table/constant';
import { formatCurrency, formatDate } from '@/lib';
import classNames from 'classnames';

export const UploadBankStatement = () => {
    const {
        uploadBankStatementCSV,
        balance,
        unsuccessfulTransactions,
        pagination,
        isLoadingBalance,
        isLoadingTransactions,
        handlePageChange,
    } = useUploadBankStatement();

    const columns: Column[] = [
        {
            key: 'timestamp',
            label: 'Date',
            sortable: true,
            width: COL_WIDTH_SMALL,
            render: (value: number) => formatDate(value),
        },
        {
            key: 'name',
            label: 'Merchant',
            align: 'left',
            sortable: true,
            width: COL_WIDTH_SMALL,
        },
        {
            key: 'type',
            label: 'Type',
            sortable: true,
            align: 'left',
            render: (value: string) => (
                <Badge variant={value === 'CREDIT' ? 'info' : 'warning'}>
                    {value}
                </Badge>
            ),
            width: COL_WIDTH_XSMALL,
        },
        {
            key: 'amount',
            label: 'Amount',
            sortable: true,
            align: 'left',
            render: (value: number) => formatCurrency(value),
            width: COL_WIDTH_SMALL,
        },
        {
            key: 'status',
            label: 'Status',
            sortable: true,
            align: 'left',
            render: (value: string) => (
                <Badge variant={value === 'FAILED' ? 'danger' : 'warning'}>
                    {value}
                </Badge>
            ),
            width: COL_WIDTH_XSMALL,
        },
        {
            key: 'description',
            label: 'Description',
            width: COL_WIDTH_LARGE,
            align: 'left',
        },
    ];

    const getRowDecoration = (row: Transaction) => {
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


    return (
        <div className={styles.wrapper}>
            <Container maxWidth="xl">
                <div className={styles.header}>
                    <h1 className={styles.title}>Bank Statement Viewer</h1>
                    <p className={styles.description}>
                        Upload your CSV bank statement to view transaction analysis
                    </p>
                </div>
                <div className={styles.uploaderSection}>
                    <CsvUploader
                        onUploadHandler={uploadBankStatementCSV}
                    />
                </div>
                <div className={styles.balanceSection}>
                    <div className={classNames(
                        styles.balanceCard,
                        balance && balance.total_balance < 0 && styles.balanceCardNegative
                    )}>
                        <h2>Total Balance</h2>
                        {isLoadingBalance ? (
                            <div className={styles.loader}>Loading...</div>
                        ) : (
                            <p className={styles.balanceAmount}>
                                {balance ? formatCurrency(balance.total_balance) : '-'}
                            </p>
                        )}
                        <span className={styles.balanceSubtext}>
                            From successful transactions only
                        </span>
                    </div>
                </div>
                <div className={styles.tableSection}>
                    <div className={styles.tableHeader}>
                        <h2>Unsuccessful Transactions</h2>
                        {pagination && (
                            <span className={styles.tableCount}>
                                {pagination.total_elements} transaction
                                {pagination.total_elements !== 1 ? 's' : ''}
                            </span>
                        )}
                    </div>
                    <Table
                        data={unsuccessfulTransactions}
                        columns={columns}
                        rowDecoration={(row) => getRowDecoration(row as Transaction)}
                        showPaginationArrow={true}
                        page={pagination?.current_page || 1}
                        setPage={handlePageChange}
                        totalData={pagination?.total_elements || 0}
                        totalPage={pagination?.total_pages || 0}
                        pageSize={DEFAULT_PAGE_SIZE}
                        paginationLayout="showItemCount"
                        isLoading={isLoadingTransactions}
                    />
                </div>
            </Container>
        </div>
    );
};
