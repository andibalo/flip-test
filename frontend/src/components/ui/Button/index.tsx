import React from 'react';
import classNames from 'classnames';
import { Spinner } from '../Spinner';
import styles from './Button.module.css';

interface ButtonProps {
    children: React.ReactNode;
    onClick?: () => void;
    variant?: 'primary' | 'secondary' | 'success' | 'danger' | 'outline';
    size?: 'sm' | 'md' | 'lg';
    disabled?: boolean;
    type?: 'button' | 'submit' | 'reset';
    fullWidth?: boolean;
    className?: string;
    style?: React.CSSProperties;
    isLoading?: boolean;
}

export const Button: React.FC<ButtonProps> = ({
    children,
    onClick,
    variant = 'primary',
    size = 'md',
    disabled = false,
    type = 'button',
    fullWidth = false,
    className,
    style,
    isLoading = false,
}) => {
    const isDisabled = disabled || isLoading;

    return (
        <button
            type={type}
            onClick={onClick}
            disabled={isDisabled}
            className={classNames(
                styles.button,
                styles[variant],
                styles[size],
                fullWidth && styles.fullWidth,
                className
            )}
            style={style}
        >
            {isLoading && <Spinner size={size === 'sm' ? 'sm' : size === 'lg' ? 'md' : 'sm'} />}
            {!isLoading && children}
        </button>
    );
};
