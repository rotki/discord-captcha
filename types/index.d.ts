declare module '#app' {
  import { type Client } from '@discordjs/core';

  interface NuxtApp {
    $discord: Client;
  }
}
