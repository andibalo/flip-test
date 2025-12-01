export interface Transaction {
    timestamp: number;
    name: string;
    type: 'CREDIT' | 'DEBIT';
    amount: number;
    status: 'SUCCESS' | 'FAILED' | 'PENDING';
    description: string;
}

export interface UploadStatementCsvRequest {
    file: File;
}

export interface UploadStatementCsvResponse {
    transactions_uploaded: number;
}

export interface GetBalanceResponse {
    total_balance: number;
}

export interface GetUnsuccessfulTransactionsRequest {
    page?: number;
    page_size?: number;
    sorts?: string
}

export interface GetUnsuccessfulTransactionSummary {
    total_count: number;
    pending_count: number;
    failed_count: number;
}

export interface GetUnsuccessfulTransactionsResponse {
    transactions: Transaction[];
    summary: GetUnsuccessfulTransactionSummary;
}

