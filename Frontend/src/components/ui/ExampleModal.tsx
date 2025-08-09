import { FC, useState } from 'react';
import { Modal, Input, TypeAheadDropdown } from '@/components/ui';

interface Option {
  id: string;
  label: string;
}

// Example usage of the Modal component for reference
interface ExampleModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const ExampleModal: FC<ExampleModalProps> = ({ isOpen, onClose }) => {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [selectedOption, setSelectedOption] = useState<Option | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const options: Option[] = [
    { id: '1', label: 'Option 1' },
    { id: '2', label: 'Option 2' },
    { id: '3', label: 'Option 3' },
  ];

  const handleSubmit = async () => {
    setError(null);
    
    if (!name || !email) {
      setError('Please fill in all required fields');
      return;
    }

    setIsSubmitting(true);
    
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      // Reset form
      setName('');
      setEmail('');
      setSelectedOption(null);
      onClose();
    } catch {
      setError('Something went wrong. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    setName('');
    setEmail('');
    setSelectedOption(null);
    setError(null);
    onClose();
  };

  return (
    <Modal 
      isOpen={isOpen} 
      onClose={handleClose} 
      title="Example Modal"
      maxWidth="500px"
    >
      <Modal.Error error={error} />
      
      <Modal.Body>
        <div className="modal-form-section">
          <Input
            label="Name"
            id="name"
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Enter your name"
            disabled={isSubmitting}
          />

          <Input
            label="Email"
            id="email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="Enter your email"
            disabled={isSubmitting}
          />

          <TypeAheadDropdown<Option>
            options={options}
            value={selectedOption}
            onChange={setSelectedOption}
            onSearch={(query) => {
              // This triggers when user types in the search box
              // The filtering is handled internally by the component
              console.log('Searching for:', query);
            }}
            placeholder="Select an option..."
            id="option"
            label="Option (Optional)"
          />
        </div>
      </Modal.Body>
      
      <Modal.Actions>
        <button 
          onClick={handleClose} 
          className="button-cancel"
          disabled={isSubmitting}
        >
          Cancel
        </button>
        <button 
          onClick={handleSubmit}
          className="button-accept"
          disabled={isSubmitting || !name || !email}
        >
          {isSubmitting ? (
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
};

export default ExampleModal;
