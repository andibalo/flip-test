import React from 'react';
import classNames from 'classnames';
import styles from './Container.module.css';

interface ContainerProps {
    children: React.ReactNode;
    maxWidth?: 'sm' | 'md' | 'lg' | 'xl' | 'full';
    className?: string;
    style?: React.CSSProperties;
}

export const Container: React.FC<ContainerProps> = ({
    children,
    maxWidth = 'lg',
    className,
    style
}) => {
    return (
        <div
            className={classNames(styles.container, styles[maxWidth], className)}
            style={style}
        >
            {children}
        </div>
    );
};
