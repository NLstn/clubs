# ODataTable Component

A table component with **server-side pagination and sorting** via OData v4 queries. This component extends the base Table component to automatically handle OData query construction and response parsing.

## Import

```tsx
import { ODataTable, ODataTableColumn } from '@/components/ui';
```

## Basic Usage

```tsx
interface News {
  ID: string;
  Title: string;
  Content: string;
  CreatedAt: string;
}

const columns: ODataTableColumn<News>[] = [
  {
    key: 'Title',
    header: 'Title',
    render: (item) => item.Title,
    sortable: true,
  },
  {
    key: 'Content',
    header: 'Content',
    render: (item) => item.Content,
  },
  {
    key: 'CreatedAt',
    header: 'Created',
    render: (item) => new Date(item.CreatedAt).toLocaleString(),
    sortable: true,
    sortField: 'CreatedAt',
  },
];

<ODataTable
  endpoint="/api/v2/News"
  filter="ClubID eq 'abc-123'"
  columns={columns}
  keyExtractor={(item) => item.ID}
  pageSize={10}
  initialSortField="CreatedAt"
  initialSortDirection="desc"
/>
```

## Props

### ODataTableProps<T>

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `endpoint` | `string` | Required | OData entity set path (e.g., `/api/v2/News`) |
| `columns` | `ODataTableColumn<T>[]` | Required | Column definitions with optional sorting |
| `keyExtractor` | `(item: T) => string` | Required | Function to extract unique key from each item |
| `filter` | `string` | `undefined` | OData filter query (e.g., `"ClubID eq 'abc'"`) |
| `expand` | `string \| string[]` | `undefined` | OData expand clause for related entities |
| `select` | `string[]` | `undefined` | OData select clause for specific fields |
| `pageSize` | `number` | `10` | Number of items per page |
| `className` | `string` | `''` | Additional CSS classes |
| `emptyMessage` | `string` | `'No data available'` | Message when no data |
| `loadingMessage` | `string` | `'Loading...'` | Loading state message |
| `errorMessage` | `string` | `'Error loading data'` | Error state message |
| `initialSortField` | `string` | `undefined` | Initial field to sort by |
| `initialSortDirection` | `'asc' \| 'desc'` | `'desc'` | Initial sort direction |

### ODataTableColumn<T>

Extends `TableColumn<T>` with additional properties:

| Prop | Type | Description |
|------|------|-------------|
| `key` | `string` | Column identifier |
| `header` | `string` | Column header text |
| `render` | `(item: T) => React.ReactNode` | Cell render function |
| `className` | `string` | Optional column CSS class |
| `sortable` | `boolean` | Enable sorting for this column |
| `sortField` | `string` | OData field name for sorting (if different from key) |

## Features

### Server-Side Pagination

The component automatically handles pagination using OData `$skip` and `$top` query parameters:

- Shows current page range (e.g., "Showing 1-10 of 45")
- Displays "Page X of Y" indicator
- Provides first/previous/next/last navigation buttons
- Fetches only the data needed for the current page

### Server-Side Sorting

Sortable columns display an interactive header:

- Click header to toggle between ascending/descending
- Visual indicator (▲/▼) shows current sort direction
- Uses OData `$orderby` parameter
- Resets to page 1 when sort changes

### OData Query Construction

The component automatically builds OData queries with:

- `$skip` - For pagination offset
- `$top` - For page size
- `$orderby` - For sorting (field + direction)
- `$filter` - From the filter prop
- `$expand` - From the expand prop
- `$select` - From the select prop
- `$count=true` - Always included to get total count

### Response Parsing

Automatically parses OData collection responses:

```typescript
{
  "@odata.context": "...",
  "@odata.count": 45,
  "value": [...] // Array of entities
}
```

Extracts the `value` array and `@odata.count` for pagination.

## Advanced Examples

### With Related Entity Expansion

```tsx
<ODataTable
  endpoint="/api/v2/Members"
  filter="ClubID eq 'abc-123'"
  expand="User"
  columns={[
    {
      key: 'Name',
      header: 'Name',
      render: (member) => member.User?.Name || 'N/A',
      sortable: true,
    },
    // ... more columns
  ]}
  keyExtractor={(item) => item.ID}
/>
```

### With Field Selection

```tsx
<ODataTable
  endpoint="/api/v2/News"
  filter="ClubID eq 'abc-123'"
  select={['ID', 'Title', 'CreatedAt']}
  columns={columns}
  keyExtractor={(item) => item.ID}
/>
```

### With Action Columns

