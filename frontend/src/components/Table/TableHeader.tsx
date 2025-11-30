import React from 'react';
import styles from './TableHeader.module.css';

interface TableHeaderProps {
    children: React.ReactNode;
}

export const TableHeader: React.FC<TableHeaderProps> = ({ children }) => {
    return <thead className={styles.tableHeader}>{children}</thead>;
};
