import React from 'react';
import styles from './TableCell.module.css';

interface TableCellProps {
    children: React.ReactNode;
    align?: 'left' | 'center' | 'right';
    className?: string;
}

export const TableCell: React.FC<TableCellProps> = ({ children, align = 'left', className }) => {
    return (
        <td className={`${styles.tableCell} ${styles[`align-${align}`]} ${className || ''}`}>
            {children}
        </td>
    );
};
