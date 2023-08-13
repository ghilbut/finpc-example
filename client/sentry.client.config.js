import * as Sentry from '@sentry/nextjs';

Sentry.init({
    dsn: "https://21225647aa80316b11f1d5c5d97666a1@o4505689921683456.ingest.sentry.io/4505690191888384",
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