import React from 'react';
import { FiChevronLeft, FiChevronRight } from 'react-icons/fi';
import { Button } from '@/components/ui/Button';
import styles from './PaginationArrow.module.css';

interface PaginationArrowProps {
    page: number;
    setPage: (page: number) => void;
    totalData?: number;
    totalCurrentData: number;
    pageSize: number;
    totalPage?: number;
    layout?: 'default' | 'showItemCount';
}

export const PaginationArrow: React.FC<PaginationArrowProps> = ({
    page,
    setPage,
    totalData = 0,
    totalPage = 0,
    totalCurrentData,
    pageSize,
    layout = 'default',
}) => {
    const isFirst = page === 1;
    const isLast = totalPage > 0 ? page >= totalPage : totalCurrentData < pageSize;

    if (layout === 'showItemCount') {
        let perPageStartItemCount = 0;
        let perPageEndItemCount = 0;

        if (totalData > 0) {
            perPageStartItemCount = (page - 1) * pageSize + 1;
            perPageEndItemCount = (page - 1) * pageSize + totalCurrentData;
        }

        return (
            <div className={styles.paginationContainer}>
                <div className={styles.itemCount}>
                    <span className={styles.countText}>
                        {perPageStartItemCount}-{perPageEndItemCount} dari {totalData} hasil
                    </span>
                </div>
                <div className={styles.buttonGroup}>
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setPage(page - 1)}
                        disabled={isFirst}
                        className={styles.arrowButton}
                    >
                        <FiChevronLeft size={16} />
                    </Button>
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setPage(page + 1)}
                        disabled={isLast}
                        className={styles.arrowButton}
                    >
                        <FiChevronRight size={16} />
                    </Button>
                </div>
            </div>
        );
    }

    return (
        <div className={styles.buttonGroupRight}>
            <Button
                variant="outline"
                size="sm"
                onClick={() => setPage(page - 1)}
                disabled={isFirst}
                className={styles.arrowButton}
            >
                <FiChevronLeft size={16} />
            </Button>
            <Button
                variant="outline"
                size="sm"
                onClick={() => setPage(page + 1)}
                disabled={isLast}
                className={styles.arrowButton}
            >
                <FiChevronRight size={16} />
            </Button>
        </div>
    );
};
