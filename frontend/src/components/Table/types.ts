import { CSSProperties } from 'react';

export interface Row {
    [key: string]: any;
}

export interface Column {
    key: string;
    label: string;
    sortable?: boolean;
    width?: string | number;
    align?: 'left' | 'center' | 'right';
    render?: (value: any, row: Row) => React.ReactNode;
}

export interface SortConfig {
    key: string;
    direction: 'asc' | 'desc';
}

export interface RowDecoration {
    backgroundColor?: string;
    textColor?: string;
    borderColor?: string;
    className?: string;
}

export interface TableProps {
    data?: Row[];
    columns?: Column[];
    containerStyle?: CSSProperties;
    actionButton?: (rowData: Row) => React.ReactNode;
    ListEmptyComponent?: React.ReactNode;
    rowDecoration?: (row: Row, index: number) => RowDecoration | undefined;
    initialSort?: SortConfig;
    onSort?: (sortConfig: SortConfig) => void | Promise<void>;
    showPaginationArrow?: boolean;
    page?: number;
    setPage?: (page: number) => void;
    totalData?: number;
    totalPage?: number;
    pageSize?: number;
    paginationLayout?: 'default' | 'showItemCount';
    isLoading?: boolean;
    useClientSideSort?: boolean;
}
