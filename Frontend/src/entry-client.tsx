import { hydrateRoot } from 'react-dom/client';
import { createBrowserRouter } from 'react-router-dom';
import routes, { AppRouter } from './router';
import './index.css';
import './i18n/index.ts';

hydrateRoot(
    document.getElementById('root')!,
    <AppRouter router={createBrowserRouter(routes)} />
);
