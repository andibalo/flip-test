import { toast } from 'react-toastify';

type ToastType = 'success' | 'error' | 'warning' | 'info';

export function useToast() {
    const showToast = ({
        type,
        title,
        description,
    }: {
        type: ToastType;
        title: string;
        description?: string;
    }) => {
        const message = description ? `${title}: ${description}` : title;

        switch (type) {
            case 'success':
                toast.success(message);
                break;
            case 'error':
                toast.error(message);
                break;
            case 'warning':
                toast.warning(message);
                break;
            case 'info':
                toast.info(message);
                break;
        }
    };

    return { showToast };
}

export default useToast;
