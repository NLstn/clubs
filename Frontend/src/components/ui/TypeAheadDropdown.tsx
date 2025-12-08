import { useState, useEffect, useRef } from 'react';
import { Input } from './Input';
import './TypeAheadDropdown.css';

interface Option {
    id: string;
    label: string;
}

interface TypeAheadDropdownProps<T extends Option> {
    options: T[];
    value: T | null;
    onChange: (value: T | null) => void;
    onSearch: (query: string) => void;
    placeholder?: string;
    id?: string;
    label?: string;
}

export const TypeAheadDropdown = <T extends Option>({
    options,
    value,
    onChange,
    onSearch,
    placeholder = 'Search...',
    id,
    label
}: TypeAheadDropdownProps<T>) => {
    const [searchQuery, setSearchQuery] = useState<string>(value?.label || '');
    const [isOpen, setIsOpen] = useState(false);
    const prevValueRef = useRef(value);

    // Sync search query with value changes
    useEffect(() => {
        // Only update if value actually changed (not just re-render)
        if (prevValueRef.current !== value) {
            prevValueRef.current = value;
            const newQuery = value?.label || '';
            if (searchQuery !== newQuery) {
                // Synchronizing with external prop value is a valid use case
                // eslint-disable-next-line react-hooks/set-state-in-effect
                setSearchQuery(newQuery);
            }
        }
    }, [value, searchQuery]);

    const handleInputChange = (query: string) => {
        setSearchQuery(query);
        setIsOpen(true);
        if (!query) {
            onChange(null);
        }
        onSearch(query);
    };

    const handleSelectOption = (option: T) => {
        onChange(option);
        setSearchQuery(option.label);
        setIsOpen(false);
    };

    return (
        <div className="typeahead-container">
            <div className="ta-select">
                <Input
                    id={id}
                    label={label}
                    value={searchQuery}
                    onChange={(e) => handleInputChange(e.target.value)}
                    placeholder={placeholder}
                    autoComplete="off"
                    onFocus={() => setIsOpen(true)}
                />
                {isOpen && searchQuery && !value && options.length > 0 && (
                    <div className="ta-dropdown">
                        {options.map(option => (
                            <div
                                key={option.id}
                                className="ta-option"
                                onClick={() => handleSelectOption(option)}
                            >
                                {option.label}
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
};

export default TypeAheadDropdown;