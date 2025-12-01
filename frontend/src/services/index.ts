import axios, { AxiosError, AxiosRequestConfig, AxiosResponse } from "axios";
import { CONFIG } from "@/config";

export interface PaginationMetadata {
    current_page: number;
    current_elements: number;
    total_pages: number;
    total_elements: number;
}

export interface ApiSuccessResponse<T> {
    data: T;
    pagination?: PaginationMetadata;
    success: string;
}

export interface ApiErrorResponse {
    metadata: {
        path: string;
        code: string;
        statusCode: number;
        status: string;
        message: string;
        error: string;
        timestamp: string;
    };
    success: string;
}

const apiInstance = axios.create({
    baseURL: CONFIG.BE_API_URL,
    timeout: 30000,
});

async function apiRequest<T = any, R = AxiosResponse<T>, D = any>(
    method: "get" | "post" | "put" | "patch" | "delete",
    url: string,
    data?: D,
    auth?: boolean,
    config: AxiosRequestConfig<D> = {},
): Promise<R> {
    if (method === "post" && data instanceof FormData) {
        config.headers = {
            ...config.headers,
        };
    } else if (method === "post") {
        config.headers = {
            ...config.headers,
            "Content-Type": "application/json",
        };
    }

    return apiInstance.request({
        method,
        url,
        data,
        ...config,
    });
}

apiInstance.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.status === 401) {
            // TODO: implement
        }

        if (error.status === 403) {
            // TODO: implement
        }

        if (error.code === "ECONNABORTED" || error.message?.includes("timeout")) {
            const timeoutError = {
                ...error,
                error: {
                    code: "TIMEOUT",
                    message: "CONNECTION_TIMEOUT",
                    errors: "Request timeout",
                },
            };
            return Promise.reject(timeoutError);
        }

        if (error.response) {
            console.error("Error response:", error.response.data);
            console.error("Error status:", error.response.status);
        } else if (error.request) {
            console.error("Error request:", error.request);
        } else {
            console.error("Error message:", error.message);
        }
        return Promise.reject(error);
    },
);

const getErrorMessageFromAxiosError = (
    error: AxiosError<ApiErrorResponse>,
): string => {
    if (error.response && error.response.data) {
        const apiErrorMsg = error.response.data.metadata.message;
        return apiErrorMsg || "An error occurred";
    }

    return error.message || "An error occurred";
};

export {
    apiRequest,
    getErrorMessageFromAxiosError,
};
