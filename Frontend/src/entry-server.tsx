import { renderToString } from 'react-dom/server';
import { createStaticHandler, createStaticRouter } from 'react-router-dom';
import routes, { AppRouter } from './router';

export async function render(url: string) {
    const handler = createStaticHandler(routes);
    const request = new Request('http://localhost' + url);
    const context = await handler.query(request);
    const router = createStaticRouter(handler.dataRoutes, context);
    const html = renderToString(<AppRouter router={router} />);
    return { html, state: router.state };
}
