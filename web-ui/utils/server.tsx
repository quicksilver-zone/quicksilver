//@ts-ignore
import { type Serve } from 'bun';

export default {
  port: 8080,
  hostname: '0.0.0.0',
  //@ts-ignore
  fetch(req) {
    console.log('request made');
    return new Response('Bun!');
  },
} satisfies Serve;
