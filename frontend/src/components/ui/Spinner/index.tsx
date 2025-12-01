import React from 'react';
import classNames from 'classnames';
import styles from './Spinner.module.css';

export type SpinnerSize = 'sm' | 'md' | 'lg';

interface SpinnerProps {
    size?: SpinnerSize;
    className?: string;
    color?: string;
}

export const Spinner: React.FC<SpinnerProps> = ({
    size = 'md',
    className,
    color,
}) => {
    const style = color ? { borderTopColor: color, borderRightColor: color } : undefined;

    return (
        <span
            className={classNames(styles.spinner, styles[size], className)}
            style={style}
            role="status"
            aria-label="Loading"
        />
    );
};
