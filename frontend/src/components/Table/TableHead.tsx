import React from 'react';
import classNames from 'classnames';
import styles from './TableHead.module.css';

interface TableHeadProps {
    children: React.ReactNode;
    className?: string;
}

export const TableHead: React.FC<TableHeadProps> = ({ children, className }) => {
    return <th className={classNames(styles.tableHead, className)}>{children}</th>;
};
