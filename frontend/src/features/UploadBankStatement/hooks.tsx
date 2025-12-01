'use client';

import { useState } from 'react';
import useToast from '@/hooks/useToast';
import {
    uploadStatementCsv,
    getBalance,
    getUnsuccessfulTransactions,
} from '@/services/transaction';
import {
    GetBalanceResponse,
    Transaction,
} from '@/services/transaction/types';
import { getErrorMessageFromAxiosError, PaginationMetadata } from '@/services';
import axios from 'axios';
import { SortConfig } from '@/components/Table/types';

export const DEFAULT_PAGE_SIZE = 10;

export function useUploadBankStatement() {
    const { showToast } = useToast();
    const [isUploadingCSV, setIsUploadingCSV] = useState(false);
    const [isLoadingBalance, setIsLoadingBalance] = useState(false);
    const [isLoadingTransactions, setIsLoadingTransactions] = useState(false);
    const [balance, setBalance] = useState<GetBalanceResponse | null>(null);
    const [unsuccessfulTransactions, setUnsuccessfulTransactions] = useState<
        Transaction[]
    >([]);
    const [pagination, setPagination] = useState<PaginationMetadata | null>(null);

    const uploadBankStatementCSV = async (file: File) => {
        setIsUploadingCSV(true);
        try {
            const result = await uploadStatementCsv({ file });

            showToast({
                type: 'success',
                title: 'Upload Successful',
                description: `Successfully uploaded ${result.transactions_uploaded} transactions`,
            });

            await Promise.all([fetchBalance(), fetchUnsuccessfulTransactions(1, DEFAULT_PAGE_SIZE, "-timestamp")]);

        } catch (error) {
            if (axios.isAxiosError(error)) {
                showToast({
                    type: 'error',
                    title: 'Upload Failed',
                    description: getErrorMessageFromAxiosError(error),
                });
                throw error;
            }

            showToast({
                type: 'error',
                title: 'Upload Failed',
                description: 'Failed to upload CSV file',
            });
            throw error;
        } finally {
            setIsUploadingCSV(false);
        }
    };

    const fetchBalance = async () => {
        setIsLoadingBalance(true);
        try {
            const result = await getBalance();
            setBalance(result);
            return result;
        } catch (error) {
            throw error;
        } finally {
            setIsLoadingBalance(false);
        }
    };

    const fetchUnsuccessfulTransactions = async (page?: number, pageSize?: number, sort?: string) => {
        setIsLoadingTransactions(true);
        try {
            const result = await getUnsuccessfulTransactions({
                page,
                page_size: pageSize,
                sorts: sort,
            });
            setUnsuccessfulTransactions(result.data.transactions);
            setPagination(result.pagination || null);
            return result;
        } catch (error) {
            throw error;
        } finally {
            setIsLoadingTransactions(false);
        }
    };

    const handlePageChange = async (page: number) => {
        await fetchUnsuccessfulTransactions(page, DEFAULT_PAGE_SIZE);
    };

    const handleSort = async (sortConfig: SortConfig) => {
        const page = 1;
        const pageSize = DEFAULT_PAGE_SIZE;
        const sort = `${sortConfig.direction === 'asc' ? '+' : '-'}${sortConfig.key}`

        await fetchUnsuccessfulTransactions(page, pageSize, sort);
    };

    return {
        uploadBankStatementCSV,
        fetchBalance,
        fetchUnsuccessfulTransactions,
        isUploadingCSV,
        isLoadingBalance,
        isLoadingTransactions,
        balance,
        unsuccessfulTransactions,
        pagination,
        handlePageChange,
        handleSort,
    };
}