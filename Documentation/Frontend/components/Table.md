# Table Component

A flexible, reusable table component with support for custom rendering, loading states, error handling, and responsive design.

## Import

```tsx
import Table, { TableColumn } from '@/components/ui/Table';
```

## Basic Usage

```tsx
interface User {
  id: string;
  name: string;
  email: string;
  role: string;
}

const users: User[] = [
  { id: '1', name: 'John Doe', email: 'john@example.com', role: 'Admin' },
  { id: '2', name: 'Jane Smith', email: 'jane@example.com', role: 'User' },
];

const columns: TableColumn<User>[] = [
  {
    key: 'name',
    header: 'Name',
    render: (user) => user.name,
  },
  {
    key: 'email',
    header: 'Email',
    render: (user) => user.email,
  },
  {
    key: 'role',
    header: 'Role',
    render: (user) => <span className={`role-${user.role.toLowerCase()}`}>{user.role}</span>,
  },
];

<Table
  columns={columns}
  data={users}
  keyExtractor={(user) => user.id}
/>
```

## Props

### TableProps<T>

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `columns` | `TableColumn<T>[]` | Required | Array of column definitions |
| `data` | `T[]` | Required | Array of data to display |
| `keyExtractor` | `(item: T) => string` | Required | Function to extract unique key from each item |
| `className` | `string` | `''` | Additional CSS classes for the table container |
| `emptyMessage` | `string` | `'No data available'` | Message to display when data array is empty |
| `footer` | `React.ReactNode` | `undefined` | Optional footer content |
| `loading` | `boolean` | `false` | Shows loading state when true |
| `error` | `string \| null` | `null` | Error message to display |
| `loadingMessage` | `string` | `'Loading...'` | Custom loading message |
| `errorMessage` | `string` | `'Error loading data'` | Custom error message |

### TableColumn<T>

| Prop | Type | Description |
|------|------|-------------|
| `key` | `string` | Unique identifier for the column |
| `header` | `string` | Column header text |
| `render` | `(item: T) => React.ReactNode` | Function to render cell content |
| `className` | `string` | Optional CSS class for column cells |

## Advanced Examples

### With Loading State

```tsx
const [loading, setLoading] = useState(true);
const [users, setUsers] = useState<User[]>([]);

<Table
  columns={columns}
  data={users}
  keyExtractor={(user) => user.id}
  loading={loading}
  loadingMessage="Fetching users..."
/>
```

### With Error Handling

```tsx
const [error, setError] = useState<string | null>(null);

<Table
  columns={columns}
  data={users}
  keyExtractor={(user) => user.id}
  error={error}
  errorMessage="Failed to load users. Please try again."
/>
```

### With Custom Actions

```tsx
const columns: TableColumn<User>[] = [
  {
    key: 'name',
    header: 'Name',
    render: (user) => user.name,
  },
  {
    key: 'actions',
    header: 'Actions',
    render: (user) => (
      <div className="table-actions">
        <button 
          className="action-button edit"
          onClick={() => handleEdit(user.id)}
        >
          Edit
        </button>
        <button 
          className="action-button remove"
          onClick={() => handleDelete(user.id)}
        >
          Remove
        </button>
      </div>
    ),
  },
];
```

### With Footer

```tsx
<Table
  columns={columns}
  data={users}
  keyExtractor={(user) => user.id}
  footer={
    <div>
      Showing {users.length} of {totalUsers} users
    </div>
  }
/>
```

### Responsive Columns

```tsx
const columns: TableColumn<User>[] = [
  {
    key: 'name',
    header: 'Name',
    render: (user) => user.name,
  },
  {
    key: 'email',
    header: 'Email',
    render: (user) => user.email,
    className: 'hide-mobile', // Hides on mobile devices
  },
  {
    key: 'created',
    header: 'Created',
    render: (user) => formatDate(user.createdAt),
    className: 'hide-small', // Hides on small screens
  },
];
```

## Styling

### CSS Classes

The Table component uses these CSS classes:

