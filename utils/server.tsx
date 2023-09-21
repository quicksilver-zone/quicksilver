import { type Serve } from 'bun';

export default {
  port: 8080,
  hostname: '0.0.0.0',
  fetch(req) {
    console.log('request made');
    return new Response('Bun!');
  },
} satisfies Serve;
