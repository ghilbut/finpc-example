import * as Sentry from '@sentry/nextjs';

Sentry.init({
    dsn: "https://8bfb9ce2ae9f33673b99ea874f755e98@o4505685110423552.ingest.sentry.io/4505685292613632",
    integrations: [
        new Sentry.BrowserTracing(),
        new Sentry.Replay(),
    ],
    // Performance Monitoring
    tracesSampleRate: 1.0, // Capture 100% of the transactions, reduce in production!
    // Session Replay
    replaysSessionSampleRate: 0.1, // This sets the sample rate at 10%. You may want to change it to 100% while in development and then sample at a lower rate in production.
    replaysOnErrorSampleRate: 1.0, // If you're not already sampling the entire session, change the sample rate to 100% when sampling sessions where errors occur.
});