- `.table-container` - Main container
- `.reusable-table` - Table element
- `.table-footer` - Footer container
- `.table-loading-text` - Loading message
- `.table-error-text` - Error message
- `.table-empty-text` - Empty state message
- `.table-actions` - Action buttons container
- `.action-button` - Individual action button

### Action Button Types

Pre-defined action button styles:

```css
.action-button.edit     /* Blue - for edit actions */
.action-button.remove   /* Red - for delete actions */
.action-button.promote  /* Green - for promote actions */
.action-button.demote   /* Orange - for demote actions */
```

### Responsive Classes

- `.hide-mobile` - Hides column on screens < 768px
- `.hide-small` - Hides column on screens < 480px

### Custom Styling

You can customize the table appearance by:

1. **Adding custom CSS classes:**
   ```tsx
   <Table className="my-custom-table" ... />
   ```

2. **Using CSS custom properties:**
   ```css
   .my-custom-table {
     --space-md: 1rem; /* Custom spacing */
   }
   ```

3. **Column-specific styling:**
   ```tsx
   {
     key: 'status',
     header: 'Status',
     render: (item) => item.status,
     className: 'status-column'
   }
   ```

## Accessibility

The Table component includes:

- Proper semantic HTML structure (`<table>`, `<thead>`, `<tbody>`)
- Keyboard navigation support
- Screen reader friendly markup
- Proper ARIA labels (can be extended)

## Performance Considerations

- Use `React.memo()` for expensive render functions
- Consider virtualization for large datasets (>1000 rows)
- Implement pagination for better performance with large datasets

## Examples in the Codebase

The Table component is used throughout the application in:
- **AdminClubMemberList.tsx** - Member management with role changes and permissions
- **AdminClubEventList.tsx** - Event listings
- **AdminClubFineList.tsx** - Fine management
- **ProfileFines.tsx** - User's personal fines
- **ProfileInvites.tsx** - User's invitations
- **ProfileSessions.tsx** - User's active sessions
- **EventRSVPList.tsx** - Event RSVP management

### Real Example from AdminClubMemberList

```tsx
interface Member {
    id: string;
    name: string;
    role: string;
    joinedAt: string;
    userId?: string;
    birthDate?: string;
}

const columns: TableColumn<Member>[] = [
    {
        key: 'name',
        header: 'Name',
        render: (member) => member.name
    },
    {
        key: 'role',
        header: 'Role',
        render: (member) => translateRole(member.role)
    },
    {
        key: 'joined',
        header: 'Joined',
        render: (member) => member.joinedAt ? new Date(member.joinedAt).toLocaleDateString() : 'N/A'
    },
    {
        key: 'birthDate',
        header: 'Birth Date',
        render: (member) => member.birthDate ? new Date(member.birthDate).toLocaleDateString() : 'Not shared'
    },
    {
        key: 'actions',
        header: 'Actions',
        render: (member) => (
            <div className="member-actions">
                {canDeleteMember(currentUserRole, member.role) && (
                    <button
                        onClick={() => deleteMember(member.id)}
                        className="action-button remove"
                        aria-label="Remove member"
                    >
                        Remove
                    </button>
                )}
                {member.role === 'member' && canChangeRole(currentUserRole, member.role, 'admin', member) && (
                    <button
                        onClick={() => handleRoleChange(member.id, 'admin')}
                        className="action-button promote"
                    >
                        Promote
                    </button>
                )}
                {/* Additional action buttons based on permissions */}
            </div>
        )
    }
];

<Table
    columns={columns}
    data={members}
    keyExtractor={(member) => member.id}
    loading={loading}
    error={error}
    emptyMessage="No members found"
    loadingMessage="Loading members..."
    errorMessage="Failed to load members"
    footer={
        members.length > 0 ? (
            <div>
                {t('clubs.totalMembers', { count: members.length }) || `Total: ${members.length} members`}
            </div>
        ) : null
    }
/>
```

## Best Practices

1. **Always provide a unique keyExtractor**
2. **Use TypeScript interfaces for type safety**
3. **Handle loading and error states appropriately**
4. **Consider responsive design when defining columns**
5. **Use semantic action button classes**
6. **Provide meaningful empty state messages**
7. **Test with screen readers for accessibility**
