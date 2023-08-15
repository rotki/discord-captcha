# @rotki/discord-captcha

## Configuration

### Create and configure

Create a new [discord application](https://discord.com/developers/applications). 

In the bot settings of your new application make sure to go under `Priviledged Gateway Intents`
and enable `SERVER_MEMBERS_INTENT` and `MESSAGE_CONTENT_INTENT`.

### Permissions

Make sure the bot has the following permissions:
- Manage Server
- Manage Roles
- Manage Channels
- Create Instant Invite
- Create Expressions
- Mention Everyone

### Add the bot to your server

Make sure that your bot is not public.
To invite the bot to your server you can use the following link after replacing the placeholders:

```
https://discord.com/oauth2/authorize?client_id={appId}&permissions={permissions}&scope=b`
```

For the permissions use the permission calculator in the application bot settings.

### Role Addition

You should ensure that the added role is below the bot's role in the Guild's 
roles settings.

Otherwise, this might cause permission errors even if you have given the `MANAGE_ROLES` 
permission to the bot.


## Development
Look at the [Nuxt 3 documentation](https://nuxt.com/docs/getting-started/introduction) to learn more.

Make sure to install the dependencies:

```bash
# npm
npm install

# pnpm
pnpm install

# yarn
yarn install
```

## Development Server

Start the development server on `http://localhost:3000`:

```bash
# npm
npm run dev

# pnpm
pnpm run dev

# yarn
yarn dev
```

## Production

Build the application for production:

```bash
# npm
npm run build

# pnpm
pnpm run build

# yarn
yarn build
```

Locally preview production build:

```bash
# npm
npm run preview

# pnpm
pnpm run preview

# yarn
yarn preview
```

Check out the [deployment documentation](https://nuxt.com/docs/getting-started/deployment) for more information.