```tsx
const columns: ODataTableColumn<News>[] = [
  {
    key: 'Title',
    header: 'Title',
    render: (item) => item.Title,
    sortable: true,
  },
  {
    key: 'actions',
    header: 'Actions',
    render: (item) => (
      <div>
        <Button
          variant="accept"
          size="sm"
          onClick={() => handleEdit(item.ID)}
        >
          Edit
        </Button>
        <Button
          variant="cancel"
          size="sm"
          onClick={() => handleDelete(item.ID)}
        >
          Delete
        </Button>
      </div>
    ),
  },
];
```

### With Custom Sorting Field

When the display field differs from the sortable field:

```tsx
{
  key: 'userName',
  header: 'User',
  render: (member) => member.User?.Name || 'Unknown',
  sortable: true,
  sortField: 'User/Name', // OData navigation path
}
```

### Refreshing Data

Use a key prop to force re-fetch:

```tsx
const [refreshKey, setRefreshKey] = useState(0);

const refreshData = () => {
  setRefreshKey(prev => prev + 1);
};

<ODataTable
  key={refreshKey}
  endpoint="/api/v2/News"
  filter="ClubID eq 'abc-123'"
  columns={columns}
  keyExtractor={(item) => item.ID}
/>

<Button onClick={refreshData}>Refresh</Button>
```

## Styling

### CSS Classes

- `.odata-table-wrapper` - Main container
- `.odata-table-sort-header` - Sortable header button
- `.sort-indicator` - Sort direction indicator (▲/▼)
- `.odata-table-pagination` - Pagination container
- `.pagination-info` - "Showing X-Y of Z" text
- `.pagination-controls` - Navigation buttons container
- `.page-indicator` - "Page X of Y" text

### Custom Styling

```tsx
<ODataTable
  className="my-custom-table"
  // ... other props
/>
```

```css
.my-custom-table .odata-table-sort-header {
  color: var(--custom-color);
}

.my-custom-table .pagination-controls button {
  background: var(--custom-bg);
}
```

## Accessibility

- Sortable headers are keyboard accessible buttons
- Buttons have proper ARIA labels
- Disabled navigation buttons are properly marked
- Screen reader friendly pagination info

## Performance Considerations

- Only fetches data for the current page (efficient for large datasets)
- Sorting and filtering done on the server (no client-side processing)
- Use `select` prop to limit fields returned by the API
- Consider adding indexes to frequently sorted/filtered database columns

## OData v2 Conventions

⚠️ **Important**: All OData v2 endpoints use **PascalCase** for field names.

```tsx
// Correct ✅
interface News {
  ID: string;
  Title: string;
  CreatedAt: string;
}

// Wrong ❌
interface News {
  id: string;
  title: string;
  createdAt: string;
}
```

## Example in the Codebase

See `src/pages/clubs/admin/news/AdminClubNewsList.tsx` for a complete working example:

- Fetches news posts for a club
- Sortable by Title, CreatedAt, UpdatedAt
- Paginated with 10 items per page
- Includes edit/delete action buttons
- Refreshes after add/edit/delete operations

## Comparison with Regular Table

| Feature | Table | ODataTable |
|---------|-------|-----------|
| Data Source | Client-side array | Server-side OData API |
| Pagination | Manual/external | Built-in, automatic |
| Sorting | Manual/external | Built-in, clickable headers |
| Large Datasets | Load all data | Load only current page |
| Filtering | Manual/external | Via OData $filter |
| Related Data | Manual joins | Via OData $expand |

## Best Practices

1. **Use for large datasets** - When you have more than 50-100 records
2. **Set appropriate page sizes** - 10-25 items is usually optimal
3. **Make commonly sorted columns sortable** - Dates, names, titles
4. **Use filters effectively** - Reduce server load with specific filters
5. **Expand sparingly** - Only expand relations you actually display
6. **Select fields** - If you only need a few fields, use `select`
7. **Handle errors gracefully** - ODataTable shows error states automatically
8. **Test with slow connections** - Loading states should be clear

## Troubleshooting

### No data showing

- Check browser console for API errors
- Verify the `endpoint` path is correct
- Check `filter` syntax (use OData filter syntax)
- Ensure backend OData service is running

### Sorting not working

- Verify column has `sortable: true`
- Check `sortField` matches backend field name (case-sensitive)
- Ensure field is actually sortable in backend

### Pagination incorrect

- Verify `@odata.count` is included in API response
- Check if backend supports `$count=true`
- Look for errors in browser console

### Type errors

- Ensure your interface matches backend PascalCase
- Check `keyExtractor` returns a string
- Verify column `render` functions return ReactNode
