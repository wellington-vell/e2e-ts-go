import { createClient } from 'authula';
import {
  AccessControlPlugin,
  AdminPlugin,
  CSRFPlugin,
  EmailPasswordPlugin,
} from 'authula/plugins';

import { env } from '@/lib/env';

export const authClient = createClient({
  url: env.VITE_SERVER_URL + '/auth',
  plugins: [
    new EmailPasswordPlugin(),
    new AdminPlugin(),
    new AccessControlPlugin(),
    new CSRFPlugin({
      cookieName: 'csrf_token',
      headerName: 'X-CSRF-TOKEN',
    }),
  ],
});
