import {credentials} from '@grpc/grpc-js';
import {z} from 'zod';
import {Stock, TradingClient} from '~/grpc/trading';
import {BoardClient, Question, Subject} from '~/grpc/board';
import {procedure, router} from '../trpc';
import * as Sentry from '@sentry/nextjs';

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
const board = new BoardClient(`${host}:${port}`, creds, opts);

export const appRouter = router({
    listSubjects: procedure.query(async (): Promise<Subject[]> => {
        const subjects: Promise<Subject[]> = new Promise((resolve, reject) => {
            board.listSubject({}, (err, subjectList) => {
                if (err) {
                    Sentry.captureException(err)
                    console.error(err);
                    reject(err);
                    return;
                }

                resolve(subjectList.subjectList);
            });
        });
        return subjects;
    }),

    getSubject: procedure.input(
        z.object({
            id: z.number(),
        })
    ).query(async ({input}): Promise<Subject> => {
        const subject: Promise<Subject> = new Promise((resolve, reject) => {
            board.getSubject(input, (err, subject) => {
                if (err) {
                    Sentry.captureException(err)
                    console.error(err);
                    reject(err);
                    return;
                }

                resolve(subject);
            });
        });
        return subject;
    }),

    listQuestions: procedure.input(
        z.object({
            id: z.number(),
        })
    ).query(async ({input}): Promise<Question[]> => {
        const questions: Promise<Question[]> = new Promise((resolve, reject) => {
            board.listQuestion(input, (err, questionList) => {
                if (err) {
                    Sentry.captureException(err)
                    console.error(err);
                    reject(err);
                    return;
                }

                resolve(questionList.questionList);
            });
        });
        return questions;
    }),

    createQuestion: procedure.input(
        z.object({
            question: z.string(),
            subjectId: z.number(),
        })
    ).mutation(async ({input: newQuestion}) => {
        new Promise((resolve, reject) => {
            board.createQuestion(newQuestion, (err, question) => {
                if (err) {
                    Sentry.captureException(err)
                    console.error(err);
                    reject(err);
                    return;
                }
            });
        });
    }),

    like: procedure.input(
        z.object({
            id: z.number(),
        })).mutation(async ({input}) => {
        new Promise((resolve, reject) => {
            board.like(input, (err, empty) => {
                if (err) {
                    Sentry.captureException(err)
                    console.error(err);
                    reject(err);
                    return;
                }
            });
        });
    }),

    getStockList: procedure.query(async (): Promise<Stock[]> => {
        const stocks: Promise<Stock[]> = new Promise((resolve, reject) => {
            trading.getStockList({}, (err, stockListResp) => {
                if (err) {
                    Sentry.captureException(err)
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
