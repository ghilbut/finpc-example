import { credentials } from '@grpc/grpc-js';
import { z } from 'zod';
import { Stock, TradingClient } from '~/grpc/trading';
import { procedure, router } from '../trpc';

const host = process.env.GRPC_HOST || '127.0.0.1';
const port = process.env.GRPC_PORT || '9095';
const creds = process.env.GRPC_INSECURE === 'true' ?
    credentials.createInsecure() :
    credentials.createSsl(Buffer.from(process.env.GRPC_CACERT || '', 'base64'));
const opts = process.env.GRPC_HOST_OVERRIDE ? {
    'grpc.ssl_target_name_override': process.env.GRPC_HOST_OVERRIDE,
    'grpc.default_authority': process.env.GRPC_HOST_OVERRIDE
} : undefined;

console.log('GRPC_HOST: ', host);
console.log('GRPC_PORT: ', port);
console.log('GRPC_INSECURE: ', process.env.GRPC_INSECURE || 'false');
console.log('GRPC_HOST_OVERRIDE: ', process.env.GRPC_HOST_OVERRIDE || '');
console.log('GRPC_OPTIONS: ', opts);

const trading = new TradingClient(`${host}:${port}`, creds, opts);

export const appRouter = router({
    listSubjects: procedure.query(() => {
        return [
            { id: 1, title: 'subject #1' },
            { id: 2, title: 'subject #2' },
            { id: 3, title: 'subject #3' },
            { id: 4, title: 'subject #4' },
        ];
    }),

    getSubject: procedure.input(
        z.object({
            id: z.number(),
        })
    ).query(async ({ input }) => {
        return { id: input.id, title: `subject #1` };
    }),

    createSubject: procedure.input(z.object({
        id: z.number(),
        title: z.string(),
    })).mutation(async ({ input }) => {

    }),

    deleteSubject: procedure.input(z.object({
        id: z.number(),
    })).mutation(async ({ input }) => {

    }),

    listQuestions: procedure.query(async ({ input }) => {
        return [
            { id: 1, question: 'question #1' },
            { id: 2, question: 'question #2' },
            { id: 3, question: 'question #3' },
            { id: 4, question: 'question #4' },
        ];
    }),

    createQuestion: procedure.input(z.object({
        id: z.number(),
        title: z.string(),
    })).mutation(async ({ input }) => {

    }),

    deleteQuestion: procedure.input(z.object({
        id: z.number(),
    })).mutation(async ({ input }) => {

    }),

    like: procedure.input(z.object({
        id: z.number(),
    })).mutation(async ({ input }) => {

    }),

    unlike: procedure.input(z.object({
        id: z.number(),
    })).mutation(async ({ input }) => {

    }),

    getStockList: procedure.query(async (): Promise<Stock[]> => {
        const stocks: Promise<Stock[]> = new Promise((resolve, reject) => {
            trading.getStockList({}, (err, stockListResp) => {
                if (err) {
                    console.error(err);
                    reject(err);
                    return;
                }

                resolve(stockListResp.stockList);
            });
        });
        return stocks;
    }),
});

// export type definition of API
export type AppRouter = typeof appRouter;
