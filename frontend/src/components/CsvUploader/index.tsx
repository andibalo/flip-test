'use client';

import React, { useState, useCallback, DragEvent, ChangeEvent } from 'react';
import { FiUploadCloud, FiFile, FiX } from 'react-icons/fi';
import useToast from '@/hooks/useToast';
import styles from './CsvUploader.module.css';

interface CsvUploaderProps {
    onUploadError?: (error: any) => void;
    onUploadHandler: (file: File) => Promise<void>;
}

export const CsvUploader: React.FC<CsvUploaderProps> = ({
    onUploadError,
    onUploadHandler,
}) => {
    const { showToast } = useToast();
    const [isDragging, setIsDragging] = useState(false);
    const [isUploading, setIsUploading] = useState(false);
    const [selectedFile, setSelectedFile] = useState<File | null>(null);

    const validateFile = (file: File): boolean => {
        const validTypes = ['text/csv', 'application/vnd.ms-excel'];
        const isValidType = validTypes.includes(file.type) || file.name.endsWith('.csv');

        if (!isValidType) {
            showToast({
                type: 'error',
                title: 'Invalid File Type',
                description: 'Please upload a CSV file',
            });
            onUploadError?.('Invalid file type');
            return false;
        }

        return true;
    };

    const handleDragEnter = useCallback((e: DragEvent<HTMLDivElement>) => {
        e.preventDefault();
        e.stopPropagation();
        setIsDragging(true);
    }, []);

    const handleDragLeave = useCallback((e: DragEvent<HTMLDivElement>) => {
        e.preventDefault();
        e.stopPropagation();
        setIsDragging(false);
    }, []);

    const handleDragOver = useCallback((e: DragEvent<HTMLDivElement>) => {
        e.preventDefault();
        e.stopPropagation();
    }, []);

    const handleDrop = useCallback(
        (e: DragEvent<HTMLDivElement>) => {
            e.preventDefault();
            e.stopPropagation();
            setIsDragging(false);

            const files = e.dataTransfer.files;
            if (files && files.length > 0) {
                const file = files[0];
                if (validateFile(file)) {
                    setSelectedFile(file);
                }
            }
        },
        [validateFile]
    );

    const handleFileSelect = useCallback(
        (e: ChangeEvent<HTMLInputElement>) => {
            const files = e.target.files;
            if (files && files.length > 0) {
                const file = files[0];
                if (validateFile(file)) {
                    setSelectedFile(file);
                }
            }
        },
        [validateFile]
    );

    const handleDropzoneClick = () => {
        if (!selectedFile) {
            document.getElementById('csv-upload')?.click();
        }
    };

    const handleUpload = async () => {
        if (!selectedFile) return;

        setIsUploading(true);

        try {
            await onUploadHandler(selectedFile);

            setSelectedFile(null);
        } catch (error) {
            onUploadError?.(error);
        } finally {
            setIsUploading(false);
        }
    };

    const handleClear = (e: React.MouseEvent) => {
        e.stopPropagation();
        setSelectedFile(null);
    };

    return (
        <div className={styles.container}>
            <div
                className={`${styles.dropzone} ${isDragging ? styles.dragging : ''} ${selectedFile ? styles.hasFile : ''
                    }`}
                onClick={handleDropzoneClick}
                onDragEnter={handleDragEnter}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
            >
                {!selectedFile ? (
                    <>
                        <div className={styles.uploadIcon}>
                            <FiUploadCloud size={48} />
                        </div>
                        <h3 className={styles.title}>Upload CSV File</h3>
                        <p className={styles.description}>
                            Drag and drop your CSV file here, or click to browse
                        </p>
                        <input
                            type="file"
                            accept=".csv,text/csv,application/vnd.ms-excel"
                            onChange={handleFileSelect}
                            className={styles.fileInput}
                            id="csv-upload"
                        />
                        <label htmlFor="csv-upload" className={styles.browseButton}>
                            Browse Files
                        </label>
                    </>
                ) : (
                    <div className={styles.fileDetails}>
                        <div className={styles.fileIcon}>
                            <FiFile size={40} />
                        </div>
                        <div className={styles.fileInfo}>
                            <h4 className={styles.fileName}>{selectedFile.name}</h4>
                            <p className={styles.fileSize}>
                                {(selectedFile.size / 1024).toFixed(2)} KB
                            </p>
                        </div>
                        <button
                            onClick={handleClear}
                            className={styles.clearButton}
                            disabled={isUploading}
                            aria-label="Clear file"
                        >
                            <FiX size={20} />
                        </button>
                    </div>
                )}
            </div>

            {selectedFile && (
                <button
                    onClick={handleUpload}
                    disabled={isUploading}
                    className={`${styles.uploadButton} ${isUploading ? styles.uploading : ''}`}
                >
                    {isUploading ? (
                        <>
                            <span className={styles.spinner}></span>
                            Uploading...
                        </>
                    ) : (
                        <>
                            <FiUploadCloud size={18} />
                            Upload Transactions
                        </>
                    )}
                </button>
            )}
        </div>
    );
};
