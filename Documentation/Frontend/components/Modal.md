# Modal Component Documentation

The Modal component is a reusable, enhanced modal system based on the style from the Add Fine modal. It provides a consistent, accessible, and beautiful modal experience across the application.

## Features

- **Enhanced Design**: Beautiful header with gradient background, improved spacing, and smooth animations
- **Compound Component Pattern**: Flexible composition with Modal.Header, Modal.Body, Modal.Actions, etc.
- **Built-in Error Handling**: Dedicated Modal.Error component with animated error display
- **Loading States**: Built-in Modal.LoadingSpinner component
- **Accessibility**: Proper ARIA labels and keyboard support
- **Mobile Responsive**: Optimized for mobile devices
- **Backdrop Blur**: Modern backdrop-filter blur effect
- **Smooth Animations**: Fade-in and slide-in animations

## Basic Usage

```tsx
import { Modal } from '@/components/ui';

function MyComponent() {
  const [isOpen, setIsOpen] = useState(false);
  
  return (
    <Modal isOpen={isOpen} onClose={() => setIsOpen(false)} title="My Modal">
      <Modal.Body>
        <p>Modal content goes here</p>
      </Modal.Body>
      <Modal.Actions>
        <button onClick={() => setIsOpen(false)} className="button-cancel">
          Cancel
        </button>
        <button className="button-accept">
          Confirm
        </button>
      </Modal.Actions>
    </Modal>
  );
}
```

## Advanced Usage with Error Handling and Loading

```tsx
import { Modal } from '@/components/ui';

function AdvancedModal() {
  const [isOpen, setIsOpen] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  
  const handleSubmit = async () => {
    setIsLoading(true);
    try {
      // API call here
      await submitData();
      setIsOpen(false);
    } catch (err) {
      setError('Failed to submit data');
    } finally {
      setIsLoading(false);
    }
  };
  
  return (
    <Modal 
      isOpen={isOpen} 
      onClose={() => setIsOpen(false)} 
      title="Advanced Modal"
      maxWidth="600px"
    >
      <Modal.Error error={error} />
      
      <Modal.Body>
        <div className="modal-form-section">
          {/* Form inputs */}
        </div>
      </Modal.Body>
      
      <Modal.Actions>
        <button 
          onClick={() => setIsOpen(false)} 
          className="button-cancel"
          disabled={isLoading}
        >
          Cancel
        </button>
        <button 
          onClick={handleSubmit}
          className="button-accept"
          disabled={isLoading}
        >
          {isLoading ? (
            <>
              <Modal.LoadingSpinner />
              Submitting...
            </>
          ) : (
            'Submit'
          )}
        </button>
      </Modal.Actions>
    </Modal>
  );
}
```

## Props

### Modal Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| isOpen | boolean | - | Controls modal visibility |
| onClose | () => void | - | Called when modal should close |
| title | string | - | Modal title displayed in header |
| children | ReactNode | - | Modal content |
| maxWidth | string | "550px" | Maximum width of the modal |
| showCloseButton | boolean | true | Whether to show the X close button |
| className | string | "" | Additional CSS classes |

### Modal.Error Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| error | string \| null | - | Error message to display |
| onClose | () => void | - | Optional close handler for error |

## Styling

The modal uses CSS custom properties (CSS variables) for theming:

- `--space-*`: Spacing values
- `--color-*`: Color values
- `--border-radius-*`: Border radius values

## Migration from Old Modals

### Before
```tsx
<div className="modal" onClick={onClose}>
  <div className="modal-content" onClick={(e) => e.stopPropagation()}>
    <h2>Title</h2>
    {error && <p style={{ color: 'red' }}>{error}</p>}
    <p>Content</p>
    <div className="modal-actions">
      <button onClick={onClose}>Cancel</button>
      <button>Confirm</button>
    </div>
  </div>
</div>
```

### After
```tsx
<Modal isOpen={isOpen} onClose={onClose} title="Title">
  <Modal.Error error={error} />
  <Modal.Body>
    <p>Content</p>
  </Modal.Body>
  <Modal.Actions>
    <button onClick={onClose}>Cancel</button>
    <button>Confirm</button>
  </Modal.Actions>
</Modal>
```

## Accessibility

The modal includes:
- Proper ARIA labels
- Focus management
- Keyboard support (ESC to close)
- Screen reader friendly structure

## Examples in the App

The following components have been updated to use the new Modal component:
- AddFine modal
- ClubDetails leave confirmation modal
- AdminClubFineList template modal  
- AddNews modal

Other modals throughout the app can be similarly migrated to use this component for a consistent user experience.
