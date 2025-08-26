import { createBrowserRouter } from 'react-router-dom';
import routes, { AppRouter } from './router';

function App() {
    const router = createBrowserRouter(routes);
    return <AppRouter router={router} />;
}

export default App;
