import { useState } from 'react';
import Layout from '../components/layout/Layout';
import { Button, ButtonState } from '@/components/ui';
import './ButtonDemo.css';

/**
 * Demo page to showcase the enhanced Button component with feedback states.
 * This page demonstrates all button states: idle, loading, success, and error.
 */
const ButtonDemo = () => {
  const [state1, setState1] = useState<ButtonState>('idle');
  const [state2, setState2] = useState<ButtonState>('idle');
  const [state3, setState3] = useState<ButtonState>('idle');
  const [cancelledCount, setCancelledCount] = useState(0);

  // Simulate async operation for state1
  const handleSave = async () => {
    setState1('loading');
    await new Promise(resolve => setTimeout(resolve, 2000));
    setState1('success');
    setTimeout(() => setState1('idle'), 3000);
  };

  // Simulate async operation with error for state2
  const handleSaveWithError = async () => {
    setState2('loading');
    await new Promise(resolve => setTimeout(resolve, 2000));
    setState2('error');
    setTimeout(() => setState2('idle'), 3000);
  };

  // Simulate cancellable operation for state3
  const handleCancellableOperation = async () => {
    setState3('loading');
    
    // Simulate a long operation
    setTimeout(() => {
      setState3('success');
      setTimeout(() => setState3('idle'), 3000);
    }, 5000);
  };

  const handleCancel = () => {
    setState3('idle');
    setCancelledCount(prev => prev + 1);
  };

  return (
    <Layout title="Button Component Demo">
      <div className="button-demo">
        <h2>Enhanced Button Component Demo</h2>
        <p className="demo-description">
          This page demonstrates the new feedback capabilities of the Button component,
          including loading states with spinners, success/error feedback, and cancellation.
        </p>

        <section className="demo-section">
          <h3>1. Success Flow</h3>
          <p>Click the button to simulate a successful save operation.</p>
          <Button
            variant="primary"
            state={state1}
            successMessage="Saved successfully!"
            onClick={handleSave}
          >
            Save Changes
          </Button>
          <div className="state-info">Current state: <strong>{state1}</strong></div>
        </section>

        <section className="demo-section">
          <h3>2. Error Flow</h3>
          <p>Click the button to simulate a failed operation with error feedback.</p>
          <Button
            variant="accept"
            state={state2}
            errorMessage="Save failed!"
            onClick={handleSaveWithError}
          >
            Save (Will Fail)
          </Button>
          <div className="state-info">Current state: <strong>{state2}</strong></div>
        </section>

        <section className="demo-section">
          <h3>3. Cancellable Operation</h3>
          <p>Click the button to start a long operation. Click the X button to cancel it.</p>
          <Button
            variant="secondary"
            state={state3}
            successMessage="Completed!"
            onCancel={handleCancel}
            onClick={handleCancellableOperation}
          >
            Start Long Process
          </Button>
          <div className="state-info">
            Current state: <strong>{state3}</strong>
            {cancelledCount > 0 && <span> | Cancelled: {cancelledCount} times</span>}
          </div>
        </section>

        <section className="demo-section">
          <h3>4. Button Variants with Different States</h3>
          <div className="button-grid">
            <div className="button-showcase">
              <h4>Primary Variants</h4>
              <Button variant="primary">Idle</Button>
              <Button variant="primary" state="loading">Loading</Button>
              <Button variant="primary" state="success" successMessage="Success!">Success</Button>
              <Button variant="primary" state="error" errorMessage="Error!">Error</Button>
            </div>
            
            <div className="button-showcase">
              <h4>Accept Variants</h4>
              <Button variant="accept">Idle</Button>
              <Button variant="accept" state="loading">Loading</Button>
              <Button variant="accept" state="success" successMessage="Done!">Success</Button>
              <Button variant="accept" state="error" errorMessage="Failed!">Error</Button>
            </div>
            
            <div className="button-showcase">
              <h4>Cancel Variants</h4>
              <Button variant="cancel">Idle</Button>
              <Button variant="cancel" state="loading">Loading</Button>
              <Button variant="cancel" state="success" successMessage="Cancelled!">Success</Button>
              <Button variant="cancel" state="error" errorMessage="Error!">Error</Button>
            </div>
          </div>
        </section>

        <section className="demo-section">
          <h3>5. Different Sizes</h3>
          <div className="button-sizes">
            <Button size="sm" state="loading">Small</Button>
            <Button size="md" state="loading">Medium</Button>
            <Button size="lg" state="loading">Large</Button>
          </div>
        </section>

        <section className="demo-section">
          <h3>6. Full Width Buttons</h3>
          <Button fullWidth variant="primary" state="idle">Full Width Idle</Button>
          <br /><br />
          <Button fullWidth variant="accept" state="loading">Full Width Loading</Button>
          <br /><br />
          <Button fullWidth variant="secondary" state="success" successMessage="Saved!">Full Width Success</Button>
        </section>

        <section className="demo-section api-section">
          <h3>API Reference</h3>
          <div className="api-table">
            <table>
              <thead>
                <tr>
                  <th>Prop</th>
                  <th>Type</th>
                  <th>Description</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td><code>state</code></td>
                  <td><code>'idle' | 'loading' | 'success' | 'error'</code></td>
                  <td>Current state of the button (default: 'idle')</td>
                </tr>
                <tr>
                  <td><code>successMessage</code></td>
                  <td><code>string</code></td>
                  <td>Message to display when state is 'success'</td>
                </tr>
                <tr>
                  <td><code>errorMessage</code></td>
                  <td><code>string</code></td>
                  <td>Message to display when state is 'error'</td>
                </tr>
                <tr>
                  <td><code>onCancel</code></td>
                  <td><code>() =&gt; void</code></td>
                  <td>Callback when cancel button is clicked during loading</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>
      </div>
    </Layout>
  );
};

export default ButtonDemo;
