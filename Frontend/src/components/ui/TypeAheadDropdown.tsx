import { useState, useEffect } from 'react';

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

    useEffect(() => {
        setSearchQuery(value?.label || '');
    }, [value]);

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
            {label && <label htmlFor={id}>{label}</label>}
            <div className="ta-select">
                <input
                    id={id}
                    type="text"
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