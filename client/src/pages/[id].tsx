import Link from "next/link";
import { useRouter } from 'next/router'
import { trpc } from '~/utils/trpc'

export default function Home() {
    const router = useRouter();

    const subject = trpc.getSubject.useQuery({ id: parseInt(router.query.id) });
    const list = trpc.listQuestions.useQuery();

    console.log(subject.data);

    return (
        <main>
            <h1>Question List</h1>
            <ul>
                {list.data && list.data.map(q => {
                    return (<li key="{q.id}">{q.question}</li>);
                })}
            </ul>
        </main>
    )
}
