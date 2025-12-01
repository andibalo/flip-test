import React from 'react';
import classNames from 'classnames';
import styles from './Badge.module.css';

export type BadgeVariant = 'success' | 'danger' | 'warning' | 'info';

interface BadgeProps {
    children: React.ReactNode;
    variant?: BadgeVariant;
    className?: string;
}

export const Badge: React.FC<BadgeProps> = ({
    children,
    variant = 'info',
    className,
}) => {
    return (
        <span className={classNames(styles.badge, styles[variant], className)}>
            {children}
        </span>
    );
};
