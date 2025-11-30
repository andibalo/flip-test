import React from 'react';
import styles from './TableRow.module.css';

interface TableRowProps {
    children: React.ReactNode;
    decoration?: {
        backgroundColor?: string;
        textColor?: string;
        borderColor?: string;
        className?: string;
    };
    onClick?: () => void;
}

export const TableRow: React.FC<TableRowProps> = ({ children, decoration, onClick }) => {
    const style: React.CSSProperties = {};

    if (decoration?.backgroundColor) style.backgroundColor = decoration.backgroundColor;
    if (decoration?.textColor) style.color = decoration.textColor;
    if (decoration?.borderColor) style.borderColor = decoration.borderColor;

    return (
        <tr
            className={`${styles.tableRow} ${decoration?.className || ''} ${onClick ? styles.clickable : ''}`}
            style={style}
            onClick={onClick}
        >
            {children}
        </tr>
    );
};
