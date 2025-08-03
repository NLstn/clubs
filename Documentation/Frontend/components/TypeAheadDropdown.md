# TypeAheadDropdown Component

A reusable type-ahead dropdown component that provides search functionality with autocomplete suggestions.

## Import

```tsx
import TypeAheadDropdown from '@/components/ui/TypeAheadDropdown';
```

## Basic Usage

```tsx
interface User {
    id: string;
    label: string; // Required for the Option interface
    name: string;
    email: string;
}

const users: User[] = [
    { id: '1', label: 'John Doe', name: 'John Doe', email: 'john@example.com' },
    { id: '2', label: 'Jane Smith', name: 'Jane Smith', email: 'jane@example.com' },
];

const [selectedUser, setSelectedUser] = useState<User | null>(null);
const [filteredUsers, setFilteredUsers] = useState<User[]>(users);

const handleSearch = (query: string) => {
    const filtered = users.filter(user => 
        user.label.toLowerCase().includes(query.toLowerCase())
    );
    setFilteredUsers(filtered);
};

<TypeAheadDropdown<User>
    options={filteredUsers}
    value={selectedUser}
    onChange={setSelectedUser}
    onSearch={handleSearch}
    placeholder="Search users..."
    id="user-select"
    label="Select User"
/>
```

## Props

### TypeAheadDropdownProps<T>

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `options` | `T[]` | Required | Array of options to display in dropdown |
| `value` | `T \| null` | Required | Currently selected option |
| `onChange` | `(value: T \| null) => void` | Required | Callback when selection changes |
| `onSearch` | `(query: string) => void` | Required | Callback when user types in search box |
| `placeholder` | `string` | `'Search...'` | Placeholder text for the input |
| `id` | `string` | `undefined` | HTML id attribute for the input |
| `label` | `string` | `undefined` | Label text displayed above the input |

### Option Interface

All options must extend the `Option` interface:

```tsx
interface Option {
    id: string;
    label: string;
}
```

The `label` property is what gets displayed in the dropdown and input field.

## Advanced Examples

### With Custom Option Types

```tsx
interface MemberOption {
    id: string;
    label: string; // This will be the display name
    userId: string;
    role: string;
}

interface FineTemplateOption {
    id: string;
    label: string; // This will be the template name
    amount: number;
    description: string;
}

// Usage
<TypeAheadDropdown<MemberOption>
    options={memberOptions}
    value={selectedMember}
    onChange={setSelectedMember}
    onSearch={handleMemberSearch}
    placeholder="Search member..."
    id="member"
    label="Member"
/>

<TypeAheadDropdown<FineTemplateOption>
    options={templateOptions}
    value={selectedTemplate}
    onChange={setSelectedTemplate}
    onSearch={handleTemplateSearch}
    placeholder="Search fine template (optional)..."
    id="template"
    label="Fine Template (Optional)"
/>
```

### With Search Filtering

```tsx
const [searchQuery, setSearchQuery] = useState('');
const [allUsers, setAllUsers] = useState<User[]>([]);
const [filteredUsers, setFilteredUsers] = useState<User[]>([]);

const handleUserSearch = (query: string) => {
    setSearchQuery(query);
    
    if (!query.trim()) {
        setFilteredUsers(allUsers);
        return;
    }
    
    const filtered = allUsers.filter(user =>
        user.label.toLowerCase().includes(query.toLowerCase()) ||
        user.email.toLowerCase().includes(query.toLowerCase())
    );
    
    setFilteredUsers(filtered);
};

// Clear selection when search query changes but no exact match
useEffect(() => {
    if (selectedUser && searchQuery && !selectedUser.label.includes(searchQuery)) {
        setSelectedUser(null);
    }
}, [searchQuery, selectedUser]);
```

### With API Integration

