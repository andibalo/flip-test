import React from 'react';
import classNames from 'classnames';
import styles from './TableCell.module.css';

interface TableCellProps {
    children: React.ReactNode;
    align?: 'left' | 'center' | 'right';
    className?: string;
}

export const TableCell: React.FC<TableCellProps> = ({ children, align = 'left', className }) => {
    return (
        <td className={classNames(styles.tableCell, styles[`align-${align}`], className)}>
            {children}
        </td>
    );
};
