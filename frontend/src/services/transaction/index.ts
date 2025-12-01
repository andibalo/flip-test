import { apiRequest, ApiSuccessResponse, PaginationMetadata } from '@/services';
import {
    UploadStatementCsvRequest,
    UploadStatementCsvResponse,
    GetBalanceResponse,
    GetUnsuccessfulTransactionsRequest,
    GetUnsuccessfulTransactionsResponse,
} from './types';

export async function uploadStatementCsv(
    req: UploadStatementCsvRequest
): Promise<UploadStatementCsvResponse> {
    const formData = new FormData();
    formData.append('file', req.file);

    const response = await apiRequest<ApiSuccessResponse<UploadStatementCsvResponse>>(
        'post',
        '/upload',
        formData,
        false
    );

    return response.data.data;
}

export async function getBalance(): Promise<GetBalanceResponse> {
    const response = await apiRequest<ApiSuccessResponse<GetBalanceResponse>>(
        'get',
        '/balance',
        undefined,
        false
    );

    return response.data.data;
}

export async function getUnsuccessfulTransactions(
    req?: GetUnsuccessfulTransactionsRequest
): Promise<{
    data: GetUnsuccessfulTransactionsResponse;
    pagination?: PaginationMetadata;
}> {
    const params = new URLSearchParams();

    if (req?.page) {
        params.append('page', req.page.toString());
    }

    if (req?.page_size) {
        params.append('page_size', req.page_size.toString());
    }

    if (req?.sorts) {
        params.append('sorts', req.sorts);
    }

    const url = params.toString() ? `/issues?${params.toString()}` : '/issues';

    const response = await apiRequest<
        ApiSuccessResponse<GetUnsuccessfulTransactionsResponse>
    >('get', url, undefined, false);

    return {
        data: response.data.data,
        pagination: response.data.pagination,
    };
}
