import Link from 'next/link';
import {trpc} from '~/utils/trpc'
import * as Sentry from '@sentry/nextjs';

export default function Home() {
    const {data} = trpc.listSubjects.useQuery();

    return (
        <main>
            <h1>Subject List</h1>
            <ul>
                {data && data.map(subject => {
                    let subjectEnabled: string;
                    if (subject.enabled) {
                        subjectEnabled = 'enabled';
                    } else {
                        subjectEnabled = "disabled";
                    }

                    return (
                        <li key="{subject.id}">
                            <Link href={`/${subject.id}`}>{subject.title}</Link>
                            <b>&nbsp;{` (${subjectEnabled})`}</b>
                        </li>
                    );
                })}
            </ul>
        </main>
    )
}
