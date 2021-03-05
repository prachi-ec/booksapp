/* eslint-disable no-console,@typescript-eslint/no-var-requires */
const express = require('express');
const next = require('next');
const { createProxyMiddleware } = require('http-proxy-middleware');

const devProxy = {
  '/api': {
    target: 'http://localhost:8088/',
    pathRewrite: { '^/api': '/' },
    changeOrigin: true,
  },
};

const port = parseInt(process.env.PORT, 10) || 3000;
const dev = true;
const app = next({
  dir: '.', // base directory where everything is, could move to src later
  dev,
});

const handle = app.getRequestHandler();

const server = express();

app
  .prepare()
  .then(() => {
    // Set up the proxy.
    Object.keys(devProxy).forEach(function (context) {
      server.use(createProxyMiddleware(context, devProxy[context]));
    });

    // Default catch-all handler to allow Next.js to handle all other routes
    server.all('*', (req, res) => handle(req, res));

    server.listen(port, err => {
      if (err) {
        console.log(err);
        throw err;
      }
      console.log(`> dev server Ready on port ${port}`);
    });
  })
  .catch(err => {
    console.log('An error occurred, unable to start the server');
    console.log(err);
  });