```tsx
const [options, setOptions] = useState<User[]>([]);
const [loading, setLoading] = useState(false);

const handleSearch = async (query: string) => {
    if (!query.trim()) {
        setOptions([]);
        return;
    }
    
    setLoading(true);
    try {
        const response = await api.get(`/api/users/search?q=${encodeURIComponent(query)}`);
        setOptions(response.data);
    } catch (error) {
        console.error('Search failed:', error);
        setOptions([]);
    } finally {
        setLoading(false);
    }
};

<TypeAheadDropdown<User>
    options={options}
    value={selectedUser}
    onChange={setSelectedUser}
    onSearch={handleSearch}
    placeholder={loading ? "Searching..." : "Search users..."}
    id="user-search"
    label="Find User"
/>
```

## Behavior

### Selection Logic
- When user types, `onSearch` is called with the current query
- If input is cleared, `onChange` is called with `null`
- When an option is selected, input shows the option's `label` and dropdown closes
- Clicking on an option calls `onChange` with the selected option

### Dropdown Visibility
- Dropdown shows when:
  - Input is focused AND
  - There's a search query AND
  - No option is currently selected AND
  - There are options available
- Dropdown hides when:
  - An option is selected
  - Input loses focus (not shown in current implementation)

## Styling

### CSS Classes

The component uses these CSS classes:

- `.typeahead-container` - Main container
- `.ta-select` - Input wrapper
- `.ta-dropdown` - Dropdown container
- `.ta-option` - Individual option items

### Custom Styling

Since the component doesn't have its own CSS file, you'll need to add styles:

```css
.typeahead-container {
    position: relative;
    margin-bottom: 1rem;
}

.typeahead-container label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: rgba(255, 255, 255, 0.9);
}

.ta-select {
    position: relative;
}

.ta-select input {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #444;
    border-radius: 4px;
    background: #1a1a1a;
    color: rgba(255, 255, 255, 0.9);
    font-size: 1rem;
}

.ta-select input:focus {
    outline: none;
    border-color: #2196f3;
    box-shadow: 0 0 0 2px rgba(33, 150, 243, 0.2);
}

.ta-dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background: #2a2a2a;
    border: 1px solid #444;
    border-top: none;
    border-radius: 0 0 4px 4px;
    max-height: 200px;
    overflow-y: auto;
    z-index: 1000;
}

.ta-option {
    padding: 0.75rem;
    cursor: pointer;
    border-bottom: 1px solid #333;
    color: rgba(255, 255, 255, 0.85);
}

.ta-option:hover {
    background: #333;
}

.ta-option:last-child {
    border-bottom: none;
}
```

## Accessibility

To improve accessibility, consider adding:

```tsx
<input
    id={id}
    type="text"
    value={searchQuery}
    onChange={(e) => handleInputChange(e.target.value)}
    placeholder={placeholder}
    autoComplete="off"
    onFocus={() => setIsOpen(true)}
    aria-expanded={isOpen}
    aria-haspopup="listbox"
    role="combobox"
    aria-describedby={`${id}-help`}
/>
```

## Examples in the Codebase

The TypeAheadDropdown component is used in:
- **AddFine.tsx** - For selecting members and fine templates

### Real Example from AddFine

```tsx
interface MemberOption {
    id: string;
    label: string;
    userId: string;
}

interface FineTemplateOption {
    id: string;
    label: string;
    amount: number;
    description: string;
}

<TypeAheadDropdown<MemberOption>
    options={memberOptions}
    value={selectedOption}
    onChange={setSelectedOption}
    onSearch={handleSearch}
    placeholder="Search member..."
    id="member"
    label="Member"
/>

<TypeAheadDropdown<FineTemplateOption>
    options={templateOptions}
    value={selectedTemplate}
    onChange={handleTemplateSelection}
    onSearch={handleTemplateSearch}
    placeholder="Search fine template (optional)..."
    id="template"
    label="Fine Template (Optional)"
/>
```

## Best Practices

1. **Always implement search filtering** in the `onSearch` callback
2. **Use descriptive label properties** for better user experience
3. **Handle empty states** appropriately
4. **Consider performance** for large option lists
5. **Provide clear placeholder text**
6. **Handle loading states** when fetching data asynchronously
7. **Clear selection when appropriate** (e.g., when search changes)
8. **Add proper CSS styling** since component has no default styles
